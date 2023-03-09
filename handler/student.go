package handler

import (
	"log"
	"math"
	"net/http"
	"practice/json-golang/storage"
	"strconv"
)

const ListLimit = 2

type studentList struct {
	Students   []storage.Student
	SearchTerm string
	Limit      int
	PageNumber int
	Total      int
	TotalPage  int
}

func (h Handler) ListOfStudent(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var err error
	pageNumber := 1
	pn := r.FormValue("page")
	if pn != "" {
		pageNumber, _ = strconv.Atoi(pn)
		if err != nil {
			pageNumber = 1
		}
	}

	offset := 0
	if pageNumber > 1 {
		offset = (pageNumber * ListLimit) - ListLimit
	}

	st := r.FormValue("SearchTerm")
	uf := storage.StudentFilter{
		SearchTerm: st,
		Offset:     offset,
		Limit:      ListLimit,
	}
	listStudent, err := h.storage.ListOfStudent(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("student-list.html")
	if t == nil {
		log.Println("unable to lookup template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	total := 0
	if len(listStudent) > 0 {
		total = listStudent[0].Total
	}
	totalPage := int(math.Ceil(float64(total) / float64(ListLimit)))

	data := studentList{
		Students:   listStudent,
		SearchTerm: st,
		Limit:      ListLimit,
		PageNumber: pageNumber,
		Total:      total,
		TotalPage:  totalPage,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
