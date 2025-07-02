package repository

import "context"

type SessionRepository interface {
	AddedSession(ctx context.Context, userId string) (string, error)
	GetSession(ctx context.Context, sessionId string)
	DeleteSession(ctx context.Context, sessionId string)
	DeleteAllUserSession(ctx context.Context, userId string)
}
