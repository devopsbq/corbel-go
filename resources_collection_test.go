package corbel

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestResourcesAddToCollection(t *testing.T) {

	var (
		client *Client
		search *Search
		err    error
	)

	endpoints := map[string]string{"iam": "https://iam-int.bqws.io", "resources": "https://resources-int.bqws.io"}
	client, err = NewClient(
		nil,
		endpoints,
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		300)

	err = client.IAM.OauthToken()
	if err != nil {
		t.Errorf("GetToken must not fail. Got: %v  Want: nil", err)
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

	_, err = client.Resources.AddToCollection("test:GoTestResource", test1)
	if err != nil {
		t.Errorf("Failed to AddToCollection a struct. Got: %v  Want: nil", err)
	}
	search = client.Resources.SearchCollection("test:GoTestResource")
	err = search.Page(0, &arrResourceForTest)
	if err != nil {
		t.Errorf("Failed to SearchCollection an array of structs. Got: %v  Want: nil", err)
	}
	err = client.Resources.DeleteFromCollection("test:GoTestResource", arrResourceForTest[0].ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection from item in an array of structs. Got: %v  Want: nil", err)
	}
}

func TestResourcesGetFromCollection(t *testing.T) {

	var (
		client *Client
		err    error
		search *Search
	)

	endpoints := map[string]string{"iam": "https://iam-int.bqws.io", "resources": "https://resources-int.bqws.io"}
	client, err = NewClient(
		nil,
		endpoints,
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		300)

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

	var arrResourceForTest []ResourceForTest

	test1 := ResourceForTest{
		Key1: "test string",
		Key2: 123456,
		Key3: 1.123456,
		Key4: true,
	}

	_, err = client.Resources.AddToCollection("test:GoTestResource", &test1)
	if err != nil {
		t.Errorf("Failed to AddFromCollection to a struct. Got: %v  Want: nil", err)
	}

	search = client.Resources.SearchCollection("test:GoTestResource")
	err = search.Page(0, &arrResourceForTest)
	if err != nil {
		t.Errorf("Failed to SearchCollection.Page to an array of structs. Got: %v  Want: nil", err)
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
	if err != nil {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: nil", err)
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

	test2.Key1 = "new string"
	test2.Key2 = 654321
	test2.Key3 = 654.321
	test2.Key4 = false

	err = client.Resources.UpdateInCollection("test:GoTestResource", test2.ID, &test2)
	if err != nil {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: nil", err)
	}

	test3 := ResourceForTest{}
	err = client.Resources.GetFromCollection("test:GoTestResource", test2.ID, &test3)

	if got, want := test3.ID, test2.ID; got != want {
		t.Errorf("Failed to GetFromCollection after UpdateInCollection to a struct. Got: %v  Want: %v", got, want)
	}
	if got, want := test3.Key1, test2.Key1; got != want {
		t.Errorf("Failed to GetFromCollection after UpdateInCollection to a struct. Got: %v  Want: %v", got, want)
	}
	if got, want := test3.Key2, test2.Key2; got != want {
		t.Errorf("Failed to GetFromCollection after UpdateInCollection to a struct. Got: %v  Want: %v", got, want)
	}
	if got, want := test3.Key3, test2.Key3; got != want {
		t.Errorf("Failed to GetFromCollection after UpdateInCollection to a struct. Got: %v  Want: %v", got, want)
	}
	if got, want := test3.Key4, test2.Key4; got != want {
		t.Errorf("Failed to GetFromCollection after UpdateInCollection to a struct. Got: %v  Want: %v", got, want)
	}

	err = client.Resources.DeleteFromCollection("test:GoTestResource", test3.ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection to a struct. Got: %v  Want: nil", err)
	}

	type ResourceWithAcl struct {
		ID   string             `json:"id,omitempty"`
		Name string             `json:"name,omitempty"`
		ACL  map[string]UserACL `json:"_acl,omitempty"`
	}

	resAcl := &ResourceWithAcl{
		Name: "Prueba ACL",
	}

	id, err := client.Resources.AddToCollection("test:GoTestResource", resAcl)
	if err != nil {
		t.Errorf("Failed to AddFromCollection to a struct. Got: %v  Want: nil", err)
	}
	s := strings.Split(id, "/")
	resAcl.ID = s[len(s)-1]
	resAcl.ACL = make(map[string]UserACL)

	resAcl.ACL["ALL"] = UserACL{Permission: "READ", Properties: make(map[string]interface{})}

	b, err := json.Marshal(resAcl.ACL)
	fmt.Println(string(b))
	err = client.Resources.UpdateResourceACL("test:GoTestResource", resAcl.ID, resAcl.ACL)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL . Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 1 || resAcl.ACL["ALL"].Permission != "READ" {
		t.Errorf("Failed to UpdateResourceACL . Got: %d/%s  Want: 1/READ", len(resAcl.ACL), resAcl.ACL["ALL"])
	}

	resAcl.ACL["user1"] = UserACL{Permission: "WRITE", Properties: make(map[string]interface{})}
	err = client.Resources.UpdateResourceACL("test:GoTestResource", resAcl.ID, resAcl.ACL)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL . Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 2 {
		t.Errorf("Failed to UpdateResourceACL . Got: %d  Want: 2", len(resAcl.ACL))
	}

	resAcl.ACL["user1"] = UserACL{Permission: "READ", Properties: make(map[string]interface{})}
	err = client.Resources.UpdateResourceACL("test:GoTestResource", resAcl.ID, resAcl.ACL)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL . Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 2 || resAcl.ACL["user1"].Permission != "READ" {
		t.Errorf("Failed to UpdateResourceACL . Got: %d/%s  Want: 2/READ", len(resAcl.ACL), resAcl.ACL["user1"])
	}

	delete(resAcl.ACL, "ALL")
	err = client.Resources.UpdateResourceACL("test:GoTestResource", resAcl.ID, resAcl.ACL)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL . Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 1 || resAcl.ACL["ALL"].Permission != "" {
		t.Errorf("Failed to UpdateResourceACL . Got: %d/%s  Want: 1/", len(resAcl.ACL), resAcl.ACL["ALL"])
	}

	err = client.Resources.DeleteFromCollection("test:GoTestResource", resAcl.ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection to a struct. Got: %v  Want: nil", err)
	}
}
