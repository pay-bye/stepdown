#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

go version | rg 'go1\.26\.3'
test "$(go env GOTOOLCHAIN)" = "local"
rg -n '^go 1\.26\.3$' go.mod
go test ./...

fixtures=()
while IFS= read -r dir; do
  fixtures+=("./${dir}")
done < <(find testdata -mindepth 1 -maxdepth 1 -type d | sort)

go run ./cmd/stepdown "${fixtures[@]}"
go run ./cmd/stepdown ./...
