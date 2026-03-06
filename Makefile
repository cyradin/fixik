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
