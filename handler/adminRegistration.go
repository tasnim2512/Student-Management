package handler

import (
	"log"
	"net/http"
	"practice/json-golang/storage"
	"strings"

	"github.com/justinas/nosurf"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AdminForm struct {
	Admin     storage.Admin
	FormError map[string]error
	CSRFToken string
}

func (h Handler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	h.parseCreateAdminTemplate(w, AdminForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreAdmin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := AdminForm{}
	admin := storage.Admin{}
	
	err := h.decoder.Decode(&admin, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form.Admin = admin
	if err := admin.Validate(); err != nil {
		formErr := make(map[string]error)
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
		}
		form.FormError=formErr
		form.CSRFToken = nosurf.Token(r)
		h.parseCreateAdminTemplate(w, form)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	_, err = h.storage.AdminRegistration(admin)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h Handler) parseCreateAdminTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("admin-registration.html")
	if t == nil {
		log.Println("unable to lookup create student template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
