package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/handler"
	"tonotdolist/pkg/config"
	"tonotdolist/pkg/server/http"
)

const (
	HTTPHostKey = "server.http.host"
	HTTPPortKey = "server.http.port"
)

func init() {
	config.RegisterRequiredKey(HTTPHostKey, HTTPPortKey)
}

func NewHTTPServer(logger zerolog.Logger, viper *viper.Viper, userHandler *handler.UserHandler) *http.Server {
	s := http.NewServer(gin.Default(), logger, http.WithHost(viper.GetString(HTTPHostKey)), http.WithPort(viper.GetUint16(HTTPPortKey)))

	v1 := s.Group("/v1")
	{
		noAuth := v1.Group("/")
		noAuth.POST("login", userHandler.Login)
		noAuth.POST("register", userHandler.Register)
	}

	return s
}
