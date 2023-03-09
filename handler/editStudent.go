package handler

import (
	"fmt"
	"log"
	"net/http"
	"practice/json-golang/storage"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) EditStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	classList, err := h.storage.GetClasses()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	editStudent, err := h.storage.GetStudentById(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := StudentForm{
		ListOfClasses: classList,
	}
	form.Student = *editStudent
	form.CSRFToken = nosurf.Token(r)
	h.parseEditTemplate(w, form)
}

func (h Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	Sid, _ := strconv.Atoi(id)

	student := storage.Student{
		ID: Sid,
	}
	var form StudentForm

	err := h.decoder.Decode(&student, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form.Student = student

	form.CSRFToken = nosurf.Token(r)
	if err := student.Validate(); err != nil {

		if vErr, ok := err.(validation.Errors); ok {
			newErr := make(map[string]error)
			for key, val := range vErr {
				newErr[strings.Title(key)] = val
			}
			form.FormError = newErr
		}

		h.parseEditTemplate(w, form)
		return
	}

	updateStudent, err := h.storage.UpdateStudent(student)

	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, fmt.Sprintf("/admin/%v/edit/student", updateStudent.ID), http.StatusSeeOther)
}

func (h Handler) parseEditTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("edit-student.html")
	if t == nil {
		log.Println("unable to look up edit student template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
