# Use the official Golang image to create a build artifact.
FROM golang:1.22-alpine AS builder

# Enable CGO and install necessary dependencies
RUN apk add --no-cache gcc musl-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd

# Install the migrate tool with SQLite support
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xz && \
    mv migrate.linux-amd64 /usr/local/bin/migrate

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

# Copy the migrations folder
COPY --from=builder /app/db/migrations ./db/migrations

# Install SQLite3
RUN apk --no-cache add sqlite

# Command to run the executable
CMD ["./main"]
