package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"

	"github.com/pkg/errors"
	"github.com/psi59/payhere-assignment/domain"
	"github.com/psi59/payhere-assignment/internal/db"
	"github.com/psi59/payhere-assignment/internal/valid"
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
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr); mysqlErr.Number == 1062 {
			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}

func (r *TokenBlacklistRepository) IsExists(c context.Context, token string) (bool, error) {
	if valid.IsNil(c) {
		return false, domain.ErrNilContext
	}
	if len(token) == 0 {
		return false, fmt.Errorf("empty token")
	}
	conn, err := db.ConnFromContext(c)
	if err != nil {
		return false, errors.WithStack(err)
	}

	var record struct {
		Exists bool `gorm:"exists"`
	}
	if err := conn.Model(&AuthToken{}).Select("1 AS exists").Where("token = ?", token).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	return record.Exists, nil
}

type AuthToken struct {
	Token     string    `gorm:"token"`
	ExpiresAt time.Time `gorm:"expires_at"`
}

func (t *AuthToken) TableName() string {
	return "token_blacklist"
}
