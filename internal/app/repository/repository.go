package repository

import "effectivemobile/internal/pkg/config"

type Repository struct {
	cfg config.Config
}

func New(cfg config.Config) *Repository {
	return &Repository{cfg: cfg}
}
