package report

import (
	"fmt"
	"io"
	"sort"
)

const (
	SectionOrder          = "section-order"
	MultiTypeInterleave   = "multi-type-interleave"
	GroupedType           = "grouped-type-declaration"
	DFSPublicRoot         = "dfs-public-root"
	OrphanMethod          = "orphan-unexported-method"
	HelperPlacement       = "helper-placement"
	ReceiverTypeMissing   = "receiver-type-not-declared"
	ParseFailure          = "parse-failure"
	TypeResolutionFailure = "type-resolution-failure"
	PackageLoadFailure    = "package-load-failure"
)

type Diagnostic struct {
	Path        string
	Line        int
	Column      int
	Rule        string
	Description string
	ToolError   bool
}

func (d Diagnostic) String() string {
	if d.Path == "" {
		return fmt.Sprintf("%s: %s", d.Rule, d.Description)
	}
	return fmt.Sprintf("%s:%d:%d: %s: %s", d.Path, d.Line, d.Column, d.Rule, d.Description)
}

func At(path string, line int, column int, rule string, description string) Diagnostic {
	return Diagnostic{
		Path:        path,
		Line:        line,
		Column:      column,
		Rule:        rule,
		Description: description,
	}
}

func Tool(rule string, description string) Diagnostic {
	item := At("", 0, 0, rule, description)
	item.ToolError = true
	return item
}

func Sort(items []Diagnostic) {
	sort.SliceStable(items, func(left int, right int) bool {
		return less(items[left], items[right])
	})
}

func Write(writer io.Writer, items []Diagnostic) error {
	for _, item := range items {
		if _, err := fmt.Fprintln(writer, item.String()); err != nil {
			return err
		}
	}
	return nil
}

func less(left Diagnostic, right Diagnostic) bool {
	if left.Path != right.Path {
		return left.Path < right.Path
	}
	if left.Line != right.Line {
		return left.Line < right.Line
	}
	if left.Column != right.Column {
		return left.Column < right.Column
	}
	return left.Rule < right.Rule
}
