package item

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/psi59/payhere-assignment/domain"
)

type Usecase interface {
	Create(c context.Context, input *CreateInput) (*CreateOutput, error)
	Get(c context.Context, input *GetInput) (*GetOutput, error)
}

const ErrNilUsecase domain.ConstantError = "nil ItemUsecase"

type CreateInput struct {
	User        *domain.User    `validate:"required"`
	Name        string          `validate:"required"`
	Description string          `validate:"required"`
	Price       int             `validate:"gt=0"`
	Cost        int             `validate:"gt=0"`
	Category    string          `validate:"required"`
	Barcode     string          `validate:"required"`
	Size        domain.ItemSize `validate:"required,oneof=small large"`
	ExpiryAt    time.Time       `validate:"required"`
}

func (i *CreateInput) Validate() error {
	if err := valid.ValidateStruct(i); err != nil {
		return errors.WithStack(err)
	}
	if err := i.Size.Validate(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type CreateOutput struct {
	Item *domain.Item
}

type GetInput struct {
	User   *domain.User `validate:"required"`
	ItemID int          `validate:"required"`
}

func (i *GetInput) Validate() error {
	if err := valid.ValidateStruct(i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type GetOutput struct {
	Item *domain.Item
}
