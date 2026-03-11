BINARY_NAME ?= aicoder
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)

-include .env

%:
	@:

build:
	go build \
		-ldflags "-X main.GitCommit=$(COMMIT)" \
		-o dist/$(BINARY_NAME) \
		./cmd/$(BINARY_NAME)

build-%:
	@$(MAKE) build BINARY_NAME=$*

test:
	go test  ./... --race

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1 run

MIGRATE_COMMAND=go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

migrate:
	${MIGRATE_COMMAND} -path $(FIXIK_POSTGRES_MIGRATIONS_DIR) -database $(FIXIK_POSTGRES_URL) up

migrate-create:
	@if [ -z "$(word 2,$(MAKECMDGOALS))" ]; then echo "Name is required: make migrate-create <name>"; exit 1; fi
	migrate create -ext sql -seq -dir $(FIXIK_POSTGRES_MIGRATIONS_DIR) $(word 2,$(MAKECMDGOALS))

migrate-down:
	${MIGRATE_COMMAND} -path $(FIXIK_POSTGRES_MIGRATIONS_DIR) -database $(FIXIK_POSTGRES_URL) down 1

migrate-down-all:
	${MIGRATE_COMMAND} -path $(FIXIK_POSTGRES_MIGRATIONS_DIR) -database $(FIXIK_POSTGRES_URL) down

migrate-version:
	${MIGRATE_COMMAND} -path $(FIXIK_POSTGRES_MIGRATIONS_DIR) -database $(FIXIK_POSTGRES_URL) version

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.6 init  -d cmd/fixik,internal/router