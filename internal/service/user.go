package service

import (
	"context"
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"tonotdolist/common"
	"tonotdolist/internal/model"
	"tonotdolist/internal/repository"
	"tonotdolist/pkg/config"
)

const bcryptCostKey = "auth.bcryptCost"

func init() {
	config.RegisterRequiredKey(bcryptCostKey)
}

type UserService interface {
	Login(context.Context, *common.UserLoginRequest) error
	Register(context.Context, *common.UserRegisterRequest) error
}

type userService struct {
	bcryptCost     int
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository, viper *viper.Viper) UserService {
	return &userService{
		userRepository: userRepository,
		bcryptCost:     viper.GetInt(bcryptCostKey),
	}
}

func (s *userService) Login(ctx context.Context, req *common.UserLoginRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 15)
	if err != nil {
		return err
	}

	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if user.Password != string(hashedPassword) {
		return common.ErrUnauthorized
	}

	return nil
}

func (s *userService) Register(ctx context.Context, req *common.UserRegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 15)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return common.ErrConflict
		}

		return err
	}

	return nil
}
