package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// MountRoute func will mount all rest routes
func MountRoute(router *gin.RouterGroup) {
	router.GET("/health", healthCheck)
	router.GET("/search/", adSearch)

	router.OPTIONS("/*any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "pass"})
	})
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "pass"})
}

func adSearch(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	lastname := c.Query("lastname")
	category := c.Query("category")
	region := c.Request.URL.Query().Get("region")

	log.Println(category, "cat")
	log.Println(region, "reg")
	log.Println(lastname, "lastname")
	log.Println(offsetStr, "off")
	log.Println(limitStr, "limitStr")

	c.JSON(handleSearchAds(c))
}
