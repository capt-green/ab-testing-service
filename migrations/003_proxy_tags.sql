-- +goose Up
-- +goose StatementBegin
-- Add tags column to proxies table
ALTER TABLE proxies
    ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}';

-- Create index for faster tag-based searches
CREATE INDEX IF NOT EXISTS idx_proxies_tags ON proxies USING GIN (tags);
-- +goose StatementEnd
