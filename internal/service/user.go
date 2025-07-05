package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
	"tonotdolist/common"
	"tonotdolist/internal/model"
	"tonotdolist/internal/repository"
	"tonotdolist/pkg/config"
)

const (
	bcryptCostKey    = "auth.bcryptCost"
	sessionLengthKey = "auth.sessionLength"
)

func init() {
	config.RegisterRequiredKey(bcryptCostKey, sessionLengthKey)
}

type UserService interface {
	GetSession(ctx context.Context, sessionId string) (string, error)
	Login(context.Context, *common.UserLoginRequest) (string, error)
	Register(context.Context, *common.UserRegisterRequest) (string, error)
}

type userService struct {
	bcryptCost        int
	sessionLength     int64
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository
}

func NewUserService(userRepository repository.UserRepository, sessionRepository repository.SessionRepository, viper *viper.Viper) UserService {
	return &userService{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		bcryptCost:        viper.GetInt(bcryptCostKey),
		sessionLength:     viper.GetInt64(sessionLengthKey),
	}
}

func (s *userService) GetSession(ctx context.Context, sessionId string) (string, error) {
	session, err := s.sessionRepository.GetSession(ctx, sessionId)
	if err != nil {
		if !errors.Is(err, common.ErrUnauthorized) {
			return "", fmt.Errorf("error fetching user session info from repo: %w", err)
		}

		return "", err
	}

	if time.Now().After(time.Unix(session.Expire, 0)) {
		return "", common.ErrUnauthorized
	}

	return session.UserID, nil
}

func (s *userService) Login(ctx context.Context, req *common.UserLoginRequest) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("error fetching user data from repo: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return "", common.ErrUnauthorized
	}

	sessionId, err := s.createSession(ctx, user.UserId)
	if err != nil {
		return "", fmt.Errorf("error create session: %w", err)
	}

	return sessionId, nil
}

func (s *userService) Register(ctx context.Context, req *common.UserRegisterRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return "", err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("error generating user id: %w", err)
	}

	user := &model.User{
		UserId:   id.String(),
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error inserting user data into db: %w", err)
	}

	sessionId, err := s.createSession(ctx, user.UserId)
	if err != nil {
		return "", fmt.Errorf("error create session: %w", err)
	}

	return sessionId, nil
}

func (s *userService) createSession(ctx context.Context, userId string) (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("error generating session id: %w", err)
	}
	sessionId := uuid.String()

	err = s.sessionRepository.AddSession(ctx, userId, sessionId, time.Now().Unix()+s.sessionLength)
	if err != nil {
		return "", fmt.Errorf("error adding session id to repo: %w", err)
	}

	return sessionId, nil
}
