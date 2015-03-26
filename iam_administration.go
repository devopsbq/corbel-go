package corbel

import (
	"fmt"
	"net/http"
)

// IAMAuthConfiguration is the representation of an AuthConfiguration object used by IAM
type IAMAuthConfiguration struct {
	// Type defined the Auth Configuration type defined
	Type        string `json:"type"`
	RedirectURL string `json:"redirectUri"`
	// ClientId used for Facebook, Google and Corbel Oauth 2.0
	ClientID     string `json:"clientID,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	// OAuthServerURL is the Oauth URL to use in Corbel Oauth
	OAuthServerURL string `json:"oAuthServerUrl,omitempty"`
	// CounsumerKey used for Twitter Oauth
	ConsumerKey    string `json:"consumerKey,omitempty"`
	ConsumerSecret string `json:"consumerSecret,omitempty"`
}

// IAMDomain is the representation of an Domain object used by IAM
type IAMDomain struct {
	ID                 string                          `json:"id"`
	Description        string                          `json:"description,omitempty"`
	AuthURL            string                          `json:"authUrl, omitempty"`
	AllowedDomains     string                          `json:"allowedDomains"`
	Scopes             []string                        `json:"scopes,omitempty"`
	DefaultScopes      []string                        `json:"defaultScopes,omitempty"`
	AuthConfigurations map[string]IAMAuthConfiguration `json:"authConfigurations,omitempty"`
	UserProfileFields  []string                        `json:"userProfileFields,omitempty"`
	CreatedDate        int                             `json:"createdDate,omitempty"`
	CreatedBy          string                          `json:"createdBy,omitempty"`
}

// IAMRule is the representation of a Rule for a Scope object used by IAM
type IAMRule struct {
	MediaTypes []string `json:"mediaTypes"`
	Methods    []string `json:"methods"`
	Type       string   `json:"type"`
	URI        string   `json:"uri"`
	TokenType  string   `json:"tokenType"`
}

// IAMClient is the representation of a Client object used by IAM
type IAMClient struct {
	ID                       string   `json:"id,omitempty"`
	Key                      string   `json:"key"`
	Name                     string   `json:"name"`
	Domain                   string   `json:"domain"`
	Version                  string   `json:"version,omitempty"`
	SignatureAlgorithm       string   `json:"signatureAlgorithm,omitempty"`
	Scopes                   []string `json:"scopes,omitempty"`
	ClientSideAuthentication bool     `json:"clientSideAuthentication"`
	ResetURL                 string   `json:"resetUrl,omitempty"`
	ResetNotificationID      string   `json:"resetNotificationId,omitempty"`
}

// IAMScope is the representation of a Scope object used by IAM
type IAMScope struct {
	ID         string            `json:"id"`
	Audience   string            `json:"audience"`
	Type       string            `json:"type"`
	Scopes     []string          `json:"scopes,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
	Rules      []IAMRule         `json:"rules,omitempty"`
}

// DomainAdd adds an Domain defined struct to the platform
func (i *IAMService) DomainAdd(domain *IAMDomain) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("POST", "iam", "/v1.0/domain/", domain)
	return returnErrorHTTPSimple(i.client, req, err, 201)
}

// DomainUpdate updates an domain by using IAMDomain
func (i *IAMService) DomainUpdate(id string, domain *IAMDomain) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("PUT", "iam", fmt.Sprintf("/v1.0/domain/%s", id), domain)
	return returnErrorHTTPSimple(i.client, req, err, 204)
}

// DomainGet gets the desired IAMUdomain from the domain by id
func (i *IAMService) DomainGet(id string, domain *IAMDomain) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/domain/%s", id), nil)
	return returnErrorHTTPInterface(i.client, req, err, domain, 200)
}

// DomainDelete deletes the desired domain from IAM by id
func (i *IAMService) DomainDelete(id string) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/domain/%s", id), nil)
	return returnErrorHTTPSimple(i.client, req, err, 204)
}

// DomainSearch gets the desired objects in base of a search query
func (i *IAMService) DomainSearch() *Search {
	return NewSearch(i.client, "iam", "/v1.0/domain")
}

// ClientAdd adds an Client defined struct to the platform
func (i *IAMService) ClientAdd(client *IAMClient) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("POST", "iam", fmt.Sprintf("/v1.0/domain/%s/client/", client.Domain), client)
	return returnErrorHTTPSimple(i.client, req, err, 201)
}

// ClientUpdate updates an client by using IAMClient
func (i *IAMService) ClientUpdate(id string, client *IAMClient) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("PUT", "iam", fmt.Sprintf("/v1.0/domain/%s/client/%s", client.Domain, id), client)
	return returnErrorHTTPSimple(i.client, req, err, 204)
}

// ClientGet gets the desired IAMClient
func (i *IAMService) ClientGet(domainName, id string, client *IAMClient) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/domain/%s/client/%s", domainName, id), nil)
	return returnErrorHTTPInterface(i.client, req, err, client, 200)
}

// ClientDelete deletes the desired client from IAM by id
func (i *IAMService) ClientDelete(domainName, id string) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/domain/%s/client/%s", domainName, id), nil)
	return returnErrorHTTPSimple(i.client, req, err, 204)
}

// ClientSearch gets the desired objects in base of a search query
func (i *IAMService) ClientSearch(domainName string) *Search {
	return NewSearch(i.client, "iam", fmt.Sprintf("/v1.0/domain/%s/client", domainName))
}

// ScopeAdd adds an Scope defined struct to the platform
func (i *IAMService) ScopeAdd(scope *IAMScope) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("POST", "iam", "/v1.0/scope/", scope)
	return returnErrorHTTPSimple(i.client, req, err, 201)
}

// ScopeUpdate updates an scope by using IAMScope
func (i *IAMService) ScopeUpdate(scope *IAMScope) error {
	return i.ScopeAdd(scope)
}

// ScopeGet gets the desired IAMScope from the scope by id
func (i *IAMService) ScopeGet(id string, scope *IAMScope) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/scope/%s", id), nil)
	return returnErrorHTTPInterface(i.client, req, err, scope, 200)
}

// ScopeDelete deletes the desired scope from IAM by id
func (i *IAMService) ScopeDelete(id string) error {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/scope/%s", id), nil)
	return returnErrorHTTPSimple(i.client, req, err, 204)
}
