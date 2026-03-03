# API Client Documentation

This directory contains the auto-generated API client and related utilities for communicating with the Envoy server.

## Files Overview

- `client.ts` - Main API client configured with openapi-fetch
- `api.ts` - Auto-generated TypeScript types from OpenAPI spec
- `custom-fetch.ts` - Custom fetch wrapper with error handling
- `README.md` - This documentation

## Setup

### 1. Generate API Types

First, make sure the server is running on `http://localhost:8080`, then generate the types:

```bash
npm run gen-api
```

This will fetch the OpenAPI spec from `http://localhost:8080/openapi` and generate TypeScript types in `api.ts`.

### 2. Usage Examples

#### Basic API Client Usage

```typescript
import client from './api/client';

// GET request
const response = await client.GET('/me');
console.log(response.data);

// POST request
const result = await client.POST('/auth/login', {
  body: {
    email: 'user@example.com',
    password: 'password123'
  }
});
```

#### Using the AuthService

```typescript
import { AuthService } from '../services/auth';

// Register a new user
await AuthService.register({
  email: 'user@example.com',
  password: 'password123'
});

// Login
await AuthService.login({
  email: 'user@example.com',
  password: 'password123'
});

// Get current user
const user = await AuthService.getMe();
```

#### Using the React Hook

```typescript
import { useAuth } from '../hooks/useAuth';

function MyComponent() {
  const { 
    user, 
    isAuthenticated, 
    login, 
    logout, 
    isLoggingIn 
  } = useAuth();

  const handleLogin = async () => {
    await login({
      email: 'user@example.com',
      password: 'password123'
    });
  };

  return (
    <div>
      {isAuthenticated ? (
        <div>
          <p>Welcome, {user?.data?.email}!</p>
          <button onClick={logout}>Logout</button>
        </div>
      ) : (
        <button onClick={handleLogin} disabled={isLoggingIn}>
          {isLoggingIn ? 'Logging in...' : 'Login'}
        </button>
      )}
    </div>
  );
}
```

## Available Endpoints

### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login user (sets HTTP-only cookie)
- `GET /me` - Get current authenticated user profile

### Response Format

All responses follow this format:

```typescript
{
  status: number;
  message: string;
  data?: T; // Response data
  success: boolean;
}
```

### Error Handling

The API client automatically handles HTTP errors and returns them in the `error` property:

```typescript
const response = await client.POST('/auth/login', { body: credentials });

if (response.error) {
  console.error('Login failed:', response.error.message);
  // Handle error (show message to user, etc.)
} else {
  // Success
  console.log('Login successful:', response.data);
}
```

## Cookie Authentication

The API uses HTTP-only cookies for authentication. The `client.ts` is configured with `credentials: 'include'` to automatically send cookies with each request.

## Type Safety

The generated types provide full type safety for:

- Request parameters and body
- Response data
- Error responses
- HTTP methods and status codes

## Development

When the server API changes, regenerate the types:

```bash
npm run gen-api
```

This ensures your client code always matches the server's API contract.
