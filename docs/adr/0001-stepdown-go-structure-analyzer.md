# ADR-0001: Stepdown Go Structure Analyzer

Status: Accepted
Date: 2026-05-28
Accepted: 2026-05-28
Amended: 2026-05-31 — rehomed to `github.com/stepdown-dev/stepdown-go`, module path `stepdown.dev/go`, owner Stinnett Holdings LLC. Locator metadata only; no rule semantics changed.
Decision owner: stepdown maintainer
Initial maintainer: John Stinnett
Owner: Stinnett Holdings LLC
Ratified by: Founder (John Stinnett), with review by External Auditor, Founder Chief Architect, and Foundation Auditor

## Context

`stepdown` is a structural Go-language source analyzer. It is a commodity Go linter in the same family as `staticcheck`, `gosec`, and `govulncheck`: a pinned external Go tool that reads Go source as input and emits pass/fail diagnostics.

The analyzer enforces a positive grammar of valid Go source-file structure based on Robert C. Martin's stepdown rule (Clean Code, 2008): code should read top-down, with high-level declarations appearing before the supporting declarations they depend on. The stepdown rule is a pedagogical convention from the 2010-era enterprise programming tradition that has aged into a structural defense against agentic source drift in 2026: AI-driven code generators do not naturally produce code humans can read top-down by call dependency, and a mechanical AST analyzer is one structural enforcement gate that addresses this drift before review.

`stepdown` is intentionally narrow. It enforces source-file structure, not semantic correctness, not security, not performance, not API design. It is one rule, expressed as a positive grammar over Go's abstract syntax tree. Deviations from the grammar fail. There are no waivers.

This ADR is canonical for `stepdown`'s tool semantics. Consumers that adopt `stepdown` in their own verification pipelines record consumer-side facts (pinned version, invocation command, runtime ceiling, version-drift checks) in their own adoption records, not in this ADR.

## Decision

### Tool identity

- Tool name: `stepdown`
- Binary name: `stepdown`
- Repository: `github.com/stepdown-dev/stepdown-go`
- Module path: `stepdown.dev/go` (vanity import path; a `go-import` meta tag resolves it to the repository above)
- Command package: `stepdown.dev/go/cmd/stepdown`
- License: Apache 2.0
- Owner: Stinnett Holdings LLC

The module path is the vanity domain `stepdown.dev/go`, owned by the steward (Stinnett Holdings LLC) and resolved to the repository by a `go-import` meta tag; it carries no steward-name token at all. The repository URL `github.com/stepdown-dev/stepdown-go` names the steward organization, which is **ownership/steward locator metadata**, not tool semantics. Repository metadata, module path, license headers, and similar steward-identifying tokens are explicitly carved out as locator metadata. The tool's **rule vocabulary, diagnostics, fixture content, error messages, configuration, and rule documentation** use only Go-language concepts and contain no steward or consumer identity. The same locator-metadata discipline applies as for other commodity Go tools whose module paths name their stewards (e.g., `dominikh/go-tools` for `staticcheck`, `securego/gosec` for `gosec`).

### Maintainer ownership

The initial maintainer is John Stinnett; the tool is owned by Stinnett Holdings LLC. Maintainer succession is recorded by an amendment to this ADR (changing this section) or by a successor ADR. An orphan maintainer is a recognized failure mode (see Failure Modes); the maintainer's responsibility is to either continue the work or arrange succession before stepping away.

### Public-facing motivation

`stepdown` enforces top-down readability for human review of Go source code. As AI-driven code generation becomes common, source structure tends to drift across editing sessions: declarations accumulate in arbitrary positions because the generator's local context window does not capture global file structure. Over many iterations, files lose top-down readability — callers appear below callees, constructors appear after the methods that depend on them, helpers interleaved with the receiver methods they support.

`stepdown` addresses this drift by making the structural rule mechanical. Code that does not pass `stepdown` fails the local verification gate. The tool does not depend on human review to catch structural drift; it catches drift before review.

### Grammar

`stepdown` enforces the following positive grammar over each non-test, non-generated Go source file in the default build configuration. Each section is optional (may be empty) but must appear in this order when present:

```
package declaration
import block
constants    (zero or more const blocks; all const before any var)
package vars (zero or more var blocks)

for each type declared in the file, in source order:
    type declaration (single-spec only — see "Single-spec type declarations only" below)
    constructors for this type
    getters for this type
    setters for this type
    non-accessor receiver methods for this type
      (with DFS-from-public-roots ordering within this section)

exported package-level functions    (non-constructor, non-receiver)
unexported package-level helper functions
```

**Per-type blocks.** Each type's section block (constructor → getters → setters → non-accessor receiver methods) is contiguous and appears in source order of the type declarations. A second type's block does not begin until the first type's block ends. The constants and package-vars blocks at the top of the file are shared across types. The exported-package-level-functions and unexported-helper sections at the end of the file are shared across types.

