package domain

import (
	"fmt"

	"github.com/psi59/payhere-assignment/internal/i18n"
	"golang.org/x/text/language"
)

const (
	ErrNilContext               ConstantError = "nil Context"
	ErrNilInput                 ConstantError = "nil Input"
	ErrExpiredToken             ConstantError = "ExpiredToken"
	ErrUserNotFound             ConstantError = "UserNotFound"
	ErrTokenBlacklistNotFound   ConstantError = "TokenBlacklistNotFound"
	ErrInvalidRequest           ConstantError = "InvalidRequest"
	ErrPasswordMismatch         ConstantError = "PasswordMismatch"
	ErrDuplicatedUser           ConstantError = "DuplicatedUser"
	ErrDuplicatedTokenBlacklist ConstantError = "DuplicatedTokenBlacklist"
)

type ConstantError string

func (e ConstantError) Error() string {
	return string(e)
}

type HTTPError struct {
	StatusCode int
	ErrorCode  string
	Internal   error
}

func NewHTTPError(statusCode int, msgID string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		ErrorCode:  msgID,
		Internal:   err,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%s:%d] %s: %v", e.ErrorCode, e.StatusCode, e.Message(), e.Internal)
}

func (e *HTTPError) Message() string {
	return i18n.T(language.English, e.ErrorCode, nil)
}
