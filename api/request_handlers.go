package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleRead(c *gin.Context) (int, interface{}) {
	var logData StructuredLog

	//category := c.Query("category")
	//region := c.Query("region")
	//area := c.Query("area")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	//sortBy := c.Query("sort")
	//searchText := c.Query("q")
	logData.ID = ""
	logData.Account = ""
	if offsetStr != "" {
		offset, errOff := strconv.ParseInt(offsetStr, 10, 64)
		if errOff != nil {
			LogEvent(logData, "error", "OffsetError", errOff.Error())
			//return JsonapiErrorResp(http.StatusNotAcceptable, "Request offset not compatible")
			offset = 0
		}
		log.Println(offset)
	}
	if limitStr != "" {
		limit, errLim := strconv.ParseInt(limitStr, 10, 64)
		if errLim != nil {
			LogEvent(logData, "error", "LimitError", errLim.Error())
			return JsonapiErrorResp(http.StatusNotAcceptable, "Request limit not compatible")
		}
		log.Println(limit)
	}

	LogEvent(logData, "info", "ReadRequestReceived", "Request received to search by params")
	return 200, logData
}
