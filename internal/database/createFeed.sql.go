// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: createFeed.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const addFeed = `-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2,
  $3
)
RETURNING id, created_at, updated_at, name, url, user_id
`

type AddFeedParams struct {
	Name   string
	Url    string
	UserID uuid.UUID
}

func (q *Queries) AddFeed(ctx context.Context, arg AddFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, addFeed, arg.Name, arg.Url, arg.UserID)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}
