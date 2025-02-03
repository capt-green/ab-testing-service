-- +goose Up
-- +goose StatementBegin
-- Create proxies table
CREATE TABLE IF NOT EXISTS proxies
(
    id         VARCHAR(255) PRIMARY KEY,
    listen_url VARCHAR(255) NOT NULL,
    mode       VARCHAR(50)  NOT NULL,
    condition  JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create targets table
CREATE TABLE IF NOT EXISTS targets
(
    id        VARCHAR(255) PRIMARY KEY,
    proxy_id  VARCHAR(255) NOT NULL REFERENCES proxies (id) ON DELETE CASCADE,
    url       VARCHAR(255) NOT NULL,
    weight    FLOAT        NOT NULL DEFAULT 1.0,
    is_active BOOLEAN      NOT NULL DEFAULT true
);

-- Create users table
CREATE TABLE IF NOT EXISTS users
(
    id            VARCHAR(255) PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create visits table
CREATE TABLE IF NOT EXISTS visits
(
    id         VARCHAR(255) PRIMARY KEY,
    proxy_id   VARCHAR(255) NOT NULL REFERENCES proxies (id),
    target_id  VARCHAR(255) NOT NULL REFERENCES targets (id),
    user_id    VARCHAR(255) NOT NULL,
    rid        VARCHAR(255) NOT NULL,
    rrid       VARCHAR(255) NOT NULL,
    ruid       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_proxy
        FOREIGN KEY (proxy_id)
            REFERENCES proxies (id)
            ON DELETE CASCADE,
    CONSTRAINT fk_target
        FOREIGN KEY (target_id)
            REFERENCES targets (id)
            ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_visits_proxy_id ON visits (proxy_id);
CREATE INDEX idx_visits_target_id ON visits (target_id);
CREATE INDEX idx_visits_user_id ON visits (user_id);
CREATE INDEX idx_targets_proxy_id ON targets (proxy_id);
CREATE INDEX idx_visits_rid ON visits (rid);
CREATE INDEX idx_visits_rrid ON visits (rrid);
CREATE INDEX idx_visits_ruid ON visits (ruid);
-- +goose StatementEnd
