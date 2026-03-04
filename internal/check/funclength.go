package check

import (
	"go/ast"
	"go/token"
)

const funcLengthThreshold = 50

func init() {
	checkers = append(checkers, checkFuncLength)
}

func checkFuncLength(fset *token.FileSet, file *ast.File, filename string) []Violation {
	var violations []Violation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		startLine := fset.Position(fn.Pos()).Line
		endLine := fset.Position(fn.End()).Line
		lineCount := endLine - startLine + 1
		if lineCount <= funcLengthThreshold {
			continue
		}
		violations = append(violations, Violation{
			File:      filename,
			Line:      startLine,
			Check:     "func_length",
			Name:      fn.Name.Name,
			Value:     lineCount,
			Threshold: funcLengthThreshold,
		})
	}
	return violations
}
