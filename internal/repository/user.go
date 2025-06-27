package repository

import (
	"context"
	"tonotdolist/internal/model"
)

type UserRepository interface {
	GetByEmail(context.Context, string) (*model.User, error)
	CreateUser(context.Context, *model.User) error
}
