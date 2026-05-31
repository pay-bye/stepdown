package grammar

import (
	"fmt"
	"go/ast"

	"stepdown.dev/go/internal/report"
)

func checkMethodOrder(model Model) []report.Diagnostic {
	var diagnostics []report.Diagnostic
	for _, subject := range model.Types {
		methods := receiverMethods(model, subject.Name)
		diagnostics = append(diagnostics, checkReceiverOrder(model, subject.Name, methods)...)
	}
	return diagnostics
}

func receiverMethods(model Model, owner string) []Declaration {
	var methods []Declaration
	for _, declaration := range model.Declarations {
		if declaration.Category == ReceiverMethod && declaration.Owner == owner {
			methods = append(methods, declaration)
		}
	}
	return methods
}

func checkReceiverOrder(model Model, owner string, methods []Declaration) []report.Diagnostic {
	if len(methods) == 0 {
		return nil
	}
	expected, owned := expectedMethodOrder(methods)
	diagnostics := orphanDiagnostics(model, methods, owned)
	if len(expected) != len(methods) {
		return diagnostics
	}
	for index, method := range methods {
		if method.Name != expected[index].Name {
			diagnostics = append(diagnostics, diagnostic(model, method.Pos, report.DFSPublicRoot, methodOrderDescription(owner, expected[index].Name)))
			return diagnostics
		}
	}
	return diagnostics
}

func expectedMethodOrder(methods []Declaration) ([]Declaration, map[string]bool) {
	byName := methodsByName(methods)
	owned := map[string]bool{}
	var expected []Declaration
	for _, method := range methods {
		if !exported(method.Name) {
			continue
		}
		expected = append(expected, method)
		for _, child := range methodCalls(method, byName) {
			appendOwned(&expected, owned, child, byName)
		}
	}
	return expected, owned
}

func appendOwned(expected *[]Declaration, owned map[string]bool, name string, methods map[string]Declaration) {
	if owned[name] || exported(name) {
		return
	}
	method, ok := methods[name]
	if !ok {
		return
	}
	owned[name] = true
	*expected = append(*expected, method)
	for _, child := range methodCalls(method, methods) {
		appendOwned(expected, owned, child, methods)
	}
}

func orphanDiagnostics(model Model, methods []Declaration, owned map[string]bool) []report.Diagnostic {
	var diagnostics []report.Diagnostic
	for _, method := range methods {
		if !exported(method.Name) && !owned[method.Name] {
			diagnostics = append(diagnostics, diagnostic(model, method.Pos, report.OrphanMethod, "expected unexported receiver method to be reached from an exported receiver method"))
		}
	}
	return diagnostics
}

func methodsByName(methods []Declaration) map[string]Declaration {
	items := map[string]Declaration{}
	for _, method := range methods {
		items[method.Name] = method
	}
	return items
}

func methodCalls(method Declaration, methods map[string]Declaration) []string {
	var calls []string
	ast.Inspect(method.Func.Body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || !selectorUsesReceiver(method.Func, selector.X) {
			return true
		}
		if _, ok := methods[selector.Sel.Name]; ok {
			calls = append(calls, selector.Sel.Name)
		}
		return true
	})
	return calls
}

func methodOrderDescription(owner string, expected string) string {
	return fmt.Sprintf("expected non-accessor receiver method order for %s to follow DFS from %s", owner, expected)
}
