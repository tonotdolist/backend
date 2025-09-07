package v1

import (
	"tonotdolist/common"
	"tonotdolist/pkg/api"
)

func init() {
	api.RegisterRequest[common.UserRegisterRequest, UserRegisterRequest](version)
	api.RegisterRequest[common.UserLoginRequest, UserLoginRequest](version)

	api.RegisterResponse[common.UserLoginResponse, UserLoginResponse](version)
	api.RegisterResponse[common.UserRegisterResponse, UserRegisterResponse](version)
}

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

type UserLoginResponse struct {
	SessionId string `json:"session_id,id"`
}

type UserRegisterResponse struct {
	SessionId string `json:"session_id,id"`
}

func (r *UserLoginResponse) FromInternalResponse(resp interface{}) {
	castedResp := resp.(*common.UserLoginResponse)

	r.SessionId = castedResp.SessionID
}

func (r *UserRegisterResponse) FromInternalResponse(resp interface{}) {
	castedResp := resp.(*common.UserRegisterResponse)

	r.SessionId = castedResp.SessionID
}

var _ api.VersionedRequest = (*UserLoginRequest)(nil)
var _ api.VersionedRequest = (*UserRegisterRequest)(nil)

var _ api.VersionedResponse = (*UserLoginResponse)(nil)
var _ api.VersionedResponse = (*UserRegisterResponse)(nil)
