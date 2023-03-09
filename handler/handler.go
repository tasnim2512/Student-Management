package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"practice/json-golang/storage"
	"strconv"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/form"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
	storage        dbStorage
	Templates      *template.Template
}

type dbStorage interface {
	AdminRegistration(a storage.Admin) (*storage.Admin, error)
	GetAdminByUsername(username string) (*storage.Admin, error)

	CreateClass(storage.Class) (*storage.Class, error)
	GetClasses() ([]storage.Class, error)
	GetClassesByName(string) (storage.Class,error)

	GetSubjects() ([]storage.Subject, error)
	GetSubjectsByClassId(int) ([]storage.Subject, error)
	GetSubjectsByClassIdEdit(int, int) ([]storage.Subject, error)
	GetSubjectByClassId(int) (*storage.Subject, error)
	GetSubjectMarkByClassId(int, int) (*storage.Subject, error)
	CreateSubject(storage.Subject) (*storage.Subject, error)

	ListOfStudent(storage.StudentFilter) ([]storage.Student, error)
	CreateStudent(u storage.Student) (*storage.Student, error)
	GetStudentDetailByStudentId(id string) (*storage.StudentDetails, error)
	UpdateStudent(u storage.Student) (*storage.Student, error)
	DeleteStudent(id string) error
	GetStudentById(id string) (*storage.Student, error)
	GetStudentList() ([]storage.Student, error)
	GetStudentByClass(string) ([]storage.Student, error)
	GetStudentByClassId(string) (*storage.Student, error)
	GetStudentByMark() ([]storage.StudentSubject, error)

	CreateStudentSubject(storage.StudentSubject) (*storage.StudentSubject, error)
	UpdateStudentSubject(storage.StudentSubject) (*storage.StudentSubject, error)

	GetMarksByStudentId(int) (*storage.TotalMarks, error)
	GetTotalMarksByStudentId() ([]storage.Student, error)
}
type ErrorPage struct {
	Code    int
	Message string
}

func (h Handler) Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	ep := ErrorPage{
		Code:    code,
		Message: error,
	}

	tf := "default"
	switch code {
	case 400, 401, 402, 403, 404:
		tf = "4xx"
	case 500, 501, 503:
		tf = "5xx"
	}

	tpl := fmt.Sprintf("templates/errors/%s.html", tf)
	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Fatalln(err)
	}

	if err := t.Execute(w, ep); err != nil {
		log.Fatalln(err)
	}
}

func NewHandler(sm *scs.SessionManager, formDecoder *form.Decoder, storage dbStorage) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        formDecoder,
		storage:        storage,
	}

	h.ParseTemplates()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(Method)

	r.Get("/", h.HOME)
	r.Get("/registration", h.CreateAdmin)

	r.Post("/admin/store", h.StoreAdmin)

	r.Get("/login", h.Login)
	r.Post("/login", h.LoginPostHandler)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	r.Route("/students", func(r chi.Router) {
		r.Use(h.Authenticator)

		r.Get("/", h.ListOfStudent)

		r.Get("/{id:[0-9]+}/delete", h.DeleteStudent)

		r.Get("/{id:[0-9]+}/add/marks/{classId:[0-9]+}", h.GetMarks)

		r.Post("/{id:[0-9]+}/store/marks/{classId:[0-9]+}", h.StoreMarks)

		r.Get("/{id:[0-9]+}/edit/marks/{classId:[0-9]+}", h.EditMarks)

		r.Post("/{id:[0-9]+}/update/marks/{classId:[0-9]+}", h.UpdateMarks)

		r.Get("/{id:[0-9]+}/detail", h.DetailStudent)

		r.Get("/{id:[0-9]+}/result", h.ShowResult)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(h.Authenticator)

		r.Get("/options", h.OPTIONS)

		r.Get("/create/class", h.CreateClass)

		r.Get("/classList", h.ListOfClass)

		r.Post("/store/class", h.StoreClass)

		r.Get("/create/subject", h.CreateSubject)

		r.Post("/store/subject", h.StoreSubject)

		r.Get("/create/student", h.CreateStudent)

		r.Post("/store/student", h.StoreStudent)

		r.Get("/{id:[0-9]+}/edit/student", h.EditStudent)

		r.Put("/{id:[0-9]+}/update/student", h.UpdateStudent)

		r.Get("/show/class/{id:[0-9]+}/result", h.GetResultSheet)

	})

	r.Group(func(r chi.Router) {
		r.Use(h.Authenticator)
		r.Get("/logout", h.Logout)
	})
	return r
}

func Method(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := h.sessionManager.GetString(r.Context(), "userId")
		uid, err := strconv.Atoi(userId)

		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if uid <= 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) ParseTemplates() error {
	templates := template.New("my-templates").Funcs(template.FuncMap{
		"calculatePrevPage": func(currentPageNumber int) int {
			if currentPageNumber == 1 {
				return 0
			}
			return currentPageNumber - 1
		},
		"calculateNextPage": func(currentPageNumber, totalPage int) int {
			if currentPageNumber == totalPage {
				return 0
			}
			return currentPageNumber + 1
		},
	}).Funcs(sprig.FuncMap())

	newFs := os.DirFS("assets/templates")
	tmpl := template.Must(templates.ParseFS(newFs, "*/*/*.html", "*/*.html", "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}
	h.Templates = tmpl
	return nil
}
