package handler

import (
	"log"
	"net/http"
)

func  (h Handler) HOME(w http.ResponseWriter, r *http.Request) {
	t := h.Templates.Lookup("home.html")
	if t == nil {
		log.Println("unable to lookup home template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w,nil); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
