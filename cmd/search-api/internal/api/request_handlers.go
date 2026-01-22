package api

import (
	"log"
	"net/http"

	"github.com/ashraful88/search-goes-on/internal/search"
	"github.com/gin-gonic/gin"
)

var contentType = "application/json; charset=utf-8"

func handleGetSearchParams(c *gin.Context) {
	var logData StructuredLog
	logData.ID = ""
	logData.Account = ""

	LogEvent(logData, "info", "FiltersRequested", "Request received for filter metadata")

	// Get filter metadata from config
	filters := search.GetFiltersMetadata()

	// Return filters in JSON API format
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"type":    "filters",
			"id":      "available-filters",
			"attributes": gin.H{
				"filters": filters,
			},
		},
		"meta": gin.H{
			"count": len(filters),
		},
	})
}

func handleSearch(c *gin.Context) {
	var logData StructuredLog
	logData.ID = ""
	logData.Account = ""

	q := c.Query("q")
	qry := c.Request.URL.Query()
	log.Println(qry, "qry")

	logData.RawInterface = qry
	LogEvent(logData, "info", "SearchRequestReceived", "Request received to search by params")

	result, err := search.QuerySearch(q, qry)
	body, _ := JsonapiSearchResultRaw(result)
	if err != nil {
		LogEvent(logData, "error", "querySearchError", err.Error())
		// send response
		c.JSON(JsonapiErrorResp(500, "Error something went wrong"))
	} else {
		// send response
		c.Data(200, contentType, []byte(body))
	}

}
