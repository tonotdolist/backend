package service

import "github.com/rs/zerolog"

type Service struct {
	logger zerolog.Logger
}

func NewService(logger zerolog.Logger) *Service {
	return &Service{
		logger: logger,
	}
}
