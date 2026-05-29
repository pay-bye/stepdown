#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

go version | grep -Eq 'go1\.26\.3'
test "$(go env GOTOOLCHAIN)" = "local"
grep -Eq '^go 1\.26\.3$' go.mod
go test ./...
go run ./cmd/stepdown -h >/dev/null
go run ./cmd/stepdown --help >/dev/null
go run ./cmd/stepdown -help >/dev/null
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.3.0 ./...

fixtures=()
while IFS= read -r dir; do
  fixtures+=("./${dir}")
done < <(find testdata -mindepth 1 -maxdepth 1 -type d | sort)

go run ./cmd/stepdown "${fixtures[@]}"
go run ./cmd/stepdown ./...