**Single-spec type declarations only.** The grammar accepts single-spec type declarations of any form — `type Foo struct { ... }`, `type ID string`, `type Handler func(...)`, `type Items []Item`, `type Reader interface { ... }`, type aliases, and any other `TypeSpec` Go supports — one type per declaration. The per-type block grammar requires every type declaration to be followed by its own block of constructors, getters, setters, and methods (sections may be empty); the grouped `type ( ... )` form cannot satisfy this because Go syntax does not permit blocks between specs inside the group. Inputs using the grouped form emit the `grouped-type-declaration` diagnostic.

Getter and setter sections require struct fields and are therefore vacuous for non-struct type declarations (named primitives, function types, map/slice types, interfaces, aliases). The constructor and non-accessor receiver method sections may still apply to non-struct types where Go semantics permit receiver methods on those types (named primitives, function types, map/slice types — interfaces and aliases cannot have receiver methods).

**Receiver method grouping.** Receiver methods are grouped by receiver type within the file. Methods on different receiver types do not interleave.

**Receiver type normalization.** Methods with value receiver `Foo` and methods with pointer receiver `*Foo` belong to the same type's block. The receiver type is normalized to its underlying named type for grouping purposes; the pointer/value distinction is preserved in the method declaration but does not affect block assignment.

**Same-file receiver type declaration.** Every receiver method in the file must have its receiver type (after pointer/value normalization) declared in the same file. The positive grammar treats per-type blocks as local to the file containing the type declaration; receiver methods on types declared in another file have no per-type block to belong to. Receiver methods whose receiver type is not declared in the same file emit the `receiver-type-not-declared` diagnostic. This is intentional: stepdown enforces top-down readability *within a file*, and a method on a type declared elsewhere defeats that property because the reader cannot see the type declaration before its methods. Projects that want to split receiver methods across files should split the type declaration's block across files only by moving the type declaration; methods follow their type.

**Getter and setter section position.** Getter and setter sections appear in their section-order position. Within-section ordering of getters and setters is not enforced.

**DFS-from-public-roots.** Within the non-accessor receiver methods section for each type, each exported receiver method is immediately followed in source order by its unexported receiver-method callees in depth-first pre-order from the root. The traversal is local to the file and to the receiver type.

**DFS ownership for shared unexported receiver methods.** When an unexported receiver method is called by more than one exported root in the same file on the same receiver type, ownership is determined by source order of the exported roots:

- Exported roots are traversed in their source order.
- The first exported root whose DFS reaches an unexported receiver method **owns its placement** — the method is placed immediately after that root's subtree in DFS pre-order.
- Later exported roots may freely call the already-placed unexported method without requiring it to be re-placed or duplicated.
- Every unexported receiver method in the section must be reachable from at least one exported root in the same file on the same receiver type. Unreachable cases emit the `orphan-unexported-method` diagnostic.

#### DFS explicit bounds

DFS in `stepdown` is AST-local, not whole-program call graph analysis. The bounds:

- **Same file only** — no cross-file analysis
- **Same receiver type only** — receiver-method DFS does not cross receiver types
- **Direct calls visible in the AST only** — `f.method()` where `f` is the receiver value or pointer, statically resolvable in the AST
- **No interface dispatch** — interface-typed call sites are not part of the DFS graph
- **No function values** — first-class function passing, callbacks, and stored function references are not part of the DFS graph
- **No reflection-based calls**
- **No cross-package calls**
- **No DFS into package-level helper functions** — unexported package-level helper functions stay in the file-end section even if called from receiver methods

#### Classification predicates

Each top-level function and each receiver method is classified into exactly one category by mechanical AST predicate. Predicates are exact; AST-level ambiguity defaults mechanically to the most general category. Classification errors are reserved for AST-level failures only (see "Classification errors" below).

**Constructor:**
- Top-level function (no receiver)
- Function name matches the exact pattern `New<Type>` where `<Type>` is the name of a type declared in the same file
- Return type is exactly one of: `<Type>`, `*<Type>`, `(<Type>, error)`, or `(*<Type>, error)`

Both value-returning and pointer-returning constructor forms are recognized to match idiomatic Go. `func NewFoo() Foo`, `func NewFoo() *Foo`, `func NewFoo() (Foo, error)`, and `func NewFoo() (*Foo, error)` are all constructors for type `Foo`.

**Getter:**
- Receiver method
- Method name corresponds to a struct field of the receiver type. Field correspondence uses Go export-case conventions: method `Bar` corresponds to a field named `Bar` (exported) or `bar` (unexported); the field name is matched case-insensitively against the method name
- Zero non-receiver parameters
- Single return value whose type exactly matches the field's type
- Body is exactly one statement: `return r.<fieldName>` where `r` is the receiver name and `<fieldName>` is the field

