-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS attachments (
    attachment_id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(task_id),
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    uploaded_by INT REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists attachments;
-- +goose StatementEnd
