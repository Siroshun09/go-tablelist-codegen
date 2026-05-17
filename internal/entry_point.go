package internal

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/go-tablelist-codegen/generator"
	"github.com/Siroshun09/serrors/v2"
)

type Options struct {
	PackageName string
	Query       string
	Output      string
}

func Run(db database.DB, opts Options) (returnErr error) {
	tables, err := database.GetTables(db, opts.Query)
	if err != nil {
		return serrors.Wrap(fmt.Errorf("failed to get table list: %w", err))
	}

	var w io.Writer
	if opts.Output != "" {
		dir := filepath.Dir(opts.Output)
		err := os.MkdirAll(dir, 0o750)
		if err != nil {
			return serrors.Wrap(fmt.Errorf("failed to create output directory: %w", err), slog.String("dir", dir))
		}

		file, err := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return serrors.Wrap(fmt.Errorf("failed to open output file: %w", err), slog.String("file", opts.Output))
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				returnErr = errors.Join(returnErr, serrors.Wrap(
					fmt.Errorf("failed to close output file: %w", err),
					slog.String("file", opts.Output),
				))
			}
		}(file)

		w = file
	} else {
		w = os.Stdout
	}

	err = generator.GenerateCode(w, generator.TemplateParam{
		PackageName: opts.PackageName,
		Tables:      tables,
	})
	if err != nil {
		return serrors.Wrap(fmt.Errorf("failed to generate code: %w", err))
	}

	return nil
}
