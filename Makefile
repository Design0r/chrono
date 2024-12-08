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
	air

dev: 
	make -j2 live/templ live/server 
