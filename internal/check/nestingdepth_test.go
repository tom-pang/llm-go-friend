package check

import (
	"go/parser"
	"go/token"
	"testing"
)

func TestCheckNestingDepth_Deep(t *testing.T) {
	// 5 levels deep: if > for > range > switch > select
	src := `package x

func deep() {
	if true {
		for i := 0; i < 1; i++ {
			for _, v := range []int{1} {
				switch v {
				case 1:
					select {
					default:
					}
				}
			}
		}
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkNestingDepth(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Check != "nesting_depth" {
		t.Errorf("check: got %q, want %q", v.Check, "nesting_depth")
	}
	if v.Name != "deep" {
		t.Errorf("name: got %q, want %q", v.Name, "deep")
	}
	if v.Value != 5 {
		t.Errorf("value: got %d, want 5", v.Value)
	}
	if v.Threshold != 4 {
		t.Errorf("threshold: got %d, want 4", v.Threshold)
	}
}

func TestCheckNestingDepth_Shallow(t *testing.T) {
	// 3 levels deep: should not trigger
	src := `package x

func shallow() {
	if true {
		for i := 0; i < 1; i++ {
			if false {
			}
		}
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkNestingDepth(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(violations))
	}
}

func TestCheckNestingDepth_TypeSwitch(t *testing.T) {
	// typeswitch nested enough to trigger
	src := `package x

func deepTypeSwitch() {
	if true {
		for i := 0; i < 1; i++ {
			for _, v := range []interface{}{1} {
				switch v.(type) {
				case int:
					if true {}
				}
			}
		}
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkNestingDepth(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Value != 5 {
		t.Errorf("value: got %d, want 5", violations[0].Value)
	}
}

func TestCheckNestingDepth_FuncLit(t *testing.T) {
	// func literal counts as a nesting level
	src := `package x

func withClosure() {
	if true {
		for i := 0; i < 1; i++ {
			f := func() {
				if true {
					for j := 0; j < 1; j++ {
					}
				}
			}
			_ = f
		}
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkNestingDepth(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Value != 5 {
		t.Errorf("value: got %d, want 5", violations[0].Value)
	}
}

func TestCheckNestingDepth_SelectStmt(t *testing.T) {
	// select nested enough to trigger
	src := `package x

func deepSelect() {
	if true {
		for i := 0; i < 1; i++ {
			select {
			default:
				if true {
					for j := 0; j < 1; j++ {
					}
				}
			}
		}
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkNestingDepth(fset, f, "test.go")
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Value != 5 {
		t.Errorf("value: got %d, want 5", violations[0].Value)
	}
}

func TestCheckNestingDepth_ElseDoesNotAddDepth(t *testing.T) {
	// else branch doesn't add depth, but nested if inside else does
	src := `package x

func elseTest() {
	if true {
		if true {
			if true {
				if true {
				}
			}
		}
	} else {
		_ = 1
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	violations := checkNestingDepth(fset, f, "test.go")
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations (max depth 4, threshold 4), got %d", len(violations))
	}
}
