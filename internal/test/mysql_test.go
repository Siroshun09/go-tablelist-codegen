package test

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQL(t *testing.T) {
	if os.Getenv("ENABLE_MYSQL_TEST") != "true" {
		t.Skip("ENABLE_MYSQL_TEST is not true, skipping MySQL test")
	}

	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")

	root, err := getProjectRoot()
	require.NoError(t, err)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	schemaSQL, err := os.ReadFile(filepath.Join(root, "internal/test/schema_mysql.sql"))
	require.NoError(t, err)

	_, err = db.ExecContext(t.Context(), string(schemaSQL))
	require.NoError(t, err)

	args := []string{
		"run",
		"./cmd/mysql",
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
