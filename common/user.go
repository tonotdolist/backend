package common

type UserRegisterRequest struct {
	Email    string
	Password string
}

type UserLoginRequest struct {
	Email    string
	Password string
}

type UserRegisterResponse struct {
	SessionID string
}

type UserLoginResponse struct {
	SessionID string
}

type UserSession struct {
	UserID string
	Expire int64
}
