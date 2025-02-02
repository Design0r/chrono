%:
	@:

.PHONY: build dev

ifeq ($(OS),Windows_NT)
  BIN_SUFFIX := .exe
else
  BIN_SUFFIX :=
endif

MIGRATION_DIR = ./db/migrations/
DB_DIR = ./db/chrono.db/
CSS_DIR = ../assets/static/css

generate:
	@echo "Generating sqlc repositoy..."
	@sqlc generate

migrate:
	@-mkdir ${MIGRATION_DIR}
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose sqlite3 ${DB_DIR} -dir=${MIGRATION_DIR} create ${args} sql

live/templ:
	templ generate --watch --proxy="http://localhost:8080"  --open-browser=true

live/server:
	air

live/tailwind:
	cd node && \
	npm install && \
	npx --yes tailwindcss -i $(CSS_DIR)/input.css -o $(CSS_DIR)/output.css --minify --watch

dev: 
	make -j3 live/templ live/server live/tailwind

build:
	@cd node && \
	npm install && \
	npx --yes tailwindcss -i $(CSS_DIR)/input.css -o $(CSS_DIR)/output.css --minify

	templ generate
	go build -o ./build/Chrono$(BIN_SUFFIX) ./cmd/main.go

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
