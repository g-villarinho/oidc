-- Tabela de clients OAuth2
CREATE TABLE oauth_clients (
    id UUID PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    client_secret VARCHAR(255) NOT NULL,
    client_name VARCHAR(255) NOT NULL,
    redirect_uris TEXT[] NOT NULL,
    grant_types TEXT[] NOT NULL,
    response_types TEXT[] NOT NULL,
    scopes TEXT[] NOT NULL,
    logo_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela de usuários
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela de authorization codes
CREATE TABLE authorization_codes (
    code VARCHAR(255) PRIMARY KEY,
    client_id varchar(255) NOT NULL REFERENCES oauth_clients(client_id),
    user_id UUID NOT NULL REFERENCES users(id),
    redirect_uri TEXT NOT NULL,
    scopes TEXT[] NOT NULL,
    code_challenge VARCHAR(255),
    code_challenge_method VARCHAR(10),
    used BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_codes_expires ON authorization_codes(expires_at);
CREATE INDEX idx_auth_codes_expires ON authorization_codes(expires_at);
CREATE INDEX idx_auth_codes_used ON authorization_codes(used) WHERE used = FALSE;

-- Tabela de tokens
CREATE TABLE tokens (
    id UUID PRIMARY KEY,
    access_token_hash VARCHAR(64) NOT NULL UNIQUE,
    refresh_token_hash VARCHAR(64) NOT NULL UNIQUE,
    authorization_code VARCHAR(255) REFERENCES authorization_codes(code) ON DELETE SET NULL,
    client_id VARCHAR(255) NOT NULL REFERENCES oauth_clients(client_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scopes TEXT[] NOT NULL,
    token_type VARCHAR(50) NOT NULL DEFAULT 'Bearer',
    access_token_expires_at TIMESTAMP NOT NULL,
    refresh_token_expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    revoked_at TIMESTAMP,
    revoked_reason VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMP
);

CREATE INDEX idx_tokens_access_hash ON tokens(access_token_hash);
CREATE INDEX idx_tokens_refresh_hash ON tokens(refresh_token_hash);
CREATE INDEX idx_tokens_access_expires ON tokens(access_token_expires_at);
CREATE INDEX idx_tokens_user_id ON tokens(user_id);
CREATE INDEX idx_tokens_client_id ON tokens(client_id);
CREATE INDEX idx_tokens_revoked ON tokens(revoked) WHERE revoked = FALSE;
CREATE INDEX idx_tokens_auth_code ON tokens(authorization_code) WHERE authorization_code IS NOT NULL;

