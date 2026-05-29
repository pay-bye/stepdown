package main

import (
	"bytes"
	"testing"
)

const expectedHelp = `stepdown - Go source structure analyzer for top-down declaration order

Usage:
  stepdown <package-pattern> [<package-pattern>...]
  stepdown -h
  stepdown --help
  stepdown -help

Examples:
  stepdown ./...
  go run github.com/pay-bye/stepdown/cmd/stepdown@v0.1.0 ./...

Exit codes:
  0  no findings, or help printed
  1  source structure findings
  2  usage, load, parse, type-resolution, or output error

ADR: https://github.com/pay-bye/stepdown/blob/main/docs/adr/0001-stepdown-go-structure-analyzer.md
README: https://github.com/pay-bye/stepdown#readme
`

func TestRunPrintsHelpForAliases(t *testing.T) {
	for _, alias := range []string{"-h", "--help", "-help"} {
		stdout, stderr, code := execute(alias)

		if code != 0 {
			t.Fatalf("%s code = %d, want 0", alias, code)
		}
		if stdout != expectedHelp {
			t.Fatalf("%s stdout = %q, want help", alias, stdout)
		}
		if stderr != "" {
			t.Fatalf("%s stderr = %q, want empty", alias, stderr)
		}
	}
}

func TestRunRequiresPackagePatterns(t *testing.T) {
	stdout, stderr, code := execute()

	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if stderr == "" {
		t.Fatal("stderr is empty")
	}
}

func TestRunRejectsMixedHelpAndPackages(t *testing.T) {
	stdout, stderr, code := execute("-h", "./...")

	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if stderr == "" {
		t.Fatal("stderr is empty")
	}
}

func execute(args ...string) (string, string, int) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run(args, &stdout, &stderr)

	return stdout.String(), stderr.String(), code
}
