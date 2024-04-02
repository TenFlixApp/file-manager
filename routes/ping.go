package routes

import "github.com/gin-gonic/gin"

func PingRoute(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}
