package corbel

import "testing"

func TestResourcesRelations(t *testing.T) {
	var (
		client *Client
		err    error
		search *Search
	)

	client, err = NewClientForEnvironment(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	err = client.IAM.OauthToken()
	if got := err; got != nil {
		t.Errorf("GetToken must not fail. Got: %v  Want: nil", got)
	}

	type ResourceForTest struct {
		ID   string  `json:"id,omitempty"`
		Key1 string  `json:"key1"`
		Key2 int     `json:"key2"`
		Key3 float64 `json:"key3"`
		Key4 bool    `json:"key4"`
	}

	var arrResourceForTest []ResourceForTest
	var arrResourceForTestDest []ResourceForTest
	test1 := ResourceForTest{
		Key1: "test string",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	err = client.Resources.AddToCollection("test:GoTestOrigin", &test1)
	if got := err; got != nil {
		t.Errorf("Failed to AddFromCollection to a struct. Got: %v  Want: nil", got)
	}

	test2 := ResourceForTest{
		Key1: "test string",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	err = client.Resources.AddToCollection("test:GoTestDestination", &test2)
	if got := err; got != nil {
		t.Errorf("Failed to AddFromCollection to a struct. Got: %v  Want: nil", got)
	}

	search = client.Resources.SearchCollection("test:GoTestOrigin")
	err = search.Page(0, &arrResourceForTest)
	if got := err; got != nil {
		t.Errorf("Failed to SearchCollection.Page to an array of structs. Got: %v  Want: nil", got)
	}
	if got, want := len(arrResourceForTest), 1; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	search = client.Resources.SearchCollection("test:GoTestDestination")
	err = search.Page(0, &arrResourceForTestDest)
	if got := err; got != nil {
		t.Errorf("Failed to SearchCollection.Page to an array of structs. Got: %v  Want: nil", got)
	}
	if got, want := len(arrResourceForTestDest), 1; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	type relationDataForTest struct {
		Field1 string `json:"field1"`
		Field2 string `json:"field2"`
	}

	relData := relationDataForTest{
		Field1: "data for field1",
		Field2: "data for field2",
	}

	err = client.Resources.AddRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[0].ID, relData)
	if got := err; got != nil {
		t.Errorf("Failed to AddRelation. Got: %v  Want: nil", got)
	}

	err = client.Resources.DeleteRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[0].ID)
	if got := err; got != nil {
		t.Errorf("Failed to DeleteRelation. Got: %v  Want: nil", got)
	}

	err = client.Resources.AddRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[0].ID, relData)
	if got := err; got != nil {
		t.Errorf("Failed to AddRelation. Got: %v  Want: nil", got)
	}

	err = client.Resources.DeleteAllRelations("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation")
	if got := err; got != nil {
		t.Errorf("Failed to DeleteAllRelations. Got: %v  Want: nil", got)
	}

	err = client.Resources.DeleteFromCollection("test:GoTestOrigin", arrResourceForTest[0].ID)
	if got := err; got != nil {
		t.Errorf("Failed to DeleteFromCollection. Got: %v  Want: nil", got)
	}
	err = client.Resources.DeleteFromCollection("test:GoTestDestination", arrResourceForTestDest[0].ID)
	if got := err; got != nil {
		t.Errorf("Failed to DeleteFromCollection. Got: %v  Want: nil", got)
	}
}
