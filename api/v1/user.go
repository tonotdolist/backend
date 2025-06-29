package v1

import (
	"tonotdolist/common"
	"tonotdolist/pkg/api"
)

func init() {
	api.RegisterRequest[common.UserRegisterRequest, UserRegisterRequest](version)
	api.RegisterRequest[common.UserLoginRequest, UserLoginRequest](version)
}

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

func (r *UserLoginRequest) ToInternalRequest() interface{} {
	return &common.UserLoginRequest{
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r *UserRegisterRequest) ToInternalRequest() interface{} {
	return &common.UserRegisterRequest{
		Email:    r.Email,
		Password: r.Password,
	}
}

var _ api.VersionedRequest = (*UserLoginRequest)(nil)
var _ api.VersionedRequest = (*UserRegisterRequest)(nil)
