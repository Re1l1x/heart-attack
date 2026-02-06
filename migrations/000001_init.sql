-- +goose Up

CREATE TYPE confirmation_state AS ENUM ('not_confirmed', 'confirmed', 'cancelled');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    is_bot BOOLEAN NOT NULL DEFAULT FALSE,
    language_code TEXT,
    is_premium BOOLEAN NOT NULL DEFAULT FALSE,
    sex TEXT,
    about TEXT NOT NULL DEFAULT '',
    state TEXT NOT NULL DEFAULT 'start',
    time_ranges TEXT NOT NULL DEFAULT '000000',
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE places (
    id BIGSERIAL PRIMARY KEY,
    description TEXT NOT NULL
);

CREATE TABLE meetings (
    id BIGSERIAL PRIMARY KEY,
    dill_id BIGINT NOT NULL REFERENCES users(id),
    doe_id BIGINT NOT NULL REFERENCES users(id),
    pair_score DOUBLE PRECISION NOT NULL,
    is_fullmatch BOOLEAN NOT NULL DEFAULT FALSE,
    place_id BIGINT REFERENCES places(id),
    time TEXT,
    dill_state confirmation_state NOT NULL DEFAULT 'not_confirmed',
    doe_state confirmation_state NOT NULL DEFAULT 'not_confirmed'
);

-- +goose Down
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS places;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS confirmation_state;
