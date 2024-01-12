package mysql

import (
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	ErrCodeDuplicateEntry = 1062
)

func IsDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == ErrCodeDuplicateEntry
	}

	return false
}
