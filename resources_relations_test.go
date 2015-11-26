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
		"int",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	err = client.IAM.OauthToken()
	if err != nil {
		t.Errorf("GetToken must not fail. Got: %v  Want: nil", err)
	}

	type ResourceForTest struct {
		ID   string  `json:"id,omitempty"`
		Key1 string  `json:"key1"`
		Key2 int     `json:"key2"`
		Key3 float64 `json:"key3"`
		Key4 bool    `json:"key4"`
	}

	var (
		arrResourceForTest     []ResourceForTest
		arrResourceForTestDest []ResourceForTest
		arrSearch              []ResourceForTest
	)
	test1 := ResourceForTest{
		Key1: "test string",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	_, err = client.Resources.AddToCollection("test:GoTestOrigin", &test1)
	if err != nil {
		t.Errorf("Failed to AddToCollection to a struct. Got: %v  Want: nil", err)
	}

	test2 := ResourceForTest{
		Key1: "test2",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	_, err = client.Resources.AddToCollection("test:GoTestDestination", &test2)
	if err != nil {
		t.Errorf("Failed to AddToCollection to a struct. Got: %v  Want: nil", err)
	}

	test3 := ResourceForTest{
		Key1: "test3",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	_, err = client.Resources.AddToCollection("test:GoTestDestination", &test3)
	if err != nil {
		t.Errorf("Failed to AddToCollection to a struct. Got: %v  Want: nil", err)
	}

	search = client.Resources.SearchCollection("test:GoTestOrigin")
	err = search.Page(0, &arrResourceForTest)
	if err != nil {
		t.Errorf("Failed to SearchCollection.Page to an array of structs. Got: %v  Want: nil", err)
	}
	if got, want := len(arrResourceForTest), 1; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	// ensure there is no error searching for a relation that not exists
	err = client.Resources.DeleteAllRelations("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation")
	if err != nil {
		t.Errorf("Failed to DeleteAllRelations. Got: %v  Want: nil", err)
	}

	search = client.Resources.SearchCollection("test:GoTestDestination")
	err = search.Page(0, &arrResourceForTestDest)
	if err != nil {
		t.Errorf("Failed to SearchCollection.Page to an array of structs. Got: %v  Want: nil", err)
	}
	if got, want := len(arrResourceForTestDest), 2; got != want {
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

	_, err = client.Resources.AddRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[0].ID, relData)
	if err != nil {
		t.Errorf("Failed to first AddRelation. Got: %v  Want: nil", err)
	}

	_, err = client.Resources.AddRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[1].ID, nil)
	if err != nil {
		t.Errorf("Failed to second AddRelation. Got: %v  Want: nil", err)
	}

	search = client.Resources.SearchRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation")
	err = search.Page(0, &arrSearch)
	if err != nil {
		t.Errorf("Failed to SearchRelation.Page to an array of structs. Got: %v  Want: nil", err)
	}
	if got, want := len(arrSearch), 2; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	var arrRelationData []RelationData

	search = client.Resources.SearchRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestOtherRelation")
	err = search.Page(0, &arrRelationData)
	if err != nil {
		t.Errorf("Failed to SearchRelation.Page to an array of structs. Got: %v  Want: nil", err)
	}
	if got, want := len(arrRelationData), 0; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	_, err = client.Resources.MoveRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[1].ID, 1)
	if err != nil {
		t.Errorf("Failed to MoveRelation. Got: %v  Want: nil", err)
	}

	type customRelationData struct {
		Order  float64                  `json:"_order,omitempty"`
		ID     string                   `json:"id,omitempty"`
		Links  []map[string]interface{} `json:"links, omitempty"`
		Field1 string                   `json:"field1"`
		Field2 string                   `json:"field2"`
	}

	var customRelData []customRelationData

	search = client.Resources.SearchRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation")
	search.Sort.Asc = []string{"_order"}
	err = search.Page(0, &customRelData)
	if err != nil {
		t.Errorf("Failed to SearchRelation.Page to an array of structs. Got: %v  Want: nil", err)
	}
	if got, want := len(customRelData), 2; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	// variables to store the result of GetFromRelationDefinition
	var (
		test4 ResourceForTest
		test5 ResourceForTest
	)

	err = client.Resources.GetFromRelationDefinition(customRelData[0].ID, &test4)
	if err != nil {
		t.Errorf("Failed to GetFromRelationDefinition to the correct struct. Got: %v  Want: nil", err)
	}

	err = client.Resources.GetFromRelationDefinition(customRelData[1].ID, &test5)
	if err != nil {
		t.Errorf("Failed to GetFromRelationDefinition to the correct struct. Got: %v  Want: nil", err)
	}

	// test2 == test5 && test3 == test4
	if got, want := test5.Key1, test2.Key1; got != want {
		t.Errorf("MoveRelation failed. Objects didn't move. Got: %v. Want: %v", got, want)
	}
	if got, want := test4.Key1, test3.Key1; got != want {
		t.Errorf("MoveRelation failed. Objects didn't move. Got: %v. Want: %v", got, want)
	}

	err = client.Resources.DeleteRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[0].ID)
	if err != nil {
		t.Errorf("Failed to DeleteRelation. Got: %v  Want: nil", err)
	}

	_, err = client.Resources.AddRelation("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation", "test:GoTestDestination", arrResourceForTestDest[0].ID, relData)
	if err != nil {
		t.Errorf("Failed to AddRelation. Got: %v  Want: nil", err)
	}

	err = client.Resources.DeleteAllRelations("test:GoTestOrigin", arrResourceForTest[0].ID, "test:GoTestRelation")
	if err != nil {
		t.Errorf("Failed to DeleteAllRelations. Got: %v  Want: nil", err)
	}

	err = client.Resources.DeleteFromCollection("test:GoTestOrigin", arrResourceForTest[0].ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection. Got: %v  Want: nil", err)
	}
	err = client.Resources.DeleteFromCollection("test:GoTestDestination", arrResourceForTestDest[0].ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection. Got: %v  Want: nil", err)
	}
	err = client.Resources.DeleteFromCollection("test:GoTestDestination", arrResourceForTestDest[1].ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection. Got: %v  Want: nil", err)
	}
}
