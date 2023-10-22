package service

import (
	"context"
	"effectivemobile/internal/app/model"
	"effectivemobile/internal/pkg/config"
)

type IRepository interface {
	DeleteUser(ctx context.Context, id int) error
	AddUser(ctx context.Context, info model.UserInfo) (int, error)
	UpdateUser(ctx context.Context, info model.UserInfo) error
	GetUser(ctx context.Context, filters model.UserFilter) ([]model.UserInfo, error)
}

type IEnrichService interface {
	Enrich(user model.UserInfo) (model.UserInfo, error)
	Agify(user model.UserInfo, apiUrl string) (model.UserInfo, error)
	Genderize(user model.UserInfo, apiUrl string) (model.UserInfo, error)
	Nationalize(user model.UserInfo, apiUrl string) (model.UserInfo, error)
}

type Service struct {
	cfg           config.Config
	repo          IRepository
	enrichService IEnrichService
}

func New(cfg config.Config, r IRepository, ec IEnrichService) *Service {
	return &Service{cfg, r, ec}
}
