package database

import (
	"database/sql"
	"errors"
)

// errNotFound is a sentinel error for not found cases
var errNotFound = errors.New("not found")

// isNoRows returns true if the error is sql.ErrNoRows
func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
