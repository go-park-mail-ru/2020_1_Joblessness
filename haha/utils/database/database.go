package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type databaseSetttings struct {
	user string
	password string
	name string
	db *sql.DB
}

var dbSettings = databaseSetttings{
	user:     "",
	password: "",
	name:     "",
	db: nil,
}

func InitDatabase(user, password, name string) {
	dbSettings = databaseSetttings{
		user: user,
		password: password,
		name: name,
	}
}

func OpenDatabase() (err error){
	if dbSettings.db != nil {
		return errors.New("Close current DB")
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		dbSettings.user, dbSettings.password, dbSettings.name)
	dbSettings.db, err = sql.Open("postgres", dbinfo)

	return err
}

func CloseDatabase() (err error) {
	if dbSettings.db == nil {
		return nil
	}

	err = dbSettings.db.Close()
	dbSettings.db = nil

	return err
}

func GetDatabase() *sql.DB {
	return dbSettings.db
}