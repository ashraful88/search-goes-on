package api

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleSearchAds(c *gin.Context) (int, interface{}) {
	var logData StructuredLog

	/*area := c.Query("area")
	searchText := c.Query("q")
	sortBy := c.Query("sort") */
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	lastname := c.Query("lastname")
	category := c.Query("category")
	region := c.Request.URL.Query().Get("region")

	log.Println(category, "cat")
	log.Println(region, "reg")
	log.Println(lastname, "lastname")
	log.Println(offsetStr, "off")

	logData.ID = ""
	logData.Account = ""

	offset, errOff := strconv.ParseInt(offsetStr, 10, 64)
	if errOff != nil {
		LogEvent(logData, "error", "OffsetError", errOff.Error())
		offset = 0
	}
	log.Println(offset)

	limit, errLim := strconv.ParseInt(limitStr, 10, 64)
	if errLim != nil {
		LogEvent(logData, "error", "LimitError", errLim.Error())
		limit = 20
	}
	log.Println(limit)

	LogEvent(logData, "info", "ReadRequestReceived", "Request received to search by params")
	//result := search.QueryES(category, region)
	return 200, logData
}
