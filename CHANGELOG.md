# Changelog

All notable changes to this project are documented in this file.

This project follows Semantic Versioning.

## [v0.1.0] - 2026-05-28

Initial release. `stepdown` enforces top-down declaration order in Go source files through a positive grammar, specified by [ADR-0001](docs/adr/0001-stepdown-go-structure-analyzer.md).

### Added

- Structural analyzer that enforces file-level section order: package, imports, constants, package vars, per-type blocks, exported package-level functions, then unexported package-level helper functions.
- Per-type blocks ordered as constructors, getters, setters, then non-accessor receiver methods, with each type's block contiguous in source order.
- Depth-first ordering of receiver methods from exported roots: each exported method is followed by the unexported methods it calls, with shared callees owned by the first exported root in source order.
- Acceptance of single-spec type declarations of any form (struct, interface, named primitive, function type, map, slice, alias). Getter and setter rules apply only to struct types.
- Same-file receiver type requirement: a receiver method's type must be declared in the same file.
- File selection via `go/packages`, skipping test files, generated files (Go generated-code marker), and files outside the default build.
- Diagnostics in `file:line:column: rule-name: description` format with stable rule names.
- Command-line interface accepting Go package patterns, with `-h`/`--help`/`-help` and exit codes `0` (clean), `1` (findings), `2` (tool or load error).
- Self-policing: `stepdown` runs against its own source as part of verification.
- Apache 2.0 license.

[v0.1.0]: https://github.com/pay-bye/stepdown/releases/tag/v0.1.0