**Setter:**
- Receiver method
- Method name matches the exact pattern `Set<FieldName>` where `<FieldName>` corresponds to a struct field of the receiver type. Field correspondence after dropping the `Set` prefix uses the same case rule as getters: method `SetBar` corresponds to a field named `Bar` (exported) or `bar` (unexported); the field name is matched case-insensitively
- Exactly one non-receiver parameter
- No return type
- Body is exactly one statement: `r.<fieldName> = <paramName>` where `r` is the receiver name, `<fieldName>` is the field, and `<paramName>` is the parameter name
- Parameter type exactly matches the field type

**Non-accessor receiver method:**
- Any receiver method that does not match the getter or setter predicates above

**Exported package-level function (non-constructor, non-receiver):**
- Top-level function (no receiver)
- Exported name (starts with uppercase)
- Does not match the constructor predicate

**Unexported package-level helper function:**
- Top-level function (no receiver)
- Unexported name (starts with lowercase)

#### Classification errors

Classification errors are reserved for cases where the analyzer cannot produce a deterministic classification. They are:

- **`parse-failure`** — the file is not valid Go syntax (the Go parser failed)
- **`type-resolution-failure`** — `go/packages` cannot resolve types referenced in the file (typically missing imports, unresolved external references, or partial package load)
- **`package-load-failure`** — `go/packages` cannot load the package at all

No other condition is a classification error. Body-shape mismatches (e.g., a method named `Set*` whose body is not a single field assignment) default to non-accessor receiver method classification. Name-pattern mismatches default to the appropriate fallback category. The analyzer does not invent classification errors for cases the grammar can resolve mechanically.

### Examples

All fixtures are positive witnesses: synthetic Go source files that conform to the grammar. Each fixture lives in its own directory under `testdata/`. The test harness loads each fixture via `go/packages`, runs the analyzer, and asserts zero diagnostics. There is no rejected-form corpus; the implementation contract is the positive grammar and the positive witness fixtures only.

Fixture file naming convention:

```
testdata/
  <case>/
    input.go
```

**`testdata/single_type/input.go`** — single type with constructor, getter, setter, DFS-ordered methods, and a package-level helper.

```go
package alpha

import "errors"

const MaxLorem = 100

var ErrInvalid = errors.New("invalid")

type Foo struct {
	id  int
	bar string
}

func NewFoo(id int, bar string) (Foo, error) {
	if err := requireNonBlank(bar); err != nil {
		return Foo{}, err
	}
	return Foo{id: id, bar: bar}, nil
}

func (f Foo) ID() int     { return f.id }
func (f Foo) Bar() string { return f.bar }

func (f *Foo) SetBar(bar string) { f.bar = bar }

func (f Foo) Lorem() (Foo, error) {
	if err := f.requireValid(); err != nil {
		return Foo{}, err
	}
	return f.applyLorem(), nil
}

func (f Foo) requireValid() error {
	if f.bar == "" {
		return ErrInvalid
	}
	return nil
}

func (f Foo) applyLorem() Foo {
	f.bar = f.bar + "_lorem"
	return f
}

func requireNonBlank(value string) error {
	if value == "" {
		return ErrInvalid
	}
	return nil
}
```

**`testdata/multi_type/input.go`** — two types, each with its own complete block, demonstrating per-type block grammar and receiver method grouping.

```go
package alpha

import "errors"

var errInvalid = errors.New("invalid")

type Foo struct {
	id int
}

func NewFoo(id int) Foo {
	return Foo{id: id}
}

func (f Foo) ID() int { return f.id }

func (f Foo) Process() error {
	return f.validate()
}

func (f Foo) validate() error {
	if f.id < 0 {
		return errInvalid
	}
	return nil
}

type Widget struct {
	name string
}

func NewWidget(name string) Widget {
	return Widget{name: name}
}

func (w Widget) Name() string { return w.name }

func (w Widget) Render() string {
	return w.format()
}

func (w Widget) format() string {
	return w.name + "_widget"
}
```

**`testdata/support_types/input.go`** — multiple type declarations including support types with no methods (empty blocks beyond the type declaration itself).

```go
package alpha

type Config struct {
	Host string
	Port int
}

type Settings struct {
	Debug bool
	Limit int
}

type Service struct {
	config Config
}

func NewService(config Config) Service {
	return Service{config: config}
}

func (s Service) Run() error {
	return nil
}
```

`Config` and `Settings` have type-only blocks (no constructor, getter, setter, or non-accessor receiver method sections). `Service` has a constructor and a non-accessor receiver method. All three type blocks appear in source order; each block ends before the next type declaration.

**`testdata/exported_package_function/input.go`** — type with an exported package-level function appearing after the type block and before unexported helpers.

