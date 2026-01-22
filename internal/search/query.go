package search

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// QueryBuilder builds Elasticsearch queries from URL parameters and YAML config
type QueryBuilder struct {
	QueryString string
	Params      map[string][]string
	Query       map[string]interface{}
	Aggs        map[string]interface{}
	Size        int
	From        int
	Sort        []map[string]interface{}
}

// NewQueryBuilder creates a new query builder instance
func NewQueryBuilder(queryString string, params map[string][]string) *QueryBuilder {
	qb := &QueryBuilder{
		QueryString: queryString,
		Params:      params,
		Query:       make(map[string]interface{}),
		Aggs:        make(map[string]interface{}),
		Size:        40,
		From:        0,
		Sort:        []map[string]interface{}{},
	}
	return qb
}

// Build constructs the complete Elasticsearch query
func (qb *QueryBuilder) Build() map[string]interface{} {
	boolQuery := map[string]interface{}{}

	// Build query context (affects scoring)
	must := qb.buildQueryContext()
	if len(must) > 0 {
		boolQuery["must"] = must
	}

	// Build filter context (does not affect scoring)
	filter := qb.buildFilterContext()
	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}

	// If no query or filters, use match_all
	if len(must) == 0 && len(filter) == 0 {
		qb.Query = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		qb.Query = map[string]interface{}{
			"bool": boolQuery,
		}
	}

	// Build pagination and sorting
	qb.buildPaginationAndSort()

	// Build aggregations
	qb.buildAggregations()

	// Construct final query
	result := map[string]interface{}{
		"query": qb.Query,
		"size":  qb.Size,
		"from":  qb.From,
	}

	if len(qb.Sort) > 0 {
		result["sort"] = qb.Sort
	}

	if len(qb.Aggs) > 0 {
		result["aggs"] = qb.Aggs
	}

	return result
}

// buildQueryContext builds the query context (must clause) for text search
func (qb *QueryBuilder) buildQueryContext() []interface{} {
	must := []interface{}{}

	if qb.QueryString != "" {
		// Get search fields from config
		searchFields := qb.getSearchFields()
		defaultOperator := viper.GetString("search.default_operator")
		if defaultOperator == "" {
			defaultOperator = "OR"
		}

		// Build multi_match query for full-text search
		multiMatch := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":    qb.QueryString,
				"fields":   searchFields,
				"type":     "best_fields",
				"operator": defaultOperator,
			},
		}

		minimumShouldMatch := viper.GetInt("search.minimum_should_match")
		if minimumShouldMatch > 0 {
			multiMatch["multi_match"].(map[string]interface{})["minimum_should_match"] = minimumShouldMatch
		}

		must = append(must, multiMatch)
	}

	return must
}

// buildFilterContext builds filter context from URL parameters and YAML config
func (qb *QueryBuilder) buildFilterContext() []interface{} {
	filter := []interface{}{}

	// Get filters configuration
	filtersConfig := viper.GetStringMap("filters")

	for filterName, filterConfigRaw := range filtersConfig {
		filterConfig, ok := filterConfigRaw.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if this filter has values in the request params
		paramValues, hasParam := qb.Params[filterName]
		if !hasParam || len(paramValues) == 0 {
			continue
		}

		filterType := getStringFromMap(filterConfig, "type")
		field := getStringFromMap(filterConfig, "field")

		switch filterType {
		case "terms":
			// Multi-value term filter (e.g., category, brand)
			filter = append(filter, map[string]interface{}{
				"terms": map[string]interface{}{
					field: paramValues,
				},
			})

		case "term":
			// Single value term filter (e.g., in_stock)
			filter = append(filter, map[string]interface{}{
				"term": map[string]interface{}{
					field: paramValues[0],
				},
			})

		case "range":
			// Range filter (e.g., price, rating)
			rangeFilter := qb.buildRangeFilter(field, paramValues[0], filterConfig)
			if rangeFilter != nil {
				filter = append(filter, rangeFilter)
			}
		}
	}

	return filter
}

// buildRangeFilter builds a range filter from the config
func (qb *QueryBuilder) buildRangeFilter(field string, rangeKey string, filterConfig map[string]interface{}) map[string]interface{} {
	rangesRaw, ok := filterConfig["ranges"]
	if !ok {
		return nil
	}

	ranges, ok := rangesRaw.([]interface{})
	if !ok {
		return nil
	}

	// Find the matching range
	for _, rangeRaw := range ranges {
		rangeMap, ok := rangeRaw.(map[string]interface{})
		if !ok {
			continue
		}

		key := getStringFromMap(rangeMap, "key")
		if key == rangeKey {
			rangeQuery := map[string]interface{}{}

			if from, ok := rangeMap["from"]; ok {
				rangeQuery["gte"] = from
			}

			if to, ok := rangeMap["to"]; ok {
				rangeQuery["lt"] = to
			}

			return map[string]interface{}{
				"range": map[string]interface{}{
					field: rangeQuery,
				},
			}
		}
	}

	return nil
}

