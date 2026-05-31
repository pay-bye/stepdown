# Repository Hardening v0.1.0

Status: Draft for implementation
Date: 2026-05-29

## Problem

`stepdown` has executable analyzer behavior and local verification, but the repository lacks the public project surfaces required before v0.1.0: CI, security disclosure, changelog, contribution templates, ownership routing, dependency update automation, command help, and static-analysis verification.

Affected readers are maintainers, contributors, security reporters, and release reviewers. The spec covers repository hardening only. It does not change analyzer semantics, release a tag, configure repository settings, enable branch protection, register external services, or prove that external badges and hosted services have indexed the repository. The only command behavior added is help output for explicit help aliases.

## Success Criteria

- Repository verification runs locally through `./scripts/verify.sh`.
- GitHub Actions runs the same verification script on every push and pull request.
- Static analysis covers `staticcheck`, `gosec`, `gofmt`, and `govulncheck`.
- `stepdown -h`, `stepdown --help`, and `stepdown -help` print short help to stdout and exit `0`.
- Public policy files exist for security disclosure, changelog, contribution conduct, ownership, issue routing, and pull request review.
- README badges render the intended external status links without claiming external service registration.
- Repository-hardening files preserve ADR-0001: no waivers, no rejected-form fixture corpus, no per-case diagnostic assertion tables, and no semantics outside source-file structure.
- Every new file has one local responsibility and a stable path.

## Assumptions

- The repository is `github.com/stepdown-dev/stepdown-go` and the module path is the vanity `stepdown.dev/go`; the steward-org token is locator metadata, not tool semantics.
- `go.mod` declares `go 1.26.3`.
- The maintainer GitHub handle for ownership routing is `@johnastinnett`.
- Private security disclosure uses GitHub Private Security Advisories.
- Contributor Covenant enforcement reports go to `john.a.stinnett@gmail.com`.
- GitHub Actions uses GitHub-hosted `ubuntu-24.04` runners.
- The implementation does not need repository-admin permission because repository settings are out of scope.

## Constraints And Non-Goals

- Do not change analyzer behavior, diagnostics, fixtures, ADR-0001 semantics, or the public invocation model except for the help aliases defined in this spec.
- Preserve zero-argument behavior: no package patterns exits `2` with usage text on stderr.
- Do not add waiver mechanisms, negative fixtures, `testdata/violations/`, `expected.txt`, bad-shape catalogs, or per-case diagnostic assertions.
- Do not configure GitHub repository settings, branch protection, repository topics, repository description, private vulnerability reporting toggles, Go Report Card registration, pkg.go.dev indexing, or release tags.
- Do not add non-Go Dependabot ecosystems.
- Do not add GitHub issue templates that contain example rejected-form source.
- Do not introduce a second verification script. `scripts/verify.sh` remains the local and CI command surface.
- Do not require unprovisioned CI runner tools in `scripts/verify.sh`.
- Do not widen lint correction beyond the exact findings needed for the pinned `golangci-lint` gate.

## Options Considered

### Issue-template shape

| Option | Value | Complexity | Change cost | Failure mode | Reversibility | Time to first result |
|---|---|---|---|---|---|---|
| Markdown templates | Simple files, easy to edit. | Low. | Low. | Required data can be omitted; new-rule proposals drift into vague requests. | Easy. | Immediate. |
| GitHub issue forms | Structured fields make minimal reproductions, expected behavior, actual behavior, and ADR justification explicit. | Moderate. | Low. | YAML can over-prescribe examples or accidentally invite rejected-form source into repository templates. | Easy. | Immediate. |

Choice: use GitHub issue forms. The forms collect structured reports without embedding bad examples in source. The new-rule proposal form routes semantic expansion through ADR discipline.

### Static-analysis installation

