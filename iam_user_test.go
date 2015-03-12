package corbel

import (
	"testing"
	"time"
)

func TestIAMUser(t *testing.T) {

	var (
		client *Client
		err    error
	)

	client, err = NewClientForEnvironment(
		nil,
		"qa",
		"22b0e55f",
		"test-client-full",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	if err != nil {
		t.Errorf("Error instancing client. Got: %v  Want: nil", err)
	}

	err = client.IAM.OauthToken()
	if err != nil {
		t.Errorf("Error getting token. Got: %v  Want: nil", err)
	}

	now := time.Now()

	anUserProperties := make(map[string]interface{})
	anUserProperties["string"] = "test string"
	anUserProperties["integer"] = 123456
	anUserProperties["float"] = 1.23
	anUserProperties["date"] = now

	anUser := IAMUser{
		Domain:      "silkroad-qa",
		Username:    "corbel-go",
		Email:       "corbel-go@corbel.org",
		FirstName:   "Corbel",
		LastName:    "Go",
		ProfileURL:  "http://corbel.org/corbel-go",
		PhoneNumber: "555-555-555",
		Scopes:      []string{},
		Properties:  anUserProperties,
		Country:     "Somewhere",
	}

	err = client.IAM.Add(&anUser)
	if err != nil {
		t.Errorf("Error creating user. Got: %v  Want: nil", err)
	}

	search := client.IAM.Search()
	search.Query.Eq["username"] = "corbel-go"

	var arrUsers []IAMUser

	err = search.Page(0, &arrUsers)
	if err != nil {
		t.Errorf("Error searching users. Got: %v  Want: nil", err)
	}
	if got, want := len(arrUsers), 1; got != want {
		t.Errorf("Error on search. Got: %v. Expect %v user.", got, want)
	}
	if arrUsers[0].Username != anUser.Username {
		t.Errorf("Error user found != user created")
	}

	anUser2 := IAMUser{}
	err = client.IAM.Get(arrUsers[0].ID, &anUser2)
	if err != nil {
		t.Errorf("Error getting users. Got: %v  Want: nil", err)
	}
	if anUser.FirstName != anUser2.FirstName {
		t.Errorf("Error user getted != user created")
	}

	anUser2.Country = "Internet"
	err = client.IAM.Update(anUser2.ID, &anUser2)
	if err != nil {
		t.Errorf("Error updating users. Got: %v  Want: nil", err)
	}

	anUser3 := IAMUser{}
	err = client.IAM.Get(anUser2.ID, &anUser3)
	if err != nil {
		t.Errorf("Error getting users. Got: %v  Want: nil", err)
	}
	if anUser2.Country != anUser3.Country {
		t.Errorf("User did not updated successfully")
	}

	err = client.IAM.Delete(anUser3.ID)
	if err != nil {
		t.Errorf("Error deleting users. Got: %v  Want: nil", err)
	}

}
