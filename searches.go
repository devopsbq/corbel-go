package silkroad

import "encoding/json"

// Search is the struct that contains the especification of a search
type Search struct {
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

// QueryString returns the query string to append to the url we are searchIng for
func (s *Search) QueryString() (string, error) {
	var (
		apiQueryString []byte
		err            error
	)

	apiQueryString, err = json.Marshal(s)
	return string(apiQueryString), err
}

// NewSearch returns a New search struct
func NewSearch() *Search {
	return &Search{
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
type Sort struct {
	Asc  []string `json:"$asc,omitempty"`
	Desc []string `json:"$desc,omitempty"`
}

// NewSort returns a Sort struct
func NewSort() *Sort {
	return &Sort{}
}

// QueryString returns the query string to append to the url we are sorting
func (s *Sort) QueryString() (string, error) {
	var (
		apiSortString []byte
		apiSortMap    = make(map[string]string)
		err           error
	)

	for _, field := range s.Asc {
		apiSortMap[field] = "asc"
	}

	for _, field := range s.Desc {
		apiSortMap[field] = "desc"
	}

	apiSortString, err = json.Marshal(apiSortMap)
	return string(apiSortString), err
}

//
// var (
//   queryString    string
//   apiQueryString []byte
//   // apiSortString string
//   // apiAggregationString string
//   err error
// )
//
// queryString = "?"
//
// apiQueryString, err = json.Marshal(s)
// if err != nil {
//   return "", err
// }
//
// if string(apiQueryString) != "{}" {
//   queryString = fmt.Sprintf("%sapi:query=%s", queryString, string(apiQueryString))
// }
//
// return string(queryString), nil
