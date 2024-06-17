# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="MullerJulia mrjuliamelnik@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd

# Install wget
RUN apt-get update && apt-get install -y wget

# Download and install migrate
RUN wget -q -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz && \
    tar -xzf migrate.tar.gz && \
    ls -la && \
    mv migrate /usr/local/bin/migrate && \
    rm migrate.tar.gz