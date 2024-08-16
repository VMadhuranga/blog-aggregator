-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, last_fetched_in, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;