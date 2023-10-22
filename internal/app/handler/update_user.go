package handler

import (
	"context"
	"fmt"
)

// UpdateUserRequest example
type UpdateUserRequest struct {
	Id         int    `json:"id" example:"2"`
	Name       string `json:"name" example:"Bill"`
	Surname    string `json:"surname" example:"McDonald Jr."`
	Patronymic string `json:"patronymic" example:""`
}

// UpdateUserHandler godoc
//
//	@Summary		Update user
//	@Description	update user by passing new values
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param 			request body 	UpdateUserRequest true "query params"
//	@Success		200	{object}	emptyResponse
//	@Router			/updateUser [post]
func (h Handler) UpdateUserHandler(ctx context.Context, req UpdateUserRequest) (*emptyResponse, error) {
	err := h.service.UpdateUser(ctx, toModelUser(req))
	if err != nil {
		return nil, fmt.Errorf("cannot add segment: %w", err)
	}
	return &emptyResponse{}, nil
}
