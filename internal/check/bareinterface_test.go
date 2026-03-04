package check

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckBareInterface_AnyParam(t *testing.T) {
	src := `package x

func DoStuff(v any) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkBareInterface(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Check != "bare_interface" {
		t.Errorf("check: got %q, want %q", v.Check, "bare_interface")
	}
	if v.Name != "DoStuff" {
		t.Errorf("name: got %q, want %q", v.Name, "DoStuff")
	}
	if v.Value != 1 {
		t.Errorf("value: got %d, want 1", v.Value)
	}
	if v.Threshold != 0 {
		t.Errorf("threshold: got %d, want 0", v.Threshold)
	}
}

func TestCheckBareInterface_EmptyInterface(t *testing.T) {
	src := `package x

func DoStuff(v interface{}) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkBareInterface(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].Check != "bare_interface" {
		t.Errorf("check: got %q, want %q", violations[0].Check, "bare_interface")
	}
}

func TestCheckBareInterface_UnexportedFunc(t *testing.T) {
	src := `package x

func doStuff(v any) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkBareInterface(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckBareInterface_NamedInterface(t *testing.T) {
	src := `package x

import "io"

func DoStuff(r io.Reader) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkBareInterface(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckBareInterface_ExportedMethod(t *testing.T) {
	src := `package x

type Server struct{}

func (s *Server) Handle(v any) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkBareInterface(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].Name != "Handle" {
		t.Errorf("name: got %q, want %q", violations[0].Name, "Handle")
	}
}

func TestCheckBareInterface_UnexportedReceiverType(t *testing.T) {
	src := `package x

type server struct{}

func (s *server) Handle(v any) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkBareInterface(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}
