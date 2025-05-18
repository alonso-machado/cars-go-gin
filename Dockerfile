# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod first
COPY go.mod .

# Download all dependencies and create go.sum
RUN go mod download
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o car-service .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/car-service .
COPY --from=builder /app/migrations ./migrations

# Copy environment file (you can override this in docker-compose)
COPY .env .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./car-service"]
