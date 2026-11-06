INFISICAL_ENV ?= dev
INFISICAL_PATH ?= /backend

.PHONY: run run-secrets build test vet security gosec govulncheck fmt db-push db-pull db-reset

run:
	go run ./cmd/api

run-secrets:
	infisical run --env="$(INFISICAL_ENV)" --path="$(INFISICAL_PATH)" -- go run ./cmd/api

# Database schema targets require the Supabase CLI login (`supabase login`),
# the project link (`supabase link`), and the database password (stored after
# linking or supplied through SUPABASE_DB_PASSWORD).
db-push:
	supabase db push

db-pull:
	supabase db pull

db-reset:
	supabase db reset

build:
	mkdir -p bin
	go build -o bin/api ./cmd/api

test:
	go test ./...

vet:
	go vet ./...

security:
	./scripts/check-security.sh

gosec:
	gosec -conf .gosec.json ./...

govulncheck:
	govulncheck ./...
	govulncheck -C cmd/api -scan module
	govulncheck -C cmd/worker -scan module

fmt:
	go fmt ./...
