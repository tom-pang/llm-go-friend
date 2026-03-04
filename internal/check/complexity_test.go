package check

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckComplexity_Over(t *testing.T) {
	// Complexity: 1 (base) + 6 ifs + 2 cases + 2 binary ops = 11
	src := `package x

func complex() {
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	switch {
	case true:
	case false:
	default:
	}
	_ = true && false || true
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkComplexity(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Check != "cyclomatic_complexity" {
		t.Errorf("check: got %q, want %q", v.Check, "cyclomatic_complexity")
	}
	if v.Name != "complex" {
		t.Errorf("name: got %q, want %q", v.Name, "complex")
	}
	if v.Value != 11 {
		t.Errorf("value: got %d, want 11", v.Value)
	}
	if v.Threshold != 10 {
		t.Errorf("threshold: got %d, want 10", v.Threshold)
	}
	if v.Line != 3 {
		t.Errorf("line: got %d, want 3", v.Line)
	}
}

func TestCheckComplexity_Under(t *testing.T) {
	src := `package x

func simple() {
	if true {}
	for i := 0; i < 1; i++ {}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkComplexity(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckComplexity_TypeSwitchAndSelect(t *testing.T) {
	// Complexity: 1 (base) + 2 type switch cases + 1 select comm + 6 ifs = 10
	// At threshold — should NOT trigger.
	src := `package x

func medium() {
	switch interface{}(nil).(type) {
	case int:
	case string:
	default:
	}
	select {
	case <-make(chan int):
	default:
	}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
	if true {}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkComplexity(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations at threshold, got %d", len(violations))
	}
}
