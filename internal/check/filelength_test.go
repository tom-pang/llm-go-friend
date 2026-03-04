package check

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestCheckFileLength_Over(t *testing.T) {
	// 501 lines: "package x" on line 1, then 500 blank lines.
	src := "package x\n" + strings.Repeat("\n", 500)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "big.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkFileLength(fset, f, "big.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.File != "big.go" {
		t.Errorf("file: got %q, want %q", v.File, "big.go")
	}
	if v.Line != 1 {
		t.Errorf("line: got %d, want 1", v.Line)
	}
	if v.Check != "file_length" {
		t.Errorf("check: got %q, want %q", v.Check, "file_length")
	}
	if v.Value != 501 {
		t.Errorf("value: got %d, want 501", v.Value)
	}
	if v.Threshold != 500 {
		t.Errorf("threshold: got %d, want 500", v.Threshold)
	}
}

func TestCheckFileLength_AtThreshold(t *testing.T) {
	// Exactly 500 lines: should not trigger.
	src := "package x\n" + strings.Repeat("\n", 499)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "ok.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkFileLength(fset, f, "ok.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations at threshold, got %d", len(violations))
	}
}

func TestCheckFileLength_Short(t *testing.T) {
	src := "package x\n"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "short.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkFileLength(fset, f, "short.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestRunAll_IncludesFileLength(t *testing.T) {
	src := "package x\n" + strings.Repeat("\n", 500)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "big.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := RunAll(fset, f, "big.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation from RunAll, got %d", len(violations))
	}
	if violations[0].Check != "file_length" {
		t.Errorf("check: got %q, want %q", violations[0].Check, "file_length")
	}
}
