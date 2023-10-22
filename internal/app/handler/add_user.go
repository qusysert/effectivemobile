package handler

import (
	"context"
	"fmt"
)

type AddUserRequest struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type AddUserResponse struct {
	Id int `json:"user"`
}

func (h Handler) AddUserHandler(ctx context.Context, req AddUserRequest) (*AddUserResponse, error) {
	id, err := h.service.AddUser(ctx, toModelUser(req))
	if err != nil {
		return nil, fmt.Errorf("cannot add user: %w", err)
	}
	return &AddUserResponse{Id: id}, nil
}
