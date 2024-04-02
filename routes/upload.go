package routes

import (
	"file-manager/background"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type BindFile struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func UploadRoute(c *gin.Context) {
	var bindFile BindFile

	// Bind file
	if err := c.ShouldBind(&bindFile); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// Save uploaded file
	file := bindFile.File
	dst := filepath.Join("uploaded", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	go background.ProcessFile(dst)
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}
