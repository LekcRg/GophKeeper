-- +goose Up
-- +goose StatementBegin
CREATE TABLE test (
  id SERIAL PRIMARY KEY,
  test varchar(72) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE test
-- +goose StatementEnd
