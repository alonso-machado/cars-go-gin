.PHONY: build run test clean migrate-up migrate-down docker-build docker-up docker-down docker-logs

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=car-service
MIGRATE=migrate -path ./migrations -database "postgres://john:doe@localhost:5432/car_service?sslmode=disable"

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Database migrations
migrate-up:
	$(MIGRATE) -verbose up

migrate-down:
	$(MIGRATE) -verbose down

# Docker commands
docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Run the application
run: build
	./$(BINARY_NAME)

# Install dependencies
deps:
	$(GOCMD) mod download
	$(GOCMD) mod verify

# Install development tools
tools:
	# Install migrate
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	# Install swag
	go install github.com/swaggo/swag/cmd/swag@latest
	# Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2

# Generate Swagger docs
swag:
	swag init -g internal/api/car_handler.go -o docs/swagger

# Lint the code
lint:
	golangci-lint run
