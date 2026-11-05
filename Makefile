INFISICAL_ENV ?= dev
INFISICAL_PATH ?= /backend

.PHONY: run run-secrets build test vet fmt db-push db-pull db-reset

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

fmt:
	go fmt ./...
