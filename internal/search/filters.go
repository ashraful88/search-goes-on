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

// GetFilterConfig return viper
func GetFilterConfig() map[string]interface{} {
	return viper.AllSettings()
}

// GetFilterConfigByName return viper
func GetFilterConfigByName(name string) map[string]interface{} {
	all := viper.AllSettings()
	return map[string]interface{}{
		name: all[name],
	}
}

// FilterMetadata represents metadata for a single filter
type FilterMetadata struct {
	Name    string                   `json:"name"`
	Label   string                   `json:"label"`
	Type    string                   `json:"type"`
	Field   string                   `json:"field"`
	Multi   bool                     `json:"multi,omitempty"`
	Options []map[string]interface{} `json:"options,omitempty"`
}

// GetFiltersMetadata returns filter metadata for the API response
func GetFiltersMetadata() []FilterMetadata {
	filtersConfig := viper.GetStringMap("filters")
	metadata := []FilterMetadata{}

	for filterName, filterConfigRaw := range filtersConfig {
		filterConfig, ok := filterConfigRaw.(map[string]interface{})
		if !ok {
			continue
		}

		filter := FilterMetadata{
			Name:  filterName,
			Label: getStringValue(filterConfig, "label"),
			Type:  getStringValue(filterConfig, "type"),
			Field: getStringValue(filterConfig, "field"),
		}

		// Check if multi-select is enabled
		if multi, ok := filterConfig["multi"].(bool); ok {
			filter.Multi = multi
		}

		// Add range options if available
		if filter.Type == "range" {
			if rangesRaw, ok := filterConfig["ranges"]; ok {
				if ranges, ok := rangesRaw.([]interface{}); ok {
					options := []map[string]interface{}{}
					for _, rangeRaw := range ranges {
						if rangeMap, ok := rangeRaw.(map[string]interface{}); ok {
							options = append(options, rangeMap)
						}
					}
					filter.Options = options
				}
			}
		}

		metadata = append(metadata, filter)
	}

	return metadata
}

// Helper function to get string value from map
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
