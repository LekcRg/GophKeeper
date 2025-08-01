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

CREATE TYPE vault_type AS ENUM('password','note','card', 'binary');

CREATE TABLE vault (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) NOT NULL,
  name VARCHAR(30) NOT NULL,
  type vault_type NOT NULL,
  binary_path VARCHAR(100),
  encrypted_data BYTEA,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE vault;
DROP TYPE vault_type;
-- +goose StatementEnd
