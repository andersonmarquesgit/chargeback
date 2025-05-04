API_BINARY=chargeback-api

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker containers..."
	docker compose up -d
	@echo "Docker containers started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_api
	@echo "Stopping any existing containers..."
	docker compose down
	@echo "Building and starting containers..."
	docker compose up --build -d
	@echo "All containers up and running!"

## down: stop docker compose
down:
	@echo "Stopping containers..."
	docker compose down
	@echo "Stopped."

## build_api: builds the API binary as a linux executable
build_api:
	@echo "Running go mod tidy..."
	cd ../chargeback-api && go mod tidy
	@echo "Generating Swagger docs..."
	cd ../chargeback-api && swag init -g internal/interfaces/http/routes/routes.go
	@echo "Building Go binary..."
	cd ../chargeback-api && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${API_BINARY} ./cmd/api
	@echo "Build complete!"
