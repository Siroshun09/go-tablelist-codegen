package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/go-tablelist-codegen/internal"
	"github.com/Siroshun09/serrors/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	internal.ParseFlags()

	if internal.Flag.PackageName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "package name is required")
		os.Exit(1)
	}

	db, err := sql.Open("pgx", generateDSN())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to close database connection: %v\n", err)
			os.Exit(1)
		}
	}(db)

	err = internal.Run(db, internal.Options{
		PackageName: internal.Flag.PackageName,
		Output:      internal.Flag.Output,
		Query:       database.QueryForPostgreSQL,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		if internal.Flag.Debug {
			_, _ = fmt.Fprintln(os.Stderr, serrors.GetStackTrace(err).String())
		}
		os.Exit(1)
	}
}

func generateDSN() string {
	info := internal.Flag.DBConnInfo
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		info.User, info.Password, info.Host, strconv.Itoa(info.Port), info.Database, info.SSLMode,
	)
}
