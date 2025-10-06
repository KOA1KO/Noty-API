package service

import (
	"Noty/internal/config"
	"Noty/internal/models"
	"Noty/internal/storage"
	"Noty/internal/storage/pg"
	"Noty/pkg/auth"
	"Noty/pkg/hasher"
	"Noty/pkg/logger/sl"
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"
)

type UserService struct {
	Config         *config.Config
	Logger         *slog.Logger
	UserRepository storage.UserStorage
	TokenManager   auth.TokenManager
}

func NewUserService(userRepo storage.UserStorage, TM auth.TokenManager,
	logger *slog.Logger, config *config.Config) *UserService {
	return &UserService{
		Config:         config,
		Logger:         logger,
		UserRepository: userRepo,
		TokenManager:   TM,
	}
}

func (s *UserService) SignIn(ctx context.Context, input SignInput) (models.Tokens, error) {
	const op = "internal.service.user_service.SignIn"
	log := s.Logger.With(
		sl.Operation(op),
	)

	var hasher hasher.Hasher
	hashedPass, err := hasher.Hash(input.Password)
	if err != nil {
		log.Error("hasher error", sl.Err(err))
		return models.Tokens{}, ErrInternal
	}

	userID, err := s.UserRepository.CreateUser(ctx, input.Username, hashedPass)
	if err != nil {
		if errors.Is(err, pg.ErrAlredyExists) {
			log.Debug("user is already sign in")
			return models.Tokens{}, ErrUserExists
		}
		log.Error(err.Error())
		return models.Tokens{}, ErrInternal
	}
	log.Info("user is writed to DB")
	return s.createSession(ctx, userID)
}

func (s *UserService) Login(ctx context.Context, input SignInput) (models.Tokens, error) {
	const op = "internal.service.user_service.Login"
	log := s.Logger.With(
		sl.Operation(op),
	)
	user, err := s.UserRepository.GetUser(input.Username)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			log.Debug("user not found")
			return models.Tokens{}, ErrUserNotFound
		}
		log.Error(err.Error())
		return models.Tokens{}, ErrInternal
	}

	var hasher hasher.Hasher
	if !hasher.Compare(user.Password, input.Password) {
		log.Debug("invalid credentials")
		return models.Tokens{}, ErrInvalidCredentials
	}

	tokens, err := s.createSession(ctx, user.ID)
	if err != nil {
		log.Error(err.Error())
		return models.Tokens{}, ErrInternal
	}

	log.Info("user founded")
	return models.Tokens{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (s *UserService) createSession(ctx context.Context, userID int) (models.Tokens, error) {
	const op = "internal.service.user_service.createSession"
	log := s.Logger.With(
		sl.Operation(op),
	)

	var (
		result models.Tokens
		err    error
	)

	result.Access, err = s.TokenManager.NewJWT(strconv.Itoa(userID), s.Config.AccessJwtTTL)
	if err != nil {
		log.Error("access JWT generation error", sl.Err(err))
		return models.Tokens{}, ErrInternal
	}
	result.Refresh, err = s.TokenManager.NewRefreshToken()
	if err != nil {
		log.Error("refresh JWT generation error", sl.Err(err))
		return models.Tokens{}, ErrInternal
	}
	log.Debug("all tokens getted")

	var hasher hasher.Hasher
	hashedRefresh, err := hasher.Hash(result.Refresh)
	if err != nil {
		log.Error("refresh hasher", sl.Err(err))
		return models.Tokens{}, ErrInternal
	}

	session := models.Session{
		RefreshTokenHash: hashedRefresh,
		ExpiresAt:        time.Now().Add(s.Config.RefreshJwtTTL),
		IssuedAt:         time.Now(),
	}
	err = s.UserRepository.SetSession(ctx, userID, session)
	return result, err
}
