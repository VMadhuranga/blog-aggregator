-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_in TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_in;
