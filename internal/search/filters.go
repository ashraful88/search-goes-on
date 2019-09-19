package search

import (
	"log"

	"github.com/spf13/viper"
)

// FilterFilePath path to config files
var FilterFilePath string

// LoadConfigFromFile read filters from yaml file
func LoadConfigFromFile() {
	//NOTE: this func should called once
	viper.SetConfigName("search_params")
	viper.AddConfigPath(FilterFilePath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println("Fatal error config file:", err.Error())
	}
}

func mergeConfigByFileName(filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(FilterFilePath)
	err := viper.MergeInConfig()
	if err != nil {
		log.Println("Fatal error config file:", err.Error())
	}
}

//GetFilterConfig return viper
func GetFilterConfig() map[string]interface{} {
	return viper.AllSettings()
}

//GetFilterConfigByName return viper
func GetFilterConfigByName(name string) map[string]interface{} {
	all := viper.AllSettings()
	return map[string]interface{}{
		name: all[name],
	}
}
