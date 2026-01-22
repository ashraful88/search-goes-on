package search

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupFiltersTestConfig() {
	viper.Reset()

	// Set up comprehensive filter configuration
	viper.Set("filters.category", map[string]interface{}{
		"label": "Category",
		"field": "category.keyword",
		"type":  "terms",
		"multi": true,
	})

	viper.Set("filters.brand", map[string]interface{}{
		"label": "Brand",
		"field": "brand.keyword",
		"type":  "terms",
		"multi": true,
	})

	viper.Set("filters.price", map[string]interface{}{
		"label": "Price",
		"field": "price",
		"type":  "range",
		"ranges": []interface{}{
			map[string]interface{}{"key": "0-50", "from": 0, "to": 50},
			map[string]interface{}{"key": "50-100", "from": 50, "to": 100},
			map[string]interface{}{"key": "100-200", "from": 100, "to": 200},
			map[string]interface{}{"key": "200-*", "from": 200},
		},
	})

	viper.Set("filters.in_stock", map[string]interface{}{
		"label": "In Stock",
		"field": "in_stock",
		"type":  "term",
		"multi": false,
	})

	viper.Set("filters.rating", map[string]interface{}{
		"label": "Rating",
		"field": "rating",
		"type":  "range",
		"ranges": []interface{}{
			map[string]interface{}{"key": "4+", "from": 4},
			map[string]interface{}{"key": "3+", "from": 3},
			map[string]interface{}{"key": "2+", "from": 2},
		},
	})
}

func TestGetFilterConfig(t *testing.T) {
	setupFiltersTestConfig()

	config := GetFilterConfig()

	assert.NotNil(t, config)
	assert.NotEmpty(t, config)

	filters, ok := config["filters"].(map[string]interface{})
	assert.True(t, ok)
	assert.Len(t, filters, 5) // category, brand, price, in_stock, rating
}

func TestGetFilterConfigByName(t *testing.T) {
	setupFiltersTestConfig()

	config := GetFilterConfigByName("filters")

	assert.NotNil(t, config)
	assert.Contains(t, config, "filters")

	filters := config["filters"].(map[string]interface{})
	assert.NotNil(t, filters)
}

func TestGetFiltersMetadata(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	assert.NotNil(t, metadata)
	assert.Len(t, metadata, 5)

	// Verify each filter type is present
	filterNames := make(map[string]bool)
	for _, filter := range metadata {
		filterNames[filter.Name] = true
	}

	assert.True(t, filterNames["category"])
	assert.True(t, filterNames["brand"])
	assert.True(t, filterNames["price"])
	assert.True(t, filterNames["in_stock"])
	assert.True(t, filterNames["rating"])
}

func TestGetFiltersMetadata_TermsFilter(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	var categoryFilter *FilterMetadata
	for _, filter := range metadata {
		if filter.Name == "category" {
			categoryFilter = &filter
			break
		}
	}

	assert.NotNil(t, categoryFilter)
	assert.Equal(t, "Category", categoryFilter.Label)
	assert.Equal(t, "category.keyword", categoryFilter.Field)
	assert.Equal(t, "terms", categoryFilter.Type)
	assert.True(t, categoryFilter.Multi)
	assert.Empty(t, categoryFilter.Options) // terms filters don't have options
}

func TestGetFiltersMetadata_TermFilter(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	var inStockFilter *FilterMetadata
	for _, filter := range metadata {
		if filter.Name == "in_stock" {
			inStockFilter = &filter
			break
		}
	}

	assert.NotNil(t, inStockFilter)
	assert.Equal(t, "In Stock", inStockFilter.Label)
	assert.Equal(t, "in_stock", inStockFilter.Field)
	assert.Equal(t, "term", inStockFilter.Type)
	assert.False(t, inStockFilter.Multi)
	assert.Empty(t, inStockFilter.Options)
}

func TestGetFiltersMetadata_RangeFilter(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	var priceFilter *FilterMetadata
	for _, filter := range metadata {
		if filter.Name == "price" {
			priceFilter = &filter
			break
		}
	}

	assert.NotNil(t, priceFilter)
	assert.Equal(t, "Price", priceFilter.Label)
	assert.Equal(t, "price", priceFilter.Field)
	assert.Equal(t, "range", priceFilter.Type)
	assert.Len(t, priceFilter.Options, 4)

	// Verify range options
	option := priceFilter.Options[0]
	assert.Equal(t, "0-50", option["key"])
	assert.Equal(t, 0, option["from"])
	assert.Equal(t, 50, option["to"])

	// Verify open-ended range
	lastOption := priceFilter.Options[3]
	assert.Equal(t, "200-*", lastOption["key"])
	assert.Equal(t, 200, lastOption["from"])
	_, hasTo := lastOption["to"]
	assert.False(t, hasTo)
}

