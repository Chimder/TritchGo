-- +goose Up
-- +goose StatementBegin
create extension if not exists "uuid-ossp";

create table stream_stats (
    id uuid primary key default uuid_generate_v4(),
    stream_id varchar(255) not null,
    user_id varchar(255) not null,
    game_id varchar(255) not null,
    date date not null,
    airtime int default 0,
    peak_viewers int default 0,
    average_viewers int default 0,
    hours_watched int default 0,
    unique (stream_id, date)
);

create index idx_stream_id on stream_stats (stream_id);
create index idx_user_id on stream_stats (user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists stream_stats;

-- +goose StatementEnd
