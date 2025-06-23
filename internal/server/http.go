package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"tonotdolist/pkg/server/http"
)

func NewHTTPServer(logger zerolog.Logger) *http.Server {
	s := http.NewServer(gin.Default(), logger, http.WithHost("0.0.0.0"), http.WithPort(8080))

	return s
}
