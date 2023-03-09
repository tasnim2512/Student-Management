package handler

import (
	"log"
	"net/http"
	"practice/json-golang/storage"
)

type ResultSheet struct {
	StudentList []storage.Student
	Total       []storage.TotalMarks
}

func (h Handler) GetResultSheet(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	studentList, err := h.storage.GetTotalMarksByStudentId()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("result-sheet.html")
	if t == nil {
		log.Println("unable to lookup template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	data := ResultSheet{
		StudentList: studentList,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
