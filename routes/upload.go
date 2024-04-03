package routes

import (
	"file-manager/data"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type BindFile struct {
	File  *multipart.FileHeader `form:"file" binding:"required"`
	Cover *multipart.FileHeader `form:"cover" binding:"required"`
}

func UploadRoute(c *gin.Context) {
	var bindFile BindFile

	// Bind file
	if err := c.ShouldBind(&bindFile); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// Generate new UUID
	id := uuid.New()

	// Save uploaded video
	file := bindFile.File
	dst := filepath.Join("uploaded", "video", id.String()+".mp4")
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	// Save uploaded cover
	cover := bindFile.File
	dst = filepath.Join("uploaded", "cover", id.String()+".png")
	if err := c.SaveUploadedFile(cover, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	//TODO: Hydrate DB
	data.CreateFileMetadata(&data.File{
		ID:    id,
		Title: "File Title",
	})
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}
