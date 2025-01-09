%:
	@:

.PHONY: build 

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# --- OS detection logic ---
UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Darwin)
  OS_NAME := macos
endif

ifneq (,$(findstring MINGW,$(UNAME_S)))
  OS_NAME := windows
endif

ifeq ($(UNAME_S),Windows_NT)
  OS_NAME := windows
endif

# If OS_NAME is still empty after the checks, assume Linux
ifeq ($(OS_NAME),)
  OS_NAME := linux
endif

# --- Suffix for the binary ---
ifeq ($(OS_NAME),windows)
  BIN_SUFFIX := .exe
else
  BIN_SUFFIX :=
endif

migrationDir = ./db/migrations/
cssDir = ./assets/static/css

generate:
	@echo "Generating sqlc repositoy..."
	@sqlc generate

migrate:
	@-mkdir ${migrationDir}
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose -dir=${migrationDir} create ${args}

live/templ:
	templ generate --watch --proxy="http://localhost:8080"  --open-browser=true

live/server:
	air

live/tailwind:
	npm install 
	npx --yes tailwindcss -i $(cssDir)/input.css -o $(cssDir)/output.css --minify --watch

dev: 
	make -j3 live/templ live/server live/tailwind

build:
	@npm install 
	@npx --yes tailwindcss -i $(cssDir)/input.css -o $(cssDir)/output.css --minify
	@templ generate
	@go build -o ./build/Chrono$(BIN_SUFFIX) ./cmd/main.go

docker-install:
	@go install github.com/a-h/templ/cmd/templ@latest

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

test: 
	go test ./tests/ -v

deploy:
	@git pull origin main
	@docker compose up --build -d
