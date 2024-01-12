package authtoken

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/psi59/payhere-assignment/repository"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/valid"
	"github.com/rs/xid"
)

type Service struct {
	secret                   []byte
	tokenBlacklistRepository repository.TokenBlacklistRepository
}

func NewService(secret string, tokenBlacklistRepository repository.TokenBlacklistRepository) (*Service, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("empty secret")
	}
	if valid.IsNil(tokenBlacklistRepository) {
		return nil, repository.ErrTokenBlacklistRepository
	}

	return &Service{
		secret:                   []byte(secret),
		tokenBlacklistRepository: tokenBlacklistRepository,
	}, nil
}

func (s *Service) Create(c context.Context, input *CreateInput) (*CreateOutput, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case valid.IsNil(input):
		return nil, domain.ErrNilInput
	}
	if err := valid.ValidateStruct(input); err != nil {
		return nil, errors.WithStack(err)
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.AddDate(0, 0, 7)
	claims := &jwt.RegisteredClaims{
		ID:        xid.New().String(),
		Issuer:    "payhere-assignment",
		Subject:   input.Identifier,
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(s.secret)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(token))

	return &CreateOutput{
		Token:     encoded,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *Service) Verify(c context.Context, input *VerifyInput) (*VerifyOutput, error) {
	switch {
	case valid.IsNil(c):
		return nil, domain.ErrNilContext
	case valid.IsNil(input):
		return nil, domain.ErrNilInput
	}
	if err := valid.ValidateStruct(input); err != nil {
		return nil, errors.WithStack(err)
	}

	tokenByte, err := base64.StdEncoding.DecodeString(input.Token)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var claims jwt.RegisteredClaims
	t, err := jwt.ParseWithClaims(string(tokenByte), &claims, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("%w: expiresAt(%s) < now(%s)", domain.ErrExpiredToken, claims.ExpiresAt.UTC(), time.Now().UTC())
		}

		return nil, errors.WithStack(err)
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	return &VerifyOutput{
		Identifier: claims.Subject,
		ExpiresAt:  claims.ExpiresAt.Time,
	}, nil
}

func (s *Service) RegisterToBlacklist(c context.Context, input *RegisterToBlacklistInput) error {
	switch {
	case valid.IsNil(c):
		return domain.ErrNilContext
	case valid.IsNil(input):
		return domain.ErrNilInput
	}
	if err := valid.ValidateStruct(input); err != nil {
		return errors.WithStack(err)
	}

	verifyOutput, err := s.Verify(c, &VerifyInput{Token: input.Token})
	if err != nil {
		if errors.Is(err, domain.ErrExpiredToken) {
			return nil
		}

		return errors.WithStack(err)
	}

	token := &domain.AuthToken{
		Token:     input.Token,
		ExpiresAt: verifyOutput.ExpiresAt,
	}
	if err := s.tokenBlacklistRepository.Create(c, token); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
