package handler

import (
	"context"
	"fmt"
)

// AddUserRequest example
type AddUserRequest struct {
	Name       string `json:"name" example:"Peter"`
	Surname    string `json:"surname" example:"Dibin"`
	Patronymic string `json:"patronymic" example:"Ivanovich"`
}

// AddUserResponse example
type AddUserResponse struct {
	Id int `json:"user" example:"1"`
}

// AddUserHandler godoc
//
//	@Summary		Add user
//	@Description	add new user with name, surname and patronymic
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param 			request body 	AddUserRequest true "query params"
//	@Success		200	{object}	AddUserResponse
//	@Router			/addUser [post]
func (h Handler) AddUserHandler(ctx context.Context, req AddUserRequest) (*AddUserResponse, error) {
	id, err := h.service.AddUser(ctx, toModelUser(req))
	if err != nil {
		return nil, fmt.Errorf("cannot add user: %w", err)
	}
	return &AddUserResponse{Id: id}, nil
}
