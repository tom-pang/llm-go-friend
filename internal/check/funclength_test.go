package check

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCheckFuncLength_Over(t *testing.T) {
	// 60-line function: func keyword on line 3, body has enough lines
	// to span 60 lines total (line 3 through line 62).
	body := strings.Repeat("\t_ = 0\n", 58)
	src := "package x\n\nfunc big() {\n" + body + "}\n"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkFuncLength(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Check != "func_length" {
		t.Errorf("check: got %q, want %q", v.Check, "func_length")
	}
	if v.Name != "big" {
		t.Errorf("name: got %q, want %q", v.Name, "big")
	}
	if v.Value != 60 {
		t.Errorf("value: got %d, want 60", v.Value)
	}
	if v.Threshold != 50 {
		t.Errorf("threshold: got %d, want 50", v.Threshold)
	}
	if v.Line != 3 {
		t.Errorf("line: got %d, want 3", v.Line)
	}
}

func TestCheckFuncLength_Under(t *testing.T) {
	body := strings.Repeat("\t_ = 0\n", 28)
	src := "package x\n\nfunc small() {\n" + body + "}\n"

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkFuncLength(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}
