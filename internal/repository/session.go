package repository

import "context"

type SessionRepository interface {
	AddSession(ctx context.Context, userId string) (string, error)
	GetSession(ctx context.Context, sessionId string) (string, error)
	DeleteSession(ctx context.Context, sessionId string) error
	DeleteAllUserSession(ctx context.Context, userId string) error
}
