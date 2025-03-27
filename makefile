# Go project variables
BINARY_NAME=tax-server
MAIN_PACKAGE=./cmd/server
COVER_FILE=coverage.out
HTML_COVER_FILE=coverage.html

# Run the Go server locally
run:
	go run $(MAIN_PACKAGE)

# Run all unit tests with coverage
test:
	go test ./... -cover

# Generate test coverage report
coverage:
	go test ./... -coverprofile=$(COVER_FILE)
	go tool cover -html=$(COVER_FILE) -o $(HTML_COVER_FILE)

# Run the full stack using Docker Compose
up:
	docker-compose up --build

# Build the Go binary
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Clean up coverage files and binary
clean:
	rm -f $(BINARY_NAME) $(COVER_FILE) $(HTML_COVER_FILE)

# Format all Go files
# Run basic style and static checks
lint:
	go fmt ./...
	go vet ./...
	golint ./...
