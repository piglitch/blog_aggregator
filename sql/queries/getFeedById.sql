-- name: GetFeedByUserId :one
SELECT * FROM feeds 
WHERE user_id = $1;