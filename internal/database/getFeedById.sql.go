// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: getFeedById.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getFeedByUserId = `-- name: GetFeedByUserId :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at FROM feeds 
WHERE user_id = $1
`

func (q *Queries) GetFeedByUserId(ctx context.Context, userID uuid.UUID) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getFeedByUserId, userID)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
		&i.LastFetchedAt,
	)
	return i, err
}
