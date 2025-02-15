package internal

import (
	"errors"
	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/go-tablelist-codegen/generator"
	"io"
	"os"
	"path/filepath"

	"github.com/Siroshun09/serrors"
)

type Options struct {
	PackageName string
	Query       string
	Output      string
}

func Run(db database.DB, opts Options) (returnErr error) {
	tables, err := database.GetTables(db, opts.Query)
	if err != nil {
		return serrors.Errorf("failed to get table list: %v", err)
	}

	var w io.Writer
	if Flag.Output != "" {
		dir := filepath.Dir(opts.Output)
		err = os.MkdirAll(dir, 0o750)
		if err != nil {
			return serrors.Errorf("failed to create output directory %s: %v", dir, err)
		}

		file, err := os.OpenFile(opts.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			return serrors.Errorf("failed to open output file: %v", err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				returnErr = errors.Join(returnErr, serrors.Errorf("failed to close output file: %v", err))
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
		return serrors.Errorf("failed to generate code: %v", err)
	}

	return nil
}
