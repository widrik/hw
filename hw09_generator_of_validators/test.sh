#!/usr/bin/env bash
set -euo pipefail

rm -f "$(command -v go-validate)"
rm -f ./models/*generated.go

go install ./go-validate
go generate models/models.go
go test -v -tags generation ./models

rm -f go-validate
echo "PASS"
