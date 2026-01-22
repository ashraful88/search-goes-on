# Testing Guide

This document describes the test suite for the Search API project.

## Test Structure

```
internal/search/
├── query_test.go           # QueryBuilder unit tests
├── filters_test.go         # Filter configuration tests

cmd/search-api/internal/api/
├── request_handlers_test.go    # API handler tests
└── response_helper_test.go     # JSON API response format tests
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test ./... -cover
```

### Run Tests with Verbose Output

```bash
go test ./... -v
```

### Run Specific Test File

```bash
go test ./internal/search -v
go test ./cmd/search-api/internal/api -v
```

### Run Specific Test Function

```bash
go test ./internal/search -run TestBuildQueryContext_WithQueryString -v
```

### Generate Coverage Report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Coverage

### Query Builder Tests (`query_test.go`)

Tests the Elasticsearch query generation functionality:

- **TestNewQueryBuilder**: Verifies QueryBuilder initialization
- **TestBuildQueryContext_WithQueryString**: Tests full-text search query generation
- **TestBuildQueryContext_EmptyQueryString**: Tests behavior without query string
- **TestBuildFilterContext_TermsFilter**: Tests multi-value term filters (category, brand)
- **TestBuildFilterContext_TermFilter**: Tests single-value term filters (in_stock)
- **TestBuildFilterContext_RangeFilter**: Tests range filters (price ranges)
- **TestBuildPaginationAndSort**: Tests pagination and sorting parameters
- **TestBuildPaginationAndSort_Defaults**: Tests default pagination/sort values
- **TestBuildAggregations**: Tests aggregation generation for faceted search
- **TestBuild_CompleteQuery**: Tests complete query with all components
- **TestBuild_MatchAll**: Tests match_all query when no parameters provided
- **TestToJSON**: Tests JSON serialization of queries
- **TestGetSearchFields**: Tests search field configuration
- **TestBuildRangeFilter_InvalidRange**: Tests invalid range handling
- **TestBuildRangeFilter_OpenEndedRange**: Tests open-ended range (e.g., "200+")
- **TestMultipleSortFields**: Tests sorting by multiple fields

### Filter Tests (`filters_test.go`)

Tests the filter configuration and metadata functionality:

- **TestGetFilterConfig**: Tests retrieving all filter configuration
- **TestGetFilterConfigByName**: Tests retrieving specific filter config
- **TestGetFiltersMetadata**: Tests filter metadata generation
- **TestGetFiltersMetadata_TermsFilter**: Tests terms filter metadata structure
- **TestGetFiltersMetadata_TermFilter**: Tests term filter metadata structure
- **TestGetFiltersMetadata_RangeFilter**: Tests range filter metadata with options
- **TestGetFiltersMetadata_RatingRangeFilter**: Tests rating filter with open-ended ranges
- **TestGetFiltersMetadata_EmptyConfig**: Tests behavior with no filters
- **TestGetFiltersMetadata_InvalidConfig**: Tests handling of invalid configuration
- **TestGetStringValue**: Tests string extraction helper function
- **TestFilterMetadata_JSONSerialization**: Tests JSON marshaling of filter metadata
- **TestGetFiltersMetadata_ConsistentOrdering**: Tests consistency across multiple calls

### API Handler Tests (`request_handlers_test.go`)

Tests the HTTP request handlers:

- **TestHandleGetSearchParams**: Tests `/v1/filters` endpoint response
- **TestHandleGetSearchParams_ContentType**: Tests proper content-type header
- **TestHandleGetSearchParams_FilterStructure**: Tests filter response structure
- **TestHandleGetSearchParams_TermsFilter**: Tests terms filter in response
- **TestHandleGetSearchParams_RangeFilter**: Tests range filter with options
- **TestHandleGetSearchParams_EmptyFilters**: Tests response with no filters configured
- **TestJsonapiErrorResp**: Tests error response format
- **TestQueryBuilderIntegration**: Tests QueryBuilder integration with handlers
- **TestResponseContentType**: Tests content type constant
- **TestHandleSearch_MultipleCategories**: Tests multi-select filter handling
- **TestPaginationParsing**: Tests pagination parameter parsing
- **TestSortingParsing**: Tests sort parameter parsing
- **TestInvalidPaginationValues**: Tests handling of invalid pagination values
- **TestAggregationsInQuery**: Tests aggregations in generated queries

### Response Helper Tests (`response_helper_test.go`)

Tests the JSON API response formatting:

- **TestJsonapiSearchResultRaw**: Tests basic Elasticsearch to JSON API conversion
- **TestJsonapiSearchResultRaw_ItemStructure**: Tests item structure in response
- **TestJsonapiSearchResultRaw_EmptyResults**: Tests empty result handling
- **TestJsonapiSearchResultRaw_MultipleItems**: Tests multiple result items
- **TestJsonapiSearchResultRaw_ComplexAttributes**: Tests nested and complex attributes
- **TestJsonapiErrorResp_Structure**: Tests error response structure
- **TestJsonapiErrorResp_DifferentCodes**: Tests various HTTP status codes
- **TestJsonapiSearchResultRaw_PreservesFieldTypes**: Tests data type preservation
- **TestJsonapiSearchResultRaw_LargeTotalValue**: Tests large result counts
- **TestJsonapiSearchResultRaw_WithAggregations**: Documents aggregation behavior
- **TestJsonapiSearchResultRaw_IDTypes**: Tests various ID formats
- **TestJsonapiSearchResultRaw_SpecialCharacters**: Tests special character handling

