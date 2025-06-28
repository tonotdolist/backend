package repository

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"tonotdolist/internal/common"
	"tonotdolist/internal/model"
)

type UserRepository interface {
	GetByEmail(context.Context, string) (*model.User, error)
	Create(context.Context, *model.User) error
	Update(context.Context, *model.User) error
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
	var user model.User
	err := up.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (up *userRepository) Create(ctx context.Context, user *model.User) error {
	return up.db.WithContext(context.WithoutCancel(ctx)).Create(user).Error
}

func (up *userRepository) Update(ctx context.Context, user *model.User) error {
	return up.db.WithContext(context.WithoutCancel(ctx)).Save(user).Error
}

var _ UserRepository = (*userRepository)(nil)
