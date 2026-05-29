# Go Structure Analyzer Implementation Spec

Status: Submitted
Date: 2026-05-28
Authority: ADR-0001
Target release: v0.1.0

## Problem

ADR-0001 defines `stepdown` as a structural Go-language analyzer, but implementation needs an executable specification that maps the accepted ADR into concrete packages, interfaces, slices, tests, review gates, and failure handling.

The implementation risk is not that `stepdown` omits a feature. The implementation risk is that the analyzer is built around named rejected shapes, negative fixtures, or per-case diagnostic tables while claiming to enforce a positive grammar. This spec makes the positive grammar walker load-bearing and rejects any implementation that switches on named bad shapes.

This spec is an implementation translation. ADR-0001 remains canonical for tool semantics. If this spec conflicts with ADR-0001, ADR-0001 wins and the implementation work returns for spec correction.

## Success Criteria

- `stepdown` v0.1.0 builds as module `github.com/pay-bye/stepdown` with command package `github.com/pay-bye/stepdown/cmd/stepdown` and Apache 2.0 license posture preserved.
- The analyzer walks Go source against the accepted positive grammar from ADR-0001 and emits diagnostics when declarations do not match the next valid grammar position.
- Analyzer source contains no rejected-shape catalog, denied-pattern list, forbidden-shape registry, bad-shape enum, or equivalent source structure.
- Analyzer source does not switch on named bad shapes or violation types such as `ConstructorBeforeType`, `GetterAfterMethod`, `MultiTypeInterleaved`, or equivalent names.
- Rule-name constants exist only as stable output vocabulary. Runtime diagnostic generation derives from positive grammar state and observed declaration facts.
- Committed fixtures are positive witnesses only: `testdata/<case>/input.go` directories discovered mechanically and asserted to produce zero diagnostics.
- The repository contains no `testdata/violations/`, no rejected-form `expected.txt`, no committed malformed structural fixture corpus, and no per-case expected diagnostic table.
- `go test ./...` passes from the repository root.
- A command-level run against positive witness fixtures passes.
- A self-policing run of `stepdown` against its own non-test, non-generated, default-build Go source passes before v0.1.0 release.
- Reviewer can reject a no-op analyzer by inspecting the positive grammar walker and confirming mismatch paths emit diagnostics from grammar state.
- Foundation Auditor can verify that source, fixtures, tests, docs, and config do not become rejected-form containers.

## Assumptions

- The accepted ADR's Go grammar is complete enough for v0.1.0 implementation.
- The first implementation uses Go packages available through the standard Go toolchain plus `golang.org/x/tools/go/packages`.
- The command accepts explicit package patterns. `go run github.com/pay-bye/stepdown/cmd/stepdown@v0.1.0 ./...` is the release invocation model.
- A zero-argument command is a tool usage error and exits `2`; ADR-0001 documents explicit package-pattern usage and does not authorize an implicit default.
- The self-policing gate can run only after the command builds and can analyze its own source without relying on an already released tag.

## Constraints And Non-Goals

In scope:

- Implementation for ADR-0001 grammar, predicates, diagnostics, file selection, fixtures, CLI behavior, exit codes, and self-policing.
- README and CONTRIBUTING updates needed for v0.1.0 implementation usage and contributor discipline.
- Minimal package structure that keeps the walker reviewable by direct inspection.

Out of scope:

- Consumer-side Agent OS adoption records.
- Pay-bye verification-gate integration.
- Vendored binary distribution.
- Container images.
- Package-manager installation.
- Stable `stepdown --version`.
- Waiver mechanisms.
- Rule configuration.
- Semantic correctness, security, performance, API-design, or domain-policy linting.
- Any Pay-bye-specific diagnostic, fixture, rule, package, or runtime vocabulary inside `stepdown`.
- Product/runtime customer-data surfaces.

## ADR Requirement Map

This map translates each load-bearing ADR-0001 requirement into implementation requirements. No row changes tool semantics.

