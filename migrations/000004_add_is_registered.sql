-- +goose Up
ALTER TABLE users ADD COLUMN is_registered BOOLEAN NOT NULL DEFAULT FALSE;
UPDATE users SET is_registered = TRUE WHERE state IN ('completed');

-- +goose Down
ALTER TABLE users DROP COLUMN is_registered;
