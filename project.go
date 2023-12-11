package main

import (
	"context"

	"pkg.package-operator.run/cardboard/run"
	"pkg.package-operator.run/cardboard/sh"
)

func main() {
	run.New().Run(&Test{})
}

type Test struct{}

// run with: "go run ./file.go test:unit"
func (t *Test) Unit(ctx context.Context, args []string) error {
	return sh.New(sh.Environment{
		"CGO_ENABLED": "1",
	}).Bash([]string{
		"set -euo pipefail",
		"go test -coverprofile=cover.txt -race -json ./... 2>&1 | tee gotest.log | gotestfmt --hide=empty-packages",
	})
}
