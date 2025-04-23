-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resource_tags (
  resource_id INT REFERENCES resources(id) ON DELETE CASCADE,
  tag_id      INT REFERENCES tags(id)       ON DELETE CASCADE,
  PRIMARY KEY (resource_id, tag_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resource_tags;
-- +goose StatementEnd
