package api

import (
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// JsonapiSearchResultRaw es to jsonapi.org format
// Elasticsearch already responding with json, and we have so far no much changes in that.
// so don't want to decode and encode json to give jsonapi.org format.
// here the solution is done without unmarshaling, using gjson and sjson
func JsonapiSearchResultRaw(result string) (string, error) {
	var i int64
	dataStr := `{"data":[], "meta":{ "total":true } }`
	itemStr := `{"type": "article", "id": true, "attributes": true }`
	total := gjson.Get(result, "hits.total.value")
	dataStr, _ = sjson.SetRaw(dataStr, "meta.total", total.Raw)
	for i = 0; i < total.Int(); i++ {
		itemStr, _ = sjson.SetRaw(itemStr, "attributes", gjson.Get(result, fmt.Sprintf("hits.hits.%v._source", i)).Raw)
		itemStr, _ = sjson.SetRaw(itemStr, "id", gjson.Get(result, fmt.Sprintf("hits.hits.%v._id", i)).Raw)
		dataStr, _ = sjson.SetRaw(dataStr, "data.-1", itemStr)
	}
	return dataStr, nil
}
