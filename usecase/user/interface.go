package user

import (
	"context"
	"fmt"

	"github.com/psi59/payhere-assignment/domain"

	"github.com/psi59/payhere-assignment/internal/valid"
)

const ErrNilUsecase domain.ConstantError = "nil UserUsecase"

type Usecase interface {
	Create(c context.Context, input *CreateInput) (*CreateOutput, error)
	Get(c context.Context, input *GetInput) (*GetOutput, error)
	GetByPhoneNumber(c context.Context, input *GetByPhoneNumberInput) (*GetOutput, error)
}

type CreateInput struct {
	Name        string
	PhoneNumber string
	Password    string
}

func (i *CreateInput) Validate() error {
	if err := valid.ValidatePhoneNumber(i.PhoneNumber); err != nil {
		return fmt.Errorf("%w: %q", err, i.PhoneNumber)
	}

	if err := valid.ValidatePassword(i.Password); err != nil {
		return fmt.Errorf("%w: %q", err, i.Password)
	}

	return nil
}

type CreateOutput struct {
	User *domain.User
}

type GetInput struct {
	UserID int `validate:"gt=0"`
}

type GetOutput struct {
	User *domain.User
}

type GetByPhoneNumberInput struct {
	PhoneNumber string `validate:"required"`
}

func (i *GetByPhoneNumberInput) Validate() error {
	if err := valid.ValidatePhoneNumber(i.PhoneNumber); err != nil {
		return fmt.Errorf("%w: %q", err, i.PhoneNumber)
	}

	return nil
}
