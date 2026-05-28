# stepdown

A structural Go-language source analyzer that enforces top-down readability through a positive grammar over Go's abstract syntax tree.

## What it does

`stepdown` checks that Go source files are organized so they read top-down by call dependency: high-level declarations appear before the supporting declarations they depend on. It addresses a 2026-era problem — AI-driven code generators do not naturally produce code humans can read top-down — by enforcing the structural rule mechanically at the verification gate level.

`stepdown` is intentionally narrow. It enforces source-file structure, not semantic correctness, not security, not performance, not API design. It is one rule, expressed as a positive grammar. Deviations from the grammar fail. There are no waivers.

## Usage

```
go run github.com/pay-bye/stepdown/cmd/stepdown@<version> ./...
```

Where `<version>` is a published git tag (e.g., `v0.1.0`).

Exit codes:

- `0` — clean (no findings)
- `1` — findings present (at least one source file does not conform to the grammar)
- `2` — tool/load error

Diagnostic output uses the standard Go diagnostic format:

```
file:line:column: <rule-name>: <description>
```

## The grammar

Each non-test, non-generated Go source file in the default build configuration must structure its declarations as:

```
package
import
constants
package vars

for each type declared in the file, in source order:
    type declaration (single-spec)
    constructors
    getters
    setters
    non-accessor receiver methods (DFS-from-public-roots ordered)

exported package-level functions
unexported package-level helper functions
```

Sections may be empty. The grammar accepts single-spec type declarations of any form (`struct`, named primitive, function type, map, slice, interface, alias). Receiver methods must be declared in the same file as their receiver type.

For the complete specification — classification predicates, DFS ownership rules, rule names, file selection, and the canonical fixture set — see:

**[ADR-0001: Stepdown Go Structure Analyzer](docs/adr/0001-stepdown-go-structure-analyzer.md)**

The ADR is canonical for the tool's behavior. This README is a summary.

## Status

`stepdown` is at the specification stage. ADR-0001 is Accepted as of 2026-05-28. v0.1.0 implementation work follows.

## License

Apache License 2.0. See [LICENSE](LICENSE).

## Contributing

`stepdown` is governed by ADRs under `docs/adr/`. See [CONTRIBUTING.md](CONTRIBUTING.md) for the evolution process for new rules and the discipline that applies to changes.