| ADR requirement | Implementation requirement | Verification or review gate |
|---|---|---|
| Tool name and binary are `stepdown`. | Command builds under `cmd/stepdown`; usage examples call `stepdown`. | `go test ./...`; command smoke from repository root. |
| Repository and module path are `github.com/pay-bye/stepdown`; command package is `github.com/pay-bye/stepdown/cmd/stepdown`; license is Apache 2.0. | Preserve `go.mod`, `LICENSE`, README, and command package path. `pay-bye` remains locator metadata only. | File review; no tool rule, fixture, diagnostic, or config uses steward identity as semantics. |
| Ownership/steward locator metadata is not tool semantics. | Module path and license headers may name steward. Rule vocabulary, diagnostics, fixtures, error messages, configuration, and docs use Go-language concepts only. | Vocabulary sweep in Reviewer and Foundation Auditor gates. |
| Public motivation is top-down readability for human review of Go source. | README and command help describe structural Go source order, not consumer-specific policy. | README/help review. |
| Strict non-goals exclude semantic correctness, security, performance, API design, waivers, and configuration. | No rules, flags, config files, or docs introduce those scopes. | Source and docs review. |
| Positive grammar ordering: package, imports, consts, vars, per-type blocks, exported package functions, unexported helpers. | Walker models the ordered grammar and compares observed top-level declarations to the next valid grammar position. | Walker review; positive fixture run; no-op analyzer rejection gate. |
| Per-type blocks are contiguous and source ordered. | Model each single `TypeSpec` as a block owner; constructors/getters/setters/non-accessor receiver methods attach only to that type block before the next type begins. | Positive fixtures `single_type`, `multi_type`, `support_types`; walker review for mismatch diagnostics. |
| Single `TypeSpec` declarations of any Go-supported type form are accepted. | Accept one type spec per type declaration across structs, named primitives, function types, maps, slices, interfaces, aliases, and other Go-supported type expressions. | Positive fixtures for representative type forms; code review of type declaration extraction. |
| Grouped `type (...)` declarations are rejected. | Any `type` declaration with a parenthesized group or more than one spec emits `grouped-type-declaration`. | Walker review; no rejected-form fixture corpus. |
| Same-file receiver type declaration is required. | Receiver methods normalize receiver type and require a matching type declaration in the same file; otherwise emit `receiver-type-not-declared`. | Walker review; declaration model review. |
| Receiver type normalization merges `Foo` and `*Foo`. | Receiver normalization returns the underlying named type for value and pointer receivers. | Positive fixture `mixed_receivers`; code review. |
| Getter section predicate. | Getter classification applies only to struct receiver types, matches field name case-insensitively, requires zero params, one return exactly matching field type, and a single `return r.field` body. | Classification review; positive fixtures with getters. |
| Setter section predicate. | Setter classification applies only to struct receiver types, requires name `Set<FieldName>`, one param matching field type, no returns, and one assignment `r.field = param`. | Classification review; positive fixtures with setters. |
| Getter and setter sections are vacuous for non-struct types. | Non-struct type blocks skip accessor expectations; methods on named non-struct types classify as non-accessor when Go permits methods. | Positive fixtures for non-struct declarations. |
| Constructor predicate accepts value and pointer returns, optionally with `error`. | Top-level `New<Type>` functions returning `<Type>`, `*<Type>`, `(<Type>, error)`, or `(*<Type>, error)` classify as constructors for same-file type declarations. | Positive fixtures `pointer_constructor` and `pointer_constructor_with_error`; code review. |
| Fallback classification covers non-accessor receiver methods, exported package functions, and unexported helpers. | Any receiver method not matching getter/setter is non-accessor; exported non-constructor top-level function is exported package-level function; unexported top-level function is helper. | Classification review; positive fixtures. |
| Classification errors are limited to parse, type-resolution, and package-load failures. | Only parser, type-resolution, and package loading failures produce exit code `2` diagnostics. Body-shape and name-pattern mismatches fall back to mechanical categories. | Error-plumbing tests; no structural negative fixtures. |
| DFS-from-public-roots orders unexported receiver callees under exported roots. | For each receiver type's non-accessor section, build a same-file, same-receiver graph of direct receiver method calls and compare observed order to DFS pre-order from exported roots. | Positive fixtures `nested_dfs` and `shared_dfs_callee`; walker review. |
| Shared unexported receiver callee ownership belongs to first exported root in source order. | DFS marks first owner by exported root source order and does not require duplicate placement under later roots. | Positive fixture `shared_dfs_callee`; graph review. |
| DFS explicit bounds. | Graph builder is same file, same receiver type, direct AST-visible receiver calls only; no interface dispatch, function values, reflection, cross-file, cross-package, or package-level helper DFS. | Graph review. |
| File selection skips tests, generated files, and non-default build files. | Use `go/packages` default build selection and test file categorization; inspect leading comments for the Go generated-code marker. | File-selection tests with positive generated/test skip fixtures; command run. |
| Generated-file detection uses Go generated-code marker, not filename guessing. | Detect `^// Code generated .* DO NOT EDIT\\.$` in the leading comment block. Do not branch on filenames like `zz_generated.go`. | File-selection review. |
| No waivers. | No inline ignore comments, per-file opt-outs, config flags, or rule toggles. | Source, docs, CLI, and help review. |
| Self-policing. | Repository verification includes a local built-command run against non-test, non-generated, default-build source. | Self-policing command in verification contract. |
| Fixture policy. | Positive fixtures only under `testdata/<case>/input.go`; fixtures are compileable, synthetic, and domain-neutral. | Fixture sweep; mechanical fixture harness. |
| Sparse fixture-driven tests. | One harness discovers fixture directories and asserts zero diagnostics. No per-case expected diagnostic table. Internal tests only cover tool/load error plumbing. | Test review. |
| Diagnostic format and rule names. | Emit `file:line:column: <rule-name>: <description>` using stable ADR rule-name constants. | Diagnostic formatting test; command run. |
| Exit codes. | `0` clean, `1` findings, `2` tool/load error. | CLI tests. |
| Pinning mechanism. | README documents `go run github.com/pay-bye/stepdown/cmd/stepdown@v0.1.0 ./...`; no installed-binary version command in v0.1.0. | README review. |
| Evolution, removal, deprecation, and maintainer risk. | README/CONTRIBUTING preserve ADR-driven new-rule process, no waivers, no configuration creep, and maintainer succession posture. | Docs review. |