```go
package alpha

import "errors"

type Foo struct {
	value int
}

func NewFoo(value int) Foo {
	return Foo{value: value}
}

func (f Foo) Value() int { return f.value }

func Parse(s string) (Foo, error) {
	if s == "" {
		return Foo{}, errors.New("empty input")
	}
	return NewFoo(len(s)), nil
}

func parseLen(s string) int {
	return len(s)
}
```

`Parse` is an exported package-level function (not a constructor — the name does not match `NewFoo`). It appears after Foo's complete block and before `parseLen`, the unexported package-level helper.

**`testdata/mixed_receivers/input.go`** — receiver methods using both value receiver `Foo` and pointer receiver `*Foo`, demonstrating receiver type normalization.

```go
package alpha

type Foo struct {
	value int
}

func NewFoo(value int) Foo {
	return Foo{value: value}
}

func (f Foo) Value() int { return f.value }

func (f *Foo) SetValue(v int) { f.value = v }

func (f *Foo) Write(v int) {
	f.write(v)
}

func (f *Foo) write(v int) {
	f.value = v
}
```

`Value()` uses value receiver `Foo`; `SetValue`, `Write`, and `write` use pointer receiver `*Foo`. All four methods belong to `Foo`'s block because the receiver type is normalized to its underlying named type. `Write` and `write` are DFS-ordered within the non-accessor receiver methods section.

**`testdata/pointer_constructor/input.go`** — pointer-returning constructor, the common idiomatic Go form.

```go
package alpha

type Foo struct {
	id int
}

func NewFoo(id int) *Foo {
	return &Foo{id: id}
}

func (f *Foo) ID() int { return f.id }
```

`NewFoo` returns `*Foo` rather than `Foo`. The constructor predicate accepts both value and pointer return forms, so `NewFoo` is correctly classified as Foo's constructor and lives in the constructor section of Foo's block.

**`testdata/pointer_constructor_with_error/input.go`** — pointer-returning constructor with error return.

```go
package alpha

import "errors"

var ErrInvalid = errors.New("invalid")

type Foo struct {
	id int
}

func NewFoo(id int) (*Foo, error) {
	if id < 0 {
		return nil, ErrInvalid
	}
	return &Foo{id: id}, nil
}

func (f *Foo) ID() int { return f.id }
```

`NewFoo` returns `(*Foo, error)`. This is the most common idiomatic Go constructor form for types whose construction can fail. The constructor predicate accepts it.

**`testdata/nested_dfs/input.go`** — multi-level DFS within the non-accessor receiver methods section: `Public → privateA → privateB`.

```go
package alpha

type Foo struct{}

func NewFoo() Foo {
	return Foo{}
}

func (f Foo) Public() error {
	return f.privateA()
}

func (f Foo) privateA() error {
	return f.privateB()
}

func (f Foo) privateB() error {
	return nil
}
```

`Public` is the exported root. `privateA` is its direct unexported callee, placed immediately after `Public` in DFS pre-order. `privateB` is `privateA`'s direct unexported callee, placed immediately after `privateA`. Both unexported methods are reachable from the exported root, satisfying the DFS-from-public-roots and orphan-unexported-method rules.

**`testdata/shared_dfs_callee/input.go`** — two exported roots share one unexported callee; the first root in source order owns its placement.

```go
package alpha

type Foo struct{}

func NewFoo() Foo {
	return Foo{}
}

func (f Foo) RootA() error {
	return f.shared()
}

func (f Foo) shared() error {
	return nil
}

func (f Foo) RootB() error {
	return f.shared()
}
```

`RootA` and `RootB` are both exported roots. Both call `shared`. Under the DFS ownership rule, `RootA` (the first exported root in source order) owns `shared`'s placement, so `shared` appears immediately after `RootA`'s subtree. `RootB` calls `shared` without requiring it to be re-placed or duplicated.


### File selection

`stepdown` applies to non-test, non-generated Go source files in the default build configuration. File selection uses two complementary sources:

- **`go/packages`** for package, test-file (`TestGoFiles`/`XTestGoFiles`), and build-tag categorization
- **Direct AST inspection** for generated-file detection — parse each candidate file's leading comment block and check for the standard `^// Code generated .* DO NOT EDIT\.$` marker per the Go toolchain convention. This is not handled by `go/packages` alone; the analyzer parses each file's comment group itself.

Files identified as test, generated, or non-default-tag are skipped. `stepdown` does not use filename-pattern matching for file selection.

### No waivers

`stepdown` does not provide inline waiver comments, per-file opt-out comments, or rule-specific exemption mechanisms. The structural skips for test, generated, and non-default-tag files are not waivers; they are out-of-scope file categories defined by Go's own toolchain.

If a Go source file produces a `stepdown` finding, the file must change to comply with the grammar. If `stepdown` produces false positives that require waivers to survive in real code, `stepdown`'s grammar is wrong — the grammar changes, not the file. The maintainer accepts the cost of designing a grammar correct enough not to need waivers.

