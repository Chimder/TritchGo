-- +goose Up
-- +goose StatementBegin
SET search_path TO public;

create extension if not exists "uuid-ossp";

create table stream_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    game_id VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    airtime INT DEFAULT 0,
    peak_viewers INT DEFAULT 0,
    average_viewers INT DEFAULT 0,
    hours_watched INT DEFAULT 0,
    UNIQUE (stream_id, date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists stream_stats;
-- +goose StatementEnd
