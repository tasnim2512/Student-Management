package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (h Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.storage.DeleteStudent(id)

	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/students", http.StatusSeeOther)
}
