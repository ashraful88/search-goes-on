package main

import (
	"log"
	"os"

	gpmiddleware "github.com/701search/gin-prometheus-middleware"
	"github.com/ashraful88/search-goes-on/cmd/search-api/internal/api"
	"github.com/ashraful88/search-goes-on/internal/search"
	"github.com/ashraful88/search-goes-on/internal/platform/searchengine"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
	log.SetFlags(0)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	search.FilterFilePath = "./search/config"
	srvPort, hasPort := os.LookupEnv("SERVICE_PORT")
	if hasPort == false {
		log.Fatal("Service port missing")
	}

	adIndex, _ := os.LookupEnv("AD_INDEX")
	catIndex, _ := os.LookupEnv("CATEGORY_INDEX")
	regionIndex, _ := os.LookupEnv("REGION_INDEX")
	esAddr, hasESInfo := os.LookupEnv("ES_ADDRESS")
	if hasESInfo == false {
		log.Fatal("Elasticsearch address missing")
	}
	_ = searchengine.OpenElasticSearchConnection(esAddr, adIndex, catIndex, regionIndex)

	router := gin.New()
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

	log.Println("Listening ", srvPort)
	router.Run(":" + srvPort)
}
