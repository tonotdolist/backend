package repository

import "context"

type SessionRepository interface {
	AddSession(ctx context.Context, userId string, sessionId string) error
	GetSession(ctx context.Context, sessionId string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
	DeleteAllUserSession(ctx context.Context, userId string) error
}

type sessionRepository struct {
	*Repository
}

func (r *sessionRepository) AddSession(ctx context.Context, userId string, sessionId string) error {
	panic("implement me")
}

func (r *sessionRepository) GetSession(ctx context.Context, sessionId string) (string, error) {
	panic("implement me")
}

func (r *sessionRepository) DeleteSession(ctx context.Context, sessionId string) error {
	panic("implement me")
}

func (r *sessionRepository) DeleteAllUserSession(ctx context.Context, userId string) error {
	panic("implement me")
}

var _ SessionRepository = (*sessionRepository)(nil)
