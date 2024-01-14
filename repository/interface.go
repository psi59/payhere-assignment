package repository

import (
	"context"

	"github.com/psi59/payhere-assignment/domain"
)

const (
	ErrNilUserRepository        domain.ConstantError = "nil UserRepository"
	ErrTokenBlacklistRepository domain.ConstantError = "nil TokenBlacklistRepository"
)

type UserRepository interface {
	Create(c context.Context, user *domain.User) error
	Get(c context.Context, userID int) (*domain.User, error)
	GetByPhoneNumber(c context.Context, phoneNumber string) (*domain.User, error)
}

type TokenBlacklistRepository interface {
	Create(c context.Context, token *domain.AuthToken) error
	Get(c context.Context, token string) (*domain.AuthToken, error)
}
