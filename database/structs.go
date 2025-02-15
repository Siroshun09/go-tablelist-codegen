package database

import (
	"strings"
	"unicode"
)

// Table holds TableName and Columns.
//
// FieldName is upper camel case of TableName.
type Table struct {
	TableName string
	FieldName string
	Columns   []Column
}

// Column holds ColumnName.
//
// FieldName is upper camel case of ColumnName.
type Column struct {
	ColumnName string
	FieldName  string
}

func toFieldName(name string) string {
	b := strings.Builder{}
	b.Grow(len(name))
	shouldUppercase := false
	for i, c := range name {
		if i == 0 {
			b.WriteRune(unicode.ToUpper(c))
			continue
		}

		if c == '_' {
			shouldUppercase = true
			continue
		}

		if shouldUppercase {
			b.WriteRune(unicode.ToUpper(c))
			shouldUppercase = false
		} else {
			b.WriteRune(unicode.ToLower(c))
		}
	}
	return b.String()
}
