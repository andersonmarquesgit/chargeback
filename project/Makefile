API_BINARY=chargeback-api
BATCH_BINARY=chargeback-batch
PROCESSOR_BINARY=chargeback-processor

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker containers..."
	docker compose up -d
	@echo "Docker containers started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_api build_batch build_processor
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
	@echo "Ensuring swag is up-to-date..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger docs..."
	cd ../chargeback-api && swag init -g internal/interfaces/http/routes/routes.go
	@echo "Building Go binary..."
	cd ../chargeback-api && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${API_BINARY} ./cmd
	@echo "Build complete!"

build_batch:
	@echo "Building Batch binary..."
	cd ../chargeback-batch && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${BATCH_BINARY} ./cmd
	@echo "Batch build complete!"

build_processor:
	@echo "Building Processor binary..."
	cd ../chargeback-processor && env GOOS=linux CGO_ENABLED=0 go build -buildvcs=false -o ${PROCESSOR_BINARY} ./cmd
	@echo "Processor build complete!"
