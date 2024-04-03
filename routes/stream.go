package routes

import (
	"github.com/gin-gonic/gin"
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