### Self-policing

`stepdown`'s own source files (excluding tests, generated files, and non-default-tag files) pass `stepdown`'s own check. `stepdown` eats its own dog food. The repository's CI includes a step that runs `stepdown` against its own source. A release version of `stepdown` that does not pass its own check is not shipped.

### Implementation discipline

Two architectural commitments are spec-driving and must not be violated as the analyzer is built.

**Positive grammar only.** The analyzer walks input against the positive grammar; mismatches emit the named diagnostic. No `forbiddenPatterns []Pattern`, `denyList []string`, `antiPatterns []RuleViolation`, or `switch violationType { ... }` cascades in analyzer source, test harness, or fixtures. Rule names exist in source as stable diagnostic identifiers, not as a runtime catalog of failure kinds the implementation switches on. The implementation contract is positive grammar plus positive witnesses; there is no rejected-form corpus to drive implementation against. If an implementor wants to verify the walker emits a specific diagnostic for some malformed input, they verify it by reading the walker code, not by adding a fixture file demonstrating the bad shape.

**Sparse fixture-driven tests.** A single Go test function walks `testdata/` and spawns one subtest per fixture directory via `t.Run`. Each subtest loads `input.go` via `go/packages`, runs the analyzer, and asserts zero diagnostics. No `expected.txt`, no programmatic per-case assertions, no test helpers beyond `go/packages` and `testing.T`, no custom assertion DSL, no test factories, no test-per-method or test-per-branch unit tests against analyzer internals, no test prose. The runner is the verification surface; the fixtures are the verification data. A reader of the test directory should feel like they are reading ARM assembly.

Analyzer-internal unit tests are permitted only for tool/load error plumbing (the `parse-failure`, `type-resolution-failure`, `package-load-failure` paths). They may not become structural rejected-form examples — no fixtures or unit tests that demonstrate "here is a malformed Go declaration; analyzer must emit diagnostic X."

Coverage gaps are addressed by adding fixtures, not by adding tests against analyzer internals.

**Proof surface and residual risk.** Because there is no rejected-form corpus, the proof that the analyzer correctly rejects invalid structure does not come from automated negative tests. The absence of negative tests is an intentional architectural commitment, not an omission. The proof rests on three combined mechanisms:

1. **Reviewable walker.** The walker must stay small, readable, and directly inspectable enough that reviewers can verify mismatch paths emit diagnostics by reading the source. This is a hard constraint on the implementation.
2. **Self-policing.** `stepdown` runs against its own source as part of local verification (see Self-policing above). A regression that silently stops a rule from firing is exposed the first time stepdown's own source touches the affected pattern.
3. **Positive witnesses.** The fixture set verifies the walker does not produce false positives on conforming source.

This tradeoff is recorded, not hidden. If the analyzer ever grows beyond "readable by direct inspection" — large enough that reviewers cannot verify mismatch paths by reading the walker — the test strategy must be revisited by an explicit new ADR decision, not by quietly adding negative fixtures or unit tests against analyzer internals. Quiet expansion of the test surface to make up for an unreviewable walker would defeat the architectural commitment.

### Fixture policy

Fixtures live under `testdata/`, one directory per case. Each `input.go` is self-contained, compileable Go using generic placeholder identifiers (`Foo`, `Bar`, `Baz`, `Widget`, `Subject`, `Config`, `Service`, or analogous neutral identifiers) with no production-system identity, business-domain vocabulary, or consumer-specific references. Fixtures contain no inline comments describing what they verify; the directory name is the description. Each fixture is the minimal Go source that demonstrates its case.

`stepdown` does not implement a meta-linter on its own fixtures. Fixture discipline is enforced by maintainer review at PR time.

### Diagnostic format

Diagnostics use standard Go diagnostic format:

```
file:line:column: <rule-name>: <description>
```

Diagnostics are deterministic, machine-readable, and editor-compatible (matches the format used by `gofmt`, `go vet`, `staticcheck`, and other Go-toolchain-family tools).

Rule names are stable identifiers:

- `section-order` — section appears out of order at the file or type-block level
- `multi-type-interleave` — a later type's declaration appears before an earlier type's block completes
- `grouped-type-declaration` — grouped/multi-spec `type ( ... )` declaration
- `dfs-public-root` — unexported receiver method declared before its exported caller
- `orphan-unexported-method` — unexported receiver method has no exported root caller in the same file
- `helper-placement` — unexported package-level helper function placed outside the file-end helper section
- `receiver-type-not-declared` — receiver method declared on a type that is not declared in the same file
- `parse-failure`, `type-resolution-failure`, `package-load-failure` — analysis cannot proceed for AST-level reasons

