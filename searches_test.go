package corbel

import "testing"

func TestSearchesQueryStrings(t *testing.T) {
	query := newQuery()

	// $eq
	query.Eq["name"] = "testName"
	query.Eq["surname"] = "testSurname"

	queryString := query.string()
	if got, want := queryString, `[{"$eq":{"name":"testName","surname":"testSurname"}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $gt
	query.Gt["age"] = 30
	queryString = query.string()
	if got, want := queryString, `[{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $gte
	query.Gte["otherValue"] = 30
	queryString = query.string()
	if got, want := queryString, `[{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30},"$gte":{"otherValue":30}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $lt
	query.Lt["age"] = 50
	queryString = query.string()
	if got, want := queryString, `[{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30},"$gte":{"otherValue":30},"$lt":{"age":50}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	// $lte
	query.Lte["otherValue"] = 50
	queryString = query.string()
	if got, want := queryString, `[{"$eq":{"name":"testName","surname":"testSurname"},"$gt":{"age":30},"$gte":{"otherValue":30},"$lt":{"age":50},"$lte":{"otherValue":50}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	query = newQuery()
	query.Ne["key"] = "value"

	queryString = query.string()
	if got, want := queryString, `[{"$ne":{"key":"value"}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	query.In["key"] = []string{"value", "value2", "value3"}
	queryString = query.string()
	if got, want := queryString, `[{"$ne":{"key":"value"},"$in":{"key":["value","value2","value3"]}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	query.All["key"] = []string{"value", "value2", "value3"}
	queryString = query.string()
	if got, want := queryString, `[{"$ne":{"key":"value"},"$in":{"key":["value","value2","value3"]},"$all":{"key":["value","value2","value3"]}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	query = newQuery()

	query.Like["title"] = "gopher"
	queryString = query.string()
	if got, want := queryString, `[{"$like":{"title":"gopher"}}]`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

}

func TestSortsQueryStrings(t *testing.T) {
	sort := newSort()
	sort.Asc = []string{"firstName", "lastName"}

	queryString := sort.string()
	if got, want := queryString, `{"firstName":"asc","lastName":"asc"}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	sort.Desc = []string{"other", "andAnother"}
	queryString = sort.string()
	if got, want := queryString, `{"andAnother":"desc","firstName":"asc","lastName":"asc","other":"desc"}`; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}
}

func TestSearchQueryString(t *testing.T) {
	client, _ := NewClient(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	search := NewSearch(client, "resources", "/v1.0/resource/test:Collection")

	if got, want := search.PageSize, 10; got != want {
		t.Errorf("Error default PerPage. Got: %v, Want: %v", got, want)
	}

	searchQueryString := search.Query.string()
	if got, want := searchQueryString, ``; got != want {
		t.Errorf("Error in search Query String: Got: %v, Want: %v", got, want)
	}

	searchSortString := search.Sort.string()
	if got, want := searchSortString, ``; got != want {
		t.Errorf("Error in search Sort String: Got: %v, Want: %v", got, want)
	}

	opts := &SearchListOptions{
		APIPageSize: 10,
	}
	if got, want := search.queryString(opts), `/v1.0/resource/test:Collection?api%3ApageSize=10`; got != want {
		t.Errorf("Error in query string for the search: Got: %v, Want: %v", got, want)
	}

	search.Query.Eq["firstName"] = "testName"
	search.Sort.Asc = []string{"firstName"}
	opts = &SearchListOptions{
		APIPageSize: 10,
		APIQuery:    search.Query.string(),
	}
	if got, want := search.queryString(opts), `/v1.0/resource/test:Collection?api%3ApageSize=10&api%3Aquery=%5B%7B%22%24eq%22%3A%7B%22firstName%22%3A%22testName%22%7D%7D%5D`; got != want {
		t.Errorf("Error in query string for the search: Got: %v, Want: %v", got, want)
	}

	opts.APIPageSize = 20
	if got, want := search.queryString(opts), `/v1.0/resource/test:Collection?api%3ApageSize=20&api%3Aquery=%5B%7B%22%24eq%22%3A%7B%22firstName%22%3A%22testName%22%7D%7D%5D`; got != want {
		t.Errorf("Error in query string for the search: Got: %v, Want: %v", got, want)
	}
}

func TestSearchPagingAndAggregation(t *testing.T) {
	var (
		err   error
		count int
		avg   float64
		sum   float64
	)

	client, _ := NewClient(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		1000)

	type ResourceForTest struct {
		ID   string  `json:"id, omitempty"`
		Key1 string  `json:"key1,omitempty"`
		Key2 int     `json:"key2,omitempty"`
		Key3 float64 `json:"key3,omitempty"`
		Key4 bool    `json:"key4,omitempty"`
	}

	_ = client.IAM.OauthToken()

	testResource1 := ResourceForTest{
		Key1: "test1",
		Key2: 123456,
		Key3: 1.23,
		Key4: true,
	}

	testResource2 := ResourceForTest{
		Key1: "test2",
		Key2: 123456,
		Key3: 1.25,
		Key4: false,
	}

	err = client.Resources.AddToCollection("test:GoCollection", &testResource1)
	if err != nil {
		t.Errorf("Error adding to collection. Expected nil, Got: %v", err)
	}

	err = client.Resources.AddToCollection("test:GoCollection", &testResource2)
	if err != nil {
		t.Errorf("Error adding to collection. Expected nil, Got: %v", err)
	}

	search := NewSearch(client, "resources", "/v1.0/resource/test:GoCollection")
	count, err = search.Count("key1")
	if err != nil {
		t.Errorf("Error searching search.Count 1. Expected nil, Got: %v", err)
	}
	if got, want := count, 2; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	count, err = search.CountAll()
	if err != nil {
		t.Errorf("Error searching search.CountAll 1. Expected nil, Got: %v", err)
	}
	if got, want := count, 2; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	avg, err = search.Average("key3")
	if err != nil {
		t.Errorf("Error searching search.Average 1. Expected nil, Got: %v", err)
	}
	if got, want := avg, 1.24; got != want {
		t.Errorf("Error with average. Got: %v, Want: %v", got, want)
	}

	sum, err = search.Sum("key2")
	if err != nil {
		t.Errorf("Error searching search.Sum 1. Expected nil, Got: %v", err)
	}
	if got, want := sum, 246912.0; got != want {
		t.Errorf("Error with sum. Got: %v, Want: %v", got, want)
	}

	search.Query.Eq["key1"] = "test1"

	count, err = search.Count("key1")
	if err != nil {
		t.Errorf("Error searching search.Count 2. Expected nil, Got: %v", err)
	}
	if got, want := count, 1; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	count, err = search.CountAll()
	if err != nil {
		t.Errorf("Error searching search.CountAll 2. Expected nil, Got: %v", err)
	}
	if got, want := count, 1; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	avg, err = search.Average("key3")
	if err != nil {
		t.Errorf("Error searching search.Average 2. Expected nil, Got: %v", err)
	}
	if got, want := avg, 1.23; got != want {
		t.Errorf("Error with average. Got: %v, Want: %v", got, want)
	}

	sum, err = search.Sum("key2")
	if err != nil {
		t.Errorf("Error searching search.Sum 2. Expected nil, Got: %v", err)
	}
	if got, want := sum, 123456.0; got != want {
		t.Errorf("Error with sum. Got: %v, Want: %v", got, want)
	}

	search.Query.Eq["key1"] = "testX"

	count, err = search.Count("key1")
	// if err.Error() != "Not Found" {
	// 	t.Errorf("Error searching search.Count 3. Expected \"Not Found\", Got: %s", err.Error())
	// }
	if got, want := count, 0; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	count, err = search.CountAll()
	if err != nil {
		t.Errorf("Error searching search.CountAll 3. Expected nil, Got: %v", err)
	}
	if got, want := count, 0; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	avg, err = search.Average("key3")
	// if err != nil {
	// 	t.Errorf("Error searching search.Average 3. Expected nil, Got: %v", err)
	// }
	if got, want := avg, 0.0; got != want {
		t.Errorf("Error with average. Got: %v, Want: %v", got, want)
	}

	sum, err = search.Sum("key2")
	// if err != nil {
	// 	t.Errorf("Error searching search.Sum 3. Expected nil, Got: %v", err)
	// }
	if got, want := sum, 0.0; got != want {
		t.Errorf("Error with sum. Got: %v, Want: %v", got, want)
	}

	var (
		arrTestResources []ResourceForTest
		searchResource   *Search
	)

	searchResource = client.Resources.SearchCollection("test:GoCollection")
	searchResource.Query.Eq["key1"] = "test1"
	count, err = searchResource.CountAll()
	if err != nil {
		t.Errorf("Error searching searchResource.CountAll 4. Expected nil, Got: %v", err)
	}
	if got, want := count, 1; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

	searchResource = client.Resources.SearchCollection("test:GoCollection")
	count, err = searchResource.CountAll()
	if err != nil {
		t.Errorf("Error searching searchResource.CountAll 4. Expected nil, Got: %v", err)
	}
	if got, want := count, 2; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}
	searchResource.Sort.Desc = []string{"key1"}

	err = searchResource.Page(0, &arrTestResources)
	if err != nil {
		t.Errorf("Error searching err = searchResource.Page(0, arrTestResources). Expected nil, Got: %v", err)
	}
	if got, want := len(arrTestResources), 2; got != want {
		t.Errorf("Error with len(arrTestResources). Got: %v, Want: %v", got, want)
	}
	if got, want := arrTestResources[0].Key1, "test2"; got != want {
		t.Errorf("Error: Desc sort failed. Got: %v, Want: %v", got, want)
	}
	if got, want := arrTestResources[1].Key1, "test1"; got != want {
		t.Errorf("Error: Desc sort failed. Got: %v, Want: %v", got, want)
	}

	err = client.Resources.DeleteFromCollection("test:GoCollection", arrTestResources[0].ID)
	if err != nil {
		t.Errorf("Error searching err = client.Resources.DeleteFromCollection. Expected nil, Got: %v", err)
	}

	err = searchResource.Page(0, &arrTestResources)
	if got, want := len(arrTestResources), 1; got != want {
		t.Errorf("Error with len(arrTestResources) 2. Got: %v, Want: %v", got, want)
	}
	if got, want := arrTestResources[0].Key1, "test1"; got != want {
		t.Errorf("Error: Desc sort failed. Got: %v, Want: %v", got, want)
	}

	err = client.Resources.DeleteFromCollection("test:GoCollection", arrTestResources[0].ID)
	if err != nil {
		t.Errorf("Error searching err = client.Resources.DeleteFromCollection. Expected nil, Got: %v", err)
	}

	searchResource = client.Resources.SearchCollection("test:GoCollection")
	count, err = searchResource.CountAll()
	if err != nil {
		t.Errorf("Error searching searchResource.CountAll 5. Expected nil, Got: %v", err)
	}
	if got, want := count, 0; got != want {
		t.Errorf("Error with count. Got: %v, Want: %v", got, want)
	}

}
