%:
	@:

.PHONY: server client

ifneq (,$(wildcard ./.env))
    include .env
    export
endif


migrationDir = ./db/migrations/

dev:
	templ generate --watch --proxy="http://localhost:8080" --cmd="air"

generate:
	@echo "Generating sqlc repositoy..."
	@sqlc generate

migrate:
	@-mkdir ${migrationDir}
	$(eval args=$(filter-out $@,$(MAKECMDGOALS)))
	@goose -dir=${migrationDir} create ${args}

up:
	@goose -dir=${migrationDir} sqlite render_box.db up

down:
	@goose -dir=${migrationDir} sqlite render_box.db down