## Test Dependencies

The test suite uses the following testing libraries:

- **testing**: Go standard library testing package
- **github.com/stretchr/testify/assert**: Assertion library for readable tests
- **github.com/spf13/viper**: Configuration management (used in tests)
- **github.com/gin-gonic/gin**: Web framework (set to test mode)
- **github.com/tidwall/gjson**: JSON parsing for response validation

## Mocking Strategy

### Configuration Mocking

Tests use `viper.Reset()` and `viper.Set()` to mock YAML configuration:

```go
func setupTestConfig() {
    viper.Reset()
    viper.Set("search.fields", []string{"name", "description"})
    viper.Set("filters.category", map[string]interface{}{
        "label": "Category",
        "field": "category.keyword",
        "type": "terms",
    })
}
```

### HTTP Request Mocking

Tests use Gin's test mode and `httptest` package:

```go
gin.SetMode(gin.TestMode)
router := gin.New()
router.GET("/v1/filters", handleGetSearchParams)

req, _ := http.NewRequest("GET", "/v1/filters", nil)
w := httptest.NewRecorder()
router.ServeHTTP(w, req)
```

### Elasticsearch Response Mocking

Tests use JSON strings to mock Elasticsearch responses:

```go
esResponse := `{
    "hits": {
        "total": {"value": 2},
        "hits": [
            {"_id": "1", "_source": {"name": "Product 1"}}
        ]
    }
}`
```

## Test Data

### Sample Product Data

Tests use the following product structure:

```json
{
    "name": "MacBook Pro",
    "description": "Powerful laptop",
    "price": 1299,
    "category": "electronics",
    "brand": "apple",
    "rating": 4.8,
    "in_stock": true
}
```

### Sample Filter Configuration

```yaml
filters:
  category:
    label: Category
    field: category.keyword
    type: terms
    multi: true

  price:
    label: Price
    field: price
    type: range
    ranges:
      - key: "0-50"
        from: 0
        to: 50
```

## Common Test Patterns

### Testing Query Generation

```go
func TestQueryGeneration(t *testing.T) {
    setupTestConfig()

    params := map[string][]string{
        "category": {"electronics"},
        "q": {"laptop"},
    }

    qb := search.NewQueryBuilder(params["q"][0], params)
    result := qb.Build()

    assert.NotNil(t, result)
    assert.Contains(t, result, "query")
}
```

### Testing API Responses

```go
func TestAPIResponse(t *testing.T) {
    setupAPITestConfig()

    router := gin.New()
    router.GET("/endpoint", handler)

    req, _ := http.NewRequest("GET", "/endpoint", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.NotNil(t, response)
}
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.25'
      - run: go test ./... -v -cover
```

## Writing New Tests

### Guidelines

1. **Use descriptive test names**: `TestFunctionName_Scenario_ExpectedResult`
2. **Setup and teardown**: Use `viper.Reset()` to clean state between tests
3. **Test one thing**: Each test should verify one specific behavior
4. **Use table-driven tests**: For testing multiple scenarios
5. **Assert clearly**: Use testify/assert for readable assertions

### Example Test Template

```go
func TestNewFeature_SpecificScenario(t *testing.T) {
    // Arrange
    setupTestConfig()
    input := prepareTestInput()

    // Act
    result := functionUnderTest(input)

    // Assert
    assert.NotNil(t, result)
    assert.Equal(t, expectedValue, result.Field)
}
```

## Debugging Tests

### Run with Race Detector

```bash
go test ./... -race
```

### Run Specific Test with Debug Output

```bash
go test ./internal/search -run TestBuildQueryContext -v
```

### Print Test Coverage per Function

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Performance Testing

### Benchmark Tests

Create benchmark tests for performance-critical code:

```go
func BenchmarkQueryBuilder(b *testing.B) {
    setupTestConfig()
    params := map[string][]string{"category": {"electronics"}}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        qb := search.NewQueryBuilder("laptop", params)
        qb.Build()
    }
}
```

Run benchmarks:

```bash
go test ./... -bench=. -benchmem
```

## Test Coverage Goals

- **Unit Tests**: > 80% coverage for core logic
- **Integration Tests**: All API endpoints tested
- **Edge Cases**: Invalid inputs, empty data, large datasets
- **Error Handling**: All error paths tested

## Troubleshooting

### Common Issues

1. **Viper state pollution**: Always call `viper.Reset()` in test setup
2. **Gin mode warnings**: Set `gin.SetMode(gin.TestMode)` at start of tests
3. **JSON parsing errors**: Verify mock data is valid JSON
4. **Import cycles**: Keep test files in same package as code under test

### Debug Tips

```bash
# Verbose output with test names
go test ./... -v

# Run only failing tests
go test ./... -failfast

# Show all output including fmt.Println
go test ./... -v | grep -A 10 "FAIL"
```

## Next Steps

1. Add integration tests with real Elasticsearch instance
2. Add end-to-end tests for complete workflows
3. Add contract tests for API responses
4. Set up automated test coverage reporting
5. Add performance regression tests