Some rule names describe the rule positively (`section-order`, `dfs-public-root`, `helper-placement`); others describe the failed condition (`multi-type-interleave`, `orphan-unexported-method`, `grouped-type-declaration`). Both forms are acceptable for diagnostic vocabulary. The rule name is a stable identifier for human-readable output; the implementation underneath remains positive-grammar driven (see Implementation discipline above).

New rules added through future ADRs introduce new rule names; existing rule names are not renamed without an ADR amendment.

### Exit codes

- `0` — clean (no findings, no errors)
- `1` — findings present (at least one source file violated the grammar)
- `2` — tool/load error (configuration failure, package loading failure, internal parser failure, type-resolution failure)

Exit code 1 versus exit code 2 distinguishes "the source has structural problems" from "the tool itself cannot proceed." Verification gates fail closed on either non-zero exit code.

### Pinning mechanism

`stepdown` supports exactly one pinning mechanism: `go run stepdown.dev/go/cmd/stepdown@<version> ./...` where `<version>` is a published git tag (`v0.1.1`, `v0.2.0`, etc.) following Go module versioning conventions. Versions `v0.1.0` and earlier were published under the predecessor module path and do not resolve at `stepdown.dev/go`; `v0.1.1` is the first version published at the vanity path.

Other distribution forms (vendored binary, container image, package-managed install) and a stable `stepdown --version` command are **deferred to a future ADR**; `stepdown` does not provide an installed-binary version check because it does not provide an installed binary. Consumers that need binary distribution build from source.

Consumers that operate under their own foundation governance may have additional constraints on pinning, on what consumer-side artifacts may be stored, and on whether analyzer source or fixtures may be vendored into the consumer's source tree. Those constraints are recorded in the consumer's own adoption record, not in this ADR.

### Evolution path

New rules and rule families require a new ADR in this repository's ADR sequence (ADR-0002, ADR-0003, etc.). Each new rule needs explicit justification:

- What structural failure mode the rule catches
- Why review and existing rules are insufficient
- Why the rule cannot be expressed under an existing rule
- What edge cases the rule handles
- What edge cases the rule does NOT handle and why

`stepdown` is intentionally a one-opinion tool: source structure should read top-down. New rules must trace back to that opinion or wait. The maintainer rejects rules that drift toward general-purpose Go style enforcement, semantic correctness, security, performance, or API design — those are out of scope.

Bug fixes, parser compatibility updates, diagnostic message improvements, and performance work do not require a new ADR; they are maintainer discretion under semantic-versioning patch releases.

### Removal and deprecation

If `stepdown` is replaced, superseded, or retired from active maintenance, the repository remains available as a read-only archive. New versions are not released. Consumers continue to pin to the last working version or migrate to a successor.

`stepdown`'s own lifecycle is independent from any specific consumer's lifecycle. Consumers may stop consuming `stepdown` without affecting `stepdown`'s repository status. `stepdown`'s deprecation does not depend on consumer coordination.

If the rule the tool enforces is subsumed by upstream Go tooling (e.g., a future version of `gofmt` or `go vet` implements equivalent enforcement), the maintainer marks `stepdown` as superseded and points consumers to the upstream tool.

## Consequences

### Positive

- Source structure remains top-down readable across edit cycles regardless of who or what authored each edit
- Mechanical enforcement does not depend on human review for structural correctness; the rule catches drift before reviewers see it
- The tool is small, fast, and predictable — single positive grammar, no configuration, no plugin model
- Other Go projects can adopt `stepdown` independently if their maintainers value top-down source structure
- The tool's vocabulary is purely Go-language; it lifts cleanly across organizations, projects, and codebases

### Costs and risks

- Strict structural grammar can produce false positives for legitimate Go idioms the grammar did not anticipate. Recovery: file an issue, propose a grammar adjustment, ship a new release if accepted by the maintainer.
- The tool is opinionated. Projects whose maintainers prefer different source ordering will not benefit from `stepdown` and should not adopt it.
- Rule-set creep: `stepdown` could grow into a general-purpose Go style enforcer if new rules are added without discipline. The evolution path requires explicit ADR authority for new rules to prevent this drift.
- Maintainer dependency: an orphaned `stepdown` is worse than no `stepdown` for consumers depending on it. Maintainer succession is recorded explicitly when the maintainer changes.
- The grammar reflects one maintainer's view of top-down readability. Disagreement with that view is legitimate; non-adoption is the expected response.
- The single-spec type declaration requirement means projects that prefer grouped `type ( ... )` declarations cannot adopt `stepdown` without restructuring their type declarations.
- The same-file receiver type declaration requirement rejects the common Go layout where a type is declared in one file and its methods live in another. Projects that split a type's methods across multiple files cannot adopt `stepdown` without consolidating the type and its methods into the same file. This is intentional (top-down readability is a per-file property), but it is a meaningful adoption cost.

### Recovery paths

