package search

import (
	"context"
	"log"
	"strings"
)

// QueryES perform a query search in ES
func QueryES(category, region string) interface{} {
	es := GetEsClient()
	//byt := []byte(`{"num":6.13,"strs":["a","b"]}`)
	// Perform the search request.
	log.Println(category, "category")
	log.Println(region, "region")
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("mudahad"),
		es.Search.WithSize(2),
		//es.Search.WithBody(&buf),
		es.Search.WithBody(strings.NewReader(`{
			"query": {
			  "bool": {
				"must": [
				  {"match": {
					"category": "`+category+`"
				  }},
				  {"match": {
					"region": "`+region+`"
				  }}
				]
			  }
			}
		  }`)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	log.Println(res)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	return res.String()
}
