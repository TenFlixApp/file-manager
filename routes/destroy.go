package routes

import (
	"file-manager/data"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DestroyRoute(c *gin.Context) {
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

	data.DeleteFileMetadata(parsedId)

	err = os.Remove(filepath.Join("uploaded", "video", id))
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to delete file")
		return
	}
	err = os.Remove(filepath.Join("uploaded", "cover", id))
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to delete file")
		return
	}
}
