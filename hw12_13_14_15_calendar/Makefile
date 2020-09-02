build:
	go build -o calendar ./cmd/calendar/main.go
	go build -o calendar ./cmd/scheduler/main.go
	go build -o calendar ./cmd/sender/main.go

run:
	go run ./cmd/calendar/main.go -config ./config/main.json

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./test

lint:  install-deps
	golangci-lint run ./...

prepare_gen:
	go install github.com/mailru/easyjson
	go get \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
        github.com/golang/protobuf/protoc-gen-go

generate: prepare_gen
	go generate ./...

install:
	go mod download

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: build run test lint clean migrate generate
