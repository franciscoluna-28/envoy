# Database Credentials

This document contains all the credentials for the PostgreSQL database setup.

## Database Connection Information

- **Host:** localhost (development only)
- **Port:** 5432 (development only)
- **Database:** envoy_dev
- **Network:** envoy-network (Docker internal)

## User Credentials

### 1. Superuser Account
- **Username:** `envoy_superuser`
- **Password:** `superuser_password_123`
- **Permissions:** 
  - Full administrative access
  - Can create/drop databases, users, tables
  - Can modify schema
  - Used for: Initial setup, migrations, administrative tasks
- **Connection String (Development):** 
  ```
  postgresql://envoy_superuser:superuser_password_123@localhost:5432/envoy_dev
  ```
- **Connection String (Docker Internal):**
  ```
  postgresql://envoy_superuser:superuser_password_123@postgres-dev:5432/envoy_dev
  ```
- **PSQL Command (Development):**
  ```bash
  psql -h localhost -p 5432 -U envoy_superuser -d envoy_dev
  ```
- **PSQL Command (Docker Internal):**
  ```bash
  docker-compose -f docker-compose.dev.yml exec postgres-dev psql -U envoy_superuser -d envoy_dev
  ```

### 2. Production User (Application User)
- **Username:** `envoy_dev_user`
- **Password:** `prod_user_password_123`
- **Permissions:**
  - SELECT, INSERT, UPDATE, DELETE on tables
  - Cannot create/drop tables or modify schema
  - Cannot create new users or databases
  - Used for: Regular application operations
- **Connection String (Development):**
  ```
  postgresql://envoy_dev_user:prod_user_password_123@localhost:5432/envoy_dev
  ```
- **Connection String (Docker Internal):**
  ```
  postgresql://envoy_dev_user:prod_user_password_123@postgres-dev:5432/envoy_dev
  ```
- **PSQL Command (Development):**
  ```bash
  psql -h localhost -p 5432 -U envoy_dev_user -d envoy_dev
  ```
- **PSQL Command (Docker Internal):**
  ```bash
  docker-compose -f docker-compose.dev.yml exec postgres-dev psql -U envoy_dev_user -d envoy_dev
  ```

### 3. Read-only User
- **Username:** `envoy_dev_readonly_user`
- **Password:** `readonly_user_password_456`
- **Permissions:**
  - SELECT only on tables
  - Cannot modify data or schema
  - Used for: Reporting, monitoring, analytics
- **Connection String (Development):**
  ```
  postgresql://envoy_dev_readonly_user:readonly_user_password_456@localhost:5432/envoy_dev
  ```
- **Connection String (Docker Internal):**
  ```
  postgresql://envoy_dev_readonly_user:readonly_user_password_456@postgres-dev:5432/envoy_dev
  ```
- **PSQL Command (Development):**
  ```bash
  psql -h localhost -p 5432 -U envoy_dev_readonly_user -d envoy_dev
  ```
- **PSQL Command (Docker Internal):**
  ```bash
  docker-compose -f docker-compose.dev.yml exec postgres-dev psql -U envoy_dev_readonly_user -d envoy_dev
  ```

## PGAdmin Credentials

- **URL:** http://localhost:5050 (development only)
- **Email:** `admin@envoy.dev`
- **Password:** `admin_password_789`

## Connection String Summary

### For Development (External Access)
```bash
# Superuser (for migrations/admin)
postgresql://envoy_superuser:superuser_password_123@localhost:5432/envoy_dev

# Production User (for application)
postgresql://envoy_dev_user:prod_user_password_123@localhost:5432/envoy_dev

# Read-only User (for reporting)
postgresql://envoy_dev_readonly_user:readonly_user_password_456@localhost:5432/envoy_dev
```

### For Docker Internal (Container-to-Container)
```bash
# Superuser (for migrations/admin)
postgresql://envoy_superuser:superuser_password_123@postgres-dev:5432/envoy_dev

# Production User (for application)
postgresql://envoy_dev_user:prod_user_password_123@postgres-dev:5432/envoy_dev

# Read-only User (for reporting)
postgresql://envoy_dev_readonly_user:readonly_user_password_456@postgres-dev:5432/envoy_dev
```

## Usage Recommendations

### For Envoy Application (Development)
Use the **Production User** credentials with localhost access:
```
Database URL: postgresql://envoy_dev_user:prod_user_password_123@localhost:5432/envoy_dev
```

### For Envoy Application (Docker Compose)
Use the **Production User** credentials with service name:
```
Database URL: postgresql://envoy_dev_user:prod_user_password_123@postgres-dev:5432/envoy_dev
```

