package domain

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
