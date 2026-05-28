package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/pay-bye/stepdown/internal/analyze"
	"github.com/pay-bye/stepdown/internal/report"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: stepdown <package-pattern> [<package-pattern>...]")
		return 2
	}
	result := analyze.Patterns(context.Background(), args)
	if err := report.Write(stdout, result.Diagnostics); err != nil {
		fmt.Fprintf(stderr, "write diagnostics: %v\n", err)
		return 2
	}
	return result.ExitCode()
}
