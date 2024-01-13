package authtoken

import (
	"context"
	"time"

	"github.com/psi59/payhere-assignment/domain"
)

type Usecase interface {
	Create(c context.Context, input *CreateInput) (*CreateOutput, error)
	Verify(c context.Context, input *VerifyInput) (*VerifyOutput, error)
	RegisterBlacklist(c context.Context, input *RegisterBlacklistInput) error
	GetBlacklist(c context.Context, input *GetBlacklistInput) (*GetBlacklistOutput, error)
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

type RegisterBlacklistInput struct {
	Token string `validate:"required"`
}

type GetBlacklistInput struct {
	Token string `validate:"required"`
}

type GetBlacklistOutput struct {
	Token *domain.AuthToken
}
