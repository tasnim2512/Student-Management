package postgres

import (
	"fmt"
	"log"
	"practice/json-golang/storage"
)

const getListQuery = `
		WITH tot AS (SELECT count(*) AS total
		FROM students
		WHERE 
			students.deleted_at IS NULL  
			AND (first_name ILIKE '%%' || $1 || '%%' OR last_name ILIKE '%%' || $1 || '%%') )
			SELECT students.* , tot.total AS total,
			CASE 
			WHEN students_subjects.marks IS NULL THEN  0
			ELSE students_subjects.marks 
			END as marks
			FROM students 
			LEFT JOIN tot ON TRUE
			LEFT JOIN students_subjects ON students.id = students_subjects.student_id
		WHERE 
			students.deleted_at IS NULL  
			AND (first_name ILIKE '%%' || $1 || '%%' OR last_name ILIKE '%%' || $1 || '%%')
			ORDER BY students.id DESC
			OFFSET $2 LIMIT $3;
	`

func (s PostgresStorage) ListOfStudent(uf storage.StudentFilter) ([]storage.Student, error) {
	var listStudent []storage.Student
	if err := s.DB.Select(&listStudent, getListQuery, uf.SearchTerm, uf.Offset, uf.Limit); err != nil {
		log.Println(err)
		return nil, err
	}
	return listStudent, nil
}

const insertQuery = `
INSERT INTO students(
	class_id,
	first_name,
	last_name,
	roll,
	username
) VALUES(
	:class_id,
	:first_name,
	:last_name,
	:roll,
	:username
) RETURNING *;
`

func (s PostgresStorage) CreateStudent(u storage.Student) (*storage.Student, error) {
	stmt, err := s.DB.PrepareNamed(insertQuery)
	if err != nil {
		log.Fatal(err)
	}
	if err := stmt.Get(&u, u); err != nil {
		return nil, err
	}
	if u.ID == 0 {
		log.Println("unable to insert user")
		return nil, fmt.Errorf("unable to insert user")
	}
	return &u, nil
}

const UpdateStudentQuery = `
UPDATE students SET
first_name=:first_name,
last_name=:last_name,
roll=:roll
WHERE id=:id AND deleted_at IS NULL
RETURNING *;
`

func (s PostgresStorage) UpdateStudent(u storage.Student) (*storage.Student, error) {
	stmt, err := s.DB.PrepareNamed(UpdateStudentQuery)
	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&u, u); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

const getStudentByIdQuery = `
		SELECT * FROM students WHERE id=$1 AND deleted_at IS NULL;
	`

func (s PostgresStorage) GetStudentById(id string) (*storage.Student, error) {
	var u storage.Student
	if err := s.DB.Get(&u, getStudentByIdQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

const DeleteStudentQuery = `
UPDATE students SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted_at IS NULL 
RETURNING id;
`

func (s PostgresStorage) DeleteStudent(id string) error {
	res, err := s.DB.Exec(DeleteStudentQuery, id)
	if err != nil {
		log.Println(err)
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if row >= 0 {
		fmt.Printf("unable to delete user")
	}
	return nil
}

const getStudentDetailByStudentIdQuery = `
		SELECT students.id, students.username , class_id, class_name FROM students LEFT JOIN classes ON students.class_id = classes.id WHERE students.id = $1 AND students.deleted_at IS NULL;
	`

func (s PostgresStorage) GetStudentDetailByStudentId(id string) (*storage.StudentDetails, error) {
	var u storage.StudentDetails
	if err := s.DB.Get(&u, getStudentDetailByStudentIdQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}

const getStudentListQuery = `
		SELECT * FROM students;
	`

func (s PostgresStorage) GetStudentList() ([]storage.Student, error) {
	var u []storage.Student
	if err := s.DB.Select(&u, getStudentListQuery); err != nil {
		log.Println(err)
		return nil, err
	}
	return u, nil
}

const getStudentByClassQuery = `
		SELECT students.* 
		FROM students 
		LEFT JOIN classes ON students.class_id=classes.id
		WHERE students.class_id=$1 AND students.deleted_at IS NULL;
	`

func (s PostgresStorage) GetStudentByClass(id string) ([]storage.Student, error) {
	var u []storage.Student
	if err := s.DB.Select(&u, getStudentByClassQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}
	return u, nil
}


const getStudentByClassIDQuery = `
		SELECT students.* 
		FROM students 
		LEFT JOIN classes ON students.class_id=classes.id
		WHERE students.class_id=$1 AND students.deleted_at IS NULL;
	`

func (s PostgresStorage) GetStudentByClassId(id string) (*storage.Student, error) {
	var u storage.Student
	if err := s.DB.Get(&u, getStudentByClassIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}
	return &u, nil
}