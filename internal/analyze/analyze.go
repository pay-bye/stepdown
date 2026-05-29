// Package analyze runs source loading and grammar checks for package patterns.
package analyze

import (
	"context"

	"github.com/pay-bye/stepdown/internal/grammar"
	"github.com/pay-bye/stepdown/internal/report"
	"github.com/pay-bye/stepdown/internal/source"
)

type Result struct {
	Diagnostics []report.Diagnostic
}

func (r Result) ExitCode() int {
	if len(r.Diagnostics) == 0 {
		return 0
	}
	if hasToolError(r.Diagnostics) {
		return 2
	}
	return 1
}

func Patterns(ctx context.Context, patterns []string) Result {
	loaded := source.Load(ctx, patterns)
	result := Result{Diagnostics: loaded.Diagnostics}
	if hasToolError(result.Diagnostics) {
		report.Sort(result.Diagnostics)
		return result
	}
	for _, file := range loaded.Files {
		model := grammar.Build(file)
		result.Diagnostics = append(result.Diagnostics, grammar.Check(model)...)
	}
	report.Sort(result.Diagnostics)
	return result
}

func hasToolError(items []report.Diagnostic) bool {
	for _, item := range items {
		if item.ToolError {
			return true
		}
	}
	return false
}
