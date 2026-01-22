package search

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"strings"

	"github.com/ashraful88/search-goes-on/internal/platform/searchengine"
)

// QuerySearch perform a query search in ES
func QuerySearch(q string, params map[string][]string) (string, error) {
	// Build the Elasticsearch query using the QueryBuilder
	qb := NewQueryBuilder(q, params)
	queryJSON, err := qb.ToJSON()
	if err != nil {
		log.Printf("Error building query: %s", err)
		return "", errors.New("error building search query")
	}

	log.Printf("Elasticsearch Query: %s", queryJSON)

	se := searchengine.GetSearchClient()
	// Perform the search request
	res, err := se.Client.Search(
		se.Client.Search.WithContext(context.Background()),
		se.Client.Search.WithIndex(se.MainIndexName),
		se.Client.Search.WithBody(strings.NewReader(queryJSON)),
		se.Client.Search.WithTrackTotalHits(true),
		se.Client.Search.WithPretty(),
	)

	if err != nil {
		log.Printf("Error executing search: %s", err)
		return "", errors.New("error executing search")
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Elasticsearch error response: %s", res.Status())
		return "", errors.New("elasticsearch returned error status")
	}

	return read(res.Body), nil
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}

const searchAll = `{
	"query" : { "match_all" : {} },
	"size" : 40,
	"sort" : { "date" : "desc"}
	}`
