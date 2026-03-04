package check

import (
	"go/ast"
	"go/token"
)

// Violation represents a single check failure.
type Violation struct {
	File      string `toon:"file"`
	Line      int    `toon:"line"`
	Check     string `toon:"check"`
	Name      string `toon:"name"`
	Value     int    `toon:"value"`
	Threshold int    `toon:"threshold"`
}

// Checker analyzes a parsed Go file and returns any violations found.
type Checker func(fset *token.FileSet, file *ast.File, filename string) []Violation

// checkers is the registry of all active checkers.
// Each checker file appends to this slice via init() — acceptable here
// since it's registering constant config, not mutable state.
var checkers []Checker

// RunAll runs every registered checker against the given file and
// returns all violations.
func RunAll(fset *token.FileSet, file *ast.File, filename string) []Violation {
	var violations []Violation
	for _, checker := range checkers {
		violations = append(violations, checker(fset, file, filename)...)
	}
	return violations
}
