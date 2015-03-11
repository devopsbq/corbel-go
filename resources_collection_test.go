package corbel

import (
	"fmt"
	"testing"
)

func TestResourcesAddToCollection(t *testing.T) {

	var (
		client            *Client
		search            *Search
		err               error
		jsonEncodedStruct []byte
	)

	client, err = NewClient(
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
		Key1 string  `json:"key1,omitempty"`
		Key2 uint64  `json:"key2,omitempty"`
		Key3 float64 `json:"key3,omitempty"`
		Key4 bool    `json:"key4,omitempty"`
	}

	test1 := ResourceForTest{
		Key1: "test string",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	var arrResourceForTest []ResourceForTest

	fmt.Println(string(jsonEncodedStruct))
	err = client.Resources.AddToCollection("test:GoTestResource", test1)
	if got := err; got != nil {
		t.Errorf("Failed to AddToCollection a struct. Got: %v  Want: nil", got)
	}
	search = client.Resources.SearchCollection("test:GoTestResource")
	err = search.Page(0, &arrResourceForTest)
	if got := err; got != nil {
		t.Errorf("Failed to SearchCollection an array of structs. Got: %v  Want: nil", got)
	}
	err = client.Resources.DeleteFromCollection("test:GoTestResource", arrResourceForTest[0].ID)
	if got := err; got != nil {
		t.Errorf("Failed to DeleteFromCollection from item in an array of structs. Got: %v  Want: nil", got)
	}
}

func TestResourcesGetFromCollection(t *testing.T) {

	var (
		client *Client
		err    error
		search *Search
	)

	client, err = NewClient(
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
		Key1 string  `json:"key1,omitempty"`
		Key2 int     `json:"key2,omitempty"`
		Key3 float64 `json:"key3,omitempty"`
		Key4 bool    `json:"key4,omitempty"`
	}

	var arrResourceForTest []ResourceForTest

	test1 := ResourceForTest{
		Key1: "test string",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	err = client.Resources.AddToCollection("test:GoTestResource", &test1)
	if got := err; got != nil {
		t.Errorf("Failed to AddFromCollection to a struct. Got: %v  Want: nil", got)
	}

	search = client.Resources.SearchCollection("test:GoTestResource")
	err = search.Page(0, &arrResourceForTest)
	if got := err; got != nil {
		t.Errorf("Failed to SearchCollection.Page to an array of structs. Got: %v  Want: nil", got)
	}
	if got, want := len(arrResourceForTest), 1; got != want {
		t.Errorf("Bad number of structs returned. Got: %v. Want: %v", got, want)
	}

	if got, want := arrResourceForTest[0].Key1, test1.Key1; got != want {
		t.Errorf("Error with search. Object0 != Crafted Object. (key1) Got: %v. Want: %v", got, want)
	}
	if got, want := arrResourceForTest[0].Key2, test1.Key2; got != want {
		t.Errorf("Error with search. Object0 != Crafted Object. (key2) Got: %v. Want: %v", got, want)
	}
	if got, want := arrResourceForTest[0].Key3, test1.Key3; got != want {
		t.Errorf("Error with search. Object0 != Crafted Object. (key3) Got: %v. Want: %v", got, want)
	}
	if got, want := arrResourceForTest[0].Key4, test1.Key4; got != want {
		t.Errorf("Error with search. Object0 != Crafted Object. (key4) Got: %v. Want: %v", got, want)
	}

	test2 := ResourceForTest{}

	err = client.Resources.GetFromCollection("test:GoTestResource", arrResourceForTest[0].ID, &test2)
	if got := err; got != nil {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: nil", got)
	}

	if got, want := test2.Key1, test1.Key1; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	if got, want := test2.Key2, test1.Key2; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	if got, want := test2.Key3, test1.Key3; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	if got, want := test2.Key4, test1.Key4; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	err = client.Resources.DeleteFromCollection("test:GoTestResource", arrResourceForTest[0].ID)

}
