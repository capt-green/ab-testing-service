-- +goose Up
-- +goose StatementBegin
-- Create proxy_changes table
CREATE TABLE IF NOT EXISTS proxy_changes
(
    id             VARCHAR(255) PRIMARY KEY,
    proxy_id       VARCHAR(255) NOT NULL REFERENCES proxies (id) ON DELETE CASCADE,
    change_type    VARCHAR(50)  NOT NULL, -- 'targets_update', 'condition_update'
    previous_state JSONB,
    new_state      JSONB,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255) REFERENCES users (id),
    CONSTRAINT fk_proxy
        FOREIGN KEY (proxy_id)
            REFERENCES proxies (id)
            ON DELETE CASCADE
);

-- Create index for faster lookups
CREATE INDEX idx_proxy_changes_proxy_id ON proxy_changes (proxy_id);
CREATE INDEX idx_proxy_changes_created_at ON proxy_changes (created_at);
-- +goose StatementEnd
