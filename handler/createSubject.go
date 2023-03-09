package handler

import (
	"log"
	"net/http"
	"practice/json-golang/storage"
	"strconv"
	"strings"

	"github.com/justinas/nosurf"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SubjectForm struct {
	ListOfClasses    []storage.Class
	Subject          storage.Subject
	SubjectFormError map[string]error
	CSRFToken        string
	ClassID          string
}

func (h Handler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	classList, err := h.storage.GetClasses()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.parseCreateSubjectTemplate(w, SubjectForm{
		ListOfClasses: classList,
		CSRFToken:     nosurf.Token(r),
	})

}

func (h Handler) StoreSubject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := SubjectForm{}
	subject := storage.Subject{}

	err := h.decoder.Decode(&subject, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form.Subject = subject

	if err := subject.Validate(); err != nil {
		formErr := make(map[string]error)
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
		}
		form.SubjectFormError = formErr
		form.CSRFToken = nosurf.Token(r)
		h.parseCreateSubjectTemplate(w, form)
		return
	}

	cl := r.FormValue("ClassID")
	classid, err := strconv.Atoi(cl)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	subject.ClassID = classid
	_, err = h.storage.CreateSubject(subject)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admin/create/subject", http.StatusSeeOther)
}

func (h Handler) parseCreateSubjectTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-subject.html")
	if t == nil {
		log.Println("unable to lookup create subject template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
