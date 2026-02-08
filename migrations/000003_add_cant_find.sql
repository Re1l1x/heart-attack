-- +goose Up
ALTER TABLE meetings ADD COLUMN dill_cant_find BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE meetings ADD COLUMN doe_cant_find BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE meetings DROP COLUMN dill_cant_find;
ALTER TABLE meetings DROP COLUMN doe_cant_find;
