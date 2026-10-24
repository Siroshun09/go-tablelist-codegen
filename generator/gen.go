package generator

import (
	_ "embed"
	"fmt"
	"go/token"
	"io"
	"text/template"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/serrors/v2"
)

//go:embed gen.go.tmpl
var codeTmpl string

// GenerateCode validates identifiers before writing source code to the given writer.
func GenerateCode(w io.Writer, param TemplateParam) error {
	if err := validateParam(param); err != nil {
		return serrors.Wrap(err)
	}

	tmpl, err := template.New("tablelist").Parse(codeTmpl)
	if err != nil {
		return serrors.Wrap(err)
	}

	if err := tmpl.Execute(w, param); err != nil {
		return serrors.Wrap(err)
	}

	return nil
}

func validateParam(param TemplateParam) error {
	if !token.IsIdentifier(param.PackageName) || param.PackageName == "_" {
		return fmt.Errorf("invalid package name %q", param.PackageName)
	}

	tableFields := make(map[string]struct{}, len(param.Tables))
	for _, table := range param.Tables {
		if !token.IsIdentifier(table.FieldName) || table.FieldName == "_" {
			return fmt.Errorf("invalid field name %q for table %q", table.FieldName, table.TableName)
		}
		if _, exists := tableFields[table.FieldName]; exists {
			return fmt.Errorf("duplicate field name %q for table %q", table.FieldName, table.TableName)
		}
		tableFields[table.FieldName] = struct{}{}

		columnFields := map[string]struct{}{"TableName": {}}
		for _, column := range table.Columns {
			if !token.IsIdentifier(column.FieldName) || column.FieldName == "_" {
				return fmt.Errorf("invalid field name %q for column %q in table %q", column.FieldName, column.ColumnName, table.TableName)
			}
			if _, exists := columnFields[column.FieldName]; exists {
				return fmt.Errorf("duplicate field name %q for column %q in table %q", column.FieldName, column.ColumnName, table.TableName)
			}
			columnFields[column.FieldName] = struct{}{}
		}
	}
	return nil
}

type TemplateParam struct {
	PackageName string
	Tables      []database.Table
}
