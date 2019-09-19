package search

import "fmt"

// ElasticQuery elastic search query
type ElasticQuery struct {
	Query struct {
		BoolQuery struct {
			Must   []interface{} `json:"must,omitempty"`
			Filter struct {
				BoolFilter struct {
					MustFilter   []interface{} `json:"must,omitempty"`
					ShouldFilter []interface{} `json:"should,omitempty"`
				} `json:"bool,omitempty"`
			} `json:"filter,omitempty"`
		} `json:"bool,omitempty"`
	} `json:"query,omitempty"`
	Size  string `json:"size"`
	From  string `json:"from"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
}
// QueryBuilder elasticsearch query builder by url query params
type QueryBuilder struct {
	SearchQuery ElasticQuery `json:"SearchQuery"`
	Params      map[string][]string
}

// add term filter to queryDSL
func (qb *QueryBuilder) addMustFilterByTerm(paramKey, termKey string) {
	if qb.Params[paramKey] != nil {
		mustFilter := map[string]interface{}{
			"term": map[string]interface{}{
				termKey: qb.Params[paramKey][0],
			},
		}
		qb.SearchQuery.Query.BoolQuery.Filter.BoolFilter.MustFilter = append(qb.SearchQuery.Query.BoolQuery.Filter.BoolFilter.MustFilter, mustFilter)
	}
}
// add must filters in es query
func (qb *QueryBuilder) buildFilters() {
	// todo: loop and yaml conf. but this one is faster as no loop, will try array switch
	searchParams := GetFilterConfig()
	for key, val := range searchParams {
		m2 := val.(map[string]interface{})
		pk := m2["filter"].(map[string]interface{})
		pkey := fmt.Sprintf("%v", pk["key"])
		log.Println(key, pkey)
		//log.Println(key, m2["filter"])
		qb.addMustFilterByTerm(key, pkey)
	}
}


func (eq *ElasticQuery) buildLimits(params map[string][]string) {
	// defaults
	eq.Size = "40"
	eq.From = "0"
	eq.Sort = `{ "date" : "desc"}`

	// override defaults if provided
	if params["limit"] != nil {
		eq.Size = params["limit"][0]
	}
	if params["from"] != nil {
		eq.From = params["from"][0]
	}
	if params["sort"] != nil {
		if params["order"] != nil {
			eq.Sort = fmt.Sprintf(` { "%s" : "%s" }`, params["sort"], params["order"])
		}
	}
}
