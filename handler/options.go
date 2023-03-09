package handler

import (
	"log"
	"net/http"
)

type OptionForm struct {
	IsSubjectValid bool
	IsStudentValid bool
}

func (h Handler) OPTIONS(w http.ResponseWriter, r *http.Request) {
	isSubValid := false
	isStValid := false

	classList, err := h.storage.GetClasses()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	subjectList, err := h.storage.GetSubjects()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if len(subjectList) > 0 {
		isStValid = true
	}

	if len(classList) > 0 {
		isSubValid = true
	}

	t := h.Templates.Lookup("options.html")
	if t == nil {
		log.Println("unable to lookup home template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, OptionForm{
		IsSubjectValid: isSubValid,
		IsStudentValid: isStValid,
	}); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
