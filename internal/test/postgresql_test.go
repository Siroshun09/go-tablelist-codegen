package test

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgreSQL(t *testing.T) {
	if os.Getenv("ENABLE_POSTGRESQL_TEST") != "true" {
		t.Skip("ENABLE_POSTGRESQL_TEST is not true, skipping PostgreSQL test")
	}

	host := os.Getenv("POSTGRESQL_HOST")
	port := os.Getenv("POSTGRESQL_PORT")
	user := os.Getenv("POSTGRESQL_USER")
	password := os.Getenv("POSTGRESQL_PASSWORD")
	database := os.Getenv("POSTGRESQL_DATABASE")

	root, err := getProjectRoot()
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	schemaSQL, err := os.ReadFile(filepath.Join(root, "internal/test/schema_postgresql.sql"))
	require.NoError(t, err)

	_, err = db.ExecContext(t.Context(), string(schemaSQL))
	require.NoError(t, err)

	args := []string{
		"run",
		"./cmd/postgresql",
		"-package-name", "tablelist",
		"-host", host,
		"-port", port,
		"-user", user,
		"-password", password,
		"-database", database,
		"-debug",
	}

	stdout := strings.Builder{}
	stderr := strings.Builder{}

	cmd := exec.CommandContext(t.Context(), "go", args...)
	cmd.Dir = root
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if !assert.NoError(t, cmd.Run()) {
		t.Log(stderr.String())
		return
	}

	assert.Equal(t, expectedOutput, stdout.String())
	assert.Empty(t, stderr.String())
}