## Positive Grammar Walker Requirement

Implementor builds a positive grammar walker.

The analyzer source describes the valid grammar and the expected next declaration or section state. It observes source declarations, classifies them using mechanical AST and type predicates, and emits diagnostics when an observed declaration does not match the next valid grammar position.

The analyzer must not contain a rejected-shape catalog, denied-pattern list, forbidden-shape registry, bad-shape enum, or equivalent source structure. It must not switch on named bad shapes or violation types such as `ConstructorBeforeType`, `GetterAfterMethod`, `MultiTypeInterleaved`, or equivalent names.

Diagnostic rule-name constants are allowed as stable output vocabulary. Runtime diagnostic generation must derive from positive grammar state and observed declaration facts, not from a failure-kind dispatch table.

Reviewer rejects an implementation that passes tests by acting as a no-op analyzer. Reviewer verifies the walker source directly: each grammar state has accepted next declarations, mismatch paths emit diagnostics, and emitted rule names come from grammar-state failure, not from a bad-shape switch.

Foundation Auditor rejects an implementation that turns source, fixtures, tests, docs, or config into a rejected-form container. Rejected examples live only in ADR/review prose when needed to explain boundaries; they do not become source artifacts.

## No Rejected-Form Source Artifacts

Committed fixtures are positive witnesses only.

The repository must not contain:

- `testdata/violations/`;
- `expected.txt` files for rejected forms;
- malformed structural Go source as a standing fixture corpus;
- per-case expected diagnostic tables in source code;
- analyzer source structures that enumerate bad shapes;
- docs that teach contributors to add rejected examples as proof.

The fixture harness discovers `testdata/<case>/input.go` directories mechanically and asserts zero diagnostics. Fixture directory names describe valid grammar coverage, not invalid examples. Fixture file content uses generic Go identifiers such as `Foo`, `Bar`, `Baz`, `Widget`, `Subject`, `Config`, and `Service`; it contains no production-system identity, business-domain vocabulary, consumer-specific references, Pay-bye vocabulary, Agent OS vocabulary, workflow vocabulary, payroll vocabulary, foundation vocabulary, watchlist vocabulary, or predecessor-story vocabulary.

Internal tests may cover tool/load error plumbing for `parse-failure`, `type-resolution-failure`, and `package-load-failure` when those conditions cannot be exercised through positive fixtures. Those tests remain focused on error plumbing and must not become structural rejected-form examples.

