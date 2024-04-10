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

type GenericForm struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
	Type string                `form:"type" binding:"required"`
}

type MediaForm struct {
	Video *multipart.FileHeader `form:"video" binding:"required"`
	Cover *multipart.FileHeader `form:"cover" binding:"required"`
}

func UploadMediaRoute(c *gin.Context) {
	var form MediaForm

	// Bind form data
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// Generate new UUID
	id := uuid.New()

	// Save uploaded video
	file := form.Video
	dst := filepath.Join("uploaded", "video", id.String())
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload video err: %s", err.Error()))
		return
	}

	// Save uploaded cover
	cover := form.Cover
	dst = filepath.Join("uploaded", "cover", id.String())
	if err := c.SaveUploadedFile(cover, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload cover err: %s", err.Error()))
		return
	}

	data.CreateFileMetadata(&data.File{
		ID: id,
		Type: data.FileType{
			Name: "media",
		},
	})

	c.String(http.StatusOK, id.String())
}

func UploadGenericRoute(c *gin.Context) {
	var form GenericForm

	// Bind form data
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// Generate new UUID
	id := uuid.New()

	// Save uploaded video
	file := form.File
	dst := filepath.Join("uploaded", form.Type, id.String())
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload video err: %s", err.Error()))
		return
	}

	data.CreateFileMetadata(&data.File{
		ID: id,
		Type: data.FileType{
			Name: form.Type,
		},
	})

	c.String(http.StatusOK, id.String())
}
