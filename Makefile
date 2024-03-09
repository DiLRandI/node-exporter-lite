.PHONY: build test

APP_NAME=node-exporter-lite

build:
	@echo "Building the project"
	@go build -o bin/$(APP_NAME) cmd/main.go

test:
	@echo "Running tests"
	@go test -v -cover ./...