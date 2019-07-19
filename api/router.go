package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MountRoute func will mount all rest routes
func MountRoute(router *gin.RouterGroup) {
	router.GET("/health", healthCheck)
	router.GET("/search", readSearchList)

	router.OPTIONS("/*any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "pass"})
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "pass"})
}

func readSearchList(c *gin.Context) {
	c.JSON(handleRead(c))
}
