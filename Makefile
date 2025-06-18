%:
	@:

.PHONY: build dev test

ifeq ($(OS),Windows_NT)
  BIN_SUFFIX := .exe
else
  BIN_SUFFIX :=
endif

MIGRATION_DIR = ./db/migrations/
DB_DIR = ./db/chrono.db/
TIMESTAMP := $(date +%Y-%m-%d_%H-%M-%S)

generate:
	@echo "Generating sqlc repositoy..."
	@sqlc generate

migrate:
	@-mkdir ${MIGRATION_DIR}
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose sqlite3 ${DB_DIR} -dir=${MIGRATION_DIR} create ${args} sql

live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=true

live/server:
	air

live/tailwind:
	npm install && \
	npm run dev

dev: 
	make -j3 live/templ live/server live/tailwind

build:
	templ generate
	go build -o ./build/chrono$(BIN_SUFFIX) -ldflags='-s -w -extldflags "-static"' ./cmd/main.go

docker-install:
	@go install github.com/a-h/templ/cmd/templ@latest

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

backup:
	@docker compose down
	@mkdir -p /home/apic/backup
	@bash -c 'timestamp=$$(date +%Y-%m-%d_%H-%M-%S); \
		echo "Backing up to chrono_$$timestamp.db"; \
		sudo cp /var/lib/docker/volumes/chrono_db/_data/chrono.db /home/apic/backup/chrono_$$timestamp.db'
	@COMPOSE_BAKE=true docker compose up -d

deploy: 
	@backup
	@git pull origin main
	@docker compose up --build -d

test:
	@go test ./... -v
