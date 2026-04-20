-- +goose Up
ALTER TABLE feeds
ADD COLUMN last_fetched_at TIMESTAMP;
-- +goose Down



-- connection string: "postgres://jeffshomefolder:@localhost:5432/gator"
-- migration syntax: goose -dir sql/schema postgres "your_connection_string" up