package item

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/internal/valid"

	"github.com/psi59/payhere-assignment/domain"
)

type Usecase interface {
	Create(c context.Context, input *CreateInput) (*CreateOutput, error)
	Get(c context.Context, input *GetInput) (*GetOutput, error)
	Delete(c context.Context, input *DeleteInput) error
	Update(c context.Context, input *UpdateInput) error
	Find(c context.Context, input *FindInput) (*FindOutput, error)
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

type DeleteInput struct {
	User   *domain.User `validate:"required"`
	ItemID int          `validate:"required"`
}

func (i *DeleteInput) Validate() error {
	if err := valid.ValidateStruct(i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type UpdateInput struct {
	User        *domain.User     `validate:"required"`
	ItemID      int              `validate:"required"`
	Name        *string          `validate:"omitempty,required"`
	Description *string          `validate:"omitempty,required"`
	Price       *int             `validate:"omitempty,gt=0"`
	Cost        *int             `validate:"omitempty,gt=0"`
	Category    *string          `validate:"omitempty,required"`
	Barcode     *string          `validate:"omitempty,required"`
	Size        *domain.ItemSize `validate:"omitempty,required,oneof=small large"`
	ExpiryAt    *time.Time       `validate:"omitempty,required"`
}

func (i *UpdateInput) Validate() error {
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

type FindInput struct {
	User        *domain.User `validate:"required"`
	Keyword     string
	SearchAfter int
}

func (i *FindInput) Validate() error {
	if err := valid.ValidateStruct(i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type FindOutput struct {
	TotalCount  int
	Items       []domain.Item
	HasNext     bool
	SearchAfter int
}
