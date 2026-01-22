package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ashraful88/search-goes-on/internal/search"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupAPITestConfig() {
	viper.Reset()

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
}

func TestHandleGetSearchParams(t *testing.T) {
	setupAPITestConfig()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test router
	router := gin.New()
	router.GET("/v1/filters", handleGetSearchParams)

	// Create a test request
	req, _ := http.NewRequest("GET", "/v1/filters", nil)
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify JSON API structure
	data, ok := response["data"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "filters", data["type"])
	assert.Equal(t, "available-filters", data["id"])

	attributes, ok := data["attributes"].(map[string]interface{})
	assert.True(t, ok)

	filters, ok := attributes["filters"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, filters, 3)

	// Verify meta
	meta, ok := response["meta"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(3), meta["count"]) // JSON numbers are float64
}

func TestHandleGetSearchParams_ContentType(t *testing.T) {
	setupAPITestConfig()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/filters", handleGetSearchParams)

	req, _ := http.NewRequest("GET", "/v1/filters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
}

func TestHandleGetSearchParams_FilterStructure(t *testing.T) {
	setupAPITestConfig()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/filters", handleGetSearchParams)

	req, _ := http.NewRequest("GET", "/v1/filters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	filters := attributes["filters"].([]interface{})

	// Verify each filter has required fields
	for _, filterRaw := range filters {
		filter := filterRaw.(map[string]interface{})
		assert.NotEmpty(t, filter["name"])
		assert.NotEmpty(t, filter["label"])
		assert.NotEmpty(t, filter["type"])
		assert.NotEmpty(t, filter["field"])
	}
}

func TestHandleGetSearchParams_TermsFilter(t *testing.T) {
	setupAPITestConfig()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/filters", handleGetSearchParams)

	req, _ := http.NewRequest("GET", "/v1/filters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	filters := attributes["filters"].([]interface{})

	// Find category filter
	var categoryFilter map[string]interface{}
	for _, filterRaw := range filters {
		filter := filterRaw.(map[string]interface{})
		if filter["name"] == "category" {
			categoryFilter = filter
			break
		}
	}

	assert.NotNil(t, categoryFilter)
	assert.Equal(t, "Category", categoryFilter["label"])
	assert.Equal(t, "terms", categoryFilter["type"])
	assert.Equal(t, "category.keyword", categoryFilter["field"])
	assert.Equal(t, true, categoryFilter["multi"])
}

func TestHandleGetSearchParams_RangeFilter(t *testing.T) {
	setupAPITestConfig()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/filters", handleGetSearchParams)

	req, _ := http.NewRequest("GET", "/v1/filters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	attributes := data["attributes"].(map[string]interface{})
	filters := attributes["filters"].([]interface{})

	// Find price filter
	var priceFilter map[string]interface{}
	for _, filterRaw := range filters {
		filter := filterRaw.(map[string]interface{})
		if filter["name"] == "price" {
			priceFilter = filter
			break
		}
	}

	assert.NotNil(t, priceFilter)
	assert.Equal(t, "Price", priceFilter["label"])
	assert.Equal(t, "range", priceFilter["type"])
	assert.Equal(t, "price", priceFilter["field"])

	// Verify options are present for range filter
	options, ok := priceFilter["options"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, options, 2)

	// Verify first option structure
	option := options[0].(map[string]interface{})
	assert.Equal(t, "0-50", option["key"])
	assert.Equal(t, float64(0), option["from"])
	assert.Equal(t, float64(50), option["to"])
}

func TestHandleGetSearchParams_EmptyFilters(t *testing.T) {
	viper.Reset() // No filters configured

	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/v1/filters", handleGetSearchParams)

	req, _ := http.NewRequest("GET", "/v1/filters", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	meta := response["meta"].(map[string]interface{})
	assert.Equal(t, float64(0), meta["count"])
}

// Test JSON API error response format
func TestJsonapiErrorResp(t *testing.T) {
	code, body := JsonapiErrorResp(500, "Test error message")

	assert.Equal(t, 500, code)
	assert.NotNil(t, body)

	// Type assert to JsonapiError
	jsonapiErr, ok := body.(JsonapiError)
	assert.True(t, ok)
	assert.NotNil(t, jsonapiErr.APIErrors)
	assert.Len(t, jsonapiErr.APIErrors, 1)
	assert.Equal(t, "Test error message", jsonapiErr.APIErrors[0].Detail)
}

// Test query builder integration
func TestQueryBuilderIntegration(t *testing.T) {
	setupAPITestConfig()

	// Set up search configuration
	viper.Set("search.fields", []string{"name", "description"})
	viper.Set("search.default_operator", "OR")

	params := map[string][]string{
		"category": {"electronics"},
		"brand":    {"apple"},
		"limit":    {"20"},
	}

	qb := search.NewQueryBuilder("laptop", params)
	result := qb.Build()

	assert.NotNil(t, result)
	assert.Equal(t, 20, result["size"])

	// Verify query structure
	query := result["query"].(map[string]interface{})
	boolQuery := query["bool"].(map[string]interface{})

	// Verify must clause
	must := boolQuery["must"].([]interface{})
	assert.Len(t, must, 1)

	// Verify filter clause
	filter := boolQuery["filter"].([]interface{})
	assert.Len(t, filter, 2) // category and brand
}

func TestResponseContentType(t *testing.T) {
	assert.Equal(t, "application/json; charset=utf-8", contentType)
}

// Test multiple query parameters
func TestHandleSearch_MultipleCategories(t *testing.T) {
	setupAPITestConfig()
	viper.Set("search.fields", []string{"name", "description"})

	params := map[string][]string{
		"category": {"electronics", "computers"},
		"q":        {"laptop"},
	}

	qb := search.NewQueryBuilder(params["q"][0], params)
	result := qb.Build()

	query := result["query"].(map[string]interface{})
	boolQuery := query["bool"].(map[string]interface{})

	// Verify filter has the terms filter with multiple categories
	filter := boolQuery["filter"].([]interface{})
	assert.Len(t, filter, 1)

	filterMap := filter[0].(map[string]interface{})
	terms := filterMap["terms"].(map[string]interface{})
	categoryValues := terms["category.keyword"].([]string)
	assert.Contains(t, categoryValues, "electronics")
	assert.Contains(t, categoryValues, "computers")
}

// Test pagination parameters
func TestPaginationParsing(t *testing.T) {
	setupAPITestConfig()

	params := map[string][]string{
		"limit": {"50"},
		"from":  {"100"},
	}

	qb := search.NewQueryBuilder("", params)
	result := qb.Build()

	assert.Equal(t, 50, result["size"])
	assert.Equal(t, 100, result["from"])
}

// Test sorting parameters
func TestSortingParsing(t *testing.T) {
	setupAPITestConfig()

	params := map[string][]string{
		"sort":  {"price"},
		"order": {"asc"},
	}

	qb := search.NewQueryBuilder("", params)
	result := qb.Build()

	sort := result["sort"].([]map[string]interface{})
	assert.Len(t, sort, 1)

	priceSort := sort[0]["price"].(map[string]interface{})
	assert.Equal(t, "asc", priceSort["order"])
}

// Test invalid pagination values
func TestInvalidPaginationValues(t *testing.T) {
	setupAPITestConfig()

	params := map[string][]string{
		"limit": {"invalid"},
		"from":  {"also-invalid"},
	}

	qb := search.NewQueryBuilder("", params)
	result := qb.Build()

	// Should fall back to defaults
	assert.Equal(t, 40, result["size"])
	assert.Equal(t, 0, result["from"])
}

// Test aggregations in response
func TestAggregationsInQuery(t *testing.T) {
	setupAPITestConfig()

	viper.Set("aggregations.category", map[string]interface{}{
		"type":  "terms",
		"field": "category.keyword",
		"size":  50,
	})

	qb := search.NewQueryBuilder("", map[string][]string{})
	result := qb.Build()

	aggs, ok := result["aggs"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, aggs)

	categoryAgg := aggs["category"].(map[string]interface{})
	terms := categoryAgg["terms"].(map[string]interface{})
	assert.Equal(t, "category.keyword", terms["field"])
	assert.Equal(t, 50, terms["size"])
}
