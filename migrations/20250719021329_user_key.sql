-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD key_hash CHAR(64);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN key_hash;
-- +goose StatementEnd
