package api

import (
	"net/http"
	"strconv"
)

func handleRead(c *gin.Context) (int, interface{}) {
	var logData StructuredLog
	
	category := c.Query("category")
	region := c.Query("region")
	area := c.Query("area")
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	sortBy := c.Query("sort")
	searchText := c.Query("q")
	logData.ID = ""
	logData.Account = ""
	if offsetStr != "" {
		offset, errOff := strconv.ParseInt(offsetStr, 10, 64)
		if errOff != nil {
			LogEvent(logData, "error", "OffsetError", errOff.Error())
			//return JsonapiErrorResp(http.StatusNotAcceptable, "Request offset not compatible")
			offset = 0
		}
	}
	if limit != "" {
		limit, errLim := strconv.ParseInt(limitStr, 10, 64)
		if errLim != nil {
			LogEvent(logData, "error", "LimitError", errLim.Error())
			return JsonapiErrorResp(http.StatusNotAcceptable, "Request limit not compatible")
		}
	}
	LogEvent(logData, "info", "ReadRequestReceived", "Request received to read by userid")
	return readDoc(user, offset, limit)
}
