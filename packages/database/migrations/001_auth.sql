CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    api_token TEXT NOT NULL DEFAULT '',
    remember_token TEXT NOT NULL DEFAULT '',
    email_verified_at TIMESTAMP NULL,
    two_factor_secret TEXT NOT NULL DEFAULT '',
    two_factor_recovery_codes TEXT NOT NULL DEFAULT '[]',
    two_factor_confirmed_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    pending_two_factor INTEGER NOT NULL,
    pending_remember INTEGER NOT NULL,
    password_confirmed_at TIMESTAMP NULL,
    authenticated_at TIMESTAMP NULL,
    last_seen_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS password_reset_tokens (
    user_id TEXT PRIMARY KEY,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL
);
