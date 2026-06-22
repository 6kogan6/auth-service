package service

import (
	"errors"
	"strings"
)

func validateEmail(email string) error {
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)
	if len(email) == 0 {
		return errors.New("Email не может быть пустым")
	}

	if !(strings.Contains(email, "@")) {
		return errors.New("Email должен содержать @")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) == 0 {
		return errors.New("Password не может быть пустым")
	}

	if len(password) < 8 {
		return errors.New("Password должен быть минимум 8 символов")
	}

	return nil
}

func validateName(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return errors.New("Name не может быть пустым")
	}

	return nil
}

func validateRefreshToken(refreshToken string) error {
	if len(refreshToken) == 0 {
		return ErrInvalidRefreshToken
	}
	return nil
}
