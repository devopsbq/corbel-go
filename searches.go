package corbel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Search is the struct used to query every searcheable api in the platform
type Search struct {
	client   *Client
	Query    *apiquery
	Sort     *sort
	PageSize int
	endpoint string
	baseURL  string
}

// aggregationCount is the json representation of Count responses
type aggregationCount struct {
	Count int `json:"count"`
}

// aggregationAvg is the json representation of Average responses
type aggregationAvg struct {
	Average float64 `json:"average"`
}

// aggregationSum is the json representation of Sum responses
type aggregationSum struct {
	Sum float64 `json:"sum"`
}

// Page fills the struct array passed as parameter as paged search by pageNumber
func (s *Search) Page(pageNumber int, result interface{}) error {
	var (
		req *http.Request
		err error
	)
	opts := &SearchListOptions{
		APIQuery:    s.Query.string(),
		APISort:     s.Sort.string(),
		APIPage:     pageNumber,
		APIPageSize: s.PageSize,
	}
	req, err = s.client.NewRequest("GET", s.endpoint, s.queryString(opts), nil)
	_, err = returnErrorHTTPInterface(s.client, req, err, result, 200)
	return err
}

// Count returns the aggregated count of an especific field in the search
func (s *Search) Count(field string) (int, error) {
	var (
		req       *http.Request
		err       error
		aggrCount aggregationCount
	)
	opts := &SearchListOptions{
		APIQuery:       s.Query.string(),
		APISort:        s.Sort.string(),
		APIAggregation: fmt.Sprintf("{\"$count\":\"%s\"}", field),
	}

	req, err = s.client.NewRequest("GET", s.endpoint, s.queryString(opts), nil)
	_, err = returnErrorHTTPInterface(s.client, req, err, &aggrCount, 200)
	if err != nil {
		return 0, err
	}

	return aggrCount.Count, nil
}

// CountAll returns the aggregated count of all items in the search.
// It's an alias of Count("*")
func (s *Search) CountAll() (int, error) {
	return s.Count("*")
}

// Average returns the average of an especific field in the search
func (s *Search) Average(field string) (float64, error) {
	var (
		req     *http.Request
		err     error
		aggrAvg aggregationAvg
	)

	opts := &SearchListOptions{
		APIQuery:       s.Query.string(),
		APISort:        s.Sort.string(),
		APIAggregation: fmt.Sprintf("{\"$avg\":\"%s\"}", field),
	}
	req, err = s.client.NewRequest("GET", s.endpoint, s.queryString(opts), nil)
	_, err = returnErrorHTTPInterface(s.client, req, err, &aggrAvg, 200)
	if err != nil {
		return 0, err
	}

	return aggrAvg.Average, nil
}

// // Sum returns the average of an especific field in the search as integer
// func (s *Search) Sum(field string) (int, error) {
// 	sum, err := s.SumFloat(field)
// 	return int(sum), err
// }

// Sum returns the average of an especific field in the search as float
func (s *Search) Sum(field string) (float64, error) {
	var (
		req     *http.Request
		err     error
		aggrSum aggregationSum
	)

	opts := &SearchListOptions{
		APIQuery:       s.Query.string(),
		APISort:        s.Sort.string(),
		APIAggregation: fmt.Sprintf("{\"$sum\":\"%s\"}", field),
	}

	req, err = s.client.NewRequest("GET", s.endpoint, s.queryString(opts), nil)
	_, err = returnErrorHTTPInterface(s.client, req, err, &aggrSum, 200)
	if err != nil {
		return 0, err
	}

	return aggrSum.Sum, nil
}

// SearchListOptions specifies the optional parameters for searches supporting
// paging and aggregation
type SearchListOptions struct {
	APIQuery       string `url:"api:query,omitempty"`
	APISort        string `url:"api:sort,omitempty"`
	APIPageSize    int    `url:"api:pageSize,omitempty"`
	APIPage        int    `url:"api:page,omitempty"`
	APIAggregation string `url:"api:aggregation,omitempty"`
}

func (s *Search) queryString(options *SearchListOptions) string {
	path, _ := addOptions(s.baseURL, options)
	return path
}

// NewSearch returns a Search struct that allows to select especific search requirements
func NewSearch(client *Client, endpoint, baseURL string) *Search {
	return &Search{
		client:   client,
		Query:    newQuery(),
		Sort:     newSort(),
		PageSize: 10,
		endpoint: endpoint,
		baseURL:  baseURL,
	}
}

// Query is the struct that contains the especification of a search
type apiquery struct {
	Eq   map[string]string   `json:"$eq,omitempty"`
	Gt   map[string]int      `json:"$gt,omitempty"`
	Gte  map[string]int      `json:"$gte,omitempty"`
	Lt   map[string]int      `json:"$lt,omitempty"`
	Lte  map[string]int      `json:"$lte,omitempty"`
	Ne   map[string]string   `json:"$ne,omitempty"`
	In   map[string][]string `json:"$in,omitempty"`
	All  map[string][]string `json:"$all,omitempty"`
	Like map[string]string   `json:"$like,omitempty"`
}

// QueryString returns the query string to append to the url we are searching for.
// api:query must be enclosed by []
func (q *apiquery) string() string {
	apiQueryString, _ := json.Marshal(q)
	if string(apiQueryString) == "{}" {
		return ""
	}
	return fmt.Sprintf("[%s]", string(apiQueryString))
}

// NewQuery returns a New search struct
func newQuery() *apiquery {
	return &apiquery{
		Eq:   make(map[string]string),
		Gt:   make(map[string]int),
		Gte:  make(map[string]int),
		Lt:   make(map[string]int),
		Lte:  make(map[string]int),
		Ne:   make(map[string]string),
		In:   make(map[string][]string),
		All:  make(map[string][]string),
		Like: make(map[string]string),
	}
}

// Sort is the struct that contains the sort to apply to a search
type sort struct {
	Asc  []string
	Desc []string
}

// NewSort returns a Sort struct
func newSort() *sort {
	return &sort{}
}

// QueryString returns the query string to append to the url we are sorting
func (s *sort) string() string {
	var (
		apiSortString []byte
		apiSortMap    = make(map[string]string)
		field         string
	)

	for _, field = range s.Asc {
		apiSortMap[field] = "asc"
	}

	for _, field = range s.Desc {
		apiSortMap[field] = "desc"
	}

	apiSortString, _ = json.Marshal(apiSortMap)
	if string(apiSortString) == "{}" {
		return ""
	}
	return string(apiSortString)
}
