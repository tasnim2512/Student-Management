package handler

import (
	"fmt"
	"log"
	"net/http"
	"practice/json-golang/storage"
	"strings"

	"github.com/justinas/nosurf"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ClassForm struct {
	Class          storage.Class
	ClassFormError map[string]error
	CSRFToken      string
}

func (h Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	h.parseCreateClassTemplate(w, ClassForm{
		CSRFToken: nosurf.Token(r),
	})

}

func (h Handler) StoreClass(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := ClassForm{}
	class := storage.Class{}

	err := h.decoder.Decode(&class, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form.Class = class

	alreadyExists, _ := h.CheckAlreadyExists(form.Class.ClassName)

	if !alreadyExists {
		formErr := make(map[string]error)
		formErr["ClassName"] = fmt.Errorf("class already exists")
		form.ClassFormError = formErr
		form.CSRFToken = nosurf.Token(r)
		h.parseCreateClassTemplate(w, form)
		return
	}

	if err := class.Validate(); err != nil {
		formErr := make(map[string]error)
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
		}
		form.ClassFormError = formErr
		form.CSRFToken = nosurf.Token(r)
		h.parseCreateClassTemplate(w, form)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	_, err = h.storage.CreateClass(class)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/admin/create/class", http.StatusSeeOther)
}

func (h Handler) parseCreateClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-class.html")
	if t == nil {
		log.Println("unable to lookup create student template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h Handler) CheckAlreadyExists(value string) (bool, error) {

	newCLass, _ := h.storage.GetClassesByName(value)
	if newCLass.ClassName == "" {
		return true, nil
	}
	return false, nil
}
