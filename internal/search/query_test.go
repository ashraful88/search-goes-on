package search

import (
	"encoding/json"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig() {
	// Reset viper
	viper.Reset()

	// Set up test search configuration
	viper.Set("search.fields", []string{"name", "description", "brand"})
	viper.Set("search.default_operator", "OR")
	viper.Set("search.minimum_should_match", 1)

	// Set up test filter configuration
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
		},
	})

	viper.Set("filters.in_stock", map[string]interface{}{
		"label": "In Stock",
		"field": "in_stock",
		"type":  "term",
		"multi": false,
	})

	// Set up test aggregation configuration
	viper.Set("aggregations.category", map[string]interface{}{
		"type":  "terms",
		"field": "category.keyword",
		"size":  50,
	})

	viper.Set("aggregations.price_ranges", map[string]interface{}{
		"type":  "range",
		"field": "price",
		"ranges": []interface{}{
			map[string]interface{}{"key": "0-50", "from": 0, "to": 50},
			map[string]interface{}{"key": "50-100", "from": 50, "to": 100},
		},
	})
}

func TestNewQueryBuilder(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"category": {"electronics"},
		"limit":    {"20"},
	}

	qb := NewQueryBuilder("laptop", params)

	assert.NotNil(t, qb)
	assert.Equal(t, "laptop", qb.QueryString)
	assert.Equal(t, params, qb.Params)
	assert.Equal(t, 40, qb.Size)
	assert.Equal(t, 0, qb.From)
	assert.NotNil(t, qb.Query)
	assert.NotNil(t, qb.Aggs)
	assert.NotNil(t, qb.Sort)
}

func TestBuildQueryContext_WithQueryString(t *testing.T) {
	setupTestConfig()

	qb := NewQueryBuilder("laptop", map[string][]string{})
	must := qb.buildQueryContext()

	assert.Len(t, must, 1)

	// Verify multi_match query structure
	mustMap, ok := must[0].(map[string]interface{})
	assert.True(t, ok)

	multiMatch, ok := mustMap["multi_match"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "laptop", multiMatch["query"])
	assert.Equal(t, "OR", multiMatch["operator"])
	assert.Equal(t, "best_fields", multiMatch["type"])

	fields, ok := multiMatch["fields"].([]string)
	assert.True(t, ok)
	assert.Contains(t, fields, "name")
	assert.Contains(t, fields, "description")
	assert.Contains(t, fields, "brand")
}

func TestBuildQueryContext_EmptyQueryString(t *testing.T) {
	setupTestConfig()

	qb := NewQueryBuilder("", map[string][]string{})
	must := qb.buildQueryContext()

	assert.Len(t, must, 0)
}

func TestBuildFilterContext_TermsFilter(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"category": {"electronics", "computers"},
		"brand":    {"apple"},
	}

	qb := NewQueryBuilder("", params)
	filters := qb.buildFilterContext()

	assert.Len(t, filters, 2)

	// Find category filter
	var categoryFilter map[string]interface{}
	for _, f := range filters {
		if filterMap, ok := f.(map[string]interface{}); ok {
			if terms, ok := filterMap["terms"].(map[string]interface{}); ok {
				if _, hasCategory := terms["category.keyword"]; hasCategory {
					categoryFilter = filterMap
					break
				}
			}
		}
	}

	assert.NotNil(t, categoryFilter)
	terms := categoryFilter["terms"].(map[string]interface{})
	categoryValues := terms["category.keyword"].([]string)
	assert.Contains(t, categoryValues, "electronics")
	assert.Contains(t, categoryValues, "computers")
}

func TestBuildFilterContext_TermFilter(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"in_stock": {"true"},
	}

	qb := NewQueryBuilder("", params)
	filters := qb.buildFilterContext()

	assert.Len(t, filters, 1)

	filterMap, ok := filters[0].(map[string]interface{})
	assert.True(t, ok)

	term, ok := filterMap["term"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "true", term["in_stock"])
}

func TestBuildFilterContext_RangeFilter(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"price": {"50-100"},
	}

	qb := NewQueryBuilder("", params)
	filters := qb.buildFilterContext()

	assert.Len(t, filters, 1)

	filterMap, ok := filters[0].(map[string]interface{})
	assert.True(t, ok)

	rangeQuery, ok := filterMap["range"].(map[string]interface{})
	assert.True(t, ok)

	priceRange, ok := rangeQuery["price"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, 50, priceRange["gte"])
	assert.Equal(t, 100, priceRange["lt"])
}

func TestBuildPaginationAndSort(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"limit": {"20"},
		"from":  {"10"},
		"sort":  {"price"},
		"order": {"asc"},
	}

	qb := NewQueryBuilder("", params)
	qb.buildPaginationAndSort()

	assert.Equal(t, 20, qb.Size)
	assert.Equal(t, 10, qb.From)
	assert.Len(t, qb.Sort, 1)

	sortMap := qb.Sort[0]
	priceSort, ok := sortMap["price"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "asc", priceSort["order"])
}

func TestBuildPaginationAndSort_Defaults(t *testing.T) {
	setupTestConfig()

	qb := NewQueryBuilder("", map[string][]string{})
	qb.buildPaginationAndSort()

	assert.Equal(t, 40, qb.Size)
	assert.Equal(t, 0, qb.From)
	assert.Len(t, qb.Sort, 1)

	// Should have default sort by created_at
	sortMap := qb.Sort[0]
	_, hasCreatedAt := sortMap["created_at"]
	assert.True(t, hasCreatedAt)
}

