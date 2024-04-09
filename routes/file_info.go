package routes

import (
	"file-manager/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func FileInfoRoute(c *gin.Context) {
	// Get file ID from URL
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

	// Get file metadata
	file := data.GetFileMetadata(parsedId)
	if file == nil {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   file.ID,
		"type": file.Type.Name,
		"_links": gin.H{
			"stream": "/files/" + file.ID.String() + "/stream",
			"cover":  "/files/" + file.ID.String() + "/cover",
		},
	})
}
