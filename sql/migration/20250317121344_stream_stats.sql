-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE stream_stats (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    stream_id varchar(255) NOT NULL,
    user_id varchar(255) NOT NULL,
    game_id varchar(255) NOT NULL,
    date date NOT NULL,
    airtime int DEFAULT 0,
    peak_viewers int DEFAULT 0,
    average_viewers int DEFAULT 0,
    hours_watched int DEFAULT 0,
    UNIQUE (stream_id, date)
);

CREATE INDEX idx_stream_id ON stream_stats (stream_id);
CREATE INDEX idx_user_id ON stream_stats (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stream_stats;

-- +goose StatementEnd
