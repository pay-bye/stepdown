package analyze

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"stepdown.dev/go/internal/report"
)

func TestReportsParseFailure(t *testing.T) {
	root := module(t)
	write(t, root, "broken.go", "package alpha\nfunc broken(\n")

	result := Patterns(context.Background(), []string{root})

	requireToolError(t, result, report.ParseFailure)
}

func TestReportsTypeResolutionFailure(t *testing.T) {
	root := module(t)
	write(t, root, "broken.go", "package alpha\nvar Value Missing\n")

	result := Patterns(context.Background(), []string{root})

	requireToolError(t, result, report.TypeResolutionFailure)
}

func TestReportsPackageLoadFailure(t *testing.T) {
	root := module(t)

	result := Patterns(context.Background(), []string{filepath.Join(root, "missing")})

	requireToolError(t, result, report.PackageLoadFailure)
}

func requireToolError(t *testing.T, result Result, rule string) {
	t.Helper()

	if result.ExitCode() != 2 {
		t.Fatalf("code = %d, want 2: %v", result.ExitCode(), result.Diagnostics)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatal("diagnostics are empty")
	}
	if result.Diagnostics[0].Rule != rule {
		t.Fatalf("rule = %s, want %s: %v", result.Diagnostics[0].Rule, rule, result.Diagnostics)
	}
}

func module(t *testing.T) string {
	t.Helper()

	root := t.TempDir()
	write(t, root, "go.mod", "module example.com/alpha\n\ngo 1.26.3\n")
	return root
}

func write(t *testing.T, root string, name string, content string) {
	t.Helper()

	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
