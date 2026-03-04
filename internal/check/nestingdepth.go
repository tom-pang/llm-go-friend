package check

import (
	"go/ast"
	"go/token"
)

const nestingDepthThreshold = 4

func init() {
	checkers = append(checkers, checkNestingDepth)
}

func checkNestingDepth(fset *token.FileSet, file *ast.File, filename string) []Violation {
	var violations []Violation
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		maxDepth := 0
		var deepestNode ast.Node
		walkNesting(fn.Body, 0, &maxDepth, &deepestNode)
		if maxDepth <= nestingDepthThreshold {
			continue
		}
		violationLine := fset.Position(fn.Pos()).Line
		if deepestNode != nil {
			violationLine = fset.Position(deepestNode.Pos()).Line
		}
		violations = append(violations, Violation{
			File:      filename,
			Line:      violationLine,
			Check:     "nesting_depth",
			Name:      fn.Name.Name,
			Value:     maxDepth,
			Threshold: nestingDepthThreshold,
		})
	}
	return violations
}

func walkNesting(node ast.Node, depth int, maxDepth *int, deepestNode *ast.Node) {
	if node == nil {
		return
	}
	ast.Inspect(node, func(n ast.Node) bool {
		if n == nil {
			return false
		}
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt,
			*ast.FuncLit:
			newDepth := depth + 1
			if newDepth > *maxDepth {
				*maxDepth = newDepth
				*deepestNode = n
			}
			// Recurse into children with incremented depth, then stop
			// default traversal for this node.
			walkNestingChildren(n, newDepth, maxDepth, deepestNode)
			return false
		}
		return true
	})
}

func walkNestingChildren(node ast.Node, depth int, maxDepth *int, deepestNode *ast.Node) {
	switch n := node.(type) {
	case *ast.IfStmt:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
		// Else can be *ast.BlockStmt or *ast.IfStmt.
		// Else blocks don't add depth — only the if itself does.
		// But a nested if inside else does add depth via the normal
		// ast.IfStmt case above.
		walkNesting(n.Else, depth, maxDepth, deepestNode)
	case *ast.ForStmt:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
	case *ast.RangeStmt:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
	case *ast.SwitchStmt:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
	case *ast.TypeSwitchStmt:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
	case *ast.SelectStmt:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
	case *ast.FuncLit:
		walkNesting(n.Body, depth, maxDepth, deepestNode)
	}
}
