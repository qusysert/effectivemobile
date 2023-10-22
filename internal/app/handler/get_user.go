package handler

import (
	"context"
	"effectivemobile/internal/app/model"
	"fmt"
)

type GetUserRequest struct {
	NameLike string `query:"name_like" json:"name_like"`
	AgeFrom  int    `query:"age_from" json:"age_from"`
	AgeTo    int    `query:"age_to" json:"age_to"`
	Gender   string `query:"gender" json:"gender"`
	Nation   string `query:"nation" json:"nation"`
	PageSize int    `query:"page_size" json:"page_size"`
	PageNum  int    `query:"page_num" json:"page_num"`
}

type GetUserResponse struct {
	Users []model.UserInfo `json:"users"`
}

func (h Handler) GetUserHandler(ctx context.Context, req GetUserRequest) (*GetUserResponse, error) {
	list, err := h.service.GetUser(ctx, toUserFilter(req))
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &GetUserResponse{Users: list}, nil
}
