package handler

import (
	"context"
	"effectivemobile/internal/app/model"
	"fmt"
)

// GetUserRequest example
type GetUserRequest struct {
	NameLike string `query:"name_like" json:"name_like" example:"Pet"`
	AgeFrom  int    `query:"age_from" json:"age_from" example:"30"`
	AgeTo    int    `query:"age_to" json:"age_to" example:"90"`
	Gender   string `query:"gender" json:"gender" example:"male"`
	Nation   string `query:"nation" json:"nation" example:"RU"`
	PageSize int    `query:"page_size" json:"page_size" example:"2"`
	PageNum  int    `query:"page_num" json:"page_num" example:"1"`
}

// GetUserResponse example
type GetUserResponse struct {
	Users []model.UserInfo `json:"users"`
}

// GetUserHandler godoc
//
// @Summary Get users
// @Description Get a list of users based on the provided filters and pagination options
// @Tags user
// @Accept json
// @Produce json
// @Param name_like query string false "Name contains"
// @Param age_from query int false "Minimum age"
// @Param age_to query int false "Maximum age"
// @Param gender query string false "Gender"
// @Param nation query string false "Nation"
// @Param page_size query int false "Page size"
// @Param page_num query int false "Page number"
// @Success 200 {object} GetUserResponse
// @Router /getUser [get]
func (h Handler) GetUserHandler(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	list, err := h.service.GetUser(ctx, toUserFilter(req))
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &GetUserResponse{Users: list}, nil
}
