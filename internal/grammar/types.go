package grammar

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/pay-bye/stepdown/internal/source"
)

const (
	ConstBlock = iota
	VarBlock
	TypeDeclaration
	GroupedTypeDeclaration
	Constructor
	Getter
	Setter
	ReceiverMethod
	ExportedFunction
	HelperFunction
	ReceiverWithoutType
)

type Category int

type Model struct {
	File         source.File
	Types        map[string]DeclaredType
	Declarations []Declaration
}

type DeclaredType struct {
	Name     string
	Type     types.Type
	Spec     *ast.TypeSpec
	Fields   map[string]Field
	IsStruct bool
	Pos      token.Pos
}

type Field struct {
	Name string
	Type types.Type
}

type Declaration struct {
	Category Category
	Name     string
	Owner    string
	Gen      *ast.GenDecl
	Type     *ast.TypeSpec
	Func     *ast.FuncDecl
	Pos      token.Pos
}
