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
		log.Fatal("Unable to create DB handle", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the DB", err)
	}
	log.Println("Connected to the database")
}

func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalln("Error closing the database connection")
	}
}

func createTransaction() *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Unable to start transaction: %v\n", err)
		return nil
	}
	return tx
}

func handleRollback(tx *sql.Tx) bool {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		log.Printf("insert files: unable to rollback: %v\n", rollbackErr)
		return false
	}
	return true
}

func handleTxCommit(tx *sql.Tx) bool {
	err := tx.Commit()
	if err != nil {
		log.Printf("Unable to commit transaction: %v\n", err)
		return false
	}
	return true
}

func handleSqlError(err error, tx *sql.Tx) bool {
	if err != nil {
		if tx != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("insert files: unable to rollback: %v\n", rollbackErr)
			}
		}
		log.Println(err)
	}
	return err != nil
}
