package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func OpenDatabase() (db *sql.DB, err error) {
	log.Print(os.Getenv("HAHA_DB_USER"), " !!! ", os.Getenv("HAHA_DB_PASSWORD"))
	dbinfo := fmt.Sprintf("host=127.0.0.1 user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("HAHA_DB_USER"), os.Getenv("HAHA_DB_PASSWORD"), os.Getenv("HAHA_DB_NAME"))
	db, err = sql.Open("postgres", dbinfo)

	return db, err
}
