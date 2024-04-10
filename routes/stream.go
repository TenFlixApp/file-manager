package routes

import (
	"file-manager/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
)

func StreamVideo(c *gin.Context) {
	id, success := c.Params.Get("id")

	if !success {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	file := filepath.Join("uploaded", "video", id)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	c.Writer.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(c.Writer, c.Request, file)
}

func StreamCover(c *gin.Context) {
	id, success := c.Params.Get("id")

	if !success {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	file := filepath.Join("uploaded", "cover", id)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	c.Writer.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	http.ServeFile(c.Writer, c.Request, file)
}

func StreamGeneric(c *gin.Context) {
	id, success := c.Params.Get("id")

	if !success {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		c.String(http.StatusBadRequest, "id must be a valid UUID")
		return
	}

	metadata := data.GetFileMetadata(parsedId)
	if metadata == nil {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	file := filepath.Join("uploaded", metadata.Type.Name, id)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	c.Writer.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	http.ServeFile(c.Writer, c.Request, file)
}
