package handler

import (
	"context"
	"fmt"
)

// DeleteUserRequest example
type DeleteUserRequest struct {
	Id int `query:"id" validate:"required" example:"3"`
}

// DeleteUserHandler godoc
//
// @Summary Delete user
// @Description Delete a user by ID
// @Tags user
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} emptyResponse
// @Router /deleteUser [delete]
func (h Handler) DeleteUserHandler(ctx context.Context, req DeleteUserRequest) (*emptyResponse, error) {
	err := h.service.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("cannot delete user: %w", err)
	}
	return &emptyResponse{}, err
}
