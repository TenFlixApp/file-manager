package data

import (
	"github.com/google/uuid"
	"log"
)

type File struct {
	ID    uuid.UUID
	Title string
}

func CreateFileMetadata(file *File) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Unable to start transaction: %v\n", err)
		return
	}

	_, err = tx.Exec(`INSERT INTO files (id, title) VALUES (?, ?)`, file.ID, file.Title)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("insert files: unable to rollback: %v", rollbackErr)
		}
		log.Fatal(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v\n", err)
		return
	}

	log.Println("Inserted file into database")
}
