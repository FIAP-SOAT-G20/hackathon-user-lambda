SHELL := /bin/bash

BIN_DIR := dist

.PHONY: build clean fmt test coverage mock package

build:
	@echo "🔨 Building Lambda function..."
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/bootstrap ./cmd/api

package: build
	@echo "📦 Packaging Lambda function..."
	@cd $(BIN_DIR) && zip -r function.zip bootstrap

fmt:
	go fmt ./...

test:
	@echo "🟢 Running tests..."
	CGO_ENABLED=1 go test ./internal/... -race -cover || go test ./internal/... -cover

coverage:
	@echo "🟢 Running tests with coverage..."
	CGO_ENABLED=1 go test ./internal/... -race -cover -coverprofile=coverage.out || go test ./internal/... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

mock:
	@echo "🟢 Generating mocks..."
	@mkdir -p internal/core/port/mocks
	@for file in internal/core/port/*.go; do \
		go run go.uber.org/mock/mockgen@latest -source=$$file -destination=internal/core/port/mocks/`basename $$file .go`_mock.go; \
	done

clean:
	rm -rf $(BIN_DIR)
