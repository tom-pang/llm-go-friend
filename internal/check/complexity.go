package check

import (
	"go/ast"
	"go/token"
)

const complexityThreshold = 10

func init() {
	checkers = append(checkers, checkComplexity)
}

func checkComplexity(fset *token.FileSet, file *ast.File, filename string) []Violation {
	var violations []Violation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		complexity := 1 // the function itself is one path
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			switch n := n.(type) {
			case *ast.IfStmt:
				complexity++
			case *ast.ForStmt:
				complexity++
			case *ast.RangeStmt:
				complexity++
			case *ast.SwitchStmt:
				for _, s := range n.Body.List {
					cc := s.(*ast.CaseClause)
					if cc.List != nil { // skip default
						complexity++
					}
				}
			case *ast.TypeSwitchStmt:
				for _, s := range n.Body.List {
					cc := s.(*ast.CaseClause)
					if cc.List != nil {
						complexity++
					}
				}
			case *ast.SelectStmt:
				for _, s := range n.Body.List {
					cc := s.(*ast.CommClause)
					if cc.Comm != nil { // nil Comm = default
						complexity++
					}
				}
			case *ast.BinaryExpr:
				if n.Op == token.LAND || n.Op == token.LOR {
					complexity++
				}
			}
			return true
		})
		if complexity <= complexityThreshold {
			continue
		}
		violations = append(violations, Violation{
			File:      filename,
			Line:      fset.Position(fn.Pos()).Line,
			Check:     "cyclomatic_complexity",
			Name:      fn.Name.Name,
			Value:     complexity,
			Threshold: complexityThreshold,
		})
	}
	return violations
}
