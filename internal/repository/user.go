package repository

import (
	"context"
	"github.com/rs/zerolog"
	"tonotdolist/internal/model"
)

type UserRepository interface {
	GetByEmail(context.Context, string) (*model.User, error)
	CreateUser(context.Context, *model.User) error
}

type userRepository struct {
	logger zerolog.Logger
	*Repository
}

func NewUserRepository(repository *Repository) UserRepository {
	return &userRepository{
		Repository: repository,
		logger:     repository.logger.With().Str("repository", "user").Logger(),
	}
}

func (up *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	panic("implement me")
}

func (up *userRepository) CreateUser(context.Context, *model.User) error {
	panic("implement me")
}
