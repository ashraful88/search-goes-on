package search

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/tidwall/sjson"
	"github.com/ashraful88/search-goes-on/internal/platform/searchengine"
)

// add must filters in es query
func (eq *ElasticQuery) buildMustFilter(params map[string][]string) {
	// todo: loop and yaml conf. but this one is faster as no loop, will try array switch

}

// add should filter if any
func (eq *ElasticQuery) buildShouldFilter(params map[string][]string) {
	// todo: loop and yaml conf
}

func (eq *ElasticQuery) buildLimits(params map[string][]string) {
	// defaults
	eq.Size = "40"
	eq.From = "0"
	eq.Sort = `{ "date" : "desc"}`

	// override defaults if provided
	if params["limit"] != nil {
		eq.Size = params["limit"][0]
	}
	if params["from"] != nil {
		eq.From = params["from"][0]
	}
	if params["sort"] != nil {
		if params["order"] != nil {
			eq.Sort = fmt.Sprintf(` { "%s" : "%s" }`, params["sort"], params["order"])
		}
	}
}

// QuerySearch perform a query search in ES
func QuerySearch(q string, params map[string][]string) (string, error) {
	var eq ElasticQuery

	eq.buildMustFilter(params)
	eq.buildShouldFilter(params)
	eq.buildLimits(params)
	sq, _ := sjson.Set("", "query", eq.Query)

	se := searchengine.GetSearchClient()
	// Perform the search request.
	res, err := se.Client.Search(
		se.Client.Search.WithContext(context.Background()),
		se.Client.Search.WithIndex(se.MainIndexName),
		se.Client.Search.WithBody(strings.NewReader(sq)),
		se.Client.Search.WithTrackTotalHits(true),
		se.Client.Search.WithPretty(),
	)
	if res.IsError() {
		log.Println("status error from es")
		return "", errors.New("status error from es")
	}

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return "", errors.New("error getting es result")
	}
	defer res.Body.Close()

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
