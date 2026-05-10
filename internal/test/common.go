package test

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/Siroshun09/serrors/v2"
)

//go:embed expected/output.go
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
