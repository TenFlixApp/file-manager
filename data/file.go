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

func GetFileMetadata(id uuid.UUID) *File {
	file := &File{}

	err := db.QueryRow(`SELECT id, title FROM files WHERE id = ?`, id).Scan(&file.ID, &file.Title)
	if err != nil {
		log.Fatalf("Unable to get file metadata: %v\n", err)
		return nil
	}

	return file
}

func DeleteFileMetadata(id uuid.UUID) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Unable to start transaction: %v\n", err)
		return
	}

	_, err = tx.Exec(`DELETE FROM files WHERE id = ?`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Fatalf("delete files: unable to rollback: %v", rollbackErr)
		}
		log.Fatal(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Unable to commit transaction: %v\n", err)
		return
	}

	log.Println("Deleted file from database")
}
