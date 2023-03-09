package postgres

import (
	"fmt"
	"log"
	"practice/json-golang/storage"
)

const cLassInsertQuery = `
INSERT INTO classes(
	class_name
) VALUES(
	:class_name
) RETURNING *;
`

func (s PostgresStorage) CreateClass(a storage.Class) (*storage.Class, error) {
	stmt, err := s.DB.PrepareNamed(cLassInsertQuery)
	if err != nil {
		log.Fatal(err)
	}
	if err := stmt.Get(&a, a); err != nil {
		return nil, err
	}
	if a.ID == 0 {
		log.Println("unable to insert user")
		return nil, fmt.Errorf("unable to insert user")
	}
	return &a, nil
}

const getClassesQuery = `
SELECT * FROM classes;
`

func (s PostgresStorage) GetClasses() ([]storage.Class, error) {
	var classes []storage.Class
	if err := s.DB.Select(&classes, getClassesQuery); err != nil {
		return nil, err
	}
	return classes, nil
}



const getClassesByNameQuery = `
SELECT class_name FROM classes WHERE class_name= $1 AND deleted_at IS NULL;
`

func (s PostgresStorage) GetClassesByName(name string) (storage.Class, error) {
	var classes storage.Class
	if err := s.DB.Get(&classes, getClassesByNameQuery,name); err != nil {
		return storage.Class{}, err
	}
	return classes, nil
}
