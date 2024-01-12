package authtoken

import (
	"context"
	"time"

	"github.com/psi59/payhere-assignment/domain"
)

type Usecase interface {
	Create(c context.Context, input *CreateInput) (*CreateOutput, error)
	Verify(c context.Context, input *VerifyInput) (*VerifyOutput, error)
	RegisterToBlacklist(c context.Context, input *RegisterToBlacklistInput) error
}

const ErrNilUsecase domain.ConstantError = "nil AuthTokenUsecase"

type CreateInput struct {
	Identifier string `validate:"required"`
}

type CreateOutput struct {
	Token     string
	ExpiresAt time.Time
}

type VerifyInput struct {
	Token string `validate:"required"`
}

type VerifyOutput struct {
	Identifier string
	ExpiresAt  time.Time
}

type RegisterToBlacklistInput struct {
	Token string `validate:"required"`
}
