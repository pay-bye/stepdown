// Package source loads Go packages and selects analyzable files.
package source

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/pay-bye/stepdown/internal/report"
)

var generatedMarker = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

type File struct {
	Path   string
	Fset   *token.FileSet
	Syntax *ast.File
	Info   *types.Info
	Types  *types.Package
}

type Result struct {
	Files       []File
	Diagnostics []report.Diagnostic
}

func Load(ctx context.Context, patterns []string) Result {
	fset := token.NewFileSet()
	dir, normalized := normalize(patterns)
	packages, err := packages.Load(config(ctx, fset, dir), normalized...)
	if err != nil {
		return Result{Diagnostics: []report.Diagnostic{report.Tool(report.PackageLoadFailure, err.Error())}}
	}
	return loadedFiles(fset, packages)
}

func config(ctx context.Context, fset *token.FileSet, dir string) *packages.Config {
	return &packages.Config{
		Context: ctx,
		Dir:     dir,
		Fset:    fset,
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
	}
}

func normalize(patterns []string) (string, []string) {
	if len(patterns) != 1 || !filepath.IsAbs(patterns[0]) {
		return "", patterns
	}
	info, err := os.Stat(patterns[0])
	if err != nil || !info.IsDir() {
		return "", patterns
	}
	return patterns[0], []string{"."}
}

func loadedFiles(fset *token.FileSet, packages []*packages.Package) Result {
	var result Result
	for _, item := range packages {
		result.Diagnostics = append(result.Diagnostics, packageDiagnostics(item.Errors)...)
		if len(item.Errors) > 0 {
			continue
		}
		result.Files = append(result.Files, selectedFiles(fset, item)...)
	}
	if len(packages) == 0 {
		result.Diagnostics = append(result.Diagnostics, report.Tool(report.PackageLoadFailure, "no packages matched"))
	}
	return result
}

func packageDiagnostics(errors []packages.Error) []report.Diagnostic {
	diagnostics := make([]report.Diagnostic, 0, len(errors))
	hasSpecific := hasSpecificPackageError(errors)
	for _, item := range errors {
		rule := packageErrorRule(item)
		if hasSpecific && rule == report.PackageLoadFailure {
			continue
		}
		diagnostics = append(diagnostics, report.Tool(rule, item.Msg))
	}
	return diagnostics
}

func hasSpecificPackageError(errors []packages.Error) bool {
	for _, item := range errors {
		switch item.Kind {
		case packages.ParseError, packages.TypeError:
			return true
		}
	}
	return false
}

func packageErrorRule(item packages.Error) string {
	switch item.Kind {
	case packages.ParseError:
		return report.ParseFailure
	case packages.TypeError:
		return report.TypeResolutionFailure
	default:
		return report.PackageLoadFailure
	}
}

func selectedFiles(fset *token.FileSet, item *packages.Package) []File {
	syntax := syntaxByPath(fset, item.Syntax)
	files := make([]File, 0, len(item.CompiledGoFiles))
	for _, path := range item.CompiledGoFiles {
		file := syntax[clean(path)]
		if file == nil || generated(fset, file) {
			continue
		}
		files = append(files, File{
			Path:   path,
			Fset:   fset,
			Syntax: file,
			Info:   item.TypesInfo,
			Types:  item.Types,
		})
	}
	return files
}

func syntaxByPath(fset *token.FileSet, files []*ast.File) map[string]*ast.File {
	items := make(map[string]*ast.File, len(files))
	for _, file := range files {
		items[clean(fset.Position(file.Package).Filename)] = file
	}
	return items
}

func clean(path string) string {
	return filepath.Clean(path)
}

func generated(fset *token.FileSet, file *ast.File) bool {
	packageLine := fset.Position(file.Package).Line
	for _, group := range file.Comments {
		if fset.Position(group.Pos()).Line >= packageLine {
			return false
		}
		if generatedComment(group) {
			return true
		}
	}
	return false
}

func generatedComment(group *ast.CommentGroup) bool {
	for _, comment := range group.List {
		if generatedMarker.MatchString(strings.TrimSpace(comment.Text)) {
			return true
		}
	}
	return false
}