func TestBuildAggregations(t *testing.T) {
	setupTestConfig()

	qb := NewQueryBuilder("", map[string][]string{})
	qb.buildAggregations()

	assert.Len(t, qb.Aggs, 2)

	// Check category aggregation
	categoryAgg, ok := qb.Aggs["category"].(map[string]interface{})
	assert.True(t, ok)

	terms, ok := categoryAgg["terms"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "category.keyword", terms["field"])
	assert.Equal(t, 50, terms["size"])

	// Check price_ranges aggregation
	priceAgg, ok := qb.Aggs["price_ranges"].(map[string]interface{})
	assert.True(t, ok)

	rangeAgg, ok := priceAgg["range"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "price", rangeAgg["field"])

	ranges, ok := rangeAgg["ranges"].([]map[string]interface{})
	assert.True(t, ok)
	assert.Len(t, ranges, 2)
}

func TestBuild_CompleteQuery(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"category": {"electronics"},
		"brand":    {"apple", "samsung"},
		"price":    {"50-100"},
		"limit":    {"10"},
		"from":     {"0"},
		"sort":     {"price"},
		"order":    {"asc"},
	}

	qb := NewQueryBuilder("laptop", params)
	result := qb.Build()

	// Verify top-level structure
	assert.NotNil(t, result["query"])
	assert.Equal(t, 10, result["size"])
	assert.Equal(t, 0, result["from"])
	assert.NotNil(t, result["sort"])
	assert.NotNil(t, result["aggs"])

	// Verify query structure
	query, ok := result["query"].(map[string]interface{})
	assert.True(t, ok)

	boolQuery, ok := query["bool"].(map[string]interface{})
	assert.True(t, ok)

	// Verify must clause (query context)
	must, ok := boolQuery["must"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, must, 1)

	// Verify filter clause (filter context)
	filter, ok := boolQuery["filter"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, filter, 3) // category, brand, price
}

func TestBuild_MatchAll(t *testing.T) {
	setupTestConfig()

	qb := NewQueryBuilder("", map[string][]string{})
	result := qb.Build()

	query, ok := result["query"].(map[string]interface{})
	assert.True(t, ok)

	_, hasMatchAll := query["match_all"]
	assert.True(t, hasMatchAll)
}

func TestToJSON(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"category": {"electronics"},
	}

	qb := NewQueryBuilder("laptop", params)
	jsonStr, err := qb.ToJSON()

	assert.NoError(t, err)
	assert.NotEmpty(t, jsonStr)

	// Verify it's valid JSON
	var result map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &result)
	assert.NoError(t, err)
	assert.NotNil(t, result["query"])
}

func TestGetSearchFields(t *testing.T) {
	setupTestConfig()

	qb := NewQueryBuilder("", map[string][]string{})
	fields := qb.getSearchFields()

	assert.Len(t, fields, 3)
	assert.Contains(t, fields, "name")
	assert.Contains(t, fields, "description")
	assert.Contains(t, fields, "brand")
}

func TestGetSearchFields_Fallback(t *testing.T) {
	viper.Reset()

	qb := NewQueryBuilder("", map[string][]string{})
	fields := qb.getSearchFields()

	assert.Len(t, fields, 1)
	assert.Equal(t, "_all", fields[0])
}

func TestGetStringFromMap(t *testing.T) {
	m := map[string]interface{}{
		"string": "value",
		"number": 42,
		"bool":   true,
	}

	assert.Equal(t, "value", getStringFromMap(m, "string"))
	assert.Equal(t, "42", getStringFromMap(m, "number"))
	assert.Equal(t, "true", getStringFromMap(m, "bool"))
	assert.Equal(t, "", getStringFromMap(m, "missing"))
}

func TestBuildRangeFilter_InvalidRange(t *testing.T) {
	setupTestConfig()

	filterConfig := map[string]interface{}{
		"ranges": []interface{}{
			map[string]interface{}{"key": "0-50", "from": 0, "to": 50},
		},
	}

	qb := NewQueryBuilder("", map[string][]string{})
	result := qb.buildRangeFilter("price", "invalid-range", filterConfig)

	assert.Nil(t, result)
}

func TestBuildRangeFilter_OpenEndedRange(t *testing.T) {
	setupTestConfig()

	filterConfig := map[string]interface{}{
		"ranges": []interface{}{
			map[string]interface{}{"key": "200+", "from": 200},
		},
	}

	qb := NewQueryBuilder("", map[string][]string{})
	result := qb.buildRangeFilter("price", "200+", filterConfig)

	assert.NotNil(t, result)

	rangeQuery := result["range"].(map[string]interface{})
	priceRange := rangeQuery["price"].(map[string]interface{})
	assert.Equal(t, 200, priceRange["gte"])
	_, hasTo := priceRange["lt"]
	assert.False(t, hasTo)
}

func TestMultipleSortFields(t *testing.T) {
	setupTestConfig()

	params := map[string][]string{
		"sort":  {"price,rating"},
		"order": {"asc"},
	}

	qb := NewQueryBuilder("", params)
	qb.buildPaginationAndSort()

	assert.Len(t, qb.Sort, 2)

	priceSort := qb.Sort[0]["price"].(map[string]interface{})
	assert.Equal(t, "asc", priceSort["order"])

	ratingSort := qb.Sort[1]["rating"].(map[string]interface{})
	assert.Equal(t, "asc", ratingSort["order"])
}
