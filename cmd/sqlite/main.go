package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/go-tablelist-codegen/internal"
	"github.com/Siroshun09/serrors/v2"
	_ "modernc.org/sqlite"
)

func main() {
	internal.ParseFlags()

	if internal.Flag.PackageName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "package name is required")
		os.Exit(1)
	}

	if internal.Flag.SchemaFile == "" {
		_, _ = fmt.Fprintln(os.Stderr, "schema file is required")
		os.Exit(1)
	}

	schema, err := os.ReadFile(internal.Flag.SchemaFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read schema file: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create memory database: %v\n", err)
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close database connection: %v\n", err)
			os.Exit(1)
		}
	}(db)

	_, err = db.Exec(string(schema))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to execute schema query: %v\n", err)
		os.Exit(1)
	}

	err = internal.Run(db, internal.Options{
		PackageName: internal.Flag.PackageName,
		Output:      internal.Flag.Output,
		Query:       database.QueryForSQLite,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		if internal.Flag.Debug {
			_, _ = fmt.Fprintln(os.Stderr, serrors.GetStackTrace(err).String())
		}
		os.Exit(1)
	}
}
