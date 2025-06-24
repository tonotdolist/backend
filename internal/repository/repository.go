package repository

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository struct {
	logger zerolog.Logger
	db     *gorm.DB
}
