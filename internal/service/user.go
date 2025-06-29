package service

import (
	"context"
	"tonotdolist/common"
)

type UserService interface {
	Login(context.Context, common.UserLoginRequest) error
	Register(context.Context, common.UserRegisterRequest) error
}
