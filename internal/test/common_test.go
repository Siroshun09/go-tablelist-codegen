package test

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Siroshun09/serrors/v2"
	"github.com/stretchr/testify/assert"
)

//go:embed expected/output.gen.go
var expectedOutput string

// getProjectRoot searches for the project root directory containing go.mod
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", serrors.Wrap(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", serrors.Wrap(os.ErrNotExist)
}

func run(t *testing.T, root string, args ...string) {
	t.Helper()

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
