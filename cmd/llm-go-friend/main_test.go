package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func binaryPath(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "llm-go-friend")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = filepath.Join(projectRoot(t), "cmd", "llm-go-friend")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %s\n%s", err, out)
	}
	return bin
}

func projectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from this test file's directory to find go.mod.
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (go.mod)")
		}
		dir = parent
	}
}

func TestIntegration_LongFile(t *testing.T) {
	bin := binaryPath(t)
	fixture := filepath.Join(projectRoot(t), "testdata", "long.go")

	cmd := exec.Command(bin, fixture)
	stdout, err := cmd.Output()

	// Should exit 1.
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("exit code: got %d, want 1\nstderr: %s", exitErr.ExitCode(), exitErr.Stderr)
	}

	out := string(stdout)
	if !strings.Contains(out, "violations[") {
		t.Errorf("expected TOON violations header, got:\n%s", out)
	}
	if !strings.Contains(out, "file_length") {
		t.Errorf("expected file_length check in output, got:\n%s", out)
	}
	if !strings.Contains(out, "500") {
		t.Errorf("expected threshold 500 in output, got:\n%s", out)
	}
}

func TestIntegration_ShortFile(t *testing.T) {
	bin := binaryPath(t)
	fixture := filepath.Join(projectRoot(t), "testdata", "short.go")

	cmd := exec.Command(bin, fixture)
	stdout, err := cmd.Output()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v", err)
	}
	if len(stdout) != 0 {
		t.Errorf("expected no output for clean file, got:\n%s", stdout)
	}
}

func TestIntegration_NoArgs(t *testing.T) {
	bin := binaryPath(t)

	cmd := exec.Command(bin)
	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %v", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("exit code: got %d, want 2", exitErr.ExitCode())
	}
}

func TestIntegration_InvalidFile(t *testing.T) {
	bin := binaryPath(t)
	fixture := filepath.Join(projectRoot(t), "testdata", "invalid.go")

	cmd := exec.Command(bin, fixture)
	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %v", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("exit code: got %d, want 2", exitErr.ExitCode())
	}
}

func TestIntegration_NonexistentFile(t *testing.T) {
	bin := binaryPath(t)

	cmd := exec.Command(bin, "/nonexistent/path/foo.go")
	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %v", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("exit code: got %d, want 2", exitErr.ExitCode())
	}
}