Reviewer-created scratch inputs used for manual inspection remain outside committed source and outside this implementation spec.

## Boundary And Ownership Model

| Boundary | Owner | Responsibility | Dependency direction |
|---|---|---|---|
| Command | `cmd/stepdown` | Parse package patterns, call analyzer, print diagnostics, return exit codes. | Command depends on internal analyzer packages. |
| Package loading and file selection | internal implementation package | Load packages with `go/packages`; select default-build, non-test, non-generated Go files. | Analyzer depends on Go toolchain metadata and parsed AST. |
| Declaration model | internal implementation package | Extract source-ordered top-level declarations, type declarations, receiver methods, fields, constructors, and package functions. | Walker depends on declaration facts, not raw scattered AST traversal. |
| Classification predicates | internal implementation package | Classify constructors, getters, setters, non-accessor receiver methods, exported package functions, helpers, and classification errors. | Walker consumes positive declaration categories. |
| Positive grammar walker | internal implementation package | Enforce ordered grammar and emit diagnostics on mismatch. | Walker depends on declaration model and predicates. |
| DFS graph | internal implementation package | Build same-file same-receiver direct-call graph and expected DFS order for non-accessor receiver methods. | Walker consumes graph order for receiver method section. |
| Diagnostics | internal implementation package | Format deterministic diagnostics and exit code classification. | Command consumes diagnostics. |
| Fixtures | `testdata/<case>/input.go` | Positive witnesses only. | Tests load fixtures through analyzer public path. |
| Documentation | README and CONTRIBUTING | Usage, ADR authority, evolution process, fixture discipline. | Docs reference ADR-0001 as canonical semantics. |

Package names are implementation details, but every package must have one responsibility expressible in one sentence. Generic drawers such as `utils`, `common`, `helpers`, `manager`, and `service` are rejected.

## Contracts And Invariants

### Command Contract

Input:

- One or more Go package patterns.

Output:

- Diagnostics on standard output or standard error in standard Go diagnostic format.
- Exit code `0` when no findings and no tool/load errors occur.
- Exit code `1` when one or more grammar findings occur.
- Exit code `2` when package load, parse, type-resolution, command usage, or internal tool errors prevent analysis.

Error semantics:

- Findings and tool/load errors are distinct.
- Package load failure, parse failure, and type-resolution failure are classification errors and produce exit code `2`.
- Source structure violations produce exit code `1`.

### Analyzer Contract

Input:

- Package patterns resolved by `go/packages`.
- Candidate files selected from default-build, non-test, non-generated Go files.

Output:

- Ordered diagnostics with file, line, column, rule name, and description.

Invariants:

- Analyzer never mutates input source.
- Analyzer never reads consumer-specific config.
- Analyzer has no waiver path.
- Analyzer accepts every valid positive fixture.
- Analyzer diagnostic order is deterministic by file path, source position, and grammar encounter order.

### Grammar State Contract

The walker exposes the ordering through data and control flow:

- current file-level section;
- current type block owner;
- current subsection inside a type block;
- observed declaration category;
- expected next categories.

Mismatch handling emits diagnostics from the current positive state. The walker does not dispatch through rejected-shape names.

## Implementation Architecture

### Package Loading And File Selection

Use `go/packages` for:

- package load;
- default build configuration;
- package file list;
- test file categorization;
- type information.

Use direct AST leading-comment inspection for generated files:

- parse leading comments of each candidate file;
- match the standard Go generated-code marker `^// Code generated .* DO NOT EDIT\\.$`;
- skip generated files.

Do not use filename guessing for generated-file detection.

### Declaration Model

For each selected file, build a source-ordered model:

- const declarations;
- var declarations;
- type declarations;
- function declarations;
- receiver methods;
- field metadata for struct types;
- source positions;
- receiver normalization facts.

Grouped type declarations emit `grouped-type-declaration`; they do not produce per-type blocks.

### Classification Predicates

Predicates follow ADR-0001 exactly.

Constructor predicate:

- top-level function;
- name exactly `New<Type>`;
- `<Type>` declared in same file;
- return is exactly `<Type>`, `*<Type>`, `(<Type>, error)`, or `(*<Type>, error)`.

