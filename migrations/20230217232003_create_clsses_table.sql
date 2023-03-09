-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS classes (
    id BIGSERIAL,
    class_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    PRIMARY KEY(id),
    UNIQUE(class_name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS classes
-- +goose StatementEnd
