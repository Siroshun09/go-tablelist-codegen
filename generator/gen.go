package generator

import (
	_ "embed"
	"github.com/Siroshun09/go-tablelist-codegen/database"
	"io"
	"text/template"

	"github.com/Siroshun09/serrors"
)

//go:embed gen.go.tmpl
var codeTmpl string

// GenerateCode writes source codes to the given writer.
func GenerateCode(w io.Writer, param TemplateParam) error {
	tmpl, err := template.New("tablelist").Parse(codeTmpl)
	if err != nil {
		return serrors.WithStackTrace(err)
	}

	if err := tmpl.Execute(w, param); err != nil {
		return serrors.WithStackTrace(err)
	}

	return nil
}

type TemplateParam struct {
	PackageName string
	Tables      []database.Table
}
