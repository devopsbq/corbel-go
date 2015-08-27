package corbel

import (
	"fmt"
	"net/http"
	"time"
)

var (
	userAgent                string
	allowedEnvironments      []string
	allowedEndpoints         []string
	allowedJTWSigningMethods []string
)

// init defines constants that will be used later
func init() {
	userAgent = fmt.Sprintf("corbel-go/%s", Version)
	allowedEnvironments = []string{"production", "staging", "current", "next", "qa", "int", "demo"}
	allowedEndpoints = []string{"iam", "oauth", "assets", "resources"}
	allowedJTWSigningMethods = []string{"HS256", "RSA"}
}

// Client is the struct that manages communication with the Corbel APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	httpClient *http.Client

	// Environment is used to define the target environment to speak with.
	Environment string

	// ClientName is the name that match the clientID
	// (Optional) The required information is the clientID
	ClientName string

	// ClientID is the application defined client on Corbel.
	ClientID string

	// ClientSecret is the application secret hash that match with clientID.
	ClientSecret string

	// ClientScopes are those scopes the client will ask for to the platform when building the client connection
	// ClientScopes is a string with the scopes delimited by spaces.
	ClientScopes string

	// ClientDomain is the SR domain where to make the operations using the provided credentials.
	// (Optional) Every clientID only maps to one SR domain.
	ClientDomain string

	// ClientJWTSigningMethod defines the signing method configured for the client.
	// Must match with the one configured on the platform since it will understand only that one.
	// Only allowed signing methods at the moment are: HS256 and RSA
	ClientJWTSigningMethod string

	// TokenExpirationTime define the amount of time in seconds that a token must be valid.
	// It must be lower than 3600 seconds, since is the imposed requisite from the platform.
	TokenExpirationTime uint64

	// UserAgent defines the UserAgent to send in the Headers for every request to the platform.
	UserAgent string

	// Token is the actual token to send as Authentication Bearer
	CurrentToken string

	// CurrentTokenExpiresAt is the unix time where the token will expire
	CurrentTokenExpiresAt int64

	// CurrentRefreshToken is the current refresh token received from the IAM service
	CurrentRefreshToken string

	// IAM endpoint struct
	IAM       *IAMService
	Resources *ResourcesService
	Assets    *AssetsService
}

// URLFor returns the formated url of the API using the actual url scheme
func (c *Client) URLFor(endpoint, uri string) string {
	if c.Environment == "production" {
		return fmt.Sprintf("https://%s.bqws.io%s", endpoint, uri)
	}
	return fmt.Sprintf("https://%s-%s.bqws.io%s", endpoint, c.Environment, uri)
}

// Token returns the token to use as bearer. If the token has already expired
// it refresh it.
// TODO: Refresh token
func (c *Client) Token() string {
	if c.CurrentTokenExpiresAt <= time.Now().Unix()*1000 {
		return ""
	}
	return c.CurrentToken
}

// NewClient returns a new Corbel API client.
// If a nil httpClient is provided, it will return a http.DefaultClient.
func NewClient(httpClient *http.Client, clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod string, tokenExpirationTime uint64) (*Client, error) {
	return NewClientForEnvironment(httpClient, "production", clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod, tokenExpirationTime)
}

// NewClientForEnvironment returns a new Corbel API client.
// If a nil httpClient is provided, it will return a http.DefaultClient.
// If a empty environment is provided, it will use production as environment.
func NewClientForEnvironment(httpClient *http.Client, environment, clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod string, tokenExpirationTime uint64) (*Client, error) {

	var thisClient *Client

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if environment == "" {
		environment = "production"
	}

	// allowedEnvironments?
	if stringInSlice(allowedEnvironments, environment) == false {
		return nil, errInvalidEnvironment
	}

	// allowedJTWSigningMethods?
	if stringInSlice(allowedJTWSigningMethods, clientJWTSigningMethod) == false {
		return nil, errInvalidJWTSigningMethod
	}

	// incorrect Token Expiration Time?
	if tokenExpirationTime > 3600 || tokenExpirationTime == 0 {
		return nil, errInvalidTokenExpirationTime
	}

	// required parameters?
	if clientID == "" || clientSecret == "" {
		return nil, errMissingClientParams
	}

	thisClient = &Client{
		httpClient:             httpClient,
		Environment:            environment,
		ClientName:             clientName,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		ClientDomain:           clientDomain,
		ClientScopes:           clientScopes,
		ClientJWTSigningMethod: clientJWTSigningMethod,
		TokenExpirationTime:    tokenExpirationTime * 1000,
		UserAgent:              userAgent,
	}

	thisClient.IAM = &IAMService{client: thisClient}
	thisClient.Resources = &ResourcesService{client: thisClient}
	thisClient.Assets = &AssetsService{client: thisClient}

	return thisClient, nil
}
