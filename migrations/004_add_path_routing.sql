-- +goose Up
-- +goose StatementBegin
-- Add path field to proxies table
ALTER TABLE proxies
    ADD COLUMN path_key VARCHAR(255) DEFAULT NULL;
CREATE UNIQUE INDEX idx_proxies_path_key ON proxies (path_key) WHERE path_key IS NOT NULL;
-- +goose StatementEnd