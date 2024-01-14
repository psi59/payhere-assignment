package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type contextKey struct{}

var (
	ctxKey = contextKey{}
)

func ConnFromContext(c context.Context) (*gorm.DB, error) {
	if c == nil {
		return nil, fmt.Errorf("nil Context")
	}

	v, ok := c.Value(ctxKey).(*gorm.DB)
	if !ok {
		return nil, ErrNilDB
	}

	return v.WithContext(c), nil
}

func ContextWithConn(c context.Context, db *gorm.DB) context.Context {
	return context.WithValue(c, ctxKey, db)
}
