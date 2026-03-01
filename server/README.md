# Envoy Server

A Go REST API server with authentication and user management using Clean Architecture principles.

## Development

### Prerequisites

- Go 1.25.0 or later
- SQLite

### Setup

1. Clone the repository
2. Navigate to the server directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Setup development environment:
   ```bash
   make setup
   ```

### Running the Application

#### Development with Hot Reload

The easiest way to develop is using the hot reload feature:

```bash
make dev
```

This will start the server with automatic reloading when you save changes.

#### Manual Development

If you prefer to run without hot reload:

```bash
make run
```

#### Build and Run

To build the application and run the binary:

```bash
make build
./bin/main.exe
```

### Available Commands

- `make help` - Show all available commands
- `make dev` - Run with hot reload (recommended for development)
- `make run` - Run the application without hot reload
- `make build` - Build the application binary
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make tidy` - Tidy Go modules
- `make setup` - Setup development environment (installs air, tidies modules)
- `make setup-simple` - Simple setup without air (if air installation fails)

### Troubleshooting

If `make setup` fails during air installation:

1. **Try the simple setup**:
   ```bash
   make setup-simple
   make run
   ```

2. **Install air manually**:
   ```bash
   go install github.com/air-verse/air@latest
   ```

3. **Check your GOPATH**:
   ```bash
   echo $GOPATH
   # Make sure $GOPATH/bin is in your PATH
   ```

4. **Alternative: Run without hot reload**:
   ```bash
   make run
   ```

5. **SQLite/CGO Issues**: If you get "Binary was compiled with 'CGO_ENABLED=0'" error:
   - Make sure you have a C compiler installed (GCC/MinGW on Windows)
   - The air configuration automatically enables CGO, but if building manually:
     ```bash
     CGO_ENABLED=1 go build -o ./bin/main.exe .
     ```

### Environment Variables

Create a `.env` file in the root directory:

```env
DATABASE_URL=./envoy.db
JWT_SECRET=your-secret-key-here
PORT=8080
```

## Project Structure

```
server/
├── internal/api/
│   ├── domain/         # Domain entities and interfaces
│   ├── handlers/       # HTTP handlers
│   ├── infra/          # Infrastructure services (auth, web)
│   ├── repositories/   # Data access layer
│   ├── routes/         # Route definitions
│   └── usecases/       # Business logic
├── pkg/
│   └── database/       # Database utilities
├── main.go             # Application entry point
├── go.mod              # Go module file
├── .air.toml           # Air hot reload configuration
├── Makefile           # Development commands
└── envoy.db           # SQLite database (created automatically)
```

## Architecture Layers

1. **Handlers** - HTTP request/response handling
2. **Use Cases** - Business logic and use cases
3. **Repositories** - Data access abstraction
4. **Domain** - Domain entities and interfaces

## API Endpoints

- `GET /health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Example Usage

Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password123"}'
```

Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "john@example.com", "password": "password123"}'
```

Health check:
```bash
curl http://localhost:8080/health
```
