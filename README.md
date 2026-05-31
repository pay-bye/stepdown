# stepdown

[![Verify](https://github.com/stepdown-dev/stepdown-go/actions/workflows/verify.yml/badge.svg)](https://github.com/stepdown-dev/stepdown-go/actions/workflows/verify.yml)
[![License](https://img.shields.io/github/license/stepdown-dev/stepdown-go)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/stepdown-dev/stepdown-go)](go.mod)
[![Go Reference](https://pkg.go.dev/badge/stepdown.dev/go.svg)](https://pkg.go.dev/stepdown.dev/go)
[![Go Report Card](https://goreportcard.com/badge/stepdown.dev/go)](https://goreportcard.com/report/stepdown.dev/go)

**A Go linter that keeps source files readable top to bottom.**

Good Go files read like a newspaper: the headline first, the details below. The type comes before the code that builds it, public methods come before the private helpers they call, and you never have to scroll up to understand what you're looking at. `stepdown` enforces that order mechanically, so it stays true no matter how many edits — human or machine — a file goes through.

That last part is the point. Code generators are good at writing correct functions and bad at placing them: a helper lands above the method that calls it, a constructor drifts below the methods that use it, declarations pile up wherever the cursor happened to be. Each edit is locally fine and the file slowly stops reading top-down. `stepdown` makes the ordering a check instead of a habit.

## Example

This file passes:

```go
package cache

import "errors"

var ErrMissing = errors.New("cache: key not found")

type Cache struct {
	entries map[string]string
}

func NewCache() *Cache {
	return &Cache{entries: map[string]string{}}
}

func (c *Cache) Lookup(key string) (string, error) {
	if err := c.require(key); err != nil {
		return "", err
	}
	return c.entries[key], nil
}

func (c *Cache) require(key string) error {
	if _, ok := c.entries[key]; !ok {
		return ErrMissing
	}
	return nil
}
```

Read it straight down: the type, then how you build it, then what it does, with `require` sitting right below the method that calls it. Move `require` above `Lookup`, drop `NewCache` beneath the methods that use it, or wedge a package-level helper between two methods, and `stepdown` reports the file with a `file:line:column` diagnostic and a non-zero exit code.

## Usage

```
go run stepdown.dev/go/cmd/stepdown@v0.1.1 ./...
```

Drop that into a CI step, or run it from a clone with `go run ./cmd/stepdown ./...`. It takes Go package patterns and analyzes the non-test, non-generated files in the default build.

Exit codes:

- `0` — clean
- `1` — one or more files do not conform
- `2` — could not analyze (usage, package load, parse, or type-resolution error)

Diagnostics use the standard Go format, so editors and CI pick them up without configuration:

```
file:line:column: rule-name: description
```

## What it enforces

Each non-test, non-generated file in the default build must order its declarations like this:

```
package
import
constants
package vars

for each type, in source order:
    type declaration
    constructors
    getters
    setters
    methods (each public method followed by the private methods it calls, depth-first)

exported package-level functions
unexported package-level helper functions
```

Sections are optional — an empty file with just a package clause passes. Type declarations of any form are accepted (struct, interface, named primitive, function type, alias, and the rest); getters and setters simply have nothing to match on non-struct types. A type's methods stay together with its declaration in the same file.

## What it doesn't do

`stepdown` checks one thing: declaration order. It does not check correctness, security, performance, or API design — `go vet`, `staticcheck`, `gosec`, and `govulncheck` already do those, and `stepdown` runs happily alongside them.

It has no configuration file, no rule toggles, and no per-line ignore comments. The order is the order. If a piece of valid Go consistently can't satisfy it, that's a bug in the grammar — [open an issue](https://github.com/stepdown-dev/stepdown-go/issues), don't reach for a waiver.

## Documentation

The complete specification — every classification rule, the depth-first ordering, file selection, and diagnostics — lives in the architecture decision record:

**[ADR-0001: Stepdown Go Structure Analyzer](docs/adr/0001-stepdown-go-structure-analyzer.md)**

The ADR is canonical for the tool's behavior; this README is the tour. `stepdown` is governed by ADRs under `docs/adr/`, and new rules arrive through new ADRs rather than configuration.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for the development setup, the verification script, and the discipline that applies to changes.

## License

[Apache License 2.0](LICENSE).

## Family

`stepdown` is the Go member of the [stepdown family](https://github.com/stepdown-dev) of structural source analyzers — all sharing one [constitution](https://github.com/stepdown-dev/.github/blob/main/PRINCIPLES.md) (positive grammar, no configuration, no waivers, self-policing). A TypeScript sibling, `stepdown-ts`, is in development.
