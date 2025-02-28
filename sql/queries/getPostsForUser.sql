-- name: GetPostsFromUser :many
SELECT posts.* FROM posts
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id 
INNER JOIN users ON users.id = feed_follows.user_id
WHERE users.id = $1
ORDER BY posts.published_at DESC;
