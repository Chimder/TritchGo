-- name: InsertStreamStats :exec
INSERT INTO stream_stats (
    stream_id, user_id, game_id, date, airtime, peak_viewers, average_viewers, hours_watched
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) ON CONFLICT (stream_id, date)
DO UPDATE SET
    airtime = EXCLUDED.airtime,
    peak_viewers = GREATEST(stream_stats.peak_viewers, EXCLUDED.peak_viewers),
    average_viewers = ROUND(stream_stats.average_viewers + EXCLUDED.average_viewers) / 2,
    hours_watched = stream_stats.hours_watched + ROUND(EXCLUDED.average_viewers * (EXCLUDED.airtime / 60.0));



-- name: GetStatsByUserId :one
SELECT * FROM stream_stats WHERE user_id = $1;

-- name: GetStatsByStreamId :one
SELECT * FROM stream_stats WHERE stream_id = $1;
