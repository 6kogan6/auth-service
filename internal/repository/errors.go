package repository

import "errors"

var ErrUserNotFound = errors.New("пользователь не найден")
var ErrRefreshTokenNotFound = errors.New("refresh token не найден")
