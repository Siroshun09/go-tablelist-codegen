package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/go-tablelist-codegen/internal"
	"github.com/Siroshun09/serrors"
	"github.com/go-sql-driver/mysql"
)

func main() {
	internal.ParseFlags()

	if internal.Flag.PackageName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "package name is required")
		os.Exit(1)
	}

	db, err := sql.Open("mysql", generateDBConfig().FormatDSN())
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
		Query:       database.QueryForMySQL,
	})
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		if internal.Flag.Debug {
			_, _ = fmt.Fprintln(os.Stderr, serrors.GetStackTrace(err).String())
		}
		os.Exit(1)
	}
}

func generateDBConfig() *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.User = internal.Flag.DBConnInfo.User
	cfg.Passwd = internal.Flag.DBConnInfo.Password
	cfg.Net = "tcp"
	cfg.Addr = internal.Flag.DBConnInfo.Host + ":" + strconv.Itoa(internal.Flag.DBConnInfo.Port)
	cfg.DBName = internal.Flag.DBConnInfo.Database
	return cfg
}