// buildPaginationAndSort builds pagination and sorting parameters
func (qb *QueryBuilder) buildPaginationAndSort() {
	// Pagination
	if limitStr, ok := qb.Params["limit"]; ok && len(limitStr) > 0 {
		if limit, err := strconv.Atoi(limitStr[0]); err == nil && limit > 0 {
			qb.Size = limit
		}
	}

	if fromStr, ok := qb.Params["from"]; ok && len(fromStr) > 0 {
		if from, err := strconv.Atoi(fromStr[0]); err == nil && from >= 0 {
			qb.From = from
		}
	}

	// Sorting
	if sortFields, ok := qb.Params["sort"]; ok && len(sortFields) > 0 {
		order := "desc"
		if orderVals, ok := qb.Params["order"]; ok && len(orderVals) > 0 {
			order = orderVals[0]
		}

		for _, sortField := range strings.Split(sortFields[0], ",") {
			sortField = strings.TrimSpace(sortField)
			if sortField != "" {
				qb.Sort = append(qb.Sort, map[string]interface{}{
					sortField: map[string]interface{}{
						"order": order,
					},
				})
			}
		}
	} else if qb.QueryString == "" {
		// Default sort by date if no query string (no relevance scoring)
		qb.Sort = append(qb.Sort, map[string]interface{}{
			"created_at": map[string]interface{}{
				"order": "desc",
			},
		})
	}
}

// buildAggregations builds aggregations for faceted search
func (qb *QueryBuilder) buildAggregations() {
	aggsConfig := viper.GetStringMap("aggregations")

	for aggName, aggConfigRaw := range aggsConfig {
		aggConfig, ok := aggConfigRaw.(map[string]interface{})
		if !ok {
			continue
		}

		aggType := getStringFromMap(aggConfig, "type")
		field := getStringFromMap(aggConfig, "field")

		switch aggType {
		case "terms":
			size := 10
			if sizeVal, ok := aggConfig["size"]; ok {
				if sizeInt, ok := sizeVal.(int); ok {
					size = sizeInt
				}
			}

			qb.Aggs[aggName] = map[string]interface{}{
				"terms": map[string]interface{}{
					"field": field,
					"size":  size,
				},
			}

		case "range":
			rangesRaw, ok := aggConfig["ranges"]
			if !ok {
				continue
			}

			ranges, ok := rangesRaw.([]interface{})
			if !ok {
				continue
			}

			rangesList := []map[string]interface{}{}
			for _, rangeRaw := range ranges {
				rangeMap, ok := rangeRaw.(map[string]interface{})
				if !ok {
					continue
				}

				r := map[string]interface{}{
					"key": getStringFromMap(rangeMap, "key"),
				}

				if from, ok := rangeMap["from"]; ok {
					r["from"] = from
				}

				if to, ok := rangeMap["to"]; ok {
					r["to"] = to
				}

				rangesList = append(rangesList, r)
			}

			qb.Aggs[aggName] = map[string]interface{}{
				"range": map[string]interface{}{
					"field":  field,
					"ranges": rangesList,
				},
			}

		case "avg":
			qb.Aggs[aggName] = map[string]interface{}{
				"avg": map[string]interface{}{
					"field": field,
				},
			}
		}
	}
}

// getSearchFields returns the list of fields to search
func (qb *QueryBuilder) getSearchFields() []string {
	fieldsRaw := viper.Get("search.fields")
	if fieldsRaw == nil {
		return []string{"_all"}
	}

	fields := []string{}
	switch v := fieldsRaw.(type) {
	case []interface{}:
		for _, field := range v {
			if fieldStr, ok := field.(string); ok {
				fields = append(fields, fieldStr)
			}
		}
	case []string:
		fields = v
	}

	if len(fields) == 0 {
		return []string{"_all"}
	}

	return fields
}

// ToJSON converts the query to JSON string
func (qb *QueryBuilder) ToJSON() (string, error) {
	query := qb.Build()
	jsonBytes, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// Helper function to safely get string from map
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}
