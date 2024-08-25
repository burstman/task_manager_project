-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
