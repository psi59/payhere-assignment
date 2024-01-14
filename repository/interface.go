package repository

import (
	"context"

	"github.com/psi59/payhere-assignment/domain"
)

const (
	ErrNilUserRepository           domain.ConstantError = "nil UserRepository"
	ErrNilTokenBlacklistRepository domain.ConstantError = "nil TokenBlacklistRepository"
	ErrNilItemRepository           domain.ConstantError = "nil ItemRepository"
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

type ItemRepository interface {
	Create(c context.Context, item *domain.Item) error
	Get(c context.Context, userID, itemID int) (*domain.Item, error)
}
