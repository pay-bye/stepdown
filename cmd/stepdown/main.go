// Command stepdown analyzes Go packages and reports source files whose
// declarations are not ordered top-down. See the README and ADR-0001 for
// the full grammar.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"stepdown.dev/go/internal/analyze"
	"stepdown.dev/go/internal/report"
)

const helpText = `stepdown - Go source structure analyzer for top-down declaration order

Usage:
  stepdown <package-pattern> [<package-pattern>...]
  stepdown -h
  stepdown --help
  stepdown -help

Examples:
  stepdown ./...
  go run stepdown.dev/go/cmd/stepdown@v0.1.1 ./...

Exit codes:
  0  no findings, or help printed
  1  source structure findings
  2  usage, load, parse, type-resolution, or output error

ADR: https://github.com/stepdown-dev/stepdown-go/blob/main/docs/adr/0001-stepdown-go-structure-analyzer.md
README: https://github.com/stepdown-dev/stepdown-go#readme
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 1 && isHelpAlias(args[0]) {
		fmt.Fprint(stdout, helpText)
		return 0
	}
	if len(args) == 0 || hasHelpAlias(args) {
		writeUsage(stderr)
		return 2
	}
	result := analyze.Patterns(context.Background(), args)
	if err := report.Write(stdout, result.Diagnostics); err != nil {
		fmt.Fprintf(stderr, "write diagnostics: %v\n", err)
		return 2
	}
	return result.ExitCode()
}

func writeUsage(stderr io.Writer) {
	fmt.Fprintln(stderr, "usage: stepdown <package-pattern> [<package-pattern>...]")
}

func hasHelpAlias(args []string) bool {
	for _, arg := range args {
		if isHelpAlias(arg) {
			return true
		}
	}
	return false
}

func isHelpAlias(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "-help"
}
