package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MountRoute func will mount all rest routes
func MountRoute(router *gin.RouterGroup) {
	router.GET("/health", healthCheck)
	router.GET("/search", adSearch)
	router.GET("/filters", searchFilters)

	router.OPTIONS("/*any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "pass"})
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "pass"})
}

func adSearch(c *gin.Context) {
	c.JSON(handleSearchAds(c))
}

func searchFilters(c *gin.Context) {
	c.JSON(handleGetFilters(c))
}
