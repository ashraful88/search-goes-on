package search

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
