package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLite(t *testing.T) {
	root, err := getProjectRoot()
	require.NoError(t, err)

	run(
		t, root,
		"run",
		"./cmd/sqlite",
		"-package-name", "tablelist",
		"-schema-file", "internal/test/schema_sqlite.sql",
		"-debug",
	)
}