- **False positive on legitimate idiom:** file an issue with a minimal reproduction; maintainer evaluates whether the grammar needs adjustment or whether the idiom is a legitimate violation. A grammar adjustment ships in a new release.
- **Rule produces too many findings on real code:** maintainer evaluates whether the rule is correctly specified. If the rule itself is wrong, the rule is revised or removed. If the rule is correct and the codebase needs to comply, the codebase is the thing that changes.
- **Tool starts enforcing semantic correctness or domain policy:** revert the rule; route the semantic concern to a different tool. `stepdown` is structural only.
- **Maintainer becomes unavailable:** identify succession; if none available, archive the repository and notify consumers via the README and CHANGELOG.

## Alternatives Considered

### Rely on human review for source structure

Rejected. Human reviewers catch some structural drift but not all. As AI-driven generation produces more code per unit time, review bandwidth becomes the bottleneck. Mechanical enforcement at the verify-gate level scales with generation volume in a way human review cannot.

### Use an existing Go linter

Rejected. `staticcheck`, `gosec`, `govulncheck`, `errcheck`, and other commodity Go linters do not enforce stepdown-style declaration order or DFS-from-public-roots. `golangci-lint` can compose linters but does not include a stepdown-style rule. Building `stepdown` as a new standalone tool is the only path that delivers the rule.

### Implement `stepdown` rules as a `golangci-lint` plugin

Rejected. `golangci-lint` plugins are version-coupled to `golangci-lint` releases. A standalone tool has independent versioning and lifecycle. The cost of being standalone (one more tool in a consumer's verify chain) is small relative to the cost of being coupled to `golangci-lint`'s release cadence.

### Use a configuration file for rule selection or customization

Rejected. `stepdown` is one opinion: source structure should read top-down. Configuration would invite divergent flavors of `stepdown` across consumers, which fragments the rule and the tool. `stepdown` is configuration-free. Consumers who want different rules use a different tool.

### Skip DFS-from-public-roots; only enforce section order

Rejected. Section order alone is well-trodden territory and other tools cover it. DFS-from-public-roots is what makes `stepdown` distinctive: it enforces top-down call-locality within a file, which is the property that addresses the agentic-drift use case. Without DFS, the tool would be a generic section-order linter and would not earn its keep relative to existing options.

### Add inline waiver / opt-out mechanism

Rejected. Waivers are entropy machines. A grammar that needs waivers to survive is a wrong grammar; revise the grammar instead. The structural file-category skips (test, generated, non-default-tag) handle the legitimate "rule does not apply to this file kind" cases. Anything else requires the source to change.

### Meta-linter on `stepdown`'s own fixtures to prevent vocabulary drift

Rejected. A meta-linter that enumerates forbidden vocabulary would itself contain the rejected vocabulary in source — a structural anti-pattern (source-side enforcement that contains the rejected story). Fixture discipline is enforced by maintainer review at PR time, not by automated meta-enforcement.

### Support installed-binary distribution and `stepdown --version`

Deferred. `stepdown` supports only `go run module@version` pinning. Binary distribution adds packaging, signing, and version-contract concerns that warrant their own ADR if consumer need surfaces.

### Permit shape-only setter classification (any single-assignment method, regardless of name)

Rejected. Shape-only classification would treat methods named `Update`, `FooMatch`, or any other name as setters if their bodies happened to be a single field assignment. The grammar uses name-and-shape classification (method named `Set<FieldName>` AND body is the field assignment) to keep the rule's intent legible. Methods that look setter-shaped but aren't named `Set*` go in the non-accessor receiver methods section.

### Permit grouped `type (...)` declarations

Rejected. The per-type block grammar requires every type declaration to be followed by its own block of constructors/getters/setters/methods; the grouped `type ( ... )` form cannot satisfy this because Go syntax does not permit blocks between specs inside the group. The grammar therefore accepts only single-spec declarations; inputs using the grouped form emit the `grouped-type-declaration` diagnostic.

### Constructor-adapter classification — not in the grammar

Considered, then removed. The intended use case was the "row mapper" or "input adapter" pattern: a function whose return type matches a constructor's input parameter. The mechanical predicate (return type matches a constructor's single non-error parameter type) over-matches on primitives — any unrelated helper returning `int` would be classified as a constructor-adapter for `NewFoo(value int)`. Narrowing the predicate would require extra name-pattern taxonomy in source, which conflicts with the positive-grammar discipline. Functions that happen to feed constructors are classified by name visibility as exported package-level functions or unexported package-level helpers, and live at the end of the file.

## Source of Truth

This ADR is canonical for:

