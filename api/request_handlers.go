package api

import (
	"log"
	"strconv"

	"github.com/ashraful88/search-goes/search"
	"github.com/gin-gonic/gin"
)

func handleRead(c *gin.Context) (int, interface{}) {
	var logData StructuredLog

	/*area := c.Query("area")
	searchText := c.Query("q")
	sortBy := c.Query("sort") */
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	category := c.Query("category")
	region := c.Query("region")
	/* log.Println(category, "cat")
	log.Println(region, "reg")
	log.Println(offsetStr, "off") */

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
	result := search.QueryES(category, region)
	return 200, result
}
