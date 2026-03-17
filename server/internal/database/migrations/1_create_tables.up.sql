CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS environments (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    project_id TEXT NOT NULL,
    connection_string_encrypted TEXT NOT NULL,
    connection_error TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    type TEXT NOT NULL DEFAULT 'development' CHECK (type IN ('production', 'staging', 'development')),
    FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
);