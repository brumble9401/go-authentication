CREATE TABLE IF NOT EXISTS roles (
    role_id UUID,
    role_name TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (role_id, created_at)
);

CREATE TABLE IF NOT EXISTS permissions (
    permission_id UUID,
    permission_name TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (permission_id, created_at)
);

CREATE TABLE IF NOT EXISTS users (
    user_id UUID,
    username TEXT,
    password_hash TEXT,
    email TEXT,
    full_name TEXT,
    role_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY ((user_id, username), created_at)
);

