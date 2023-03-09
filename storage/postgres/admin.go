package postgres

import (
	"fmt"
	"log"
	"practice/json-golang/storage"

	"golang.org/x/crypto/bcrypt"
)

const adminInsertQuery = `
INSERT INTO admin(
	first_name,
	last_name,
	email,
	username,
	password
) VALUES(
	:first_name,
	:last_name,
	:email,
	:username,
	:password
) RETURNING *;
`

func (s PostgresStorage) AdminRegistration(a storage.Admin) (*storage.Admin, error) {
	stmt, err := s.DB.PrepareNamed(adminInsertQuery)
	if err != nil {
		log.Fatal(err)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	a.Password = string(hashPass)
	if err := stmt.Get(&a, a); err != nil {
		return nil, err
	}
	if a.ID == 0 {
		log.Println("unable to insert user")
		return nil, fmt.Errorf("unable to insert user")
	}
	return &a, nil
}

const getAdminByUsernameQuery = `
		SELECT * FROM admin WHERE username=$1 AND deleted_at IS NULL;
	`

func (s PostgresStorage) GetAdminByUsername(username string) (*storage.Admin, error) {
	var a storage.Admin
	if err := s.DB.Get(&a, getAdminByUsernameQuery, username); err != nil {
		log.Println(err)
	}

	return &a, nil
}
