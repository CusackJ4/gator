-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    "name" TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;


-- connection string: "postgres://jeffshomefolder:@localhost:5432/gator"
-- migration syntax: goose -dir sql/schema postgres "your_connection_string" up

-- activate psql shell: psql "postgres://jeffshomefolder:@localhost:5432/gator"