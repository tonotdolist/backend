package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Option func(s *Server)
type Server struct {
	*gin.Engine
	logger          zerolog.Logger
	httpSrv         *http.Server
	host            string
	port            int16
	shutdownTimeout int
}

func NewServer(engine *gin.Engine, logger zerolog.Logger, options ...Option) *Server {
	s := &Server{
		Engine: engine,
		logger: logger,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func (s *Server) Start() {
	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s,
	}

	go func() {
		if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Panic().Err(err).Msg("error starting server")
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.shutdownTimeout))
	defer cancel()
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		s.logger.Fatal().Msg("server forced to shutdown")
	}

	s.logger.Info().Msg("server stopped")
	return nil
}

func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithPort(port int16) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithShutdownTimeout(timeout int) Option {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
