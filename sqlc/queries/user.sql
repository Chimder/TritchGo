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

-- -- name: InsertStreamStats :exec
-- INSERT INTO stream_stats (
--     stream_id, user_id, game_id, date, airtime, peak_viewers, average_viewers, hours_watched
-- ) VALUES (
--     $1,$2,$3,NOW(),$5,$6,$7,$8
-- ) ON CONFLICT (stream_id, date)
-- DO UPDATE SET
--     airtime = stream_stats.airtime + EXCLUDED.airtime,
--     peak_viewers = GREATEST(stream_stats.peak_viewers, EXCLUDED.peak_viewers),
--     average_viewers = (stream_stats.average_viewers + EXCLUDED.average_viewers) / 2,
--     hours_watched = stream_stats.hours_watched + (EXCLUDED.average_viewers * (15.0 / 60.0));