Getter predicate:

- receiver method on struct receiver type;
- method name corresponds to a struct field by ADR case rule;
- zero non-receiver parameters;
- single return exactly matching field type;
- body exactly `return r.field`.

Setter predicate:

- receiver method on struct receiver type;
- method name exactly `Set<FieldName>`;
- one non-receiver parameter;
- no returns;
- parameter type exactly matches field type;
- body exactly `r.field = param`.

Fallback predicates:

- receiver method not matching getter/setter becomes non-accessor receiver method;
- exported top-level non-constructor function becomes exported package-level function;
- unexported top-level function becomes unexported package-level helper.

No other classification errors are invented.

### Positive Section-Order Walker

For each file:

1. Start in the file pre-type section.
2. Accept const declarations before package var declarations.
3. Accept package var declarations before any type declaration.
4. For each type declaration in source order, enter a type block.
5. Inside a type block, accept constructors, then getters, then setters, then non-accessor receiver methods.
6. Close the current type block only when the next source declaration starts the next type block or the file-level package function sections.
7. After all type blocks, accept exported package-level functions, then unexported package-level helper functions.
8. Emit diagnostics when an observed declaration does not match the next valid grammar position.

The walker may use small data structures for valid next states. It must not use data structures whose entries are named rejected forms.

### DFS Graph

For each receiver type and file:

1. Collect non-accessor receiver methods in source order.
2. Identify exported roots in source order.
3. For each root, walk direct calls to unexported methods on the same receiver type, same file, and receiver value or pointer.
4. DFS child order follows call-site source order.
5. The first exported root that reaches an unexported method owns its placement.
6. Already-owned unexported methods are not duplicated under later roots.
7. Unexported methods that no exported root reaches emit `orphan-unexported-method`.
8. Observed non-accessor receiver method order must match the expected exported-root plus owned-DFS order.

The graph does not include interface dispatch, function values, reflection, cross-file calls, cross-package calls, or package-level helper calls.

### Diagnostics

Diagnostic rule names:

- `section-order`
- `multi-type-interleave`
- `grouped-type-declaration`
- `dfs-public-root`
- `orphan-unexported-method`
- `helper-placement`
- `receiver-type-not-declared`
- `parse-failure`
- `type-resolution-failure`
- `package-load-failure`

The implementation may store those names as constants. It must not store a failure taxonomy that drives analyzer behavior.

Descriptions state the positive expectation, for example:

- expected constructor/getter/setter/method section for current type;
- expected next type block after current block completes;
- expected helper after exported package-level functions;
- expected receiver type declaration in same file.

## Implementation Slice Plan

Each unit is sized for a focused Implementor session. The sequence is dependency ordered.

| Unit | Size | Surfaces | Completion evidence |
|---|---|---|---|
| `U1` module and command scaffold | S | `cmd/stepdown`, minimal internal run path, README usage touch-up | Command accepts package patterns, returns `2` for zero args, and builds under `go test ./...`. |
| `U2` package loading and file selection | M | package loading/file selection package plus fixture/review evidence | Default-build non-test files are selected; test files and generated files are skipped by Go toolchain category and generated marker. |
| `U3` declaration model | M | declaration extraction package plus fixture/review evidence | Source-ordered const, var, type, function, receiver, field, and position facts are available to the walker. |
| `U4` classification predicates | M | predicate package and positive fixtures | Constructors, getters, setters, non-accessor receiver methods, exported package functions, helpers, and classification errors follow ADR predicates. |
| `U5` positive section-order walker | M | walker package and diagnostics | File/type section order is enforced from positive grammar state; no rejected-shape catalog exists. |
| `U6` DFS-from-public-roots | M | DFS graph and walker integration | Receiver method order follows same-file, same-receiver DFS rules with first-root ownership and orphan detection. |
| `U7` deterministic diagnostics and exit codes | S | diagnostic formatting and command integration | Diagnostic order and format are stable; exit codes `0`, `1`, and `2` are implemented. |
| `U8` positive witness fixture harness | S | `testdata/<case>/input.go`, one fixture harness | Positive fixtures are discovered mechanically and assert zero diagnostics. |
| `U9` self-policing gate | S | local verification script or documented command | Built command runs against repository non-test, non-generated, default-build source. |
| `U10` README and CONTRIBUTING minimum updates | XS | README, CONTRIBUTING | Docs point to ADR-0001, explain usage, no waivers, fixture policy, and ADR-driven evolution. |

