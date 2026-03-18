CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS environments (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    project_id TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'development' CHECK (
        type IN (
            'production',
            'staging',
            'development'
        )
    ),
    connection_string_encrypted TEXT NOT NULL,
    connection_error TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS environment_db_users (
    id TEXT PRIMARY KEY,
    environment_id TEXT NOT NULL,
    name TEXT NOT NULL,
    connection_string_encrypted TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (environment_id) REFERENCES environments (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS environment_migrations (
    id TEXT PRIMARY KEY,
    environment_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    sql_content TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (
        status IN (
            'pending',
            'running',
            'completed',
            'failed'
        )
    ),
    executed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    duration INTEGER,
    executed_by TEXT,
    error_message TEXT,
    checksum TEXT NOT NULL,
    client_id TEXT UNIQUE,
    FOREIGN KEY (environment_id) REFERENCES environments (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_env_migrations_env_id ON environment_migrations (environment_id);

CREATE INDEX IF NOT EXISTS idx_env_migrations_status ON environment_migrations (status);

CREATE INDEX IF NOT EXISTS idx_db_users_env_id ON environment_db_users (environment_id);