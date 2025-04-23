-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS resource_subjects (
  resource_id INT REFERENCES resources(id) ON DELETE CASCADE,
  subject_id  INT REFERENCES subjects(id)  ON DELETE CASCADE,
  PRIMARY KEY (resource_id, subject_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS resource_subjects;
-- +goose StatementEnd
