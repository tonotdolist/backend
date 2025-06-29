package handler

import "github.com/gin-gonic/gin"

type UserHandler struct {
	*Handler
}

func NewUserHandler(handler *Handler) *UserHandler {
	return &UserHandler{
		Handler: handler,
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	panic("implement me")
}

func (h *UserHandler) Register(ctx *gin.Context) {
	panic("implement me")
}