No unit other than `U8` and focused tool/load error plumbing authorizes analyzer-internal tests; `U2`, `U3`, `U4`, `U5`, and `U6` are verified by positive fixtures, command runs, and direct source inspection.

Do not merge slices by adding a god package that owns loading, classification, walking, diagnostics, CLI, and fixtures in one unit. The implementation may keep packages small, but responsibilities remain separate and reviewable.

## Acceptance Criteria Per Unit

| Unit | Acceptance criterion |
|---|---|
| `U1` | `go test ./...` passes; command usage accepts package patterns; zero args exits `2` with usage text. |
| `U2` | File selection uses `go/packages` for default-build and test categorization; generated detection uses the Go marker; filename guessing is absent. |
| `U3` | Declaration model preserves source order and positions for every top-level declaration; receiver type normalization returns the same named type for `Foo` and `*Foo`. |
| `U4` | Predicate source matches ADR-0001 constructor/getter/setter/fallback definitions; non-struct accessor sections are vacuous. |
| `U5` | Walker source shows positive grammar states and next valid declaration categories; no denied-pattern list, bad-shape enum, or violation-type switch exists. |
| `U6` | DFS is same-file, same receiver type, direct AST-visible receiver calls only; first exported root owns shared unexported callees. |
| `U7` | Diagnostics match `file:line:column: <rule-name>: <description>` and exit codes distinguish clean, findings, and tool/load errors. |
| `U8` | Fixture harness discovers `testdata/<case>/input.go`; all committed fixtures are positive and domain-neutral; no `testdata/violations/` or `expected.txt` exists. |
| `U9` | Self-policing run analyzes the tool's own source and passes. |
| `U10` | README and CONTRIBUTING reflect ADR-0001 semantics, usage, no waivers, fixture policy, and ADR-driven evolution. |

## Verification Contract

Verification proves:

- the command builds and tests pass;
- positive witness fixtures produce zero diagnostics;
- the command can run against repository source;
- the analyzer source is a positive grammar walker by review inspection;
- no rejected-form source corpus exists;
- a no-op analyzer is not accepted because review confirms mismatch paths emit diagnostics from grammar state.

Verification does not prove:

- positive witness fixtures alone prove rejection behavior;
- `stepdown` enforces semantic correctness, security, performance, API design, or domain policy;
- `stepdown` is adopted by any consumer verification gate;
- installed binary distribution exists;
- a stable `stepdown --version` contract exists.

Toolchain contract:

- `stepdown/go.mod` declares `go 1.26.3`.
- Implementation verification runs under Go 1.26.3.
- Local verification uses the environment-resolved `go` command only when `go version` reports `go1.26.3` and `go env GOTOOLCHAIN` reports `local`.
- If either local toolchain check fails, verification stops until the canonical local toolchain path is used or repaired. Verification does not proceed through toolchain auto-download or a different Go version.

Required verification surfaces:

```bash
cd /path/to/stepdown
go version
go env GOTOOLCHAIN
rg -n '^go 1\.26\.3$' go.mod
go test ./...
go run ./cmd/stepdown ./testdata/...
go run ./cmd/stepdown ./...
```

The `go run ./cmd/stepdown ./...` command is the self-policing run after implementation. If the positive witness fixture command `go run ./cmd/stepdown ./testdata/...` needs a package-pattern adjustment because `testdata` is not a package pattern accepted by `go/packages`, Implementor records the exact equivalent command in the handoff without changing the proof property: command-level run against positive witness fixtures and self-policing run against the tool's own source.

Reviewer inspection:

- direct inspection of the positive grammar walker;
- direct inspection that diagnostic generation derives from grammar state and observed facts;
- direct inspection that a no-op analyzer would fail the review gate;
- direct inspection that no rejected-form source artifacts exist.

Foundation Auditor inspection:

- lift witness for module path, rule names, fixture names, diagnostics, docs, and integration tokens;
- confirmation that steward identity remains locator metadata;
- confirmation that the source does not become a rejected-form container;
- confirmation that no Pay-bye, Agent OS, workflow, foundation, watchlist, payroll, or product vocabulary appears in tool semantics, fixtures, diagnostics, config, or tests.

