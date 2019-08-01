package search

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

// EngineClient instance of search engine
var EngineClient EngineConfig

//EngineConfig client, name/alias of elsticsearch index in our cluster
type EngineConfig struct {
	Client            *elasticsearch.Client
	AdsIndexName      string
	CategoryIndexName string
	RegionIndexName   string
}

// OpenElasticSearchConnection create new ES client
func OpenElasticSearchConnection(esAddr, ads, cat, region string) *EngineConfig {
	var r map[string]interface{}
	cfg := elasticsearch.Config{
		Addresses: []string{
			esAddr,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
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

	EngineClient.Client = es
	EngineClient.AdsIndexName = ads
	EngineClient.CategoryIndexName = cat
	EngineClient.RegionIndexName = region

	return &EngineClient
}

// GetSearchClient get search engine instance
func GetSearchClient() *EngineConfig {
	return &EngineClient
}
