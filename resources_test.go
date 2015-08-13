package corbel

import "testing"

func TestResourcesAcl(t *testing.T) {
	client, err := NewClientForEnvironment(
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
	if err != nil {
		t.Errorf("GetToken must not fail. Got: %v  Want: nil", err)
	}

	type TestAclResource struct {
		*ResourceBasic
		Name string `json:"name,omitempty"`
	}

	resAcl := &TestAclResource{
		ResourceBasic: &ResourceBasic{},
		Name:          "Prueba ACL",
	}

	_, err = client.Resources.AddToCollection("test:GoTestResource", resAcl)
	if err != nil {
		t.Errorf("Failed to AddFromCollection to a struct. Got: %v  Want: nil", err)
	}

	err = resAcl.AllowAll(AclRead)
	if err != nil {
		t.Errorf("Failed to AllowAll. Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 1 || resAcl.ACL["ALL"] != "READ" {
		t.Errorf("Failed to UpdateResourceACL . Got: %d/%s  Want: 1/READ", len(resAcl.ACL), resAcl.ACL["ALL"])
	}

	err = resAcl.AllowUser("user1", AclWrite)
	if err != nil {
		t.Errorf("Failed to AllowUser. Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 2 {
		t.Errorf("Failed to UpdateResourceACL . Got: %d  Want: 2", len(resAcl.ACL))
	}

	err = resAcl.AllowUser("user1", AclRead)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL . Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 2 || resAcl.ACL["user1"] != "READ" {
		t.Errorf("Failed to UpdateResourceACL . Got: %d/%s  Want: 2/READ", len(resAcl.ACL), resAcl.ACL["user1"])
	}

	err = resAcl.RevokeAllPermission()
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL . Got: %v  Want: nil", err)
	}
	err = client.Resources.GetFromCollection("test:GoTestResource", resAcl.ID, resAcl)
	if err != nil {
		t.Errorf("Failed to UpdateResourceACL (GetResource) . Got: %v  Want: nil", err)
	}
	if len(resAcl.ACL) != 1 || resAcl.ACL["ALL"] != "" {
		t.Errorf("Failed to UpdateResourceACL . Got: %d/%s  Want: 1/", len(resAcl.ACL), resAcl.ACL["ALL"])
	}

	err = client.Resources.DeleteFromCollection("test:GoTestResource", resAcl.ID)
	if err != nil {
		t.Errorf("Failed to DeleteFromCollection to a struct. Got: %v  Want: nil", err)
	}
}
