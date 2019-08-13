package search

import (
	"log"

	"github.com/spf13/viper"
)

// FilterFilePath path to config files
var FilterFilePath string

// GetFiltersFromES get search filters from Elasticsearch for clients to search
func GetFiltersFromES() {

}

// GetFiltersFromConfig read filters from yaml file
func GetFiltersFromConfig() interface{} {
	viper.SetConfigName("categories")
	viper.AddConfigPath(FilterFilePath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println("Fatal error config file:", err.Error())
	}
	viper.SetConfigName("locations")
	viper.AddConfigPath(FilterFilePath)
	err2 := viper.MergeInConfig()
	if err2 != nil {
		log.Println("Fatal error config file:", err2.Error())
	}

	//return viper.Get("category")
	return viper.AllSettings()

}
