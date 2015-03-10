package silkroad

// Search is the struct that contains the especification of a search
type Search struct {
	eq   map[string]string
	gt   map[string]string
	gte  map[string]string
	lt   map[string]string
	lte  map[string]string
	ne   map[string]string
	in   map[string]string
	all  map[string]string
	like map[string]string
}

// NewSearch returns a new search struct
func NewSearch() *Search {
	return &Search{
		eq:   make(map[string]string),
		gt:   make(map[string]string),
		gte:  make(map[string]string),
		lt:   make(map[string]string),
		lte:  make(map[string]string),
		ne:   make(map[string]string),
		in:   make(map[string]string),
		all:  make(map[string]string),
		like: make(map[string]string),
	}
}

// Eq adds a $eq search string (Equal)
func (s *Search) Eq(key, value string) {
	s.eq[key] = value
}

// Gt adds a $gt search string (Greater than)
func (s *Search) Gt(key, value string) {
	s.gt[key] = value
}

// Gte adds a $gte search string (Greater than or equal)
func (s *Search) Gte(key, value string) {
	s.gte[key] = value
}

// Lt adds a $lt search string (Less than)
func (s *Search) Lt(key, value string) {
	s.lt[key] = value
}

// Lte adds a $lte search string (Less than or equal)
func (s *Search) Lte(key, value string) {
	s.lte[key] = value
}

// Ne adds a $ne search string (Not equal)
func (s *Search) Ne(key, value string) {
	s.ne[key] = value
}

// In adds a $in search string that will make a search query of items
// in an array (Any item In array)
func (s *Search) In(key, value string) {
	s.in[key] = value
}

// All adds a $all search string that will make a search query that will
// match if all items in the array are in (All items in array)
func (s *Search) All(key, value string) {
	s.all[key] = value
}

// Like adds a $like search string (Like this string)
func (s *Search) Like(key, value string) {
	s.like[key] = value
}

// QueryString returns the query string to append to the url we are searchin for
func (s *Search) QueryString() string {
	return ""
}
