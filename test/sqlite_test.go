package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLite(t *testing.T) {
	root, err := getProjectRoot()
	if err != nil {
		t.Fatal(err)
	}

	args := []string{
		"run",
		"./cmd/sqlite",
		"-package-name", "sqlitedb",
		"-schema-file", "test/schema_sqlite.sql",
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
