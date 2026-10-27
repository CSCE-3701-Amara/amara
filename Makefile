INFISICAL_ENV ?= dev
INFISICAL_PATH ?= /

.PHONY: run run-secrets build test vet fmt

run:
	go run ./cmd/api

run-secrets:
	infisical run --env="$(INFISICAL_ENV)" --path="$(INFISICAL_PATH)" -- go run ./cmd/api

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...
