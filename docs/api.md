# API Documentation

## Overview

This document provides detailed information about the API endpoints and how to use them.

## Endpoints

### Create a Client

**Endpoint:**
POST /client/create

**Description:**
Creates a new client with the provided details.

Request:
```json
{
  "name": "Client Name",
  "priority": 1,
  "leadCapacity": 100,
  "currentLeadCount": 0,
  "workingHoursStart": "09:00",
  "workingHoursEnd": "17:00"
}

Example:
curl -X POST http://localhost:8080/client/create -d '{
  "name": "Test Client",
  "priority": 1,
  "leadCapacity": 100,
  "currentLeadCount": 0,
  "workingHoursStart": "09:00",
  "workingHoursEnd": "17:00"
}' -H "Content-Type: application/json"



###  Assign a Lead

Endpoint:
GET /client/assign

Description:
Assigns a lead to an eligible client based on their working hours and lead capacity.

Example:
curl -X GET http://localhost:8080/client/assign


### Get Client By ID

Endpoint:
GET /client/{id}

Description:
Retrieves a specific client by ID from the database.

Example:
curl -X GET http://localhost:8080/client/1


## Usage Examples

### Create Multiple Clients

Client A:
curl -X POST http://localhost:8080/client/create -d '{
  "name": "Client A",
  "priority": 2,
  "leadCapacity": 50,
  "currentLeadCount": 10,
  "workingHoursStart": "08:00",
  "workingHoursEnd": "16:00"
}' -H "Content-Type: application/json"

Client B:
curl -X POST http://localhost:8080/client/create -d '{
  "name": "Client B",
  "priority": 3,
  "leadCapacity": 75,
  "currentLeadCount": 5,
  "workingHoursStart": "10:00",
  "workingHoursEnd": "18:00"
}' -H "Content-Type: application/json"

Client C:
curl -X POST http://localhost:8080/client/create -d '{
  "name": "Client C",
  "priority": 1,
  "leadCapacity": 200,
  "currentLeadCount": 50,
  "workingHoursStart": "07:00",
  "workingHoursEnd": "15:00"
}' -H "Content-Type: application/json"

### Assign Leads

Request:
curl -X GET http://localhost:8080/client/assign


###List All Clients

Request:
curl -X GET http://localhost:8080/client/all

### Get Client By ID

Request:
curl -X GET http://localhost:8080/client/1