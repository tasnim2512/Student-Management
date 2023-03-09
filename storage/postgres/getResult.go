package postgres

import (
	"practice/json-golang/storage"
)


const getResultByStudentID = `
	SELECT SUM(marks) as total_marks , student_id from students_subjects WHERE student_id = $1 AND deleted_at IS NULL;
`

func (s PostgresStorage) GetResultByStudentId(id int) (*storage.TotalMarks, error) {
	var total storage.TotalMarks
	if err := s.DB.Get(&total, getResultByStudentID, id); err != nil {
		return nil, err
	}
	return &total, nil
}