package search

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

// ElasticSearchClient open client for elasticsearch
var ElasticSearchClient *elasticsearch.Client

// OpenElasticSearchConnection create new ES client
func OpenElasticSearchConnection(esAddr string) *elasticsearch.Client {
	var r map[string]interface{}
	cfg := elasticsearch.Config{
		Addresses: []string{
			esAddr,
		},
	}
	es, _ := elasticsearch.NewClient(cfg)
	log.Println(elasticsearch.Version)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	ElasticSearchClient = es

	return es
}

// GetEsClient get elasticsearch client
func GetEsClient() *elasticsearch.Client {
	return ElasticSearchClient
}
