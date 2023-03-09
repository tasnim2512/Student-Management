package postgres

import (
	"log"
	"practice/json-golang/storage"
)

const subjectInsertQuery = `
INSERT INTO subjects(
	subject_name,
	class_id
) VALUES(
	:subject_name,
	:class_id
) RETURNING *;
`

func (s PostgresStorage) CreateSubject(sb storage.Subject) (*storage.Subject, error) {
	stmt, err := s.DB.PrepareNamed(subjectInsertQuery)
	if err != nil {
		log.Fatal(err)
	}
	if err := stmt.Get(&sb, sb); err != nil {
		return nil, err
	}
	return &sb, nil
}

const getSubjectsQuery = `
SELECT * FROM subjects;
`

func (s PostgresStorage) GetSubjects() ([]storage.Subject, error) {
	var subjects []storage.Subject
	if err := s.DB.Select(&subjects, getSubjectsQuery); err != nil {
		return nil, err
	}
	return subjects, nil
}

const getSubjectsByClassIdQuery = `
	SELECT subjects.id, subject_name, class_id
	FROM subjects
	WHERE subjects.class_id = $1 AND subjects.deleted_at IS NULL;
	
`

func (s PostgresStorage) GetSubjectsByClassId(classId int) ([]storage.Subject, error) {
	var subjects []storage.Subject
	if err := s.DB.Select(&subjects, getSubjectsByClassIdQuery,classId); err != nil {
		return nil, err
	}
	return subjects, nil
}

const getSubjectsByClassIdEditQuery = `
	SELECT students_subjects.id as ss_id, subjects.id, subject_name, class_id,
	CASE 
	WHEN students_subjects.marks IS NULL THEN  0
	ELSE students_subjects.marks 
	END AS mark
	FROM subjects 
	LEFT JOIN students_subjects ON students_subjects.subject_id = subjects.id
	WHERE students_subjects.student_id = $1 AND subjects.class_id = $2 AND subjects.deleted_at IS NULL;
	
`

func (s PostgresStorage) GetSubjectsByClassIdEdit(studentID int, classId int) ([]storage.Subject, error) {
	var subjects []storage.Subject
	if err := s.DB.Select(&subjects, getSubjectsByClassIdEditQuery, studentID,classId); err != nil {
		return nil, err
	}
	return subjects, nil
}

const getSubjectsMarksByClassIdQuery = `
	SELECT subjects.id, subject_name, class_id, students_subjects.marks
	FROM subjects 
	LEFT JOIN students_subjects ON subjects.id = students_subjects.subject_id
	WHERE students_subjects.subject_id = $1 AND students_subjects.student_id = $2 AND subjects.deleted_at IS NULL;
`

func (s PostgresStorage) GetSubjectMarkByClassId(sbid, stid int) (*storage.Subject, error) {
	var subjectMarks storage.Subject

	if err := s.DB.Get(&subjectMarks, getSubjectsMarksByClassIdQuery, sbid, stid); err != nil {
		return nil, err
	}
	return &subjectMarks, nil
}

const getSubjectByClassIdQuery = `
	SELECT subjects.id, subject_name, class_id
	FROM subjects 
	LEFT JOIN classes ON subjects.class_id = classes.id
	WHERE subjects.class_id = $1  AND subjects.deleted_at IS NULL;
`

func (s PostgresStorage) GetSubjectByClassId(id int) (*storage.Subject, error) {
	var subject storage.Subject

	if err := s.DB.Get(&subject, getSubjectByClassIdQuery, id); err != nil {
		return nil, err
	}
	return &subject, nil
}
