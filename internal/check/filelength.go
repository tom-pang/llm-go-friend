package check

import (
	"go/ast"
	"go/token"
)

const fileLengthThreshold = 500

func init() {
	checkers = append(checkers, checkFileLength)
}

func checkFileLength(fset *token.FileSet, file *ast.File, filename string) []Violation {
	lineCount := fset.File(file.Pos()).LineCount()
	if lineCount <= fileLengthThreshold {
		return nil
	}
	return []Violation{{
		File:      filename,
		Line:      1,
		Check:     "file_length",
		Value:     lineCount,
		Threshold: fileLengthThreshold,
	}}
}
