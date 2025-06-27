package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
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

func NewHTTPServer(logger zerolog.Logger, viper *viper.Viper) *http.Server {
	s := http.NewServer(gin.Default(), logger, http.WithHost(viper.GetString(HTTPHostKey)), http.WithPort(viper.GetUint16(HTTPPortKey)))

	return s
}
