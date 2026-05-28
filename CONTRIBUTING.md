# Contributing to stepdown

`stepdown` is governed by Architecture Decision Records (ADRs) under `docs/adr/`. Before contributing, read the canonical specification:

- **[ADR-0001: Stepdown Go Structure Analyzer](docs/adr/0001-stepdown-go-structure-analyzer.md)**

## ADR-driven evolution

The grammar `stepdown` enforces is described entirely in ADRs. The principle:

- **Bug fixes, parser compatibility updates, diagnostic message improvements, and performance work** do not require new ADRs. They are maintainer discretion under semantic-versioning patch releases.
- **Anything that changes what `stepdown` accepts or rejects as conforming Go source** requires a new ADR in this repository's `docs/adr/` sequence (ADR-0002, ADR-0003, etc.) that explicitly cites the ADR it amends or supersedes.

New rules must trace back to `stepdown`'s one opinion: source structure should read top-down. New rules need explicit justification: what structural failure mode the rule catches, why review and existing rules are insufficient, why the rule cannot be expressed under an existing rule, what edge cases it handles, and what edge cases it does NOT handle.

The maintainer rejects rules that drift toward general-purpose Go style enforcement, semantic correctness, security, performance, or API design. Those are out of scope.

## Reporting issues

**`stepdown` rejects a legitimate Go idiom the grammar did not anticipate:** file an issue with a minimal reproduction. The maintainer evaluates whether the grammar needs adjustment.

**`stepdown` accepts a structural shape that should fail:** file an issue with a minimal example. Same evaluation.

**`stepdown` is too slow on real code:** file an issue with profiling data. Performance work is maintainer discretion.

## Pull request requirements

- **Reviewability.** Changes that touch the analyzer's enforcement logic must keep the walker small enough to verify by direct inspection. The proof model in ADR-0001 § Implementation discipline depends on this; if a change makes the walker too large to read and verify, the test strategy itself must be revisited in a new ADR before merge.
- **Self-policing.** The repository's CI step that runs `stepdown` against its own source must pass before merge. A release that does not pass its own check is not shipped.
- **Local verification.** Run `./scripts/verify.sh` before submitting analyzer changes. The script runs the tests, every positive witness fixture, and the self-policing command.
- **No waivers.** PRs that introduce inline waiver comments, per-file opt-out mechanisms, or configuration flags to silence findings are rejected. If a class of valid code consistently fails, propose a grammar adjustment via new ADR — do not add a waiver.
- **Fixtures.** Fixture additions follow the discipline in ADR-0001 § Fixture policy: self-contained Go using generic placeholder identifiers (`Foo`, `Bar`, `Baz`, `Widget`, etc.), no production-system identity, no business-domain vocabulary, no consumer-specific references, no inline comments describing what the fixture verifies.
- **No negative fixtures.** `testdata/` contains only positive witnesses (valid Go conforming to the grammar). There is no `testdata/violations/` corpus and no `expected.txt` mechanism. See ADR-0001 § Implementation discipline > Proof surface and residual risk for the rationale.

## Licensing

Contributions are licensed under the Apache License 2.0. See [LICENSE](LICENSE).

By submitting a contribution, you certify that you have the right to license it under Apache 2.0 and that you intend to do so.
