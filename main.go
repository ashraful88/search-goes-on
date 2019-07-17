package main

import (
	"log"
	"strings"
	"search-goes/api"

	gpmiddleware "github.com/701search/gin-prometheus-middleware"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
)

func main() {
	log.SetFlags(0)

	/* var (
		r  map[string]interface{}
		wg sync.WaitGroup
	) */
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
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


	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// Prometheus metric setup
	p := gpmiddleware.NewPrometheus("")
	p.Use(router)

	router.Use(func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Max-Age", "10000")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Authorization,Content-Type,Accept")
		context.Next()
	})

	v1 := router.Group("/v1")
	api.MountRoute(v1)
	/* wg.Add(1)
	wg.Done()
	wg.Wait() */

}

