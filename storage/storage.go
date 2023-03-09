package storage

import (
	"database/sql"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type StudentFilter struct {
	SearchTerm string
	Offset     int
	Limit      int
}
type Student struct {
	ID        int           `form:"-" db:"id"`
	ClassID   int           `db:"class_id"`
	FirstName string        `db:"first_name"`
	LastName  string        `db:"last_name"`
	Roll      string        `db:"roll"`
	UserName  string        `db:"username"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	DeletedAt sql.NullTime  `db:"deleted_at"`
	Total     int           `db:"total"`
	Marks     int `db:"marks"`
}
type Admin struct {
	ID        int          `form:"-" db:"id"`
	FirstName string       `db:"first_name"`
	LastName  string       `db:"last_name"`
	Email     string       `db:"email"`
	UserName  string       `db:"username"`
	Password  string       `db:"password"`
	Status    bool         `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type Class struct {
	ID        int          `form:"-" db:"id"`
	ClassName string       `db:"class_name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type StudentSubject struct {
	ID        int          `form:"-" db:"id"`
	StudentID int          `db:"student_id"`
	SubjectID int          `db:"subject_id"`
	Marks     int          `db:"marks"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type TotalMarks struct {
	TotalMarks  int    `db:"total_marks"`
	StudentName string `db:"username"`
	StudentID   string `db:"student_id"`
}
type Subject struct {
	ID           int          `form:"-" db:"id"`
	StudentSubID int          `db:"ss_id"`
	ClassID      int          `db:"class_id"`
	SubjectName  string       `db:"subject_name"`
	Mark         int          `db:"mark"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

type StudentDetails struct {
	ID          int    `db:"id"`
	ClassID     int    `db:"class_id"`
	ClassName   string `db:"class_name"`
	StudentName string `db:"username"`
}

func (s Student) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.FirstName,
			validation.Required.Error("the First name field is required"),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("the Last name field is required"),
		),
		validation.Field(&s.UserName,
			validation.Required.Error("the User name field is required"),
		),
		validation.Field(&s.Roll,
			validation.Required.Error("the roll field is required"),
			is.Int.Error("please put a integer value as roll"),
			validation.Match(regexp.MustCompile(`^[1-9]\d*$`)).Error("the roll can not be negative"),
		),
	)
}
func (a Admin) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.FirstName,
			validation.Required.Error("the First name field is required"),
		),
		validation.Field(&a.LastName,
			validation.Required.Error("the Last name field is required"),
		),
		validation.Field(&a.UserName,
			validation.Required.Error("the User ame field is required"),
		),
		validation.Field(&a.Email,
			validation.Required.Error("the email field is required"),
			is.Email.Error("Please put a Valid email"),
		),
		validation.Field(&a.Password,
			validation.Required.Error("the password field is required"),
		),
	)
}
func (c Class) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClassName,
			validation.Required.Error("the class field is required"),
			validation.Match(regexp.MustCompile("^([1-9]|1[0]])$")).Error("the class should between 1 to 10"),
			validation.By(func(interface{}) error {
				return nil
			}),
		),
	)
}

func (sb Subject) Validate() error {
	return validation.ValidateStruct(&sb,
		validation.Field(&sb.SubjectName,
			validation.Required.Error("the name field is required"),
		),
	)
}
