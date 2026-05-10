package test

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	run(
		t, root,
		"run",
		"./cmd/postgresql",
		"-package-name", "tablelist",
		"-host", host,
		"-port", port,
		"-user", user,
		"-password", password,
		"-database", database,
		"-debug",
	)
}
