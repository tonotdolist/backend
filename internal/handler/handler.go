package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger zerolog.Logger
}

func NewHandler(logger zerolog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func getUserId(ctx *gin.Context) string {
	return ctx.GetString("user_id")
}
