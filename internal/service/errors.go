package service

import "errors"

var ErrEmailAlreadyExists = errors.New("пользователь с таким email уже существует")
var ErrEmailCheckFailed = errors.New("не удалось проверить email")
var ErrUserNotFound = errors.New("пользователь не найден")
var ErrUserCreateFailed = errors.New("не удалось создать пользователя")
var ErrUserCheckFailed = errors.New("не удалось проверить пользователя")
var ErrUserDeactivateFailed = errors.New("не удалось деактивировать пользователя")
var ErrTokenCreateFailed = errors.New("не удалось создать access token")
var ErrRefreshTokenCreateFailed = errors.New("не удалось создать refresh token")
var ErrRefreshTokenCheckFailed = errors.New("не удалось проверить refresh токен")
var ErrRefreshTokenRevokeFailed = errors.New("не удалось отозвать refresh tokens")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")
var ErrInvalidCredentials = errors.New("неверный email или пароль")
var ErrPasswordHashFailed = errors.New("не удалось захешировать пароль")
