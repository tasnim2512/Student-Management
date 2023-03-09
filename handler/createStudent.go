package handler

import (
	"fmt"
	"log"
	"net/http"
	"practice/json-golang/storage"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type StudentForm struct {
	ListOfClasses []storage.Class
	Student       storage.Student
	FormError     map[string]error
	CSRFToken     string
}

func (h Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	classList, err := h.storage.GetClasses()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.parseCreateTemplate(w, StudentForm{
		ListOfClasses: classList,
		CSRFToken:     nosurf.Token(r),
	})

}

func (h Handler) StoreStudent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := StudentForm{}
	student := storage.Student{}

	err := h.decoder.Decode(&student, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.Student = student
	if err := student.Validate(); err != nil {
		formErr := make(map[string]error)
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
		}
		form.FormError = formErr
		form.CSRFToken = nosurf.Token(r)
		h.parseCreateTemplate(w, form)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	cl := r.FormValue("ClassID")
	classid, err := strconv.Atoi(cl)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	student.ClassID = classid

	newStudent, err := h.storage.CreateStudent(student)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/admin/%v/edit/student", newStudent.ID), http.StatusSeeOther)
}

func (h Handler) parseCreateTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-student.html")
	if t == nil {
		log.Println("unable to lookup create student template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
