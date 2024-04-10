package data

import (
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type FileType struct {
	ID   int
	Name string
}

type File struct {
	ID   uuid.UUID
	Type FileType
}

func getFileType(tx *sql.Tx, name string, create bool) *FileType {
	fileType := &FileType{
		Name: name,
	}

	err := tx.QueryRow(`SELECT id FROM file_types WHERE name = ?`, name).Scan(&fileType.ID)
	if err != nil {
		log.Printf("Unable to get file type: %v\n", err)
		return nil
	}

	if fileType.ID == 0 && create {

		res, err := tx.Exec(`INSERT INTO file_types (name) VALUES (?)`, name)
		if handleSqlError(err, tx) {
			log.Printf("Unable to create file type: %v\n", err)
			return nil
		}

		lastInsertId, err := res.LastInsertId()
		if err != nil {
			log.Println("Unable to get last insert ID")
			return nil
		}

		fileType.ID = int(lastInsertId)
	}

	return fileType
}

func CreateFileMetadata(file *File) bool {
	tx := createTransaction()
	if tx == nil {
		return false
	}

	if file.Type.ID == 0 {
		tempType := getFileType(tx, file.Type.Name, true)
		if tempType == nil {
			handleRollback(tx)
			return false
		}

		file.Type = *tempType
	}

	_, err := tx.Exec(`INSERT INTO files (id, type) VALUES (?, ?)`, file.ID, file.Type.ID)
	if handleSqlError(err, tx) {
		return false
	}

	return !handleTxCommit(tx)
}

func GetFileMetadata(id uuid.UUID) *File {
	file := &File{
		Type: FileType{},
	}

	err := db.QueryRow(`SELECT f.id, ft.id, ft.name FROM files f JOIN file_types ft ON (f.type = ft.id) WHERE f.id = ?`, id).Scan(&file.ID, &file.Type.ID, &file.Type.Name)
	handleSqlError(err, nil)

	return file
}

func GetRandomFileMetadata(cnt int, typ string) []*File {
	var (
		videos = make([]*File, 0)
		query  *sql.Rows
		err    error
	)

	if typ == "" {
		query, err = db.Query(`SELECT f.id, ft.id, ft.name FROM files f JOIN file_types ft ON (f.type = ft.id) ORDER BY RAND() LIMIT ?`, cnt)
	} else {
		query, err = db.Query(`SELECT f.id, ft.id, ft.name FROM files f JOIN file_types ft ON (f.type = ft.id) WHERE ft.name = ? ORDER BY RAND() LIMIT ?`, typ, cnt)
	}

	if err != nil {
		return nil
	}

	for query.Next() {
		file := &File{
			Type: FileType{},
		}
		err := query.Scan(&file.ID, &file.Type.ID, &file.Type.Name)
		if handleSqlError(err, nil) {
			return nil
		}

		videos = append(videos, file)
	}

	return videos
}

func DeleteFileMetadata(id uuid.UUID) bool {
	tx := createTransaction()
	if tx == nil {
		return false
	}

	_, err := tx.Exec(`DELETE FROM files WHERE id = ?`, id)
	if handleSqlError(err, tx) {
		return false
	}

	return !handleTxCommit(tx)
}
