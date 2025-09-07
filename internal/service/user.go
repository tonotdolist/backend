package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
	"tonotdolist/common"
	"tonotdolist/internal/model"
	"tonotdolist/internal/repository"
	"tonotdolist/internal/util"
	"tonotdolist/pkg/config"
)

const (
	bcryptCostKey    = "auth.bcryptCost"
	sessionLengthKey = "auth.sessionLength"

	// 1+ numbers
	// 1+ uppercase characters
	// 1+ lowercase characters
	// 1+ special characters
	regexString = "(?=.*[^a-zA-Z0-9]+)(?=.*[0-9]+)(?=.*[a-z]+)(?=.*[A-Z]+).*"
)

func init() {
	config.RegisterRequiredKey(bcryptCostKey, sessionLengthKey)
}

type UserService interface {
	GetSession(ctx context.Context, sessionId string) (string, error)
	Login(context.Context, *common.UserLoginRequest) (string, error)
	Register(context.Context, *common.UserRegisterRequest) (string, error)
	Logout(ctx context.Context, sessionId string, userId string) error
	LogoutAll(ctx context.Context, userId string) error
}

type userService struct {
	bcryptCost        int
	sessionLength     int64
	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	regex *regexp.Regexp
}

func NewUserService(userRepository repository.UserRepository, sessionRepository repository.SessionRepository, viper *viper.Viper, logger zerolog.Logger) UserService {
	regex, err := regexp.Compile(regexString)
	if err != nil {
		logger.Fatal().Msgf("error compiling regex: %v", err)
	}

	return &userService{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		bcryptCost:        viper.GetInt(bcryptCostKey),
		sessionLength:     viper.GetInt64(sessionLengthKey),

		regex: regex,
	}
}

func (s *userService) GetSession(ctx context.Context, sessionId string) (string, error) {
	session, err := s.sessionRepository.GetSession(ctx, sessionId)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return "", common.ErrUnauthorized
		}

		return "", fmt.Errorf("error fetching user session info from repo: %w", err)
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

	id, err := util.NewID()
	if err != nil {
		return "", fmt.Errorf("error generating user id: %w", err)
	}

	user := &model.User{
		UserId:   id,
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

func (s *userService) validatePassword(password string) error {
	length := len(password)
	if length < 8 {
		return common.ErrPasswordTooShort
	}

	if length > 72 { // 72 because of bcrypt limit
		return common.ErrPasswordTooLong
	}

	if !s.regex.Match([]byte(password)) {
		return common.ErrBadPassword
	}

	return nil
}

func (s *userService) createSession(ctx context.Context, userId string) (string, error) {
	sessionId, err := util.NewID()
	if err != nil {
		return "", fmt.Errorf("error generating session id: %w", err)
	}

	err = s.sessionRepository.AddSession(ctx, userId, sessionId, time.Now().Unix()+s.sessionLength)
	if err != nil {
		return "", fmt.Errorf("error adding session id to repo: %w", err)
	}

	return sessionId, nil
}

func (s *userService) Logout(ctx context.Context, sessionId string, userId string) error {
	if err := s.sessionRepository.DeleteSession(ctx, sessionId, userId); err != nil {
		return fmt.Errorf("error deleting session from repo: %w", err)
	}

	return nil
}

func (s *userService) LogoutAll(ctx context.Context, userId string) error {
	if err := s.sessionRepository.DeleteAllUserSession(ctx, userId); err != nil {
		return fmt.Errorf("error deleting all sessions from repo: %w", err)
	}

	return nil
}
