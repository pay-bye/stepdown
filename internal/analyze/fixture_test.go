package analyze

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestPositiveWitnesses(t *testing.T) {
	for _, dir := range fixtureDirs(t) {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			result := Patterns(context.Background(), []string{dir})

			if result.ExitCode() != 0 {
				t.Fatalf("code = %d, diagnostics = %v", result.ExitCode(), result.Diagnostics)
			}
		})
	}
}

func fixtureDirs(t *testing.T) []string {
	t.Helper()

	root := filepath.Join(repoRoot(t), "testdata")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(root, entry.Name()))
		}
	}
	if len(dirs) == 0 {
		t.Fatal("no fixture directories found")
	}
	return dirs
}

func repoRoot(t *testing.T) string {
	t.Helper()

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
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
