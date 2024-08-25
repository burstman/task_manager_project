-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions (
    permission_id SERIAL PRIMARY KEY ,
    permission_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS permissions;
-- +goose StatementEnd
