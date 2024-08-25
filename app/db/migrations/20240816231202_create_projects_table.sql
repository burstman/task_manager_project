-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects (
    project_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type_project TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deadline DATE NOT NULL,
    created_by INT REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists projects;
-- +goose StatementEnd
