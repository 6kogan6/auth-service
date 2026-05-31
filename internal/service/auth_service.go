package service

import (
	"auth-service/internal/dto"
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("Password не может быть пустым")
	}

	if len(password) < 8 {
		return "", errors.New("Password должен быть минимум 8 символов")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	return string(passwordHash), nil
}

func (s *AuthService) RegisterNewUser(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	var user model.User

	if len(req.Name) == 0 {
		return dto.RegisterResponse{}, errors.New("Name не может быть пустым")
	}

	if len(req.Email) == 0 {
		return dto.RegisterResponse{}, errors.New("Email не может быть пустым")
	}

	if !(strings.Contains(req.Email, "@")) {
		return dto.RegisterResponse{}, errors.New("Email должен содержать @")
	}

	passHash, err := HashPassword(req.Password)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("неверный формат пароля: %v", err)
	}

	exists, err := s.userRepo.ExistsByEmail(req.Email)

	if err != nil {
		return dto.RegisterResponse{}, errors.New("Не удалось проверить email")
	}

	if exists {
		return dto.RegisterResponse{}, errors.New("Пользователь с таким email уже существует")
	}

	newUser := model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passHash,
		Role:         "user",
		IsActive:     true,
	}

	user, err = s.userRepo.CreateUser(newUser)
	if err != nil {
		return dto.RegisterResponse{}, errors.New("Не удалось создать пользователя")
	}

	newUserResponse := dto.RegisterResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Message: "Аккаунт создан",
	}

	return newUserResponse, nil
}
