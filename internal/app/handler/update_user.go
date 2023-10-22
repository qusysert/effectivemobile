package handler

import (
	"context"
	"fmt"
)

type UpdateUserRequest struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func (h Handler) UpdateUserHandler(ctx context.Context, req UpdateUserRequest) (*emptyResponse, error) {
	err := h.service.UpdateUser(ctx, req.Id, toModelUser(req))
	if err != nil {
		return nil, fmt.Errorf("cannot add segment: %w", err)
	}
	return &emptyResponse{}, nil
}
