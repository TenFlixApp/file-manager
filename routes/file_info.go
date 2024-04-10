package routes

import (
	"file-manager/data"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func FileInfoRoute(c *gin.Context) {
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

	file := data.GetFileMetadata(parsedId)
	if file == nil {
		c.String(http.StatusNotFound, "file not found")
		return
	}

	var linksPayload gin.H
	if file.Type.Name == "media" {
		linksPayload = gin.H{
			"stream": "/files/" + file.ID.String() + "/stream",
			"cover":  "/files/" + file.ID.String() + "/cover",
		}
	} else {
		linksPayload = gin.H{
			"stream": "/storage/" + file.ID.String(),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     file.ID,
		"type":   file.Type.Name,
		"_links": linksPayload,
	})
}