| Option | Value | Complexity | Change cost | Failure mode | Reversibility | Time to first result |
|---|---|---|---|---|---|---|
| Install binaries globally in CI | Familiar to contributors. | Moderate. | Moderate. | CI and local versions drift; PATH setup becomes a second authority. | Moderate. | Fast. |
| Use pinned `go run` commands in `scripts/verify.sh` | Versions live in one script and run the same locally and in CI. | Low. | Low. | First run downloads tools and is slower. | Easy. | Fast. |

Choice: use pinned `go run` commands. This keeps tool versions explicit and avoids a separate install surface.

### Contact-channel handling

| Option | Value | Complexity | Change cost | Failure mode | Reversibility | Time to first result |
|---|---|---|---|---|---|---|
| Invent generic contacts | Unblocks file creation superficially. | Low. | High. | Reports route nowhere; the spec violates founder direction. | Hard. | Immediate but invalid. |
| Return missing context | Preserves authority when no grounded contact exists. | Low. | Low. | Work pauses until founder supplies contacts. | Easy. | Slower but correct. |
| Use supplied principal context | Implements grounded contacts without changing scope. | Low. | Low. | External repository settings still need separate admin work. | Easy. | Immediate. |

Choice: use supplied principal context. `SECURITY.md` points to GitHub Private Security Advisories. `CODE_OF_CONDUCT.md` uses `john.a.stinnett@gmail.com` for enforcement reports.

### `govulncheck` placement

| Option | Value | Complexity | Change cost | Failure mode | Reversibility | Time to first result |
|---|---|---|---|---|---|---|
| Force `govulncheck` into `golangci-lint` | One config file appears to cover every requested tool. | High. | High. | Unsupported linter name breaks verification. | Easy. | Invalid. |
| Run `govulncheck` separately from `scripts/verify.sh` | Matches the available tool surface while preserving one verification script. | Low. | Low. | Reviewer may expect all tools inside `.golangci.yml` unless the choice is recorded. | Easy. | Immediate. |

Choice: run `govulncheck` separately from `scripts/verify.sh`. `.golangci.yml` covers `staticcheck`, `gosec`, and `gofmt`; the script covers `govulncheck` with an explicit version pin.

## Recommended Approach

Add focused repository-hardening files, keep `scripts/verify.sh` as the only local and CI verification command, and route every public contribution surface back to ADR-0001. This is the smallest design that gives maintainers reviewable governance without changing tool semantics.

External references used for tool and policy shape:

- GitHub private vulnerability reporting: <https://docs.github.com/en/code-security/how-tos/report-and-fix-vulnerabilities/configure-vulnerability-reporting/configuring-private-vulnerability-reporting-for-a-repository>
- GitHub Dependabot file reference: <https://docs.github.com/en/code-security/concepts/supply-chain-security/about-the-dependabot-yml-file>
- golangci-lint v2 configuration: <https://golangci-lint.run/docs/configuration/file/>
- golangci-lint linter list: <https://golangci-lint.run/docs/linters/>
- Contributor Covenant adoption guide: <https://www.contributor-covenant.org/adopt/>

## Boundary And Ownership Model

