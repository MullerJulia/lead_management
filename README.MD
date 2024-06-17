# Lead Management API

## Overview

This project is a Lead Management API built with Go, SQLite, and Docker. The API allows you to create clients, assign leads, and manage client information.

## Prerequisites

- Docker
- Docker-Compose

## Getting Started

### 1. Clone the Repository

```sh
git clone https://github.com/MullerJulia/lead_management.git
cd lead_management

### 2. Build and Run the Docker Containers
```sh
docker-compose up --build
```
This command will build and run the Docker containers required for the Lead Management API.

### 3. Access the API
Once the containers are up and running, you can access the API at `http://localhost:8080`.

### 4. API Documentation
For detailed information on the API endpoints and how to use them, refer to the API documentation (docs/api.md)

### 5. Testing
To run the tests for the Lead Management API, use the following command:
go test ./...

