package handler

import (
	"github.com/rs/zerolog"
	"tonotdolist/pkg/api"
)

type Handler struct {
	logger    zerolog.Logger
	responder *api.RequestResponder
}

func NewHandler(logger zerolog.Logger, responder *api.RequestResponder) *Handler {
	return &Handler{
		logger:    logger,
		responder: responder,
	}
}
