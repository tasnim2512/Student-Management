package handler

import (
	"fmt"
	"log"
	"net/http"
	"practice/json-golang/storage"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) EditMarks(w http.ResponseWriter, r *http.Request) {
	studentId := chi.URLParam(r, "id")
	sId, _ := strconv.Atoi(studentId)
	classID := chi.URLParam(r, "classId")
	cId, _ := strconv.Atoi(classID)

	studentDetail, err := h.storage.GetStudentDetailByStudentId(studentId)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	subjectList, err := h.storage.GetSubjectsByClassIdEdit(sId, cId)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	h.parseEditMarksTemplate(w, StudentSubjectTemplate{
		StudentDetail: studentDetail,
		ListOfSubject: subjectList,
		CSRFToken:     nosurf.Token(r),
	})
}

func (h Handler) UpdateMarks(w http.ResponseWriter, r *http.Request) {
	studentId := chi.URLParam(r, "id")
	sId,err := strconv.Atoi(studentId)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	classID := chi.URLParam(r, "classId")
	cId, err := strconv.Atoi(classID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	studentDetail, err := h.storage.GetStudentDetailByStudentId(studentId)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	subjectList, err := h.storage.GetSubjectsByClassIdEdit(sId, cId)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var sb StudentSubjectForm

	err = h.decoder.Decode(&sb, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := sb.Validate(); err != nil {
		formErr := make(map[string]error)
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
		}
		h.parseAddMarksTemplate(w, StudentSubjectTemplate{
			StudentDetail: studentDetail,
			ListOfSubject: subjectList,
			FormError:     formErr,
		})
	}

	if len(sb.SubjectMarks) > 0 {
		for key, val := range sb.SubjectMarks {
			k, _ := strconv.Atoi(key)
			_, err = h.storage.UpdateStudentSubject(storage.StudentSubject{
				ID:    k,
				Marks: val,
			})
			if err != nil {
				log.Println(err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/students/%v/edit/marks/%v", studentId, classID), http.StatusSeeOther)
}

func (h Handler) parseEditMarksTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("edit-marks.html")
	if t == nil {
		log.Println("unable to lookup create subject template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