### For Database Migrations
Use the **Superuser** credentials:
```
Database URL: postgresql://envoy_superuser:superuser_password_123@localhost:5432/envoy_dev
```

### For Reporting/Analytics
Use the **Read-only User** credentials:
```
Database URL: postgresql://envoy_dev_readonly_user:readonly_user_password_456@localhost:5432/envoy_dev
```

## Network Configuration

### Development Access (Ports Open)
- **PostgreSQL:** `localhost:5432` → `postgres-dev:5432`
- **PGAdmin:** `localhost:5050` → `pgadmin:80`

### Production Considerations (Ports Closed)
In production environments:
- **Close external ports** - Only allow internal Docker network access
- **Use service names** - `postgres-dev:5432` instead of `localhost:5432`
- **Implement firewall rules** - Restrict database access to application containers only

### Docker Network Configuration
```yaml
networks:
  envoy-network:
    driver: bridge
    internal: false  # Set to true in production for isolation
```

## Security Notes

⚠️ **Important:** These credentials are for development use only. In production:

1. **Change all passwords** before deploying to production
2. **Use environment variables** to store credentials securely
3. **Use SSL/TLS** for database connections
4. **Consider using secrets management** tools like HashiCorp Vault or AWS Secrets Manager
5. **Implement IP whitelisting** for database access
6. **Use stronger passwords** with a mix of uppercase, lowercase, numbers, and symbols
7. **Close external ports** - Use internal Docker network only
8. **Enable network isolation** - Set `internal: true` for production networks

## Testing Credentials

You can test the credentials using the provided PowerShell script:

```bash
.\test-databases.ps1
```

Or manually test with psql:

```bash
# Test superuser (development)
psql -h localhost -p 5432 -U envoy_superuser -d envoy_dev -c "SELECT current_user;"

# Test production user (development)
psql -h localhost -p 5432 -U envoy_dev_user -d envoy_dev -c "SELECT current_user;"

# Test read-only user (development)
psql -h localhost -p 5432 -U envoy_dev_readonly_user -d envoy_dev -c "SELECT current_user;"

# Test superuser (Docker internal)
docker-compose -f docker-compose.dev.yml exec postgres-dev psql -U envoy_superuser -c "SELECT current_user;"
```

## Environment Variables

Add these to your `.env` file:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=envoy_dev
DB_SUPERUSER=envoy_superuser
DB_SUPERUSER_PASSWORD=superuser_password_123
DB_PROD_USER=envoy_dev_user
DB_PROD_USER_PASSWORD=prod_user_password_123
DB_READONLY_USER=envoy_dev_readonly_user
DB_READONLY_USER_PASSWORD=readonly_user_password_456

# Docker Internal Configuration
DB_INTERNAL_HOST=postgres-dev
DB_INTERNAL_PORT=5432

# PGAdmin
PGADMIN_DEFAULT_EMAIL=admin@envoy.dev
PGADMIN_DEFAULT_PASSWORD=admin_password_789

# Connection Strings
DB_SUPERUSER_URL=postgresql://envoy_superuser:superuser_password_123@localhost:5432/envoy_dev
DB_PROD_USER_URL=postgresql://envoy_dev_user:prod_user_password_123@localhost:5432/envoy_dev
DB_READONLY_USER_URL=postgresql://envoy_dev_readonly_user:readonly_user_password_456@localhost:5432/envoy_dev

# Docker Internal Connection Strings
DB_SUPERUSER_INTERNAL_URL=postgresql://envoy_superuser:superuser_password_123@postgres-dev:5432/envoy_dev
DB_PROD_USER_INTERNAL_URL=postgresql://envoy_dev_user:prod_user_password_123@postgres-dev:5432/envoy_dev
DB_READONLY_USER_INTERNAL_URL=postgresql://envoy_dev_readonly_user:readonly_user_password_456@postgres-dev:5432/envoy_dev
```

## Quick Reference

| User | Username | Password | Dev Connection | Docker Connection | Use Case |
|------|----------|----------|----------------|-------------------|----------|
| Superuser | envoy_superuser | superuser_password_123 | `localhost:5432` | `postgres-dev:5432` | Migrations, Admin |
| Production | envoy_dev_user | prod_user_password_123 | `localhost:5432` | `postgres-dev:5432` | Application |
| Read-only | envoy_dev_readonly_user | readonly_user_password_456 | `localhost:5432` | `postgres-dev:5432` | Reporting |
| PGAdmin | admin@envoy.dev | admin_password_789 | `localhost:5050` | N/A | Web UI |
