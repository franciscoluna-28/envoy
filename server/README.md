# Go Clean Architecture API

A basic REST API built with Go using Clean Architecture principles.

## Project Structure

```
server/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration
│   ├── handlers/         # HTTP handlers (controllers)
│   ├── models/           # Domain models
│   ├── repositories/     # Data access layer
│   └── services/         # Business logic
├── pkg/                   # Public library code
│   └── database/         # Database setup
├── go.mod
├── go.sum
└── .env
```

## Architecture Layers

1. **Handlers** - HTTP request/response handling
2. **Services** - Business logic and use cases
3. **Repositories** - Data access abstraction
4. **Models** - Domain entities and DTOs

## Getting Started

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

- `GET /health` - Health check
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## Example Usage

Create a user:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

Get all users:
```bash
curl http://localhost:8080/api/v1/users
```
