#!/bin/bash
set -e

# Function to create database and users
setup_database() {
    local db_name=$1
    local readonly_db_name=$2
    local superuser=$3
    local prod_user="${db_name}_user"
    local readonly_user="${readonly_db_name}_user"
    
    echo "Setting up database: $db_name"
    
    psql -v ON_ERROR_STOP=1 --username "$superuser" --dbname "$db_name" <<-EOSQL
        -- IDEMPOTENT USER CREATION
        DO \$_$ 
        BEGIN
          IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$prod_user') THEN
            CREATE USER $prod_user WITH PASSWORD 'prod_user_password_123';
          END IF;

          IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$readonly_user') THEN
            CREATE USER $readonly_user WITH PASSWORD 'readonly_user_password_456';
          END IF;
        END
        \$_$;

        -- SCHEMAS & PERMISSIONS
        GRANT USAGE ON SCHEMA public TO $prod_user;
        GRANT USAGE ON SCHEMA public TO $readonly_user;
        
        GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO $prod_user;
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO $readonly_user;
        
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO $prod_user;
        ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO $readonly_user;

        -- SYSTEM TABLES
        CREATE TABLE IF NOT EXISTS environment_migrations (
            id VARCHAR(255) PRIMARY KEY,
            environment_id VARCHAR(255) NOT NULL,
            name VARCHAR(255) NOT NULL,
            status VARCHAR(50) DEFAULT 'pending',
            executed_at TIMESTAMP,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        GRANT SELECT, INSERT, UPDATE, DELETE ON environment_migrations TO $prod_user;
        GRANT SELECT ON environment_migrations TO $readonly_user;
EOSQL

    echo "✅ Setup for $db_name finished."
}

# --- MAIN EXECUTION ---

if [ ! -z "$POSTGRES_MULTIPLE_DATABASES" ]; then
    IFS=',' read -ra DBS <<< "$POSTGRES_MULTIPLE_DATABASES"
    
    # 1. Create extra databases
    for db in "${DBS[@]}"; do
        if [ "$db" != "$POSTGRES_DB" ]; then
            echo "Creating database: $db"
            psql --username "$POSTGRES_USER" --dbname "postgres" -c "CREATE DATABASE $db" || echo "Database $db already exists"
        fi
    done

    # 2. Run setup for each
    for db in "${DBS[@]}"; do
        readonly_db="${db}_readonly"
        setup_database "$db" "$readonly_db" "$POSTGRES_USER"
    done
else
    setup_database "$POSTGRES_DB" "${POSTGRES_DB}_readonly" "$POSTGRES_USER"
fi

echo "🚀 Envoy environment is ready."