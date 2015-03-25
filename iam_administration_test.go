package corbel

import (
	"os"
	"testing"
)

func TestIAMDomain(t *testing.T) {

	if os.Getenv("IAM_CLIENTID") == "" || os.Getenv("IAM_CLIENTSECRET") == "" || os.Getenv("IAM_CLIENT_DOMAIN") == "" {
		t.Skip("Skipping test since no valid keys passed to the test.")
	}

	var (
		client       *Client
		err          error
		sourceDomain IAMDomain
		targetDomain IAMDomain
		arrDomains   []IAMDomain
	)

	client, err = NewClientForEnvironment(
		nil,
		"qa",
		os.Getenv("IAM_CLIENTID"),
		"test-client-full",
		os.Getenv("IAM_CLIENTSECRET"),
		"iam:comp:root",
		os.Getenv("IAM_CLIENT_DOMAIN"),
		"HS256",
		10)

	if err != nil {
		t.Errorf("Error instancing client. Got: %v  Want: nil", err)
	}

	err = client.IAM.OauthToken()
	if err != nil {
		t.Errorf("Error getting token. Got: %v  Want: nil", err)
	}

	sourceDomain = IAMDomain{
		ID:                "corbel-go-test-domain",
		Description:       "Domain for test on corbel-go",
		Scopes:            []string{"corbel-go-scope-1", "corbel-go-scope2"},
		DefaultScopes:     []string{"corbel-go-scope-1"},
		AuthURL:           "http://corbel-go.org",
		UserProfileFields: []string{"firstName", "lastName"},
	}

	err = client.IAM.DomainAdd(&sourceDomain)
	if err != nil {
		t.Errorf("Error adding domain. Got: %v  Want: nil", err)
	}

	searchDomain := client.IAM.DomainSearch()
	searchDomain.Query.Eq["id"] = "corbel-go-test-domain"
	err = searchDomain.Page(0, &arrDomains)
	if err != nil {
		t.Errorf("Error searching domains. Got: %v  Want: nil", err)
	}

	if got, want := len(arrDomains), 1; got != want {
		t.Errorf("Wrong number of domains returned on the search. Got: %v. Want: %v.", got, want)
	}

	if got, want := sourceDomain.ID, arrDomains[0].ID; got != want {
		t.Errorf("Data returned on search does not match with the data inserted. Got: %v. Want: %v.", got, want)
	}

	sourceDomain.Description = "new description"

	err = client.IAM.DomainUpdate(sourceDomain.ID, &sourceDomain)
	if err != nil {
		t.Errorf("Error updating domain. Got: %v  Want: nil", err)
	}

	err = client.IAM.DomainGet(sourceDomain.ID, &targetDomain)
	if err != nil {
		t.Errorf("Error getting domain. Got: %v  Want: nil", err)
	}

	if got, want := sourceDomain.Description, targetDomain.Description; got != want {
		t.Errorf("Error with data returned. Got: %v. Want: %v", got, want)
	}

	err = client.IAM.DomainDelete(sourceDomain.ID)
	if err != nil {
		t.Errorf("Error deletting domain. Got: %v  Want: nil", err)
	}

}
