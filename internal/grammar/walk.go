package grammar

import (
	"go/token"

	"stepdown.dev/go/internal/report"
)

const (
	beforeVars = iota
	beforeTypes
	inTypeBlock
	inExportedFunctions
	inHelpers
)

const (
	beforeConstructors = iota
	beforeGetters
	beforeSetters
	beforeMethods
)

type fileSection int

type typeSection int

type state struct {
	model       Model
	fileSection fileSection
	typeSection typeSection
	currentType string
	earlyHelper token.Pos
	diagnostics []report.Diagnostic
}

func Check(model Model) []report.Diagnostic {
	current := state{
		model:       model,
		fileSection: beforeVars,
		typeSection: beforeConstructors,
	}
	for _, declaration := range model.Declarations {
		checkDeclaration(&current, declaration)
	}
	current.diagnostics = append(current.diagnostics, checkMethodOrder(model)...)
	return current.diagnostics
}

func checkDeclaration(current *state, declaration Declaration) {
	switch declaration.Category {
	case ConstBlock:
		checkConst(current, declaration)
	case VarBlock:
		checkVar(current, declaration)
	case TypeDeclaration:
		startType(current, declaration)
	case GroupedTypeDeclaration:
		add(current, declaration.Pos, report.GroupedType, "expected one type declaration per type block")
	case Constructor, Getter, Setter, ReceiverMethod:
		checkTypeMember(current, declaration)
	case ReceiverWithoutType:
		add(current, declaration.Pos, report.ReceiverTypeMissing, "expected receiver type declaration in the same file")
	case ExportedFunction:
		checkExportedFunction(current, declaration)
	case HelperFunction:
		checkHelper(current, declaration)
	}
}

func checkConst(current *state, declaration Declaration) {
	if current.fileSection != beforeVars {
		add(current, declaration.Pos, report.SectionOrder, "expected constants before package vars and type blocks")
	}
}

func checkVar(current *state, declaration Declaration) {
	if current.fileSection > beforeTypes {
		add(current, declaration.Pos, report.SectionOrder, "expected package vars before type blocks")
		return
	}
	current.fileSection = beforeTypes
}

func startType(current *state, declaration Declaration) {
	if current.fileSection > inTypeBlock {
		if !addEarlyHelper(current, "expected helper functions after type blocks") {
			add(current, declaration.Pos, report.SectionOrder, "expected type block before package functions")
		}
	}
	current.fileSection = inTypeBlock
	current.typeSection = beforeConstructors
	current.currentType = declaration.Name
}

func checkTypeMember(current *state, declaration Declaration) {
	if current.fileSection != inTypeBlock {
		addEarlyHelper(current, "expected helper functions after receiver methods")
		if current.fileSection != inTypeBlock {
			add(current, declaration.Pos, report.SectionOrder, "expected receiver or constructor inside a type block")
			return
		}
	}
	if declaration.Owner != current.currentType {
		add(current, declaration.Pos, report.MultiTypeInterleave, "expected declarations for current type block before another receiver type")
		return
	}
	checkTypeSection(current, declaration)
}

func checkTypeSection(current *state, declaration Declaration) {
	switch declaration.Category {
	case Constructor:
		requireTypeSection(current, declaration, beforeConstructors, "expected constructors before accessors and methods")
	case Getter:
		requireTypeSection(current, declaration, beforeGetters, "expected getters before setters and methods")
	case Setter:
		requireTypeSection(current, declaration, beforeSetters, "expected setters before methods")
	case ReceiverMethod:
		current.typeSection = beforeMethods
	}
}

func requireTypeSection(current *state, declaration Declaration, latest typeSection, description string) {
	if current.typeSection > latest {
		add(current, declaration.Pos, report.SectionOrder, description)
		return
	}
	current.typeSection = latest
}

func checkExportedFunction(current *state, declaration Declaration) {
	if current.fileSection == inHelpers {
		current.earlyHelper = token.NoPos
		add(current, declaration.Pos, report.HelperPlacement, "expected helper functions after exported package-level functions")
		return
	}
	current.fileSection = inExportedFunctions
}

func checkHelper(current *state, declaration Declaration) {
	if current.fileSection == inTypeBlock {
		current.earlyHelper = declaration.Pos
	}
	current.fileSection = inHelpers
}

func addEarlyHelper(current *state, description string) bool {
	if !current.earlyHelper.IsValid() {
		return false
	}
	add(current, current.earlyHelper, report.SectionOrder, description)
	current.earlyHelper = token.NoPos
	current.fileSection = inTypeBlock
	return true
}

func add(current *state, pos token.Pos, rule string, description string) {
	current.diagnostics = append(current.diagnostics, diagnostic(current.model, pos, rule, description))
}

func diagnostic(model Model, pos token.Pos, rule string, description string) report.Diagnostic {
	position := model.File.Fset.Position(pos)
	if position.Filename == "" {
		return report.Tool(rule, description)
	}
	return report.At(position.Filename, position.Line, position.Column, rule, description)
}
