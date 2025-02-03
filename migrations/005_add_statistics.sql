-- +goose Up
-- +goose StatementBegin
CREATE TABLE proxy_stats
(
    id            SERIAL PRIMARY KEY,
    proxy_id      VARCHAR(255) NOT NULL,
    target_id     VARCHAR(255) NOT NULL,
    timestamp     TIMESTAMP    NOT NULL,
    request_count INTEGER      NOT NULL DEFAULT 0,
    error_count   INTEGER      NOT NULL DEFAULT 0,
    unique_users  JSONB,
    FOREIGN KEY (proxy_id) REFERENCES proxies (id) ON DELETE CASCADE
);

CREATE INDEX idx_proxy_stats_proxy_id ON proxy_stats (proxy_id);
CREATE INDEX idx_proxy_stats_timestamp ON proxy_stats (timestamp);
CREATE INDEX idx_proxy_stats_target_id ON proxy_stats (target_id);

-- +goose StatementEnd
