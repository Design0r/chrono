%:
	@:

.PHONY: build 

ifneq (,$(wildcard ./.env))
    include .env
    export
endif


migrationDir = ./db/migrations/

generate:
	@echo "Generating sqlc repositoy..."
	@sqlc generate

migrate:
	@-mkdir ${migrationDir}
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose -dir=${migrationDir} create ${args}

up:
	@goose -dir=${migrationDir} sqlite chrono.db up

down:
	@goose -dir=${migrationDir} sqlite chrono.db down


live/templ:
	templ generate --watch --proxy="http://localhost:8080"  --open-browser=true

live/server:
	air

live/tailwind:
	npm install 
	npx --yes tailwindcss -i ./assets/static/css/input.css -o ./assets/static/css/output.css --minify --watch

dev: 
	make -j3 live/templ live/server live/tailwind

build:
	@npm install 
	@npx --yes tailwindcss -i ./assets/static/css/input.css -o ./assets/static/css/output.css --minify
	@templ generate
	@go build -o ./build/Chrono.exe ./cmd/main.go

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go mod tidy

test: 
	go test ./tests/ -v