| Boundary | Target path | Owner | Responsibility | Canonical fact |
|---|---|---|---|---|
| Local verification | `scripts/verify.sh` | Maintainer | Run Go toolchain preflight, tests, positive fixtures, self-policing, static analysis, and vulnerability scan. | Verification command sequence lives in this script. |
| CI verification | `.github/workflows/verify.yml` | Maintainer | Run `./scripts/verify.sh` on push and pull request. | CI delegates to local verification. |
| Static-analysis config | `.golangci.yml` | Maintainer | Configure `staticcheck`, `gosec`, and `gofmt`. | Enabled linter/formatter set lives in this file. |
| Security policy | `SECURITY.md` | Maintainer | Direct private vulnerability reports to GitHub Private Security Advisories. | Disclosure channel lives in this file. |
| Changelog | `CHANGELOG.md` | Maintainer | Maintain release notes in Keep a Changelog shape. | Release-note text lives in this file. |
| Issue routing | `.github/ISSUE_TEMPLATE/*.yml` | Maintainer | Collect structured reports without storing rejected-form examples in repository source. | Issue intake shape lives in issue forms. |
| Pull request review | `.github/PULL_REQUEST_TEMPLATE.md` | Maintainer | Preserve CONTRIBUTING review requirements. | PR checklist lives in this file. |
| Conduct policy | `CODE_OF_CONDUCT.md` | Maintainer | Adopt Contributor Covenant English 3.0 and name enforcement contact. | Conduct policy lives in this file. |
| Ownership routing | `.github/CODEOWNERS` | Maintainer | Request maintainer review for every change. | Ownership routing lives in this file. |
| Dependency updates | `.github/dependabot.yml` | Maintainer | Open Go module dependency update pull requests. | Dependency update schedule lives in this file. |
| Project badges | `README.md` | Maintainer | Show build, license, Go version, Go reference, and Go Report Card links. | Badge markup lives in README. |
| Command help | `cmd/stepdown/main.go`, `cmd/stepdown/main_test.go` | Maintainer | Print explicit command help for help aliases without running analysis. | Help behavior lives in command code and command tests. |

## Contracts

### Command help

Inputs: command arguments and standard output/error writers.

Help aliases: `-h`, `--help`, and `-help`.

Contract:

- If the command receives exactly one argument and it is a help alias, it writes help text to stdout, writes nothing to stderr, and exits `0`.
- If the command receives no arguments, existing behavior is preserved: usage text is written to stderr and exit code is `2`.
- Mixed help aliases and package patterns are usage errors: write usage text to stderr and exit `2`. This preserves explicit package-pattern invocation and avoids implicit flag parsing.
- Help output does not run package loading or analysis.
- Help output contains no waiver, configuration, or rule-toggle language.

Required help content:

```text
stepdown - Go source structure analyzer for top-down declaration order

Usage:
  stepdown <package-pattern> [<package-pattern>...]
  stepdown -h
  stepdown --help
  stepdown -help

Examples:
  stepdown ./...
  go run stepdown.dev/go/cmd/stepdown@v0.1.1 ./...

Exit codes:
  0  no findings, or help printed
  1  source structure findings
  2  usage, load, parse, type-resolution, or output error

ADR: https://github.com/stepdown-dev/stepdown-go/blob/main/docs/adr/0001-stepdown-go-structure-analyzer.md
README: https://github.com/stepdown-dev/stepdown-go#readme
```

### `scripts/verify.sh`

Inputs: repository working tree, Go toolchain, standard shell utilities available on Ubuntu runners such as `grep`, network access for pinned `go run` tools on first use.

Outputs: exit `0` only when every gate passes. Any failing command exits nonzero.

Required command order:

```bash
go version | grep -Eq 'go1\.26\.3'
test "$(go env GOTOOLCHAIN)" = "local"
grep -Eq '^go 1\.26\.3$' go.mod
go test ./...
go run ./cmd/stepdown -h >/dev/null
go run ./cmd/stepdown --help >/dev/null
go run ./cmd/stepdown -help >/dev/null
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...
go run ./cmd/stepdown "${fixtures[@]}"
go run ./cmd/stepdown ./...
```

The positive fixture discovery block remains mechanical: sorted `testdata/*` directories become `./testdata/<case>` package patterns. No negative fixture discovery is added.

The script does not call `rg` because the GitHub-hosted runner contract does not provision ripgrep as part of this repository's declared tool surface.

### `.golangci.yml`

Required shape:

```yaml
version: "2"

run:
  timeout: 5m

linters:
  default: none
  enable:
    - gosec
    - staticcheck

formatters:
  enable:
    - gofmt
```

`gofmt` is under `formatters`, not `linters`, because golangci-lint v2 moved Go formatters into the formatter section. `govulncheck` is absent from this file and runs separately through `scripts/verify.sh`.

### `.github/workflows/verify.yml`

Policy property: every pushed commit and pull request executes the same local verification contract under Go 1.26.3.

Required shape:

