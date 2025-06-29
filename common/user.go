package common

type UserRegisterRequest struct {
	Email    string
	Password string
}

type UserLoginRequest struct {
	Email    string
	Password string
}
