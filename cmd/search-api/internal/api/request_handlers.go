package api

import (
	"log"
	"net/http"

	"github.com/ashraful88/search-goes-on/internal/search"
	"github.com/gin-gonic/gin"
)

var contentType = "application/json; charset=utf-8"

func handleGetSearchParams(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "pass"})
}

func handleSearch(c *gin.Context) {
	var logData StructuredLog
	logData.ID = ""
	logData.Account = ""

	q := c.Query("q")
	qry := c.Request.URL.Query()
	log.Println(qry, "qry")

	logData.RawInterface = qry
	LogEvent(logData, "info", "ReadRequestReceived", "Request received to search by params")

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
