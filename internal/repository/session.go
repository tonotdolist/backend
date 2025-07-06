package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"tonotdolist/common"
)

type SessionRepository interface {
	AddSession(ctx context.Context, userId string, sessionId string, expire int64) error
	GetSession(ctx context.Context, sessionId string) (*common.UserSession, error)
	DeleteSession(ctx context.Context, sessionId string, userId string) error
	DeleteAllUserSession(ctx context.Context, userId string) error
}

type sessionRepository struct {
	*Repository
}

func NewSessionRepository(repository *Repository) SessionRepository {
	return &sessionRepository{
		Repository: repository,
	}
}

func (r *sessionRepository) AddSession(ctx context.Context, userId string, sessionId string, expire int64) error {
	session := &common.UserSession{
		UserID: userId,
		Expire: expire,
	}

	if err := r.rdb.Set(ctx, sessionId, formatContent(session), 0).Err(); err != nil {
		return err
	}

	if err := r.rdb.SAdd(ctx, userId, sessionId).Err(); err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) GetSession(ctx context.Context, sessionId string) (*common.UserSession, error) {
	res := r.rdb.Get(ctx, sessionId)

	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, common.ErrNotFound
		}

		return nil, fmt.Errorf("error fetching user session: %w", err)
	}

	session, err := parseContent(res.Val())
	if err != nil {
		return nil, fmt.Errorf("error parsing session data from value: %w", err)
	}

	return session, nil
}

func (r *sessionRepository) DeleteSession(ctx context.Context, sessionId string, userId string) error {
	newCtx := context.WithoutCancel(ctx)
	pipeline := r.rdb.Pipeline()
	cmd1 := pipeline.Del(newCtx, sessionId)
	cmd2 := pipeline.SRem(newCtx, userId, sessionId)
	_, err := pipeline.Exec(newCtx)
	if err != nil {
		return fmt.Errorf("unable to execute delete session pipeline: %w", err)
	}

	return errors.Join(cmd1.Err(), cmd2.Err())
}

func (r *sessionRepository) DeleteAllUserSession(ctx context.Context, userId string) error {
	userSessionList, err := r.rdb.SMembers(ctx, userId).Result()
	if err != nil {
		return fmt.Errorf("error fetching user session list: %w", err)
	}

	if len(userSessionList) > 0 {
		pipeline := r.rdb.Pipeline()
		cmd1 := pipeline.Del(ctx, userSessionList...)
		cmd2 := pipeline.Del(ctx, userId)

		_, err = pipeline.Exec(ctx)
		if err != nil {
			return fmt.Errorf("unable to execute clear user sessions pipeline: %w", err)
		}

		return errors.Join(cmd1.Err(), cmd2.Err())
	}

	return nil
}

func formatContent(session *common.UserSession) string {
	return fmt.Sprintf("%s:%d", session.UserID, session.Expire)
}

func parseContent(content string) (*common.UserSession, error) {
	parts := strings.Split(content, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format: %s", content)
	}

	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse expiration timestamp: %w", err)
	}

	return &common.UserSession{
		UserID: parts[0],
		Expire: exp,
	}, nil
}

var _ SessionRepository = (*sessionRepository)(nil)
