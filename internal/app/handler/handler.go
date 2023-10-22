package handler

import (
	"context"
	"effectivemobile/internal/app/model"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
)

type IService interface {
	AddUser(ctx context.Context, user model.UserInfo) (int, error)
	DeleteUser(ctx context.Context, id int) error
	UpdateUser(ctx context.Context, id int, info model.UserInfo) error
}

type Handler struct {
	service IService
}

func New(srv IService) *Handler {

	return &Handler{service: srv}
}

func (h Handler) RegisterHandlers(router *mux.Router, mw ...func(next http.HandlerFunc) http.HandlerFunc) {
	for _, rec := range [...]struct {
		route   string
		handler http.HandlerFunc
	}{
		{route: "/swagger.json", handler: func(w http.ResponseWriter, r *http.Request) {
			cwd, _ := os.Getwd()
			http.ServeFile(w, r, path.Join(cwd, "docs/swagger.json"))
		}},
		{route: "/swagger/{any:.+}", handler: httpSwagger.Handler(httpSwagger.URL("/swagger.json"))},
		{route: "/addUser", handler: handle(h.AddUserHandler)},
		{route: "/deleteUser", handler: handle(h.DeleteUserHandler)},
		{route: "/updateUser", handler: handle(h.UpdateUserHandler)},
	} {
		router.HandleFunc(rec.route, middlewareChain(rec.handler, mw...))
	}
}

var validate = validator.New()

func handle[REQ any, RESP any](fn func(ctx context.Context, req REQ) (*RESP, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", "*")
		headers.Set("Access-Control-Allow-Methods", "POST")
		headers.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Content-Length, Accept")

		if r.Method == http.MethodOptions {
			return
		}

		var req REQ
		headers.Set("Content-Type", "application/json")
		if err := parsePathParams(r, &req); err != nil {
			sendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				sendErrorResponse(w, err, http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(body, &req)
			if err != nil {
				sendErrorResponse(w, err, http.StatusBadRequest)
				return
			}
		}

		if err := validate.Struct(req); err != nil {
			sendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		resp, err := fn(r.Context(), req)
		if err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		respJson, err := json.Marshal(resp)
		if err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}

		headers.Set("Content-Length", strconv.Itoa(len(respJson)))
		_, err = w.Write(respJson)
		if err != nil {
			sendErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
	}
}

func sendErrorResponse(w http.ResponseWriter, respErr error, respCode int) {
	respJson, err := json.Marshal(errorResponse{Error: respErr.Error()})
	if err != nil {
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(respJson)))
	w.WriteHeader(respCode)
	_, _ = w.Write(respJson)
	log.Println(respErr)
}

type emptyRequest struct{}
type emptyResponse struct{}
type errorResponse struct {
	Error string `json:"error"`
}

func parsePathParams[REQ any](r *http.Request, req *REQ) error {

	setField := func(field reflect.StructField, val reflect.Value, value string) error {
		switch field.Type.Kind() {
		case reflect.String:
			val.SetString(value)
			return nil
		case reflect.Int:
			v, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			val.SetInt(int64(v))
			return nil
		default:
			return fmt.Errorf("Unsupported path value type %v\n", field.Type.Kind())
		}
	}

	pathParams := mux.Vars(r)
	queryParams := r.URL.Query()
	val := reflect.ValueOf(req).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		name, ok := field.Tag.Lookup("path")
		if ok {
			if value, ok := pathParams[name]; ok {
				decoded, err := url.PathUnescape(value)
				if err != nil {
					return fmt.Errorf("cannot decode special symbol: %w", err)
				}
				if err := setField(field, val.Field(i), decoded); err != nil {
					return err
				}
			}
		}
		name, ok = field.Tag.Lookup("query")
		if ok {
			if value, ok := queryParams[name]; ok {
				if len(value) == 0 {
					return fmt.Errorf("query parameters %s is empty", name)
				}
				decoded, err := url.QueryUnescape(value[0])
				if err != nil {
					return fmt.Errorf("cannot decode special symbol: %w", err)
				}
				if err := setField(field, val.Field(i), decoded); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func middlewareChain(h http.HandlerFunc, m ...func(next http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	if len(m) == 0 {
		return h
	}
	return m[0](middlewareChain(h, m[1:cap(m)]...))
}
