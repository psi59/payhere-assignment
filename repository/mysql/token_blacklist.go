package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/psi59/payhere-assignment/internal/valid"
	"gorm.io/gorm"
)

type TokenBlacklistRepository struct{}

func NewTokenBlacklistRepository() *TokenBlacklistRepository {
	return &TokenBlacklistRepository{}
}

func (r *TokenBlacklistRepository) Create(c context.Context, token *domain.AuthToken) error {
	if valid.IsNil(c) {
		return domain.ErrNilContext
	}
	if valid.IsNil(token) {
		return domain.ErrNilAuthToken
	}
	if err := token.Validate(); err != nil {
		return errors.WithStack(err)
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := conn.Create(&AuthToken{Token: token.Token, ExpiresAt: token.ExpiresAt}).Error; err != nil {
		if IsDuplicateEntry(err) {
			return errors.Wrap(domain.ErrDuplicatedTokenBlacklist, err.Error())
		}

		return errors.WithStack(err)
	}

	return nil
}

func (r *TokenBlacklistRepository) Get(c context.Context, token string) (*domain.AuthToken, error) {
	if valid.IsNil(c) {
		return nil, domain.ErrNilContext
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("empty token")
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	record := AuthToken{
		Token: token,
	}
	if err := conn.Where("token = ?", token).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(domain.ErrTokenBlacklistNotFound, err.Error())
		}

		return nil, errors.WithStack(err)
	}

	return &domain.AuthToken{
		Token:     record.Token,
		ExpiresAt: record.ExpiresAt,
	}, nil
}

type AuthToken struct {
	Token     string    `gorm:"token"`
	ExpiresAt time.Time `gorm:"expires_at"`
}

func (t *AuthToken) TableName() string {
	return "token_blacklist"
}
