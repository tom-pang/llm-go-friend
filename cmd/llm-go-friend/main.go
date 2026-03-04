package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"

	toon "github.com/toon-format/toon-go"

	"github.com/tom-pang/llm-go-friend/internal/check"
)

// output wraps violations for TOON tabular encoding.
type output struct {
	Violations []check.Violation `toon:"violations"`
}

func main() {
	paths := os.Args[1:]
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "usage: llm-go-friend <file.go> [file.go ...]")
		os.Exit(2)
	}

	fset := token.NewFileSet()
	var violations []check.Violation

	for _, path := range paths {
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %s\n", err)
			os.Exit(2)
		}
		violations = append(violations, check.RunAll(fset, f, path)...)
	}

	if len(violations) == 0 {
		os.Exit(0)
	}

	b, err := toon.Marshal(output{Violations: violations})
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal error: %s\n", err)
		os.Exit(2)
	}
	os.Stdout.Write(b)
	os.Stdout.Write([]byte{'\n'})
	os.Exit(1)
}
