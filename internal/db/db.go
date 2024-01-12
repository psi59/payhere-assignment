package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	ErrNilDB  = fmt.Errorf("nil connection")
	ErrNilKey = fmt.Errorf("nil key")
)

func Connect(c Config) (*gorm.DB, error) {
	masterDialector := c.Dialector()
	db, err := gorm.Open(masterDialector)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	stdDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get standard db object: %w", err)
	}

	maxOpenConn := 10
	if c.MaxOpenConns > 0 {
		maxOpenConn = c.MaxOpenConns
	}
	maxIdleConn := maxOpenConn / 2
	if c.MaxIdleConns > 0 {
		maxIdleConn = c.MaxIdleConns
	}
	connMaxLifetime := 10 * time.Minute
	if c.ConnMaxLifetime > 0 {
		connMaxLifetime = time.Duration(c.ConnMaxLifetime) * time.Second
	}
	connMaxIdletime := 3 * time.Minute
	if c.ConnMaxIdleTime > 0 {
		connMaxIdletime = time.Duration(c.ConnMaxIdleTime) * time.Second
	}

	stdDB.SetMaxOpenConns(maxOpenConn)
	stdDB.SetConnMaxLifetime(connMaxLifetime)
	stdDB.SetMaxIdleConns(maxIdleConn)
	stdDB.SetConnMaxIdleTime(connMaxIdletime)

	return db, nil
}

func Transaction(c context.Context, fn func(c context.Context) error, opts ...*sql.TxOptions) error {
	conn, err := ConnFromContext(c)
	if err != nil {
		return errors.Wrap(err, "failed to get connection")
	}

	if err = conn.Transaction(func(tx *gorm.DB) error {
		return fn(ContextWithConn(c, tx))
	}, opts...); err != nil {
		return errors.Wrap(err, "failed to transaction")
	}

	return nil
}
