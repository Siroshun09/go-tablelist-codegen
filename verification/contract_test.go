// Regression tests for database reads and generated source.
package verification_test

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Siroshun09/go-tablelist-codegen/database"
	"github.com/Siroshun09/go-tablelist-codegen/generator"
	entry "github.com/Siroshun09/go-tablelist-codegen/internal"
	_ "modernc.org/sqlite"
)

var errReadFailure = errors.New("injected rows.Next failure")

type connector struct{ fail bool }
type connection struct{ fail bool }
type testDriver struct{}
type testRows struct {
	index int
	fail  bool
}

func (c connector) Connect(context.Context) (driver.Conn, error) { return connection(c), nil }
func (connector) Driver() driver.Driver                          { return testDriver{} }
func (testDriver) Open(string) (driver.Conn, error)              { return connection{}, nil }
func (connection) Prepare(string) (driver.Stmt, error)           { return nil, errors.New("unused") }
func (connection) Close() error                                  { return nil }
func (connection) Begin() (driver.Tx, error)                     { return nil, errors.New("unused") }
func (c connection) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &testRows{fail: c.fail}, nil
}
func (*testRows) Columns() []string { return []string{"table_name", "column_name"} }
func (*testRows) Close() error      { return nil }
func (r *testRows) Next(dest []driver.Value) error {
	if r.index == 0 {
		r.index++
		dest[0], dest[1] = "users", "id"
		return nil
	}
	if r.fail {
		return errReadFailure
	}
	return io.EOF
}

func faultDB(t *testing.T, fail bool) *sql.DB {
	t.Helper()
	db := sql.OpenDB(connector{fail})
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func sqliteDB(t *testing.T, schema string) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("invalid test schema: %v", err)
	}
	return db
}

func checkedFile(src string) (*ast.File, error) {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, "generated.go", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	_, err = new(types.Config).Check("generated", fs, []*ast.File{f}, nil)
	return f, err
}

func fromSchema(t *testing.T, schema string) string {
	t.Helper()
	db := sqliteDB(t, schema)
	tables, err := database.GetTables(db, database.QueryForSQLite)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := generator.GenerateCode(&out, generator.TemplateParam{PackageName: "tables", Tables: tables}); err != nil {
		t.Fatal(err)
	}
	return out.String()
}

func TestValidSchema(t *testing.T) {
	src := fromSchema(t, "CREATE TABLE users (id INTEGER, display_name TEXT);")
	if _, err := checkedFile(src); err != nil {
		t.Fatal(err)
	}
}

func TestNormalEOF(t *testing.T) {
	tables, err := database.GetTables(faultDB(t, false), "test")
	if err != nil || len(tables) != 1 || len(tables[0].Columns) != 1 {
		t.Fatalf("normal EOF: tables=%#v err=%v", tables, err)
	}
}

func TestReadError(t *testing.T) {
	tables, err := database.GetTables(faultDB(t, true), "test")
	if !errors.Is(err, errReadFailure) {
		t.Fatalf("must propagate iteration error; got err=%v, partial tables=%#v", err, tables)
	}
}

func TestOutputPreservedOnReadError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tables.go")
	before := []byte("// previously correct output\n")
	if err := os.WriteFile(path, before, 0600); err != nil {
		t.Fatal(err)
	}
	err := entry.Run(faultDB(t, true), entry.Options{PackageName: "tables", Query: "test", Output: path})
	after, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !errors.Is(err, errReadFailure) || !bytes.Equal(before, after) {
		t.Fatalf("failed query must not overwrite output: err=%v, preserved=%v, output=%s", err, bytes.Equal(before, after), after)
	}
}

func TestInvalidSchemaRejectedBeforeWriting(t *testing.T) {
	for _, tc := range []struct{ name, schema string }{
		{"reserved_field", "CREATE TABLE users (table_name TEXT);"},
		{"column_collision", "CREATE TABLE users (a_b TEXT, a__b TEXT);"},
		{"table_collision", "CREATE TABLE a_b (id INTEGER); CREATE TABLE a__b (id INTEGER);"},
		{"digit_start", "CREATE TABLE \"1users\" (id INTEGER);"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := sqliteDB(t, tc.schema)
			tables, err := database.GetTables(db, database.QueryForSQLite)
			if err != nil {
				t.Fatal(err)
			}
			var out bytes.Buffer
			err = generator.GenerateCode(&out, generator.TemplateParam{PackageName: "tables", Tables: tables})
			if err == nil || out.Len() != 0 {
				t.Fatalf("expected validation error before writing: err=%v output=%q", err, out.String())
			}

			path := filepath.Join(t.TempDir(), "tables.go")
			before := []byte("// previously correct output\n")
			if err := os.WriteFile(path, before, 0600); err != nil {
				t.Fatal(err)
			}
			if err := entry.Run(db, entry.Options{PackageName: "tables", Query: database.QueryForSQLite, Output: path}); err == nil {
				t.Fatal("Run accepted invalid schema")
			}
			after, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatalf("invalid schema changed output: read err=%v output=%q", err, after)
			}
		})
	}
}

