-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    task_id SERIAL PRIMARY KEY,
    project_id INT REFERENCES projects(project_id),
    title VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'open', 
    priority INT NOT NULL DEFAULT 0, 
    due_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id),
    assigned_to INT REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tasks;
-- +goose StatementEnd
