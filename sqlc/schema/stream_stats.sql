CREATE TABLE stream_stats (
    id SERIAL PRIMARY KEY,
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