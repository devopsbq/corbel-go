package corbel

import (
	"strings"
	"testing"
)

func TestIAMGroup(t *testing.T) {
	client, err := NewClientForEnvironment(
		nil,
		"int",
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

	g := &IAMGroup{Name: "Prueba", Scopes: []string{}}
	loc, err := client.IAM.GroupAdd(g)
	if err != nil {
		t.Errorf("Error creating group. Got: %v  Want: nil", err)
	}
	id := strings.Split(loc, "/")

	getg := new(IAMGroup)
	err = client.IAM.GroupGet(id[len(id)-1], getg)
	if err != nil || getg.Name != "Prueba" {
		t.Errorf("Error retrieving group. Got: %v  Want: nil", err)
	}
	if err = client.IAM.GroupSetScopes(id[len(id)-1], []string{"silkroad:comp:base"}); err != nil {
		t.Errorf("Error adding scope to group. Got: %v  Want: nil", err)
	}
	if err = client.IAM.GroupDeleteScope(id[len(id)-1], "silkroad:comp:base"); err != nil {
		t.Errorf("Error adding scope to group. Got: %v  Want: nil", err)
	}
	if err = client.IAM.GroupDelete(id[len(id)-1]); err != nil {
		t.Errorf("Error deleting group. Got: %v  Want: nil", err)
	}

	//var groups []*IAMGroup
	//err = client.IAM.GroupGetAll(groups)
}
