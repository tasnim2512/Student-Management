package postgres

import (
	"log"
	"practice/json-golang/storage"
)

const marksInsertQuery = `
INSERT INTO students_subjects(
	student_id,
	subject_id,
	marks
) VALUES(
	:student_id,
	:subject_id,
	:marks
) RETURNING *;
`

func (s PostgresStorage) CreateStudentSubject(a storage.StudentSubject) (*storage.StudentSubject, error) {
	stmt, err := s.DB.PrepareNamed(marksInsertQuery)
	if err != nil {
		log.Fatal(err)
	}
	if err := stmt.Get(&a, a); err != nil {
		return nil, err
	}
	return &a, nil
}

const updateMarksQuery = `UPDATE students_subjects SET
	marks=:marks
	WHERE id=:id AND deleted_at IS NULL
	RETURNING *;
`

func (s PostgresStorage) UpdateStudentSubject(u storage.StudentSubject) (*storage.StudentSubject, error) {
	stmt, err := s.DB.PrepareNamed(updateMarksQuery)
	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

const getStudentByMarksQuery = `
		SELECT marks FROM students_subjects WHERE deleted_at IS NULL;
	`

func (s PostgresStorage) GetStudentByMark() ([]storage.StudentSubject, error) {
	var u []storage.StudentSubject
	if err := s.DB.Select(&u, getStudentByMarksQuery); err != nil {
		log.Println(err)
		return nil, err
	}
	return u, nil
}
