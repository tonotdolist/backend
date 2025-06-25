package repository

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository struct {
	logger zerolog.Logger
	db     *gorm.DB
}

func NewRepository(logger zerolog.Logger, db *gorm.DB) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
	}
}
