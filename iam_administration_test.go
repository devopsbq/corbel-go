package corbel

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestIAM(t *testing.T) {

	if os.Getenv("IAM_CLIENTID") == "" || os.Getenv("IAM_CLIENTSECRET") == "" || os.Getenv("IAM_CLIENT_DOMAIN") == "" {
		fmt.Println("Skipping TestIAM")
		t.Skip("Skipping test since no valid keys passed to the test.")
	}

	var (
		client       *Client
		err          error
		sourceDomain IAMDomain
		targetDomain IAMDomain
		arrDomains   []IAMDomain
		sourceClient IAMClient
		targetClient IAMClient
		arrClients   []IAMClient
		sourceScope  IAMScope
		targetScope  IAMScope
		location     string
	)

	endpoints := map[string]string{"iam": "https://iam-int.bqws.io", "resources": "https://resources-int.bqws.io"}
	client, err = NewClient(
		nil,
		endpoints,
		os.Getenv("IAM_CLIENTID"),
		"iam-client",
		os.Getenv("IAM_CLIENTSECRET"),
		"silkroad-qa:root",
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

	location, err = client.IAM.DomainAdd(&sourceDomain)
	if err != nil {
		t.Errorf("Error adding domain. Got: %v  Want: nil", err)
	}

	if location == "" {
		t.Errorf("Location is empty: '%s'", location)
	}

	searchDomain := client.IAM.DomainSearch()
	searchDomain.Query.Eq["id"] = "silkroad-qa:corbel-go-test-domain"
	err = searchDomain.Page(0, &arrDomains)
	if err != nil {
		t.Errorf("Error searching domains. Got: %v  Want: nil", err)
	}

	if got, want := len(arrDomains), 1; got != want {
		t.Errorf("Wrong number of domains returned on the search. Got: %v. Want: %v.", got, want)
	}

	// sourceDomain new id = parentDomain + domainID
	sourceDomain.ID = "silkroad-qa:corbel-go-test-domain"

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

	// Clients
	sourceClient = IAMClient{
		Name:               "corbel-go-test-client",
		Domain:             sourceDomain.ID,
		Key:                "abcdefabcdefabcdefabcdefabcdefabcdefabcdef",
		SignatureAlgorithm: "HS256",
	}

	location, err = client.IAM.ClientAdd(&sourceClient)
	if err != nil {
		t.Errorf("Error creating client. Got: %v  Want: nil", err)
	}

	searchClient := client.IAM.ClientSearch(sourceDomain.ID)
	searchClient.Query.Eq["name"] = "corbel-go-test-client"
	err = searchClient.Page(0, &arrClients)
	if err != nil {
		t.Errorf("Error searching clients. Got: %v  Want: nil", err)
	}

	if got, want := len(arrClients), 1; got != want {
		t.Errorf("Wrong number of domains returned on the search. Got: %v. Want: %v.", got, want)
	}

	if got, want := extractID(location), arrClients[0].ID; got != want {
		t.Errorf("Data returned on search does not match with the data inserted. Got: %v. Want: %v.", got, want)
	}

	sourceClient.Scopes = []string{"corbel:go:test"}

	err = client.IAM.ClientUpdate(extractID(location), &sourceClient)
	if err != nil {
		t.Errorf("Error updating client. Got: %v  Want: nil", err)
	}

	err = client.IAM.ClientGet(sourceDomain.ID, extractID(location), &targetClient)
	if err != nil {
		t.Errorf("Error getting client. Got: %v  Want: nil", err)
	}

	if got, want := sourceClient.Key, targetClient.Key; got != want {
		t.Errorf("Error with data returned. Got: %v. Want: %v", got, want)
	}

	// Scopes

	sourceSopeRule1 := IAMRule{}
	sourceSopeRule1.Methods = []string{"GET", "DELETE"}
	sourceSopeRule1.MediaTypes = []string{"application/json"}
	sourceSopeRule1.URI = "v.*/resource/cober:go:resource/{{id}}"

	sourceScope = IAMScope{
		ID:       "corbel:go:test",
		Audience: "http://resources.bqws.io",
		Type:     "http_access",
		Rules:    []IAMRule{sourceSopeRule1},
	}

	_, err = client.IAM.ScopeAdd(&sourceScope)
	if err != nil {
		t.Errorf("Error creating scope. Got: %v  Want: nil", err)
	}

	sourceScope.Scopes = []string{"corbel:go:test1"}

	_, err = client.IAM.ScopeUpdate(&sourceScope)
	if err != nil {
		t.Errorf("Error updating scope. Got: %v  Want: nil", err)
	}

	err = client.IAM.ScopeGet(sourceScope.ID, &targetScope)
	if err != nil {
		t.Errorf("Error getting scope. Got: %v  Want: nil", err)
	}

	if got, want := sourceScope.ID, targetScope.ID; got != want {
		t.Errorf("Error with data returned. Got: %v. Want: %v", got, want)
	}

	// Deletes
	err = client.IAM.ScopeDelete(sourceScope.ID)
	if err != nil {
		t.Errorf("Error deletting scope. Got: %v  Want: nil", err)
	}

	err = client.IAM.ClientDelete("corbel-go-test-domain", extractID(location))
	if err != nil {
		t.Errorf("Error deletting client. Got: %v  Want: nil", err)
	}

	err = client.IAM.DomainDelete(sourceDomain.ID)
	if err != nil {
		t.Errorf("Error deletting domain. Got: %v  Want: nil", err)
	}

}

func extractID(location string) string {
	var locationSplit []string

	locationSplit = strings.Split(location, "/")
	return locationSplit[len(locationSplit)-1]
}
