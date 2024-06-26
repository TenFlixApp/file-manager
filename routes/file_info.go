package routes

import (
	"file-manager/data"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func RandomFileRoute(c *gin.Context) {
	count, success := c.GetQuery("count")

	if !success {
		c.String(http.StatusBadRequest, "id is required")
		return
	}

	parsedCount, err := strconv.Atoi(count)
	if err != nil {
		c.String(http.StatusBadRequest, "count must be a valid integer")
		return
	}

	typ, success := c.GetQuery("type")
	if !success {
		typ = "media"
	}

	files := data.GetRandomFileMetadata(parsedCount, typ)
	var payload = make([]gin.H, len(files))
	for i, file := range files {
		payload[i] = gin.H{
			"id":   file.ID,
			"type": file.Type.Name,
			"_links": gin.H{
				"stream": "/files/" + file.ID.String() + "/stream",
				"cover":  "/files/" + file.ID.String() + "/cover",
			},
		}
	}

	c.JSON(http.StatusOK, gin.H{"medias": payload})
}
