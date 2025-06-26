package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNoRows(t *testing.T) {
	require.True(t, isNoRows(sql.ErrNoRows))
	require.False(t, isNoRows(errors.New("other error")))
}

func TestErrNotFound(t *testing.T) {
	require.EqualError(t, errNotFound, "not found")
}
