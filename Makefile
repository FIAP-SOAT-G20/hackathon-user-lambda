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


.PHONY: install-deps
install-deps: ## 📦 Install dependencies
	@echo "🟢 Installing dependencies..."
	go mod download
	@go install github.com/blmayer/awslambdarpc@latest
	@echo

.PHONY: start-lambda
start-lambda:  build  ## ▶  Start the lambda application locally to prepare to receive requests
	@echo "🟢 Starting lambda ..."
	_LAMBDA_SERVER_PORT=3300 AWS_LAMBDA_RUNTIME_API=http://localhost:3300 go run ./cmd/api/main.go
	@echo

.PHONY: trigger-lambda
trigger-lambda: ## ⚡  Trigger lambda with the input file stored in variable ./event.json
	@echo "🟢 Triggering lambda with event: ./event.json"
	@PATH="$(shell go env GOPATH)/bin:$$PATH" \
		awslambdarpc -a localhost:3300 -e ./event.json
	@echo
