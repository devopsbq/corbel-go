package silkroad

import (
	"fmt"
	"testing"
)

func TestResourcesAddToCollection(t *testing.T) {

	var (
		client            *Client
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

	fmt.Println(string(jsonEncodedStruct))
	err = client.Resources.AddToCollection("test:GoTestResource", test1)
	if got := err; got != nil {
		t.Errorf("Failed to AddToCollection a struct. Got: %v  Want: nil", got)
	}
}

func TestResourcesGetFromCollection(t *testing.T) {

	var (
		client *Client
		err    error
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
		Key1 string  `json:"key1,omitempty"`
		Key2 int     `json:"key2,omitempty"`
		Key3 float64 `json:"key3,omitempty"`
		Key4 bool    `json:"key4,omitempty"`
	}

	test1 := ResourceForTest{}

	err = client.Resources.GetFromCollection("test:GoTestResource", "54fdda20e4b09ed6fc32fa82", &test1)
	if got := err; got != nil {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: nil", got)
	}

	if got, want := test1.Key1, "test string"; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	if got, want := test1.Key2, 123456; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	if got, want := test1.Key3, 1.123456; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

	if got, want := test1.Key4, true; got != want {
		t.Errorf("Failed to GetFromCollection to a struct. Got: %v  Want: %v", got, want)
	}

}
