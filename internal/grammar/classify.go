package grammar

import (
	"go/ast"
	"go/types"
	"strings"
	"unicode"
)

func classifyFunctions(model *Model) {
	for index := range model.Declarations {
		if model.Declarations[index].Func == nil {
			continue
		}
		model.Declarations[index] = classifyFunction(*model, model.Declarations[index])
	}
}

func classifyFunction(model Model, declaration Declaration) Declaration {
	function := declaration.Func
	receiver := receiverType(function)
	if receiver != "" {
		return classifyReceiver(model, declaration, receiver)
	}
	if owner := constructorOwner(model, function); owner != "" {
		declaration.Category = Constructor
		declaration.Owner = owner
		return declaration
	}
	if exported(function.Name.Name) {
		declaration.Category = ExportedFunction
		return declaration
	}
	declaration.Category = HelperFunction
	return declaration
}

func classifyReceiver(model Model, declaration Declaration, receiver string) Declaration {
	declaration.Owner = receiver
	if _, ok := model.Types[receiver]; !ok {
		declaration.Category = ReceiverWithoutType
		return declaration
	}
	switch {
	case isGetter(model, declaration.Func, receiver):
		declaration.Category = Getter
	case isSetter(model, declaration.Func, receiver):
		declaration.Category = Setter
	default:
		declaration.Category = ReceiverMethod
	}
	return declaration
}

func constructorOwner(model Model, function *ast.FuncDecl) string {
	if !strings.HasPrefix(function.Name.Name, "New") {
		return ""
	}
	owner := strings.TrimPrefix(function.Name.Name, "New")
	subject, ok := model.Types[owner]
	if !ok || constructorReturns(model, function, subject.Type) == false {
		return ""
	}
	return owner
}

func constructorReturns(model Model, function *ast.FuncDecl, subject types.Type) bool {
	results := resultTypes(model, function.Type.Results)
	if len(results) == 1 {
		return sameType(results[0], subject) || samePointer(results[0], subject)
	}
	if len(results) == 2 && isError(results[1]) {
		return sameType(results[0], subject) || samePointer(results[0], subject)
	}
	return false
}

func isGetter(model Model, function *ast.FuncDecl, receiver string) bool {
	subject := model.Types[receiver]
	field, ok := subject.Fields[strings.ToLower(function.Name.Name)]
	if !ok || !subject.IsStruct {
		return false
	}
	if parameterCount(function.Type.Params) != 0 || resultCount(function.Type.Results) != 1 {
		return false
	}
	if !sameType(resultTypes(model, function.Type.Results)[0], field.Type) {
		return false
	}
	return returnsField(function, field.Name)
}

func isSetter(model Model, function *ast.FuncDecl, receiver string) bool {
	if !strings.HasPrefix(function.Name.Name, "Set") {
		return false
	}
	subject := model.Types[receiver]
	field, ok := subject.Fields[strings.ToLower(strings.TrimPrefix(function.Name.Name, "Set"))]
	if !ok || !subject.IsStruct {
		return false
	}
	if parameterCount(function.Type.Params) != 1 || resultCount(function.Type.Results) != 0 {
		return false
	}
	if !sameType(parameterTypes(model, function.Type.Params)[0], field.Type) {
		return false
	}
	return assignsField(function, field.Name)
}

func resultTypes(model Model, results *ast.FieldList) []types.Type {
	if results == nil {
		return nil
	}
	items := make([]types.Type, 0, results.NumFields())
	for _, field := range results.List {
		fieldType := model.File.Info.TypeOf(field.Type)
		for range field.Names {
			items = append(items, fieldType)
		}
		if len(field.Names) == 0 {
			items = append(items, fieldType)
		}
	}
	return items
}

func parameterTypes(model Model, params *ast.FieldList) []types.Type {
	if params == nil {
		return nil
	}
	items := make([]types.Type, 0, params.NumFields())
	for _, field := range params.List {
		fieldType := model.File.Info.TypeOf(field.Type)
		for range field.Names {
			items = append(items, fieldType)
		}
		if len(field.Names) == 0 {
			items = append(items, fieldType)
		}
	}
	return items
}

func parameterCount(params *ast.FieldList) int {
	if params == nil {
		return 0
	}
	return params.NumFields()
}

func resultCount(results *ast.FieldList) int {
	if results == nil {
		return 0
	}
	return results.NumFields()
}

func sameType(left types.Type, right types.Type) bool {
	return left != nil && right != nil && types.Identical(left, right)
}

func samePointer(left types.Type, right types.Type) bool {
	pointer, ok := left.(*types.Pointer)
	return ok && sameType(pointer.Elem(), right)
}

func isError(item types.Type) bool {
	return item != nil && item.String() == "error"
}

func returnsField(function *ast.FuncDecl, field string) bool {
	if function.Body == nil || len(function.Body.List) != 1 {
		return false
	}
	statement, ok := function.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(statement.Results) != 1 {
		return false
	}
	selector, ok := statement.Results[0].(*ast.SelectorExpr)
	return ok && selector.Sel.Name == field && selectorUsesReceiver(function, selector.X)
}

func assignsField(function *ast.FuncDecl, field string) bool {
	if function.Body == nil || len(function.Body.List) != 1 {
		return false
	}
	statement, ok := function.Body.List[0].(*ast.AssignStmt)
	if !ok || statement.Tok.String() != "=" || len(statement.Lhs) != 1 || len(statement.Rhs) != 1 {
		return false
	}
	selector, ok := statement.Lhs[0].(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != field || !selectorUsesReceiver(function, selector.X) {
		return false
	}
	name, ok := statement.Rhs[0].(*ast.Ident)
	return ok && parameterNamed(function, name.Name)
}

func selectorUsesReceiver(function *ast.FuncDecl, expression ast.Expr) bool {
	name, ok := expression.(*ast.Ident)
	return ok && name.Name == receiverName(function)
}

func parameterNamed(function *ast.FuncDecl, name string) bool {
	for _, field := range function.Type.Params.List {
		for _, item := range field.Names {
			if item.Name == name {
				return true
			}
		}
	}
	return false
}

func receiverName(function *ast.FuncDecl) string {
	if function.Recv == nil || len(function.Recv.List) == 0 || len(function.Recv.List[0].Names) == 0 {
		return ""
	}
	return function.Recv.List[0].Names[0].Name
}

func receiverType(function *ast.FuncDecl) string {
	if function.Recv == nil || len(function.Recv.List) == 0 {
		return ""
	}
	return typeName(function.Recv.List[0].Type)
}

func typeName(expression ast.Expr) string {
	switch item := expression.(type) {
	case *ast.Ident:
		return item.Name
	case *ast.StarExpr:
		return typeName(item.X)
	case *ast.IndexExpr:
		return typeName(item.X)
	case *ast.IndexListExpr:
		return typeName(item.X)
	default:
		return ""
	}
}

func exported(name string) bool {
	for _, item := range name {
		return unicode.IsUpper(item)
	}
	return false
}
