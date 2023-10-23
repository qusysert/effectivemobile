package main

import (
	_ "effectivemobile/docs"
	"effectivemobile/internal/app/handler"
	"effectivemobile/internal/app/repository"
	"effectivemobile/internal/app/service"
	"effectivemobile/internal/pkg/config"
	"effectivemobile/internal/pkg/enrichclient"
	db "effectivemobile/pkg/gopkg-db"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

var conn db.IClient

//	@title          Backend Trainee Assignment 2023
//	@version		1.0
//	@description	Swagger documentation fo Backend Trainee Assignment 2023 service

// @contact.name	Ivan Demchuk
// @contact.email	is.demchuk@gmail.com
// @host		localhost:8080
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load cfg")
	}

	// Creating connection to DB
	conn, err = db.New(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal(fmt.Errorf("cant create connection to db: %v", err))
	}

	repo := repository.New()

	enrichClient := enrichclient.New(cfg.AgifyHost, cfg.GenderizeHost, cfg.NationalizeHost)

	srv := service.New(cfg, repo, enrichClient)
	hdl := handler.New(srv)

	corsMw := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})

	router := mux.NewRouter()

	// Setting timeout for the server
	server := &http.Server{
		Addr:         "0.0.0.0:" + cfg.HttpPort,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      handlers.LoggingHandler(os.Stdout, corsMw.Handler(handlers.CompressHandler(handlers.RecoveryHandler()(router)))),
	}

	hdl.RegisterHandlers(router, DbMiddleware)

	http.Handle("/", router)

	log.Printf("Server started on port %s \n", cfg.HttpPort)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}

func DbMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(db.AddToContext(ctx, conn))
		next.ServeHTTP(w, r)
	}
}
