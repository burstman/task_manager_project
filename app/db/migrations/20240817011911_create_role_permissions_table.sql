-- +goose Up
-- +goose StatementBegin
CREATE TABLE role_permissions (
    role_id INT REFERENCES roles(role_id) ON DELETE CASCADE,
    permission_id INT REFERENCES permissions(permission_id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_permissions;
-- +goose StatementEnd