func TestInvalidExplicitIdentifiersRejectedBeforeWriting(t *testing.T) {
	for _, tc := range []struct {
		name  string
		param generator.TemplateParam
	}{
		{"package_keyword", generator.TemplateParam{PackageName: "package"}},
		{"package_blank", generator.TemplateParam{PackageName: "_"}},
		{"table_keyword", generator.TemplateParam{PackageName: "tables", Tables: []database.Table{{TableName: "users", FieldName: "for"}}}},
		{"column_blank", generator.TemplateParam{PackageName: "tables", Tables: []database.Table{{
			TableName: "users", FieldName: "UsersTable",
			Columns: []database.Column{{ColumnName: "id", FieldName: "_"}},
		}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := generator.GenerateCode(&out, tc.param); err == nil || out.Len() != 0 {
				t.Fatalf("expected validation error before writing: err=%v output=%q", err, out.String())
			}
		})
	}
}

func TestLiteralRoundTrip(t *testing.T) {
	// Explicit, valid FieldName isolates literal encoding from name normalization.
	for _, raw := range []string{`a\nb`, "a\"b", "a\nb"} {
		t.Run(strconv.Quote(raw), func(t *testing.T) {
			var out bytes.Buffer
			p := generator.TemplateParam{PackageName: "tables", Tables: []database.Table{
				{TableName: "users", FieldName: "UsersTable", Columns: []database.Column{
					{ColumnName: raw, FieldName: "Value"},
				}},
			}}
			if err := generator.GenerateCode(&out, p); err != nil {
				t.Fatal(err)
			}
			f, err := checkedFile(out.String())
			if err != nil {
				t.Fatalf("literal broke Go syntax: %v\n%s", err, out.String())
			}
			found := false
			ast.Inspect(f, func(n ast.Node) bool {
				kv, ok := n.(*ast.KeyValueExpr)
				if !ok {
					return true
				}
				key, ok := kv.Key.(*ast.Ident)
				if !ok || key.Name != "Value" {
					return true
				}
				lit, ok := kv.Value.(*ast.BasicLit)
				if !ok {
					t.Fatal("not literal")
				}
				value, err := strconv.Unquote(lit.Value)
				if err != nil || value != raw {
					t.Errorf("name changed: want=%q got=%q err=%v", raw, value, err)
				}
				found = true
				return true
			})
			if !found {
				t.Fatal("no value assertion executed")
			}
		})
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestWriterError(t *testing.T) {
	if err := generator.GenerateCode(failedWriter{}, generator.TemplateParam{PackageName: "tables"}); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("writer error lost: %v", err)
	}
}

// GenerateCode must preserve arbitrary Go string bytes when explicit identifiers
// are valid. Every successful iteration checks both the table and column value.
func FuzzLiteralRoundTrip(f *testing.F) {
	for _, s := range []string{"users", `a\nb`, "a\"b", "a\nb", "", string([]byte{0xff, 0})} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		var out bytes.Buffer
		param := generator.TemplateParam{PackageName: "tables", Tables: []database.Table{
			{TableName: raw, FieldName: "UsersTable", Columns: []database.Column{
				{ColumnName: raw, FieldName: "Value"},
			}},
		}}
		if err := generator.GenerateCode(&out, param); err != nil {
			t.Fatal(err)
		}
		node, err := checkedFile(out.String())
		if err != nil {
			t.Fatalf("invalid generated source: %v", err)
		}
		count := 0
		ast.Inspect(node, func(n ast.Node) bool {
			kv, ok := n.(*ast.KeyValueExpr)
			if !ok {
				return true
			}
			lit, ok := kv.Value.(*ast.BasicLit)
			if !ok {
				t.Fatal("expected literal")
			}
			got, err := strconv.Unquote(lit.Value)
			if err != nil || got != raw {
				t.Fatalf("literal not preserved: want=%q got=%q err=%v", raw, got, err)
			}
			count++
			return true
		})
		if count != 2 {
			t.Fatalf("expected two round-trip assertions, got %d", count)
		}
	})
}
