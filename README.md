# FIPE APP API

A small Golang API that fetches vehicle information (FIPE) by plate, manages users and their search history, and stores data in PostgreSQL. The app uses Gin for HTTP routing, GORM for database access and migrations, and includes JWT-based authentication.

## Table of contents

- Overview
- Features
- Tech stack
- Prerequisites
- Environment variables (.env)
- Run (local)
- Run with Docker
- API Endpoints
- Project structure
- Migrations & database
- Contributing
- License

## Overview

This service exposes endpoints to:

- Create and manage users
- Authenticate users and issue JWT tokens
- Query the FIPE external API for vehicle details by plate and return a filtered response
- Store user search history in the database

Migrations are applied automatically on startup.

## Features

- JWT authentication (middleware)
- GORM models + automatic AutoMigrate for User and History
- FIPE integration via configurable external API URL and token
- Simple health endpoint

## Tech stack

- Go 1.24
- Gin (HTTP router)
- GORM (ORM)
- PostgreSQL (database)
- Docker & docker-compose

## Prerequisites

- Go 1.24+ (if running locally)
- Docker & Docker Compose (if using containers)
- A PostgreSQL instance (the provided docker-compose starts one)

## Environment variables (.env)

Create a `.env` file in the project root (the application uses `github.com/joho/godotenv` to load it). Example:

```env
# HTTP port the server listens on
PORT=8080

# Postgres connection string
# Example: postgres://user:pass@localhost:5432/dbname?sslmode=disable
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable

# FIPE external API configuration (optional)
FIPE_EXTERNAL_API_URL=https://external-fipe.example/api
FIPE_API_TOKEN=your_fipe_api_token

# Secret key used to sign JWT tokens
SECRET_KEY=replace-with-a-secure-secret

# Gin mode: release|debug|test (defaults to release)
GIN_MODE=release
```

Notes:

- If no `.env` file is present, the app will try to read environment variables from the system (see `internal/config/config.go`).

## Run (local)

1. Ensure `DATABASE_URL` points to a running Postgres instance.
2. Set environment variables (or create a `.env` file as above).
3. From project root:

```bash
# run directly
go run ./cmd/api

# or build and run
go build -o server ./cmd/api
./server
```

The server listens on the port defined by `PORT` (default `8080`). A health check is available at `GET /healthz`.

## Run with Docker

The repository includes a `docker-compose.yml` which starts Postgres and the API (the API image is built from the `Dockerfile`).

Start everything:

```bash
docker compose up --build
```

The API will be available on the host port configured in the `.env` or `docker-compose.yml` mapping (defaults to `8080`).

## API Endpoints

All API routes are prefixed by `/api`.

- GET /healthz

  - Simple health check (no auth)

- POST /api/login

  - Body: { "email": "user@example.com", "password": "plaintext" }
  - Response: JWT token (string) or token payload
  - Public

- POST /api/users

  - Create a user
  - Body: { "name": "Name", "email": "email@example.com", "password": "secret" }
  - Response: JWT token (string)
  - Public

- GET /api/me

  - Get profile for the currently authenticated user
  - Auth: Authorization: Bearer <token>

- GET /api/users/:userID

  - Get user by id (requires auth)

- PUT /api/users/:userID

  - Update user name/email (requires auth and same user)

- PUT /api/users/:userID/update-password

  - Update password (requires auth and same user)
  - Body example: { "current": "oldpass", "new": "newpass" }

- DELETE /api/users/:userID

  - Delete user (requires auth and same user)

- GET /api/users/:userID/get-history

  - Get search history for user (requires auth and same user)

- POST /api/fipe
  - Body: { "plate": "AAA1234" }
  - Auth required
  - Returns filtered vehicle data (see `internal/models/vehicleDataFiltered.go` for the JSON shape)

Notes about authentication:

- Endpoints marked as requiring auth use the `Authorization: Bearer <token>` header. Tokens are created on login or user creation.

## Request & response shapes (examples)

- FIPE request:

```json
{ "plate": "ABC1D23" }
```

- Vehicle response (abridged):

```json
{
  "brand": "FIAT",
  "model": "Uno",
  "year": "2010",
  "model_year": "2010",
  "plate": "ABC1234",
  "fipe": [
    {
      "brand": "FIAT",
      "model": "Uno",
      "year_model": "2010",
      "reference_month": "06/2020",
      "fuel": "Flex",
      "value": "R$ 20.000,00"
    }
  ]
}
```

## Project structure

Key folders and files:

- `cmd/api` - application entrypoint
- `internal/config` - configuration and DB connection
- `internal/http` - router setup
- `internal/controllers` - HTTP handlers
- `internal/middlewares` - auth middleware
- `internal/models` - GORM models and request/response structs
- `internal/repositories` - DB access for models
- `internal/services` - external API interactions (FIPE)
- `internal/migrations` - database migrations

## Migrations & database

On startup the app:

1. Loads configuration (env variables / .env)
2. Connects to the database using the `DATABASE_URL` connection string
3. Applies migrations via the internal migrations package (see `internal/migrations`)
4. Calls `AutoMigrate` for `models.User` and `models.History`

If you need to run migrations separately or manage them, inspect `internal/migrations` which uses `gorm`/`gormigrate`.

## Contributing

1. Fork the repo and open a branch for your change.
2. Follow existing code style; run `gofmt`.
3. Add tests where relevant.
4. Open a pull request describing the change.

## Troubleshooting

- If the app can't connect to the DB, confirm `DATABASE_URL` and that Postgres is reachable.
- Check logs for migration errors printed on startup.

## License

This project does not include a license file in the repository. Add a `LICENSE` file if you want to make the project open-source and choose a license.

---
