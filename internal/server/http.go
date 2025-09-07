package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"tonotdolist/internal/handler"
	"tonotdolist/internal/middleware"
	"tonotdolist/internal/service"
	"tonotdolist/internal/util"
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

func NewHTTPServer(logger zerolog.Logger, viper *viper.Viper, userHandler *handler.UserHandler, activityHandler *handler.ActivityHandler, userService service.UserService) *http.Server {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("id", util.ValidID)
	}

	engine := gin.New()
	s := http.NewServer(engine, logger, http.WithHost(viper.GetString(httpHostKey)), http.WithPort(viper.GetUint16(httpPortKey)))

	s.Use(middleware.LogMiddleware(logger)).Use(middleware.VersionMiddleware()).Use(gin.Recovery())

	v1 := s.Group("/v1")
	{
		noAuth := v1.Group("/")
		noAuth.POST("login", userHandler.Login)
		noAuth.POST("register", userHandler.Register)
	}

	{
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware(userService))

		{
			auth.POST("logout", userHandler.Logout)
			auth.POST("logoutall", userHandler.LogoutAll)
		}

		activity := auth.Group("activity/")
		{
			activity.POST("create", activityHandler.CreateActivity)
		}
	}

	return s
}
