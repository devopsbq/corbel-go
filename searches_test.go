package silkroad

import "testing"

func TestSearchesQueryStrings(t *testing.T) {
	search := NewSearch()

	// $eq
	search.Eq["name"] = "testName"
	search.Eq["surname"] = "testSurname"

	queryString, err := search.QueryString()

	if err != nil {
		t.Errorf("Unexpected error marshalling search query string. Got: %v", err)
	}

	if got, want := queryString, `{"$eq":{"name":"testName","surname":"testSurname"}}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $gt
	search.Gt["age"] = 30
	queryString, err = search.QueryString()

	if err != nil {
		t.Errorf("Unexpected error marshalling search query string. Got: %v", err)
	}

	if got, want := queryString, `{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30}}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $gte
	search.Gte["otherValue"] = 30
	queryString, err = search.QueryString()

	if err != nil {
		t.Errorf("Unexpected error marshalling search query string. Got: %v", err)
	}

	if got, want := queryString, `{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30},"$gte":{"otherValue":30}}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $lt
	search.Lt["age"] = 50
	queryString, err = search.QueryString()

	if err != nil {
		t.Errorf("Unexpected error marshalling search query string. Got: %v", err)
	}

	if got, want := queryString, `{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30},"$gte":{"otherValue":30},"$lt":{"age":50}}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $lte
	search.Lte["otherValue"] = 50
	queryString, err = search.QueryString()

	if err != nil {
		t.Errorf("Unexpected error marshalling search query string. Got: %v", err)
	}

	if got, want := queryString, `{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30},"$gte":{"otherValue":30},"$lt":{"age":50},"$lte":{"otherValue":50}}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}
}
