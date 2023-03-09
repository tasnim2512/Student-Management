package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Result struct {
	Total       int
	StudentID   int
	StudentName string
}

func (h Handler) ShowResult(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	Sid, _ := strconv.Atoi(id)

	total, err := h.storage.GetMarksByStudentId(Sid)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.parseShowResultTemplate(w, Result{
		Total: total.TotalMarks,
		StudentName: total.StudentName,
	})
}

func (h Handler) parseShowResultTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("student-result.html")
	if t == nil {
		log.Println("unable to lookup create subject template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