```yaml
name: Verify

on:
  push:
  pull_request:

permissions:
  contents: read

jobs:
  verify:
    runs-on: ubuntu-24.04
    steps:
      - uses: actions/checkout@v6.0.2
      - uses: actions/setup-go@v6.4.0
        with:
          go-version-file: go.mod
          cache: true
      - run: ./scripts/verify.sh
        env:
          GOTOOLCHAIN: local
```

No upload, release, deploy, code-scanning SARIF, or repository-setting step is added.

### `SECURITY.md`

Policy property: reporters have one private disclosure path and are told not to disclose vulnerabilities publicly before coordination.

Required content facts:

- Supported version: `v0.1.0` once released; unreleased `main` receives maintainer review before v0.1.0.
- Private disclosure path: GitHub Private Security Advisories for this repository.
- Reporter instruction: use the repository Security tab and its private vulnerability reporting flow.
- Do not include a direct email address for security reports.
- State that repository administrators must enable GitHub private vulnerability reporting outside this implementation when the repository is ready for public reporting. This statement is an external-setting dependency, not an implementation task.

### `CHANGELOG.md`

Policy property: release notes have one durable home.

Required shape:

```markdown
# Changelog

All notable changes to this project are documented in this file.

This project follows Semantic Versioning.

## [v0.1.0] - Unreleased
```

The `v0.1.0` section stays empty until release content is known. The file does not claim a release date or tag.

### Issue Templates

Create exactly these issue forms:

- `.github/ISSUE_TEMPLATE/legitimate-idiom-rejected.yml`
- `.github/ISSUE_TEMPLATE/structural-shape-accepted.yml`
- `.github/ISSUE_TEMPLATE/new-rule-proposal.yml`

Shared constraints:

- Each form uses GitHub issue-form YAML, not freeform Markdown.
- Each form asks for `stepdown` version, Go version, command run, observed result, expected result, and minimal reproduction.
- No template embeds sample malformed source.
- Minimal reproduction fields are empty textareas for reporter input.
- Reporter prompt text says not to include production-system identity, secrets, customer data, or private code.

`legitimate-idiom-rejected.yml` captures valid Go source rejected by the grammar. It asks which ADR-0001 rule appears too narrow.

`structural-shape-accepted.yml` captures source accepted by the tool that violates top-down readability. It asks which top-down readability failure escaped detection.

`new-rule-proposal.yml` captures new rule proposals. It requires:

- structural failure mode;
- why human review and existing rules do not catch it;
- why the rule traces to top-down source readability;
- why the rule is not semantic correctness, security, performance, API design, or domain policy;
- ADR requirement acknowledgement.

### Pull Request Template

Target path: `.github/PULL_REQUEST_TEMPLATE.md`.

Checklist items:

- `./scripts/verify.sh` passed locally.
- The walker remains small enough for direct review.
- Self-policing passes.
- No waiver mechanism was added.
- Fixtures are positive witnesses only.
- No `testdata/violations/`, `expected.txt`, negative fixture corpus, or per-case diagnostic table was added.
- Any semantics change is backed by a new ADR.
- Documentation changes preserve ADR-0001 scope.

### Code Of Conduct

Target path: `CODE_OF_CONDUCT.md`.

Policy property: project participation has an explicit conduct standard and enforcement contact.

Required shape:

- Adopt Contributor Covenant English 3.0.
- Set reporting/enforcement contact to `john.a.stinnett@gmail.com`.
- Do not add a separate committee, anonymous form, chat channel, or organization-specific enforcement structure.
- Keep the upstream attribution and version link.

### CODEOWNERS

Target path: `.github/CODEOWNERS`.

Required content:

```text
* @johnastinnett
```

This file routes every pull request to John Stinnett for review. It does not configure branch protection.

### Dependabot

Target path: `.github/dependabot.yml`.

Policy property: Go module dependency update pull requests are created on a predictable cadence without adding unrelated ecosystems.

