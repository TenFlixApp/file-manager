package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

	// Save uploaded video
	file := bindFile.File
	dst := filepath.Join("uploaded", "video", filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	// Save uploaded cover
	cover := bindFile.File
	dst = filepath.Join("uploaded", "cover", filepath.Base(cover.Filename))
	if err := c.SaveUploadedFile(cover, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	//TODO: Hydrate DB
	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully.", file.Filename))
}
