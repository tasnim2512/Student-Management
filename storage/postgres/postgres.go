package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	NotFound = "sql: no rows in result set"
)

type PostgresStorage struct {
	DB *sqlx.DB
}

func NewPostgresStorage(cfg *viper.Viper) (*PostgresStorage, error) {
	db, err := connectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{
		DB: db,
	}, nil
}

func connectDatabase(cfg *viper.Viper) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s  user=%s dbname=%s sslmode=%s password=%s",
		cfg.GetString("database.host"),
		cfg.GetString("database.port"),
		cfg.GetString("database.username"),
		cfg.GetString("database.dbname"),
		cfg.GetString("database.sslmode"),
		cfg.GetString("database.password"),
	))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}
