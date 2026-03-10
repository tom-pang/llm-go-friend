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
	if isTestFile(filename) {
		return nil
	}
	var violations []Violation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		complexity := funcComplexity(fn)
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

func funcComplexity(fn *ast.FuncDecl) int {
	complexity := 1
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		complexity += nodeComplexity(n)
		return true
	})
	return complexity
}

func nodeComplexity(n ast.Node) int {
	switch n := n.(type) {
	case *ast.IfStmt:
		return 1
	case *ast.ForStmt:
		return 1
	case *ast.RangeStmt:
		return 1
	case *ast.SwitchStmt:
		return countCaseClauses(n.Body)
	case *ast.TypeSwitchStmt:
		return countCaseClauses(n.Body)
	case *ast.SelectStmt:
		return countCommClauses(n.Body)
	case *ast.BinaryExpr:
		if n.Op == token.LAND || n.Op == token.LOR {
			return 1
		}
	}
	return 0
}

func countCaseClauses(body *ast.BlockStmt) int {
	count := 0
	for _, stmt := range body.List {
		cc := stmt.(*ast.CaseClause)
		if cc.List != nil {
			count++
		}
	}
	return count
}

func countCommClauses(body *ast.BlockStmt) int {
	count := 0
	for _, stmt := range body.List {
		cc := stmt.(*ast.CommClause)
		if cc.Comm != nil {
			count++
		}
	}
	return count
}
