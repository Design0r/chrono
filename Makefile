.PHONY: build dev test generate migrate install backup deploy live/server live/frontend

ifeq ($(OS),Windows_NT)
  BIN_SUFFIX := .exe
  SET_CGO := set CGO_ENABLED=1&&
else
  BIN_SUFFIX :=
  SET_CGO := CGO_ENABLED=1
endif

MIGRATION_DIR := ./db/migrations
DB_DIR := ./db/chrono.db
GOBIN := $(shell go env GOPATH)/bin

generate:
	@echo "Generating sqlc repository..."
	@sqlc generate

migrate:
	@mkdir -p $(MIGRATION_DIR)
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose sqlite3 $(DB_DIR) -dir=$(MIGRATION_DIR) create $(args) sql

live/server:
	$(SET_CGO) $(GOBIN)/air

live/frontend:
	cd frontend && npm install && npm run dev

dev:
	$(MAKE) -j2 live/server live/frontend

build:
	$(SET_CGO) go build -o ./build/chrono$(BIN_SUFFIX) ./cmd/main.go

install:
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

test:
	@go test ./... -v