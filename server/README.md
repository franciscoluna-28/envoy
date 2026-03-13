# NewServer Project

A newborn Go web server project using Chi router.

## Getting Started

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run .
```

3. Visit the endpoints:
- `http://localhost:8080/` - Welcome message
- `http://localhost:8080/health` - Health check

## Project Structure

```
newServer/
├── main.go          # Main application entry point
├── go.mod           # Go module file
└── README.md        # This file
```

## Features

- Chi router for HTTP routing
- Basic middleware (logging, recovery, CORS)
- Health check endpoint
- Simple welcome endpoint
