// Package grammar models Go declarations and checks top-down ordering rules.
package grammar

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"stepdown.dev/go/internal/source"
)

func Build(file source.File) Model {
	model := Model{
		File:  file,
		Types: map[string]DeclaredType{},
	}
	model.Declarations = rawDeclarations(file)
	model.Types = declaredTypes(file, model.Declarations)
	classifyFunctions(&model)
	return model
}

func rawDeclarations(file source.File) []Declaration {
	var declarations []Declaration
	for _, decl := range file.Syntax.Decls {
		declarations = append(declarations, rawDeclaration(decl))
	}
	return declarations
}

func rawDeclaration(decl ast.Decl) Declaration {
	switch item := decl.(type) {
	case *ast.GenDecl:
		return rawGenDeclaration(item)
	case *ast.FuncDecl:
		return Declaration{Name: item.Name.Name, Func: item, Pos: item.Pos()}
	default:
		return Declaration{Pos: decl.Pos()}
	}
}

func rawGenDeclaration(decl *ast.GenDecl) Declaration {
	switch decl.Tok {
	case token.CONST:
		return Declaration{Category: ConstBlock, Gen: decl, Pos: decl.Pos()}
	case token.VAR:
		return Declaration{Category: VarBlock, Gen: decl, Pos: decl.Pos()}
	case token.TYPE:
		return rawTypeDeclaration(decl)
	default:
		return Declaration{Pos: decl.Pos()}
	}
}

func rawTypeDeclaration(decl *ast.GenDecl) Declaration {
	if decl.Lparen.IsValid() || len(decl.Specs) != 1 {
		return Declaration{Category: GroupedTypeDeclaration, Gen: decl, Pos: decl.Pos()}
	}
	spec, ok := decl.Specs[0].(*ast.TypeSpec)
	if !ok {
		return Declaration{Category: GroupedTypeDeclaration, Gen: decl, Pos: decl.Pos()}
	}
	return Declaration{
		Category: TypeDeclaration,
		Name:     spec.Name.Name,
		Gen:      decl,
		Type:     spec,
		Pos:      decl.Pos(),
	}
}

func declaredTypes(file source.File, declarations []Declaration) map[string]DeclaredType {
	items := map[string]DeclaredType{}
	for _, declaration := range declarations {
		if declaration.Category != TypeDeclaration {
			continue
		}
		items[declaration.Name] = declaredType(file, declaration.Type)
	}
	return items
}

func declaredType(file source.File, spec *ast.TypeSpec) DeclaredType {
	item := DeclaredType{
		Name:   spec.Name.Name,
		Type:   definedType(file, spec),
		Spec:   spec,
		Fields: map[string]Field{},
		Pos:    spec.Pos(),
	}
	if body, ok := spec.Type.(*ast.StructType); ok {
		item.IsStruct = true
		item.Fields = fields(file, body)
	}
	return item
}

func definedType(file source.File, spec *ast.TypeSpec) types.Type {
	object := file.Info.Defs[spec.Name]
	if object == nil {
		return nil
	}
	return object.Type()
}

func fields(file source.File, body *ast.StructType) map[string]Field {
	items := map[string]Field{}
	for _, field := range body.Fields.List {
		fieldType := file.Info.TypeOf(field.Type)
		for _, name := range field.Names {
			items[strings.ToLower(name.Name)] = Field{Name: name.Name, Type: fieldType}
		}
	}
	return items
}
