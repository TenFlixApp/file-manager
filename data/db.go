package data

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB

func ConnectToDB() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_CONN_STRING"))
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalln("Error closing the database connection")
	}
}
