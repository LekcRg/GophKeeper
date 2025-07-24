-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	login VARCHAR(30) UNIQUE NOT NULL,
	passhash VARCHAR(72) NOT NULL,
	key_hash CHAR(64),
	encrypted_tag TEXT NOT NULL,
	salt TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
