package database

import (
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/Siroshun09/serrors"
)

const (
	QueryForMySQL = "SELECT TABLE_NAME AS table_name, COLUMN_NAME AS column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE();"
)

// DB is an interface to execute queries.
type DB interface {
	// Query executes the given query with args.
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// GetTables gets all Table from the database by the given query.
//
// The returning slice of Table is sorted by table name.
func GetTables(db DB, query string) (tables []Table, returnErr error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, serrors.WithStackTrace(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			returnErr = serrors.WithStackTrace(errors.Join(returnErr, err))
		}
	}(rows)

	type column struct {
		TableName  string `db:"table_name"`
		ColumnName string `db:"column_name"`
	}

	var columns []column
	for rows.Next() {
		var col column
		err = rows.Scan(&col.TableName, &col.ColumnName)
		if err != nil {
			return nil, serrors.WithStackTrace(err)
		}
		columns = append(columns, col)
	}

	byTableName := make(map[string][]column)
	for _, col := range columns {
		byTableName[col.TableName] = append(byTableName[col.TableName], col)
	}

	for name, cols := range byTableName {
		cs := make([]Column, 0, len(cols))
		for _, col := range cols {
			cs = append(cs, Column{
				ColumnName: col.ColumnName,
				FieldName:  toFieldName(col.ColumnName),
			})
		}
		tables = append(tables, Table{
			TableName: name,
			FieldName: toFieldName(name) + "Table",
			Columns:   cs,
		})
	}

	slices.SortFunc(tables, func(a, b Table) int {
		return strings.Compare(a.TableName, b.TableName)
	})

	return tables, nil
}