Required shape:

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
      day: monday
      time: "09:00"
      timezone: America/Chicago
    open-pull-requests-limit: 5
```

Do not add `github-actions`, `docker`, `npm`, or other ecosystems in this unit.

### README Badges

Add badges at the top of `README.md`, below the H1 and before descriptive prose.

Required badge set:

```markdown
[![Verify](https://github.com/stepdown-dev/stepdown-go/actions/workflows/verify.yml/badge.svg)](https://github.com/stepdown-dev/stepdown-go/actions/workflows/verify.yml)
[![License](https://img.shields.io/github/license/stepdown-dev/stepdown-go)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/stepdown-dev/stepdown-go)](go.mod)
[![Go Reference](https://pkg.go.dev/badge/stepdown.dev/go.svg)](https://pkg.go.dev/stepdown.dev/go)
[![Go Report Card](https://goreportcard.com/badge/stepdown.dev/go)](https://goreportcard.com/report/stepdown.dev/go)
```

The badges prove only that README links to the intended status surfaces. External indexing and successful hosted runs remain out of scope.

## Progression And Invariants

Repository-hardening progression:

1. Add command help behavior and tests.
2. Extend local verification with help smoke checks, static analysis, and vulnerability scanning.
3. Add CI workflow that delegates to local verification.
4. Implement policy and template files.
5. Add README badges after target workflow path exists.

Invariants:

- `scripts/verify.sh` is the only verification command CI calls.
- Go version remains 1.26.3.
- `GOTOOLCHAIN=local` remains required.
- Tool versions are pinned: `golangci-lint` `v2.12.2`, `govulncheck` `v1.3.0`, `actions/checkout` `v6.0.2`, `actions/setup-go` `v6.4.0`.
- External service activation remains outside the repository diff.

## Risks And Failure Modes

| Failure mode | Detection | Containment | Recovery |
|---|---|---|---|
| CI diverges from local verification | Reviewer compares workflow command to `./scripts/verify.sh`. | Workflow runs only the script. | Replace inline CI commands with script call. |
| `govulncheck` is treated as a golangci-lint linter | `golangci-lint run` fails on unknown linter or config review finds unsupported name. | `govulncheck` runs separately in the script. | Remove it from `.golangci.yml`; keep pinned `go run`. |
| Verification script depends on unprovisioned runner tooling | CI fails before Go gates execute. | Script uses `grep` for text checks. | Replace unprovisioned tool calls with runner-provided shell utilities or provision the tool inside `scripts/verify.sh` with a pin. |
| Existing code fails newly pinned lint gates | `golangci-lint v2.12.2 run ./...` reports the exact diagnostics. | Lint correction scope is named in `U2`. | Apply only the named mechanical corrections, then rerun the full script. |
| Command help changes analysis semantics | Command tests show help aliases enter package analysis or zero-argument behavior changes. | Help aliases are handled before analysis only when they are the sole argument. | Restore zero-argument error behavior and isolate help handling from package analysis. |
| Content-wide negative grep rejects legitimate policy text | Reviewer command fails on ADR, README, CONTRIBUTING, or spec prose before implementation changes exist. | Verification uses path checks, behavior-bearing source scans, and inspection instead of repository-wide forbidden-word sweeps. | Replace content-wide scans with policy-specific checks. |
| Templates invite rejected-form source artifacts | Reviewer scans template text for embedded malformed examples. | Templates use empty textareas and reporter instructions. | Remove examples from templates. |
| Conduct or security contacts drift | Reviewer compares files to supplied principal context. | Contacts live in one policy file each. | Patch policy files and record the source of authority. |
| Badges imply external registration | README review checks badge section language. | Badges are links only; no prose claims external service activation. | Remove unsupported claims. |
| Dependabot widens ecosystems | Reviewer checks `.github/dependabot.yml`. | Only `gomod` ecosystem is allowed. | Remove unrelated ecosystems. |

## Delivery Plan

Each unit is one focused Implementor session.

| Unit | Size | Files | Acceptance |
|---|---|---|---|
| `U1` command help | S | `cmd/stepdown/main.go`, `cmd/stepdown/main_test.go` | Help aliases print required help to stdout and exit `0`; zero args still exit `2` with stderr usage. |
| `U2` verification tooling and lint correction | M | `.golangci.yml`, `scripts/verify.sh`, `internal/analyze/error_test.go`, `internal/grammar/classify.go` | Local script includes pinned `golangci-lint` and `govulncheck`; existing Go/test/fixture/self-policing gates remain; help aliases are smoke-checked; named lint corrections are applied without analyzer semantic change. |
| `U3` CI workflow | XS | `.github/workflows/verify.yml` | Workflow runs on push and pull request, uses pinned checkout/setup-go actions, and delegates to `./scripts/verify.sh`. |
| `U4` public policy files | S | `SECURITY.md`, `CODE_OF_CONDUCT.md`, `CHANGELOG.md` | Security, conduct, and changelog files satisfy their contracts without external-setting claims. |
| `U5` contribution routing | S | `.github/ISSUE_TEMPLATE/*.yml`, `.github/PULL_REQUEST_TEMPLATE.md` | Templates preserve ADR-0001 scope and do not include rejected-form examples. |
| `U6` ownership and dependency automation | XS | `.github/CODEOWNERS`, `.github/dependabot.yml` | CODEOWNERS routes all files to `@johnastinnett`; Dependabot watches only Go modules. |
| `U7` README badges | XS | `README.md` | Five badges appear below the H1 and link to the intended status surfaces. |

Do not merge unrelated concerns into a single package or script rewrite. Do not move existing analyzer code.

## Acceptance Criteria Per Unit

| Unit | Criteria |
|---|---|
| `U1` | `-h`, `--help`, and `-help` each print the required help content to stdout, write empty stderr, and return `0`. Zero args still write stderr usage and return `2`. A mixed help alias and package pattern returns `2` without analysis. |
| `U2` | `./scripts/verify.sh` includes the pinned static-analysis and vulnerability commands, uses no unprovisioned runner tools, smoke-checks all help aliases, and preserves existing positive fixture and self-policing commands. The unit also fixes only these known lint findings: `internal/analyze/error_test.go` directory and file permissions satisfy `gosec` G301/G306 with least-permissive test fixture permissions; `internal/grammar/classify.go` simplifies the boolean comparison reported by `staticcheck` S1002 without changing constructor classification behavior. |
| `U3` | Workflow file uses `actions/checkout@v6.0.2`, `actions/setup-go@v6.4.0`, `ubuntu-24.04`, `GOTOOLCHAIN=local`, and no deploy/release/upload steps. |
| `U4` | Policy files contain the supplied contact channels; changelog has an unreleased `v0.1.0` section with no fake date. |
| `U5` | Issue forms collect structured data without embedding example rejected source; PR template mirrors CONTRIBUTING requirements. |
| `U6` | CODEOWNERS contains exactly the default owner rule; Dependabot contains only the `gomod` update configuration. |
| `U7` | README badge markup matches the contract and does not state that external services have already run. |

## Verification Plan

Verification properties:

- Local and CI verification use the same command surface.
- Static analysis and vulnerability checks are version-pinned.
- Public hardening files do not weaken ADR-0001.
- No external admin task is represented as complete by repository files.

Reviewer commands:

```bash
cd /path/to/stepdown
GOTOOLCHAIN=local ./scripts/verify.sh
test ! -d testdata/violations
test -z "$(find testdata -name expected.txt -print -quit)"
! grep -RInE --exclude-dir=.git --include='*.go' 'ignore stepdown|stepdown:ignore|waiver' cmd internal testdata
grep -RInE --exclude-dir=.git 'john\.a\.stinnett@gmail\.com|Private Security Advisories|@johnastinnett|golangci-lint@v2\.12\.2|govulncheck@v1\.3\.0|actions/checkout@v6\.0\.2|actions/setup-go@v6\.4\.0' .
```

Reviewer inspection:

- `.golangci.yml` enables `gosec` and `staticcheck`, and enables `gofmt` under `formatters`.
- `scripts/verify.sh` runs `govulncheck` separately.
- command tests cover all help aliases, zero-argument usage, and mixed help/package usage.
- issue and PR templates preserve positive-witness and no-waiver discipline.
- legitimate ADR, README, CONTRIBUTING, and spec policy language remains allowed even when it names forbidden implementation residues.
- README badges are only status links.

Foundation Auditor inspection:

- Lift the repository to a neutral Go project and verify that tool semantics remain Go-language semantics.
- Confirm steward identity (`stepdown-dev`, Stinnett Holdings LLC) appears only in locator metadata such as repository and badge URLs, never in tool semantics; the vanity module path `stepdown.dev/go` carries no steward token.
- Confirm no policy file claims external GitHub settings were changed.
- Confirm supplied contact values are grounded by the principal decision artifact, not inferred.

## Release And Recovery Plan

Rollout is additive: all new files can be reverted independently. Existing behavior-bearing edits are limited to command help, verification tooling, README badges, and the named mechanical lint corrections.

Rollback:

- Revert the hardening diff if verification blocks unrelated analyzer work.
- If the CI workflow fails because external action versions change, keep local `./scripts/verify.sh` as canonical and update action pins through a follow-up repository-hardening spec.
- If `govulncheck` reports a vulnerability, stop release work and fix or explicitly risk-accept through maintainer decision before v0.1.0.

## Migration Or Cleanup Plan

No existing files are retired.

New files establish durable homes for hardening facts:

- `SECURITY.md` for security reporting policy.
- `CODE_OF_CONDUCT.md` for conduct policy.
- `CHANGELOG.md` for release notes.
- `.github/*` for repository automation and contributor routing.
- `.golangci.yml` for static-analysis configuration.

## Translation Choices Left To Implementor

| Choice | Why left open | Handoff requirement |
|---|---|---|
| Exact Contributor Covenant generated prose around the required contact | The official builder owns the full prose shape; this spec owns version and contact. | Handoff states English 3.0 was used and names the enforcement contact. |
| YAML field descriptions in issue forms | The implementation can word prompts compactly while preserving required fields. | Handoff lists the form paths and confirms no embedded rejected-form examples. |
| `scripts/verify.sh` helper function names | Local shell factoring is implementation detail. | Handoff shows the final command order. |

No target path, version pin, contact channel, verification property, badge URL, ownership rule, Dependabot ecosystem, command help content, command help behavior, or policy file is left to Implementor judgment.

## Predicate Checks

| Predicate | Result |
|---|---|
| Canonical home for verification command sequence | `scripts/verify.sh`. |
| Canonical home for CI trigger | `.github/workflows/verify.yml`. |
| Canonical home for security disclosure path | `SECURITY.md`. |
| Canonical home for conduct enforcement contact | `CODE_OF_CONDUCT.md`. |
| Canonical home for release notes | `CHANGELOG.md`. |
| Canonical home for owner routing | `.github/CODEOWNERS`. |
| Canonical home for command help behavior | `cmd/stepdown/main.go`; command tests protect the aliases and output contract. |
| Non-canonical duplicate reasons | README badges duplicate status URLs as projection for repository front-page status links; issue and PR templates project CONTRIBUTING requirements into GitHub UI forms for contributor intake; command help duplicates ADR and README URLs as projection for command-line usage. |
| External dependency named | GitHub Actions, GitHub Private Security Advisories, Dependabot, pkg.go.dev, Go Report Card, Shields.io, Contributor Covenant, golangci-lint, and Go vulnerability tooling. |
| Unsupported deferral language | None. External admin actions are out of scope rather than deferred. |
| Tests protect requirements | Verification commands protect local/CI equivalence, command help aliases, static analysis, vulnerability scanning, and ADR-0001 preservation. |
