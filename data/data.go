package data

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", "postgres://nihuo:linux123@localhost:5432/?dbname=chitchat&sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func CreateUUID() (string) {
	return uuid.NewString()
}