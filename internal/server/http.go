package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/handler"
	"tonotdolist/internal/middleware"
	"tonotdolist/pkg/config"
	"tonotdolist/pkg/server/http"
)

const (
	httpHostKey = "server.http.host"
	httpPortKey = "server.http.port"
)

func init() {
	config.RegisterRequiredKey(httpHostKey, httpPortKey)
}

func NewHTTPServer(logger zerolog.Logger, viper *viper.Viper, userHandler *handler.UserHandler) *http.Server {
	s := http.NewServer(gin.New(), logger, http.WithHost(viper.GetString(httpHostKey)), http.WithPort(viper.GetUint16(httpPortKey)))

	s.Use(middleware.RequestLogMiddleware(logger)).Use(middleware.VersionMiddleware()).Use(gin.Recovery())

	v1 := s.Group("/v1")
	{
		noAuth := v1.Group("/")
		noAuth.POST("login", userHandler.Login)
		noAuth.POST("register", userHandler.Register)
	}

	return s
}
