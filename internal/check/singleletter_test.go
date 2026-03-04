package check

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckSingleLetterVar_OutsideLoop(t *testing.T) {
	src := `package x

func foo() {
	x := 5
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkSingleLetterVar(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Check != "single_letter_var" {
		t.Errorf("check: got %q, want %q", v.Check, "single_letter_var")
	}
	if v.Value != 1 {
		t.Errorf("value: got %d, want 1", v.Value)
	}
	if v.Threshold != 0 {
		t.Errorf("threshold: got %d, want 0", v.Threshold)
	}
}

func TestCheckSingleLetterVar_ForLoopInit(t *testing.T) {
	src := `package x

func foo() {
	for i := 0; i < 10; i++ {
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkSingleLetterVar(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckSingleLetterVar_RangeStmt(t *testing.T) {
	// k is allowed in range key position, v is NOT in the allowed set
	src := `package x

func foo() {
	m := map[string]int{"a": 1}
	for k, v := range m {
		_ = k
		_ = v
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkSingleLetterVar(fset, f, "test.go")
	// m is single-letter outside loop → violation
	// k is allowed in range key → no violation
	// v is single-letter but not i/j/k → violation
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations (m and v), got %d", len(violations))
	}
}

func TestCheckSingleLetterVar_BlankIdentifier(t *testing.T) {
	src := `package x

func foo() {
	_ = 5
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkSingleLetterVar(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckSingleLetterVar_MultiLetterName(t *testing.T) {
	src := `package x

func foo() {
	count := 5
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkSingleLetterVar(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckSingleLetterVar_VarDecl(t *testing.T) {
	src := `package x

func foo() {
	var x int
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkSingleLetterVar(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].Check != "single_letter_var" {
		t.Errorf("check: got %q, want %q", violations[0].Check, "single_letter_var")
	}
}
