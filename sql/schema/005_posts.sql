-- +goose up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;


-- connection string: "postgres://jeffshomefolder:@localhost:5432/gator"
-- migration syntax: goose -dir sql/schema postgres "your_connection_string" up
-- goose -dir sql/schema postgres "postgres://jeffshomefolder:@localhost:5432/gator" up