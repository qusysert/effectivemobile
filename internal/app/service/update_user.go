package service

import (
	"context"
	"effectivemobile/internal/app/model"
	"fmt"
)

func (s *Service) UpdateUser(ctx context.Context, info model.UserInfo) error {
	var err error
	info, err = s.enrichService.Enrich(info)
	if err != nil {
		return fmt.Errorf("cannot enrich user info: %w", err)
	}
	return s.repo.UpdateUser(ctx, info)
}
