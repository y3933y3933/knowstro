-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resource_types (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resource_types;
-- +goose StatementEnd
