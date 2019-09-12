package main

import (
	"log"
	"testing"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func TestMain(t *testing.T) {
	es, _ := elasticsearch.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
}
