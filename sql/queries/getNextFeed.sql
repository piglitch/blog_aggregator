-- name: GetNextFeed :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC
LIMIT 1;
