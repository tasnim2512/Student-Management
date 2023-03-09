package main

import (
	"fmt"
	"log"
	"net/http"
	"practice/json-golang/handler"
	"practice/json-golang/storage/postgres"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
)

var sessionManager *scs.SessionManager

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("%v", err)
	}
	
	p := config.GetInt("server.port")
	decoder := form.NewDecoder()

	postGresStorage, err := postgres.NewPostgresStorage(config)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalln(err)
	}
	if err := goose.Up(postGresStorage.DB.DB, "migrations"); err != nil {
		log.Fatalln(err)
	}

	lt := config.GetDuration("session.lifetime")
	it := config.GetDuration("session.idletime")
	sessionManager = scs.New()
	sessionManager.Lifetime = lt * time.Hour
	sessionManager.IdleTimeout = it * time.Minute
	sessionManager.Cookie.Name = "student_session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = true
	sessionManager.Store = NewSqlxStore(postGresStorage.DB)

	chi := handler.NewHandler(sessionManager, decoder, postGresStorage)
	newChi := nosurf.New(chi)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", p), sessionManager.LoadAndSave(newChi)); err != nil {
		log.Fatalf("%#v", err)
	}
}
