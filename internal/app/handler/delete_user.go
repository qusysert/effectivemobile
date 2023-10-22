package handler

import (
	"context"
	"fmt"
)

type DeleteUserRequest struct {
	Id int `query:"id" validate:"required"`
}

func (h Handler) DeleteUserHandler(ctx context.Context, req DeleteUserRequest) (*emptyResponse, error) {
	err := h.service.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("cannot delete user: %w", err)
	}
	return &emptyResponse{}, err
}
