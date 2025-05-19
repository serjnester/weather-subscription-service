.PHONY:

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

BINARY_NAME=weather-app
TOOLS_DIR=$(shell pwd)/.environment
POSTGRES_URI=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable
RUN_GOOSE=$(TOOLS_DIR)/goose -dir ./internal/storage/migrations postgres "$(POSTGRES_URI)"

install-tools:
	GOBIN=${TOOLS_DIR} go install "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0" \
	&& GOBIN=${TOOLS_DIR} go install "github.com/pressly/goose/v3/cmd/goose@v3.24.3" \
	&& GOBIN=${TOOLS_DIR} go install "github.com/swaggo/swag/cmd/swag@v1.16.4"

test:
	@echo "Running tests..."
	go test ./...


up:
	docker-compose up -d --build

down:
	docker-compose down

create-migration:
	@read -p "Enter migration name: " name; \
	$(RUN_GOOSE) create $$name sql

migrate-up:
	docker-compose run --rm migrate

sqlc-gen:
	$(TOOLS_DIR)/sqlc generate

swag-init: swag-fmt generate-swagger

generate-swagger:
	$(TOOLS_DIR)/swag init --parseDependency -g router.go -d internal/handlers -o docs

swag-fmt:
	$(TOOLS_DIR)/swag fmt ./internal/handlers