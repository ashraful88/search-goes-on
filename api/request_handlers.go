package api

import (
	"log"
	"strconv"

	"github.com/ashraful88/search-goes/search"
	"github.com/gin-gonic/gin"
)

func handleGetFilters(c *gin.Context) (int, interface{}) {
	category := c.Query("category")
	if category == "" {
		search.GetSearchClient()
		return 200, search.GetFiltersFromConfig()
	}

	return 200, search.GetFiltersFromConfig()

}

func handleSearchAds(c *gin.Context) (int, interface{}) {
	var logData StructuredLog

	/*area := c.Query("area")
	searchText := c.Query("q")
	sortBy := c.Query("sort") */
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

	q := c.Query("q")
	category := c.Query("category")
	region := c.Query("region")
	area := c.Query("area")
	//region := c.Request.URL.Query().Get("region")
	qry := c.Request.URL.Query()

	log.Println(category, "category")
	log.Println(region, "region")
	log.Println(area, "area")
	log.Println(offsetStr, "off")
	log.Println(qry, "qry")

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
	result := search.QuerySearch(q, category, region)
	return 200, result
}
