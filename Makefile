%:
	@:

.PHONY: server client

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
	@goose -dir=${migrationDir} sqlite calendar.db up

down:
	@goose -dir=${migrationDir} sqlite calendar.db down


live/templ:
	templ generate --watch --proxy="http://localhost:8080"  --open-browser=false -v

live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o tmp/bin/main ./cmd/main.go" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

dev: 
	make -j2 live/templ live/server 
