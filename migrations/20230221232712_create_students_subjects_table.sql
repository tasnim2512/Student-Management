-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students_subjects (
    id BIGSERIAL,
    student_id BIGINT,
    subject_id BIGINT,
    marks INT DEFAULT 0 ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,

    PRIMARY KEY(id),
    CONSTRAINT fk_id
    FOREIGN KEY(student_id)
	REFERENCES students (id) on delete SET NULL,
    FOREIGN KEY(subject_id)
	REFERENCES subjects (id) on delete SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students_subjects
-- +goose StatementEnd
