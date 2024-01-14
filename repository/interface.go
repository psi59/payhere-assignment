package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/valid"
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
	Delete(c context.Context, userID, itemID int) error
	Update(c context.Context, userID, itemID int, input *UpdateItemInput) error
}

type UpdateItemInput struct {
	Name        *string          `validate:"omitnil,gt=0"`
	Description *string          `validate:"omitnil,gt=0"`
	Price       *int             `validate:"omitnil,gt=0"`
	Cost        *int             `validate:"omitnil,gt=0"`
	Category    *string          `validate:"omitnil,gt=0"`
	Barcode     *string          `validate:"omitnil,gt=0"`
	Size        *domain.ItemSize `validate:"omitnil,gt=0"`
	ExpiryAt    *time.Time       `validate:"omitnil,gt=0"`
}

func (i *UpdateItemInput) Validate() error {
	if valid.IsNil(i.Name) &&
		valid.IsNil(i.Description) &&
		valid.IsNil(i.Price) &&
		valid.IsNil(i.Cost) &&
		valid.IsNil(i.Category) &&
		valid.IsNil(i.Barcode) &&
		valid.IsNil(i.Size) &&
		valid.IsNil(i.ExpiryAt) {
		return fmt.Errorf("invalid input")
	}
	if err := valid.ValidateStruct(i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
