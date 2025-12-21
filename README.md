# Casdoor Test API

A Go-based REST API service that integrates with Casdoor for authentication and role-based access control (RBAC).

## Prerequisites

- Docker and Docker Compose
- Go 1.24 or later
- Make

## Getting Started

### 1. Start Dependencies with Docker Compose

Start the PostgreSQL database and Casdoor authentication service:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Casdoor web interface on http://localhost:8000

### 2. Run Database Migrations

Apply migrations to set up casdoor:

```bash
make migrate-up
```

### 3. Start the API Service

Run the Go application server:

```bash
go run main.go
```

The API will be available at http://localhost:9000

## API Endpoints

### Public Routes
- `GET /login` - Get Casdoor login URL
- `GET /callback` - Handle OAuth callback from Casdoor

### Protected Routes (Require Authentication)
All routes under `/api` require valid Casdoor authentication and appropriate permissions:

#### User Management
- `GET /api/me` - Get current user information
- `GET /api/users` - List all users (requires permission)
- `POST /api/users` - Add new user (requires permission)
- `PUT /api/users/:username` - Update user (requires permission)
- `DELETE /api/users/:username` - Delete user (requires permission)
