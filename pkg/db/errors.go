package db

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrUniqueValueViolation = errors.New("duplicate entry")
)

var mysqlDialectErrorMapping = map[uint16]error{
	1062: ErrUniqueValueViolation,
}

func IsError(dialectError error, genericError error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(dialectError, &mysqlErr) {
		return errors.Is(mysqlDialectErrorMapping[mysqlErr.Number], genericError)
	}

	return false
}
