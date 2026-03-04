package check

import (
	"go/ast"
	"go/token"
)

const paramCountThreshold = 5

func init() {
	checkers = append(checkers, checkParamCount)
}

func checkParamCount(fset *token.FileSet, file *ast.File, filename string) []Violation {
	var violations []Violation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Type.Params == nil {
			continue
		}
		count := 0
		for _, field := range fn.Type.Params.List {
			if len(field.Names) == 0 {
				count++ // unnamed parameter
			} else {
				count += len(field.Names)
			}
		}
		if count <= paramCountThreshold {
			continue
		}
		violations = append(violations, Violation{
			File:      filename,
			Line:      fset.Position(fn.Pos()).Line,
			Check:     "param_count",
			Name:      fn.Name.Name,
			Value:     count,
			Threshold: paramCountThreshold,
		})
	}
	return violations
}
