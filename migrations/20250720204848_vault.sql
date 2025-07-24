-- +goose Up
-- +goose StatementBegin
CREATE TYPE vault_type AS ENUM('password','note','card', 'binary');

CREATE TABLE vault (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) NOT NULL,
  name VARCHAR(30) NOT NULL,
  type vault_type NOT NULL,
  encrypted_data TEXT,
  created_at timestamp default current_timestamp,
  updated_at timestamp default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE vault;
DROP TYPE vault_type;
-- +goose StatementEnd
