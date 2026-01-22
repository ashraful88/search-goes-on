package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestJsonapiSearchResultRaw(t *testing.T) {
	// Mock Elasticsearch response
	esResponse := `{
		"hits": {
			"total": {
				"value": 2,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "1",
					"_source": {
						"name": "MacBook Pro",
						"price": 1299,
						"category": "electronics"
					}
				},
				{
					"_id": "2",
					"_source": {
						"name": "iPhone",
						"price": 999,
						"category": "electronics"
					}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	// Verify it's valid JSON
	var jsonResult map[string]interface{}
	err = json.Unmarshal([]byte(result), &jsonResult)
	assert.NoError(t, err)

	// Verify JSON API structure
	data, ok := jsonResult["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)

	meta, ok := jsonResult["meta"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(2), meta["total"])
}

func TestJsonapiSearchResultRaw_ItemStructure(t *testing.T) {
	esResponse := `{
		"hits": {
			"total": {
				"value": 1,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "123",
					"_source": {
						"name": "Test Product",
						"price": 99.99,
						"description": "Test description"
					}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	// Parse using gjson for easier access
	dataArray := gjson.Get(result, "data")
	assert.True(t, dataArray.IsArray())
	assert.Equal(t, int64(1), dataArray.Get("#").Int())

	// Verify first item structure
	item := dataArray.Get("0")
	assert.Equal(t, "article", item.Get("type").String())
	assert.Equal(t, "123", item.Get("id").String())

	// Verify attributes
	attributes := item.Get("attributes")
	assert.Equal(t, "Test Product", attributes.Get("name").String())
	assert.Equal(t, 99.99, attributes.Get("price").Float())
	assert.Equal(t, "Test description", attributes.Get("description").String())
}

func TestJsonapiSearchResultRaw_EmptyResults(t *testing.T) {
	esResponse := `{
		"hits": {
			"total": {
				"value": 0,
				"relation": "eq"
			},
			"hits": []
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	data := gjson.Get(result, "data")
	assert.True(t, data.IsArray())
	assert.Equal(t, int64(0), data.Get("#").Int())

	total := gjson.Get(result, "meta.total")
	assert.Equal(t, int64(0), total.Int())
}

func TestJsonapiSearchResultRaw_MultipleItems(t *testing.T) {
	esResponse := `{
		"hits": {
			"total": {
				"value": 3,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "1",
					"_source": {"name": "Product 1"}
				},
				{
					"_id": "2",
					"_source": {"name": "Product 2"}
				},
				{
					"_id": "3",
					"_source": {"name": "Product 3"}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	data := gjson.Get(result, "data")
	assert.Equal(t, int64(3), data.Get("#").Int())

	// Verify all items are present
	for i := 0; i < 3; i++ {
		item := data.Get(string(rune('0' + i)))
		assert.True(t, item.Exists())
		assert.Equal(t, "article", item.Get("type").String())
	}
}

func TestJsonapiSearchResultRaw_ComplexAttributes(t *testing.T) {
	esResponse := `{
		"hits": {
			"total": {
				"value": 1,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "complex-1",
					"_source": {
						"name": "Complex Product",
						"price": 299.99,
						"rating": 4.5,
						"in_stock": true,
						"tags": ["electronics", "laptop"],
						"metadata": {
							"brand": "Apple",
							"model": "MacBook"
						}
					}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	// Verify complex nested structure is preserved
	item := gjson.Get(result, "data.0")
	assert.Equal(t, "complex-1", item.Get("id").String())

	attributes := item.Get("attributes")
	assert.Equal(t, "Complex Product", attributes.Get("name").String())
	assert.Equal(t, 299.99, attributes.Get("price").Float())
	assert.Equal(t, 4.5, attributes.Get("rating").Float())
	assert.True(t, attributes.Get("in_stock").Bool())

	// Verify array
	tags := attributes.Get("tags")
	assert.True(t, tags.IsArray())
	assert.Contains(t, tags.String(), "electronics")

	// Verify nested object
	metadata := attributes.Get("metadata")
	assert.Equal(t, "Apple", metadata.Get("brand").String())
	assert.Equal(t, "MacBook", metadata.Get("model").String())
}

func TestJsonapiErrorResp_Structure(t *testing.T) {
	code, body := JsonapiErrorResp(404, "Resource not found")

	assert.Equal(t, 404, code)
	assert.NotNil(t, body)

	// Convert to JSON to verify structure
	jsonBytes, err := json.Marshal(body)
	assert.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	assert.NoError(t, err)

	errors, ok := result["errors"]
	assert.True(t, ok)
	assert.NotNil(t, errors)
}

func TestJsonapiErrorResp_DifferentCodes(t *testing.T) {
	testCases := []struct {
		code    int
		message string
	}{
		{400, "Bad Request"},
		{401, "Unauthorized"},
		{403, "Forbidden"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
		{503, "Service Unavailable"},
	}

	for _, tc := range testCases {
		code, body := JsonapiErrorResp(tc.code, tc.message)
		assert.Equal(t, tc.code, code)
		assert.NotNil(t, body)
	}
}

func TestJsonapiSearchResultRaw_PreservesFieldTypes(t *testing.T) {
	esResponse := `{
		"hits": {
			"total": {
				"value": 1,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "type-test",
					"_source": {
						"string_field": "text",
						"number_field": 42,
						"float_field": 3.14,
						"bool_field": true
					}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	attributes := gjson.Get(result, "data.0.attributes")

	// Verify types are preserved
	assert.Equal(t, "text", attributes.Get("string_field").String())
	assert.Equal(t, int64(42), attributes.Get("number_field").Int())
	assert.Equal(t, 3.14, attributes.Get("float_field").Float())
	assert.True(t, attributes.Get("bool_field").Bool())
}

func TestJsonapiSearchResultRaw_LargeTotalValue(t *testing.T) {
	// Note: The function iterates based on total.value, so we need to provide
	// a reasonable total value that matches the hits array length
	// Testing with a small array to avoid performance issues
	esResponse := `{
		"hits": {
			"total": {
				"value": 100,
				"relation": "gte"
			},
			"hits": []
		}
	}`

	// This test documents current behavior - the function will iterate
	// total.value times even if hits array is smaller/empty
	// For production, total.value should match actual hits returned
	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	total := gjson.Get(result, "meta.total")
	assert.Equal(t, int64(100), total.Int())
}

func TestJsonapiSearchResultRaw_WithAggregations(t *testing.T) {
	// Elasticsearch response with aggregations
	esResponse := `{
		"hits": {
			"total": {
				"value": 2,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "1",
					"_source": {"name": "Product 1", "category": "electronics"}
				},
				{
					"_id": "2",
					"_source": {"name": "Product 2", "category": "electronics"}
				}
			]
		},
		"aggregations": {
			"category": {
				"buckets": [
					{"key": "electronics", "doc_count": 2}
				]
			}
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	// The function currently doesn't include aggregations in the response
	// This test documents current behavior
	data := gjson.Get(result, "data")
	assert.Equal(t, int64(2), data.Get("#").Int())

	// Note: Aggregations are not included in the JSON API response
	// This is expected behavior based on the current implementation
}

func TestJsonapiSearchResultRaw_IDTypes(t *testing.T) {
	// Test different ID formats
	esResponse := `{
		"hits": {
			"total": {
				"value": 3,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "string-id",
					"_source": {"name": "String ID"}
				},
				{
					"_id": "123",
					"_source": {"name": "Numeric String ID"}
				},
				{
					"_id": "uuid-f47ac10b-58cc-4372-a567-0e02b2c3d479",
					"_source": {"name": "UUID ID"}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	data := gjson.Get(result, "data")

	assert.Equal(t, "string-id", data.Get("0.id").String())
	assert.Equal(t, "123", data.Get("1.id").String())
	assert.Equal(t, "uuid-f47ac10b-58cc-4372-a567-0e02b2c3d479", data.Get("2.id").String())
}

func TestJsonapiSearchResultRaw_SpecialCharacters(t *testing.T) {
	esResponse := `{
		"hits": {
			"total": {
				"value": 1,
				"relation": "eq"
			},
			"hits": [
				{
					"_id": "special-chars",
					"_source": {
						"name": "Product with \"quotes\"",
						"description": "Line 1\nLine 2",
						"unicode": "Hello 世界 🌍"
					}
				}
			]
		}
	}`

	result, err := JsonapiSearchResultRaw(esResponse)
	assert.NoError(t, err)

	// Verify special characters are properly handled
	attributes := gjson.Get(result, "data.0.attributes")
	assert.Contains(t, attributes.Get("name").String(), "quotes")
	assert.Contains(t, attributes.Get("description").String(), "Line 1")
	assert.Contains(t, attributes.Get("unicode").String(), "世界")
}
