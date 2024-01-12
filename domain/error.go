package domain

import (
	"fmt"
	"net/http"
)

const (
	ErrNilContext       ConstantError = "nil Context"
	ErrNilInput         ConstantError = "nil Input"
	ErrExpiredToken     ConstantError = "ExpiredToken"
	ErrUserNotFound     ConstantError = "UserNotFound"
	ErrInvalidRequest   ConstantError = "InvalidRequest"
	ErrPasswordMismatch ConstantError = "PasswordMismatch"
	ErrDuplicatedUser   ConstantError = "DuplicatedUser"
)

type ConstantError string

func (e ConstantError) Error() string {
	return string(e)
}

type HTTPError struct {
	ErrorCode string `json:"errorCode"`
	Internal  error  `json:"-"`
}

func NewHTTPError(msgID string, err error) *HTTPError {
	return &HTTPError{
		ErrorCode: msgID,
		Internal:  err,
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("[%s:%d] %s: %v", e.ErrorCode, e.StatusCode(), e.Message(), e.Internal)
}

func (e *HTTPError) StatusCode() int {
	//	TODO: 에러코드 정리
	return http.StatusInternalServerError
}

func (e *HTTPError) Message() string {
	//	TODO: 메시지 추가
	return "error message"
}
