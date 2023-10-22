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
