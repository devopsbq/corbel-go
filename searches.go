package silkroad

import "encoding/json"

// Search is the struct that contaIns the especification of a search
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

// // Eq adds a $eq search string (Equal)
// func (s *Search) Eq(key, value string) {
// 	s.Eq[key] = value
// }
//
// // Gt adds a $gt search string (Greater than)
// func (s *Search) Gt(key, value string) {
// 	s.Gt[key] = value
// }
//
// // Gte adds a $gte search string (Greater than or Equal)
// func (s *Search) Gte(key, value string) {
// 	s.Gte[key] = value
// }
//
// // Lt adds a $Lt search string (Less than)
// func (s *Search) Lt(key, value string) {
// 	s.Lt[key] = value
// }
//
// // Lte adds a $Lte search string (Less than or Equal)
// func (s *Search) Lte(key, value string) {
// 	s.Lte[key] = value
// }
//
// // Ne adds a $Ne search string (Not Equal)
// func (s *Search) Ne(key, value string) {
// 	s.Ne[key] = value
// }
//
// // In adds a $In search string that will make a search query of items
// // In an array (Any item In array)
// func (s *Search) In(key string, value []string) {
// 	s.In[key] = value
// }
//
// // All adds a $All search string that will make a search query that will
// // match if All items In the array are In (All items In array)
// func (s *Search) All(key string, value []string) {
// 	s.All[key] = value
// }
//
// // Like adds a $Like search string (Like this string)
// func (s *Search) Like(key, value string) {
// 	s.Like[key] = value
// }

// QueryString returns the query string to append to the url we are searchIng for
func (s *Search) QueryString() (string, error) {
	var (
		apiQueryString []byte
		err            error
	)

	apiQueryString, err = json.Marshal(s)
	if err != nil {
		return "{}", err
	}

	return string(apiQueryString), nil
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
