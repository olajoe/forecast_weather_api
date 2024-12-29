run:
	@echo "Starting server..."
	go run cmd/main.go

test:
	@echo "Running unit tests"
	go test ./...

test-coverage:
	@echo "Generating test coverage..."
	go test ./... -coverprofile=coverage.out
	@echo "Open test coverage"
	go tool cover -html coverage.out

lint:
	@echo "Running lint..."
	golangci-lint run -v --timeout 5m0s --max-issues-per-linter 10 ./...