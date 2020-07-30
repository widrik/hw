build:
	go build -o calendar ./cmd/calendar/main.go

run:
	go run ./cmd/calendar/main.go -config ./config/main.json

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./test

lint:
    go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...

.PHONY: build
