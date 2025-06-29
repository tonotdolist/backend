package service

import (
	"context"
	"tonotdolist/common"
	"tonotdolist/internal/repository"
)

type UserService interface {
	Login(context.Context, *common.UserLoginRequest) error
	Register(context.Context, *common.UserRegisterRequest) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Login(ctx context.Context, req *common.UserLoginRequest) error {
	panic("implement me")
}

func (s *userService) Register(ctx context.Context, req *common.UserRegisterRequest) error {
	panic("implement me")
}
