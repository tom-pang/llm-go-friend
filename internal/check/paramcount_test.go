package check

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckParamCount_Over(t *testing.T) {
	src := `package x

func tooMany(a, b, c int, d string, e float64, f int) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkParamCount(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Check != "param_count" {
		t.Errorf("check: got %q, want %q", v.Check, "param_count")
	}
	if v.Name != "tooMany" {
		t.Errorf("name: got %q, want %q", v.Name, "tooMany")
	}
	if v.Value != 6 {
		t.Errorf("value: got %d, want 6", v.Value)
	}
	if v.Threshold != 5 {
		t.Errorf("threshold: got %d, want 5", v.Threshold)
	}
	if v.Line != 3 {
		t.Errorf("line: got %d, want 3", v.Line)
	}
}

func TestCheckParamCount_Under(t *testing.T) {
	src := `package x

func fewParams(a int, b string, c float64) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkParamCount(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckParamCount_Unnamed(t *testing.T) {
	// Unnamed params (common in interface implementations)
	src := `package x

func unnamed(int, string, float64, int, int, int) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkParamCount(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Value != 6 {
		t.Errorf("value: got %d, want 6", violations[0].Value)
	}
}

func TestCheckParamCount_Mixed(t *testing.T) {
	// Mix of named and grouped params: a, b share int; c is alone; d, e, f share string = 6
	src := `package x

func mixed(a, b int, c float64, d, e, f string) {}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkParamCount(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Value != 6 {
		t.Errorf("value: got %d, want 6", violations[0].Value)
	}
}
