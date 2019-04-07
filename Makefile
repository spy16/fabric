all: clean build tidy test lint

clean:
	@echo "Cleaning up..."
	@rm -rf ./bin/

test:
	@echo "Running tests..."
	@go test -cover ./...

lint:
	@echo "Running linter..."
	@golint ./...

build:
	@echo "Building..."
	@mkdir ./bin/
	@go build -o bin/fabric ./cmd/fabric/*.go

tidy:
	@echo "Tidy up go mod files..."
	@go mod tidy -v
