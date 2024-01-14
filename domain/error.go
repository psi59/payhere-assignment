package domain

const (
	ErrNilContext                  ConstantError = "nil Context"
	ErrNilInput                    ConstantError = "nil Input"
	ErrInvalidRequest              ConstantError = "InvalidRequest"
	ErrExpiredToken                ConstantError = "ExpiredToken"
	ErrUserNotFound                ConstantError = "UserNotFound"
	ErrItemNotFound                ConstantError = "ItemNotFound"
	ErrTokenBlacklistNotFound      ConstantError = "TokenBlacklistNotFound"
	ErrPasswordMismatch            ConstantError = "PasswordMismatch"
	ErrUserAlreadyExists           ConstantError = "UserAlreadyExists"
	ErrTokenBlacklistAlreadyExists ConstantError = "TokenBlacklistAlreadyExists"
	ErrItemAlreadyExists           ConstantError = "ItemAlreadyExists"
)

type ConstantError string

func (e ConstantError) Error() string {
	return string(e)
}
