# Car Service API

A RESTful API for managing cars, built with Go, Gin, and PostgreSQL.

## Features

- CRUD operations for cars
- Filter cars by brand, price range, and name
- Pagination support
- Request validation
- Structured logging
- Swagger documentation
- Docker support
- Database migrations

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker and Docker Compose (optional)

## Getting Started

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/username/go-car-service.git
   cd go-car-service
   ```

2. Start the application with Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`
4. Access Swagger UI at `http://localhost:8080/swagger/index.html`

### Manual Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/username/go-car-service.git
   cd go-car-service
   ```

2. Create a `.env` file in the root directory with the following variables:
   ```env
   # Server
   SERVER_PORT=8080
   ENVIRONMENT=development

   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=john
   DB_PASSWORD=doe
   DB_NAME=car_service
   DB_SSLMODE=disable
   ```

3. Start PostgreSQL and create a database named `car_service`

4. Install dependencies and run migrations:
   ```bash
   # Install migrate tool
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

   # Run migrations
   make migrate-up
   ```

5. Run the application:
   ```bash
   go run main.go
   ```

6. The API will be available at `http://localhost:8080`

## API Documentation

Once the application is running, you can access the Swagger UI at `http://localhost:8080/swagger/index.html` for interactive API documentation.

## API Endpoints

### Cars

- `GET /api/v1/cars` - Get all cars (with pagination)
- `GET /api/v1/cars/:id` - Get a car by ID
- `GET /api/v1/cars/name/:name` - Get a car by name
- `GET /api/v1/cars/brand/:brand` - Get cars by brand
- `GET /api/v1/cars/price-range?startPrice=X&finalPrice=Y` - Get cars by price range
- `POST /api/v1/cars` - Create a new car
- `PUT /api/v1/cars/:id` - Update a car
- `DELETE /api/v1/cars/:id` - Delete a car

## Development

### Running Tests

```bash
go test -v ./...
```

### Linting

```bash
# Install golangci-lint if not installed
# https://golangci-lint.run/usage/install/

golangci-lint run
```

### Building for Production

```bash
# Build the application
go build -o bin/car-service

# Run the binary
./bin/car-service
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port the server will listen on | `8080` |
| `ENVIRONMENT` | Application environment (development, production) | `development` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `john` |
| `DB_PASSWORD` | Database password | `doe` |
| `DB_NAME` | Database name | `car_service` |
| `DB_SSLMODE` | Database SSL mode | `disable` |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
