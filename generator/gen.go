package generator

import (
	_ "embed"
	"io"
	"text/template"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/serrors/v2"
)

//go:embed gen.go.tmpl
var codeTmpl string

// GenerateCode writes source codes to the given writer.
func GenerateCode(w io.Writer, param TemplateParam) error {
	tmpl, err := template.New("tablelist").Parse(codeTmpl)
	if err != nil {
		return serrors.Wrap(err)
	}

	if err := tmpl.Execute(w, param); err != nil {
		return serrors.Wrap(err)
	}

	return nil
}

type TemplateParam struct {
	PackageName string
	Tables      []database.Table
}
