package service

import (
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/token"
	"errors"
	"strings"
	"time"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwtSecret        string
}

func NewAuthService(userRepo *repository.UserRepository, refreshTokenRepo *repository.RefreshTokenRepository, secretKey string) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        secretKey,
	}
}

func (s *AuthService) RegisterNewUser(req dto.RegisterRequest) (dto.RegisterResponse, error) {

	req.Name = strings.TrimSpace(req.Name)
	if err := validateName(req.Name); err != nil {
		return dto.RegisterResponse{}, err
	}

	req.Email = strings.TrimSpace(req.Email)
	if err := validateEmail(req.Email); err != nil {
		return dto.RegisterResponse{}, err
	}

	if err := validatePassword(req.Password); err != nil {
		return dto.RegisterResponse{}, err
	}

	passHash, err := hashPassword(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, ErrPasswordHashFailed
	}

	exists, err := s.userRepo.ExistsByEmail(req.Email)

	if err != nil {
		return dto.RegisterResponse{}, ErrEmailCheckFailed
	}

	if exists {
		return dto.RegisterResponse{}, ErrEmailAlreadyExists
	}

	newUser := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passHash,
		Role:         "user",
		IsActive:     true,
	}

	user, err := s.userRepo.CreateUser(newUser)
	if err != nil {
		return dto.RegisterResponse{}, ErrUserCreateFailed
	}

	newUserResponse := dto.RegisterResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Message: "Аккаунт создан",
	}

	return newUserResponse, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	req.Email = strings.TrimSpace(req.Email)
	if err := validateEmail(req.Email); err != nil {
		return dto.LoginResponse{}, err
	}

	if err := validatePassword(req.Password); err != nil {
		return dto.LoginResponse{}, err
	}

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return dto.LoginResponse{}, ErrInvalidCredentials
		}

		return dto.LoginResponse{}, ErrUserCheckFailed
	}

	if !user.IsActive {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	if err := comparePassword(user.PasswordHash, req.Password); err != nil {
		return dto.LoginResponse{}, err
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return dto.LoginResponse{}, ErrTokenCreateFailed
	}

	newRefreshToken := token.GenerateRefreshToken()
	tokenHash := token.HashRefreshToken(newRefreshToken)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	refreshTokenRequest := model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	_, err = s.refreshTokenRepo.CreateRefreshToken(refreshTokenRequest)
	if err != nil {
		return dto.LoginResponse{}, ErrRefreshTokenCreateFailed
	}

	return dto.LoginResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		Message:      "Вход выполнен успешно",
	}, nil
}

func (s *AuthService) GetMe(userID uint64) (dto.MeResponse, error) {

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return dto.MeResponse{}, ErrUserNotFound
		}
		return dto.MeResponse{}, ErrUserCheckFailed
	}

	if !user.IsActive {
		return dto.MeResponse{}, ErrUserNotFound
	}

	return dto.MeResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
	}, nil
}

func (s *AuthService) Refresh(req dto.RefreshRequest) (dto.RefreshResponse, error) {
	if req.RefreshToken == "" {
		return dto.RefreshResponse{}, ErrInvalidRefreshToken
	}

	tokenHash := token.HashRefreshToken(req.RefreshToken)
	refreshToken, err := s.refreshTokenRepo.GetRefreshTokenByHash(tokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return dto.RefreshResponse{}, ErrInvalidRefreshToken
		}

		return dto.RefreshResponse{}, ErrRefreshTokenCheckFailed
	}

	if refreshToken.RevokedAt != nil {
		return dto.RefreshResponse{}, ErrInvalidRefreshToken
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return dto.RefreshResponse{}, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetUserByID(refreshToken.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return dto.RefreshResponse{}, ErrInvalidRefreshToken
		}
		return dto.RefreshResponse{}, ErrUserCheckFailed
	}

	if !user.IsActive {
		return dto.RefreshResponse{}, ErrInvalidRefreshToken
	}

	accessToken, err := token.GenerateAccessToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return dto.RefreshResponse{}, ErrTokenCreateFailed
	}

	return dto.RefreshResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		Message:     "Создан новый токен",
	}, nil
}

func (s *AuthService) Logout(req dto.LogoutRequest) (dto.LogoutResponse, error) {
	if err := validateRefreshToken(req.RefreshToken); err != nil {
		return dto.LogoutResponse{}, err
	}

	hashRefreshToken := token.HashRefreshToken(req.RefreshToken)
	_, err := s.refreshTokenRepo.RevokeRefreshToken(hashRefreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			return dto.LogoutResponse{}, ErrInvalidRefreshToken
		}
		return dto.LogoutResponse{}, ErrRefreshTokenRevokeFailed
	}

	return dto.LogoutResponse{
		Message: "Вы вышли из аккаунта",
	}, nil
}

func (s *AuthService) DeactivateAccount(userID uint64) (dto.DeleteAccountResponse, error) {
	user, err := s.userRepo.DeactivateUser(userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return dto.DeleteAccountResponse{}, ErrUserNotFound
		}

		return dto.DeleteAccountResponse{}, ErrUserDeactivateFailed
	}

	err = s.refreshTokenRepo.RevokeAllRefreshTokens(user.ID)
	if err != nil {
		return dto.DeleteAccountResponse{}, ErrRefreshTokenRevokeFailed
	}

	return dto.DeleteAccountResponse{
		Message: "Аккаунт удалён",
	}, nil
}
