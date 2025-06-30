package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return common.ErrUnauthorized
	}

	return nil
}

func (s *userService) Register(ctx context.Context, req *common.UserRegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 15)
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	user := &model.User{
		UserId:   id.String(),
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
