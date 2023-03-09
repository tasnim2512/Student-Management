package handler

import (
	"fmt"
	"log"
	"net/http"
	"practice/json-golang/storage/postgres"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"

	"golang.org/x/crypto/bcrypt"
)

type LoginAdmin struct {
	UserName  string
	Password  string
	FormError map[string]error
	CSRFToken string
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.parseLoginTemplate(w, LoginAdmin{
		CSRFToken: nosurf.Token(r),
	})
}

func (ls LoginAdmin) Validate() error {
	return validation.ValidateStruct(&ls,
		validation.Field(&ls.UserName,
			validation.Required.Error("the username field is required"),
		),
		validation.Field(&ls.Password,
			validation.Required.Error("the password field is required"),
		),
	)
}

func (h Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var lf LoginAdmin
	if err := h.decoder.Decode(&lf, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
	}

	if err := lf.Validate(); err != nil {
		formErr := make(map[string]error)
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
		}
		lf.FormError = formErr
		lf.Password = ""
		lf.CSRFToken = nosurf.Token(r)
		h.parseLoginTemplate(w, lf)
		return
	}

	admin, err := h.storage.GetAdminByUsername(lf.UserName)
	if err != nil {
		if err.Error() == postgres.NotFound {
			formErr := make(map[string]error)
			formErr["UserName"] = fmt.Errorf("credential doesn't match")
			lf.FormError = formErr
			lf.Password = ""
			lf.CSRFToken = nosurf.Token(r)
			h.parseLoginTemplate(w, lf)
			return
		}

		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(lf.Password)); err != nil {
		formErr := make(map[string]error)
		formErr["UserName"] = fmt.Errorf("credential doesn't match")
		lf.FormError = formErr
		lf.Password = ""
		lf.CSRFToken = nosurf.Token(r)
		h.parseLoginTemplate(w, lf)
		return
	}

	h.sessionManager.Put(r.Context(), "userId", strconv.Itoa(admin.ID))
	http.Redirect(w, r, "/admin/options", http.StatusSeeOther)

}

func (h Handler) parseLoginTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("login.html")
	if t == nil {
		log.Println("unable to lookup login template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
