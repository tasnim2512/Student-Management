package postgres

import (
	"practice/json-golang/storage"
)

const getMarksByStudentIdQuery = `
	SELECT SUM(marks) as total_marks from students_subjects WHERE student_id = $1 AND deleted_at IS NULL;
`

func (s PostgresStorage) GetMarksByStudentId(id int) (*storage.TotalMarks, error) {
	var total storage.TotalMarks
	if err := s.DB.Get(&total, getMarksByStudentIdQuery, id); err != nil {
		return nil, err
	}
	return &total, nil
}

const getTotalMarksByStudentIdQuery = `
	SELECT SUM(marks) as marks , students.id ,roll, username
	FROM students
		LEFT JOIN students_subjects ON students.id=students_subjects.student_id
	WHERE students.deleted_at is null
	Group By 
		students_subjects.student_id ,roll, username,students.id
	ORDER BY
		SUM (marks) DESC;
`

func (s PostgresStorage) GetTotalMarksByStudentId() ([]storage.Student, error) {
	var list []storage.Student
	if err := s.DB.Select(&list, getTotalMarksByStudentIdQuery); err != nil {
		return nil, err
	}
	return list, nil
}

