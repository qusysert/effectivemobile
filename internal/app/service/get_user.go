package service

import (
	"context"
	"effectivemobile/internal/app/model"
)

func (s *Service) GetUser(ctx context.Context, filter model.UserFilter) ([]model.UserInfo, error) {
	return s.repo.GetUser(ctx, filter)
}
