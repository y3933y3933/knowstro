-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
  ADD CONSTRAINT users_name_unique UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
  DROP CONSTRAINT users_name_unique;
-- +goose StatementEnd
