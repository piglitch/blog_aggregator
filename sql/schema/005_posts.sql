-- +goose Up
CREATE TABLE posts (
  id uuid PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL,
  feed_id uuid NOT NULL
);

-- +goose Down
DROP TABLE posts;
