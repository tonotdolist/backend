package repository

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository struct {
	logger zerolog.Logger
	db     *gorm.DB
	rdb    *redis.Client
}

func NewRepository(logger zerolog.Logger, db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{
		logger: logger,
		db:     db,
		rdb:    rdb,
	}
}
