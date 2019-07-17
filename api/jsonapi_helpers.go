package api

import (
	"log"
	"net/http"
)
// SourceName for error response
var SourceName = "SearchAPI"

// ItemDataGen General
type ItemDataGen struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
}

// JsonapiDataGen jsonapi.org top wrap
type JsonapiDataGen struct {
	Data []*ItemDataGen `json:"data"`
	Meta interface{}    `json:"meta"`
}

// JsonapiError jsonapi.org error item
type JsonapiError struct {
	APIErrors []*JsonapiErrorItem `json:"errors"`
}

// JsonapiErrorItem jsonapi.org error wrap
type JsonapiErrorItem struct {
	Status string `json:"status"`
	Source string `json:"source"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// JsonapiResourcesResp a generic jsonapi multiple resource response with "type" and "id"
func JsonapiResourcesResp(status int, resType, resID string, metaData interface{}, attrs ...interface{}) (int, interface{}) {
	var jsonapiDataItem []*ItemDataGen
	for _, attr := range attrs {
		jsonapiDataItem = append(jsonapiDataItem, &ItemDataGen{
			Type:       resType,
			ID:         resID,
			Attributes: attr,
		})
	}
	return status, JsonapiDataGen{
		Data: jsonapiDataItem,
		Meta: metaData,
	}
}

// JsonapiResourceResp a generic jsonapi resource response with "type" and "id"
func JsonapiResourceResp(status int, resType, resID string, attr interface{}) (int, interface{}) {
	return status, JsonapiDataGen{
		Data: []*ItemDataGen{{
			Type:       resType,
			ID:         resID,
			Attributes: attr,
		}},
	}
}

// JsonapiCustomError single error
func JsonapiCustomError(statusCode int, status, source, title, detail *string) (int, interface{}) {
	statusStr := *status
	sourceStr := *source
	titleStr := *title
	detailStr := *detail

	if status == nil {
		statusStr = http.StatusText(statusCode)
	}
	if source == nil {
		sourceStr = SourceName
	}
	if title == nil {
		titleStr = http.StatusText(statusCode)
	}
	if detail == nil {
		detailStr = http.StatusText(statusCode)
	}
	return statusCode, JsonapiError{
		APIErrors: []*JsonapiErrorItem{
			{Status: statusStr,
				Source: sourceStr,
				Title:  titleStr,
				Detail: detailStr}},
	}
}

// JsonapiErrorResp single error json api body
func JsonapiErrorResp(status int, msg string) (int, interface{}) {
	log.Printf("%s : API error response, status %d error: %s \n", SourceName, status, msg)
	return status, JsonapiError{
		APIErrors: []*JsonapiErrorItem{
			{
				Status: http.StatusText(status),
				Source: SourceName,
				Title:  http.StatusText(status),
				Detail: msg,
			},
		},
	}

}

// JsonapiErrorsResp multiple error json api body
func JsonapiErrorsResp(status int, msges ...string) (int, interface{}) {
	var jsonapiErrItem []*JsonapiErrorItem
	for _, msg := range msges {
		jsonapiErrItem = append(jsonapiErrItem, &JsonapiErrorItem{
			Status: http.StatusText(status),
			Source: SourceName,
			Title:  http.StatusText(status),
			Detail: msg})
		log.Printf("%s : API error response, status %d error: %s \n", SourceName, status, msg)
	}
	return status, JsonapiError{
		APIErrors: jsonapiErrItem,
	}

}
