// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getPostsForUser.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getPostsFromUser = `-- name: GetPostsFromUser :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.description, posts.published_at, posts.feed_id FROM posts
INNER JOIN feed_follows ON feed_follows.feed_id = posts.feed_id 
INNER JOIN users ON users.id = feed_follows.user_id
WHERE users.id = $1
ORDER BY posts.published_at DESC
`

func (q *Queries) GetPostsFromUser(ctx context.Context, id uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsFromUser, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
