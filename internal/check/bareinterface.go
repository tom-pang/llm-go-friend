package check

import (
	"go/ast"
	"go/token"
)

func init() {
	checkers = append(checkers, checkBareInterface)
}

func checkBareInterface(fset *token.FileSet, file *ast.File, filename string) []Violation {
	var violations []Violation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if !ast.IsExported(fn.Name.Name) {
			continue
		}
		if fn.Recv != nil && !receiverIsExported(fn.Recv) {
			continue
		}
		if fn.Type.Params == nil {
			continue
		}
		for _, field := range fn.Type.Params.List {
			if isBareInterface(field.Type) {
				violations = append(violations, Violation{
					File:      filename,
					Line:      fset.Position(field.Pos()).Line,
					Check:     "bare_interface",
					Name:      fn.Name.Name,
					Value:     1,
					Threshold: 0,
				})
			}
		}
	}
	return violations
}

func receiverIsExported(recv *ast.FieldList) bool {
	if recv == nil || len(recv.List) == 0 {
		return false
	}
	typ := recv.List[0].Type
	// Unwrap pointer receiver: (*T) → T
	if star, ok := typ.(*ast.StarExpr); ok {
		typ = star.X
	}
	if ident, ok := typ.(*ast.Ident); ok {
		return ast.IsExported(ident.Name)
	}
	return false
}

func isBareInterface(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.InterfaceType:
		return t.Methods == nil || len(t.Methods.List) == 0
	case *ast.Ident:
		return t.Name == "any"
	}
	return false
}