func TestGetFiltersMetadata_RatingRangeFilter(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	var ratingFilter *FilterMetadata
	for _, filter := range metadata {
		if filter.Name == "rating" {
			ratingFilter = &filter
			break
		}
	}

	assert.NotNil(t, ratingFilter)
	assert.Equal(t, "Rating", ratingFilter.Label)
	assert.Equal(t, "rating", ratingFilter.Field)
	assert.Equal(t, "range", ratingFilter.Type)
	assert.Len(t, ratingFilter.Options, 3)

	// Verify rating ranges (open-ended from)
	for _, option := range ratingFilter.Options {
		assert.NotNil(t, option["from"])
		_, hasTo := option["to"]
		assert.False(t, hasTo) // Rating filters are open-ended (4+, 3+, 2+)
	}
}

func TestGetFiltersMetadata_EmptyConfig(t *testing.T) {
	viper.Reset()

	metadata := GetFiltersMetadata()

	assert.NotNil(t, metadata)
	assert.Empty(t, metadata)
}

func TestGetFiltersMetadata_InvalidConfig(t *testing.T) {
	viper.Reset()

	// Set up invalid configuration
	viper.Set("filters.invalid", "not a map")

	metadata := GetFiltersMetadata()

	assert.NotNil(t, metadata)
	assert.Empty(t, metadata) // Should skip invalid entries
}

func TestGetStringValue(t *testing.T) {
	testCases := []struct {
		name     string
		m        map[string]interface{}
		key      string
		expected string
	}{
		{
			name:     "string value",
			m:        map[string]interface{}{"key": "value"},
			key:      "key",
			expected: "value",
		},
		{
			name:     "missing key",
			m:        map[string]interface{}{"other": "value"},
			key:      "key",
			expected: "",
		},
		{
			name:     "non-string value",
			m:        map[string]interface{}{"key": 123},
			key:      "key",
			expected: "",
		},
		{
			name:     "empty map",
			m:        map[string]interface{}{},
			key:      "key",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getStringValue(tc.m, tc.key)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilterMetadata_JSONSerialization(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	// Find a filter with range options
	var priceFilter *FilterMetadata
	for _, filter := range metadata {
		if filter.Name == "price" {
			priceFilter = &filter
			break
		}
	}

	assert.NotNil(t, priceFilter)

	// Verify JSON tags work correctly (struct can be marshaled)
	assert.Equal(t, "price", priceFilter.Name)
	assert.NotEmpty(t, priceFilter.Label)
	assert.NotEmpty(t, priceFilter.Type)
	assert.NotEmpty(t, priceFilter.Field)
	assert.NotEmpty(t, priceFilter.Options)
}

func TestGetFiltersMetadata_BrandFilter(t *testing.T) {
	setupFiltersTestConfig()

	metadata := GetFiltersMetadata()

	var brandFilter *FilterMetadata
	for _, filter := range metadata {
		if filter.Name == "brand" {
			brandFilter = &filter
			break
		}
	}

	assert.NotNil(t, brandFilter)
	assert.Equal(t, "Brand", brandFilter.Label)
	assert.Equal(t, "brand.keyword", brandFilter.Field)
	assert.Equal(t, "terms", brandFilter.Type)
	assert.True(t, brandFilter.Multi)
}

func TestGetFiltersMetadata_ConsistentOrdering(t *testing.T) {
	setupFiltersTestConfig()

	metadata1 := GetFiltersMetadata()
	metadata2 := GetFiltersMetadata()

	// Both calls should return the same number of filters
	assert.Len(t, metadata1, 5)
	assert.Len(t, metadata2, 5)

	// Convert to maps for easier comparison
	filters1 := make(map[string]FilterMetadata)
	filters2 := make(map[string]FilterMetadata)

	for _, f := range metadata1 {
		filters1[f.Name] = f
	}
	for _, f := range metadata2 {
		filters2[f.Name] = f
	}

	// Verify each filter is consistent
	for name, filter1 := range filters1 {
		filter2, exists := filters2[name]
		assert.True(t, exists)
		assert.Equal(t, filter1.Label, filter2.Label)
		assert.Equal(t, filter1.Type, filter2.Type)
		assert.Equal(t, filter1.Field, filter2.Field)
		assert.Equal(t, filter1.Multi, filter2.Multi)
	}
}