## Positive Witness Fixtures

Implementor creates the ADR-0001 positive witness cases, unless a case is made redundant by an identical smaller fixture and the handoff records the equivalence:

- `single_type`
- `multi_type`
- `support_types`
- `exported_package_function`
- `mixed_receivers`
- `pointer_constructor`
- `pointer_constructor_with_error`
- `nested_dfs`
- `shared_dfs_callee`

Additional fixtures are allowed only when they are positive witnesses for an ADR-0001 grammar fact not covered above. Additional fixtures must stay synthetic and domain-neutral.

## Review And Audit Gates

Reviewer gate:

- implementation readiness;
- ADR fidelity;
- positive-walker enforceability;
- no slop tests;
- no rejected-form source artifacts;
- no silent scope invention;
- no no-op analyzer acceptance.

Foundation Auditor gate:

- mandatory audit because this tool defines a foundation-relevant verification discipline;
- lift witness across source, fixtures, diagnostics, docs, and config;
- no organization identity in tool semantics;
- no rejected-form source corpus;
- no source-side rejected-shape catalog;
- no consumer adoption facts inside the tool.

Founder Chief Architect closeout:

- runs after Reviewer and Foundation Auditor acceptance;
- verifies that the `Positive Grammar Walker Requirement` remains load-bearing;
- verifies no rejected-form source corpus became load-bearing for Implementor;
- authorizes only the next implementation stage, not consumer adoption.

## Release And Recovery Plan

Release shape:

- v0.1.0 implementation completes only after all implementation units pass review and audit.
- Public usage remains `go run github.com/pay-bye/stepdown/cmd/stepdown@v0.1.0 ./...`.
- Consumer adoption is separate and requires consumer-side records.

Recovery:

- If the walker becomes too large to inspect directly, stop and route a new ADR decision for test strategy.
- If legitimate Go code consistently needs waivers, revise the grammar through ADR authority; do not add waivers.
- If a fixture or diagnostic leaks consumer vocabulary, replace it with Go-language vocabulary before release.
- If the tool drifts into semantic, security, performance, API-design, or domain-policy linting, remove the drift and route that concern to another tool.
- If maintainer succession becomes unclear, amend ADR-0001 or publish a successor ADR before release dependency expands.

## Translation Choices Left To Implementor

| Choice | Allowed decision space | Required Implementor record |
|---|---|---|
| Internal package split | Any small package split that keeps loading, modeling, classification, walking, diagnostics, and CLI responsibilities reviewable. | Handoff names packages and one-sentence responsibility for each. |
| Diagnostic descriptions | Human-readable descriptions may vary as long as rule names, positions, exit code class, and positive expectation are stable. | Handoff records examples for each rule name and confirms no consumer vocabulary. |
| Fixture command pattern | Use the exact command that loads positive witness fixtures through the command path. | Handoff records command and explains why it proves positive witness fixture behavior. |
| Self-policing command before release tag exists | Use `go run ./cmd/stepdown ./...` or an equivalent built-command invocation. | Handoff records command, output, and runtime. |
| Tool/load error tests | Use focused temp-package or in-memory setup for parse/type/package failures. | Handoff records why each test is error plumbing and not structural rejected-form evidence. |

No tool semantics, grammar rule, fixture policy, diagnostic rule name, exit code, file selection rule, waiver behavior, pinning mechanism, consumer adoption fact, or distribution form is left to Implementor judgment.

## Author Checklist

- [x] ADR-0001 remains semantic authority.
- [x] The required `Positive Grammar Walker Requirement` section is present.
- [x] The spec rejects implementation by named bad-shape switching.
- [x] The spec rejects committed rejected-form fixtures, `testdata/violations/`, rejected-form `expected.txt`, and per-case diagnostic tables.
- [x] The spec defines implementation slices, verification contract, review gates, and foundation audit gates.
- [x] Positive witness fixtures are supporting evidence, not proof that invalid structures are rejected.
- [x] Rejection behavior is proven by positive grammar walker design and review inspection, not a bad-shape fixture corpus.
- [x] Workflow locators and control-plane facts stay in the workflow wrapper, not in this durable tool spec.
