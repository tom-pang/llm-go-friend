package check

import (
	"go/ast"
	"go/token"
)

func init() {
	checkers = append(checkers, checkSingleLetterVar)
}

func checkSingleLetterVar(fset *token.FileSet, file *ast.File, filename string) []Violation {
	loopVars := collectLoopVars(file)

	var violations []Violation
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch node := n.(type) {
		case *ast.AssignStmt:
			if node.Tok != token.DEFINE {
				return true
			}
			for _, lhs := range node.Lhs {
				ident, ok := lhs.(*ast.Ident)
				if !ok {
					continue
				}
				if isSingleLetterViolation(ident, loopVars) {
					violations = append(violations, Violation{
						File:      filename,
						Line:      fset.Position(node.Pos()).Line,
						Check:     "single_letter_var",
						Value:     1,
						Threshold: 0,
					})
				}
			}
		case *ast.ValueSpec:
			for _, ident := range node.Names {
				if isSingleLetterViolation(ident, loopVars) {
					violations = append(violations, Violation{
						File:      filename,
						Line:      fset.Position(node.Pos()).Line,
						Check:     "single_letter_var",
						Value:     1,
						Threshold: 0,
					})
				}
			}
		case *ast.RangeStmt:
			if node.Tok != token.DEFINE {
				return true
			}
			for _, expr := range []ast.Expr{node.Key, node.Value} {
				ident, ok := expr.(*ast.Ident)
				if !ok {
					continue
				}
				if isSingleLetterViolation(ident, loopVars) {
					violations = append(violations, Violation{
						File:      filename,
						Line:      fset.Position(node.Pos()).Line,
						Check:     "single_letter_var",
						Value:     1,
						Threshold: 0,
					})
				}
			}
		}
		return true
	})
	return violations
}

// collectLoopVars walks the AST and records variable names declared in
// for-loop init statements and range-statement key/value positions,
// keyed by their token.Pos so we can distinguish the same name used
// in different contexts.
func collectLoopVars(file *ast.File) map[token.Pos]struct{} {
	vars := make(map[token.Pos]struct{})
	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch node := n.(type) {
		case *ast.ForStmt:
			if assign, ok := node.Init.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
				for _, lhs := range assign.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok {
						vars[ident.Pos()] = struct{}{}
					}
				}
			}
		case *ast.RangeStmt:
			if node.Tok == token.DEFINE {
				if ident, ok := node.Key.(*ast.Ident); ok {
					vars[ident.Pos()] = struct{}{}
				}
				if node.Value != nil {
					if ident, ok := node.Value.(*ast.Ident); ok {
						vars[ident.Pos()] = struct{}{}
					}
				}
			}
		}
		return true
	})
	return vars
}

var loopIndexNames = map[string]struct{}{
	"i": {},
	"j": {},
	"k": {},
}

func isSingleLetterViolation(ident *ast.Ident, loopVars map[token.Pos]struct{}) bool {
	if ident.Name == "_" {
		return false
	}
	if len(ident.Name) != 1 {
		return false
	}
	// i, j, k are allowed in loop init/range positions
	if _, isLoopIndex := loopIndexNames[ident.Name]; isLoopIndex {
		if _, inLoop := loopVars[ident.Pos()]; inLoop {
			return false
		}
	}
	return true
}
