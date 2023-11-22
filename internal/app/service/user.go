package service

import (
	"context"
	"effectivemobile/internal/app/model"
	"fmt"
)

func (s *Service) AddUser(ctx context.Context, user model.UserInfo) (int, error) {
	var err error
	user, err = s.enrichService.Enrich(user)
	if err != nil {
		return 0, fmt.Errorf("cannot enrich user: %w", err)
	}
	id, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("cannot add user: %w", err)
	}
	return id, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *Service) GetUser(ctx context.Context, filter model.UserFilter) ([]model.UserInfo, error) {
	return s.repo.GetUser(ctx, filter)
}

func (s *Service) UpdateUser(ctx context.Context, info model.UserInfo) error {
	var err error
	info, err = s.enrichService.Enrich(info)
	if err != nil {
		return fmt.Errorf("cannot enrich user info: %w", err)
	}
	return s.repo.UpdateUser(ctx, info)
}
