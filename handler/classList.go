package handler

import (
	"log"
	"net/http"
	"practice/json-golang/storage"
)

type ClassList struct {
	Classes    []storage.Class
	SearchTerm string
	Limit      int
	PageNumber int
	Total      int
	TotalPage  int
}

func (h Handler) ListOfClass(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	classList, err := h.storage.GetClasses()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	t := h.Templates.Lookup("class-List.html")
	if t == nil {
		log.Println("unable to lookup template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	data := ClassList{
		Classes: classList,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