- `stepdown`'s name, repository location, license, and module path
- The public-facing motivation
- The grammar including section order, per-type blocks, receiver grouping, receiver type normalization, same-file receiver type declaration requirement, getter/setter section position, DFS-from-public-roots, DFS ownership for shared callees, and DFS explicit bounds
- Single-spec-only type declaration requirement
- Classification predicates for constructor, getter, setter, non-accessor receiver method, exported package-level function, and unexported package-level helper function
- Classification error definition (parse-failure, type-resolution-failure, package-load-failure)
- The positive witness fixtures (per-case directory layout under `testdata/`)
- File selection rules
- Fixture policy
- Diagnostic format, rule names, and exit codes
- Pinning mechanism (`go run module@version` only)
- Evolution path for new rules
- Removal and deprecation posture
- Self-policing requirement
- Implementation discipline: positive-grammar-only enforcement (no failure enumeration in analyzer source) and sparse fixture-driven test code (no per-method or per-branch unit tests)
- Initial maintainer identity and succession discipline

Consumer-side adoption details (consumer-specific pinned version, invocation command, runtime ceiling, version-drift test) are recorded by each consumer in their own adoption record. `stepdown` does not maintain a registry of adopters.

No other source of truth is permitted for the facts above.

## Failure Modes

| ID | Failure | Meaning | Recovery |
|---|---|---|---|
| `false-positive-on-idiom` | `stepdown` rejects a legitimate Go idiom the grammar did not anticipate | The grammar is wrong, not the source | File an issue with minimal reproduction; maintainer evaluates; grammar adjustment ships in new release. |
| `domain-leakage-in-fixture` | A test fixture contains production-system or business-domain vocabulary | Fixture discipline violated | Replace with synthetic Go-grammar example; review the fixture-review checklist; maintainer adds the check to PR review. |
| `rule-creep` | A new rule is added without ADR authority | Evolution path bypassed | Revert the rule; route through proper ADR authority before re-adding. |
| `self-policing-failure` | `stepdown`'s own source does not pass `stepdown`'s own check | The tool fails its own grammar | Fix the source or fix the grammar; do not ship a version that fails self-check; CI must block release. |
| `orphan-maintainer` | The tool's maintainer becomes unavailable without succession | Tool risk for consumers | Identify successor maintainer; if none, mark repository as archived and notify consumers via README and CHANGELOG. |
| `waiver-pressure` | Consumers or contributors request waiver mechanisms to silence findings | Pressure to weaken the rule | If a class of valid code consistently fails, the grammar is wrong — revise the grammar. If individual files fail and the grammar is correct, the files change. Do not add waivers. |
| `semantic-rule-creep` | A new proposed rule enforces semantic correctness, security, performance, or domain policy rather than source structure | Tool scope exceeded | Reject the rule. `stepdown` is structural only. Route the concern to a different tool. |
| `configuration-creep` | Pressure to add configuration flags, rule toggles, or per-project customization | Tool philosophy violated | Reject. `stepdown` is configuration-free by design. Consumers who want different rules use a different tool. |
| `classification-ambiguity-exploit` | Source uses an edge case the classification predicates do not handle cleanly, producing inconsistent classification | Predicate gap exposed | Add predicate clarification in a new ADR or amendment; update the canonical examples if needed. |
| `negative-enforcement-pattern` | Analyzer implementation uses a denied-list, forbidden-pattern catalog, or failure enumeration instead of positive-grammar walking | Implementation discipline violated; tool source has become a failure-pattern container | Refactor to positive-grammar walker; remove the negative list; let mismatches against positive grammar emit diagnostics. |
| `test-coverage-creep` | Implementation grows dense per-method or per-branch tests instead of sparse fixture-driven tests | Test discipline violated; the tool's verification surface has expanded beyond fixtures | Prune unit tests against analyzer internals; rely on the fixture set; add new fixtures for any uncovered grammar cases. |

## State and Lifecycle

This ADR is Accepted as of 2026-05-28. `stepdown` v0.1.0 may now be released against this specification.

Future ADRs in the `stepdown` sequence (ADR-0002, ADR-0003, etc.) amend, supersede, or extend this one. The ADR sequence is internal to this repository.

Bug fixes, parser compatibility updates, diagnostic message improvements, and performance work do not require ADRs; they are maintainer discretion under semantic-versioning patch releases.

This ADR is superseded only by a later ADR in this repository that explicitly cites it.

## Follow-On Artifacts

Expected after this ADR is accepted:

- `stepdown` v0.1.0 implementation release covering the grammar
- Self-policing CI step in this repository running `stepdown` against its own source
- Initial fixture catalog under `testdata/`, one directory per positive witness (`testdata/<case>/input.go`), each fixture asserting zero diagnostics
- Test harness implementing one subtest per fixture directory via `t.Run`, discovered by a single minimal runner, per the Implementation discipline section above (sparse, mechanical, no test helpers beyond `go/packages` and `testing.T`)
- Repository README pointing at this ADR as the authoritative tool specification
- Repository CONTRIBUTING.md describing the ADR-driven evolution process for new rules

Consumer-side adoption records are authored independently by each consuming project and are not tracked here.
