package silkroad

import (
	"fmt"
	"net/http"
)

var (
	userAgent                string
	allowedEnvironments      []string
	allowedEndpoints         []string
	allowedJTWSigningMethods []string
)

// init defines constants that will be used later
func init() {
	userAgent = fmt.Sprintf("go-silkroad/%s", Version)
	allowedEnvironments = []string{"production", "staging", "current", "next", "qa", "integration"}
	allowedEndpoints = []string{"iam", "oauth", "assets", "resources"}
	allowedJTWSigningMethods = []string{"HS256", "RSA"}
}

// Client is the struct that manages communication with the Silkroad APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	client *http.Client

	// Environment is used to define the target environment to speak with.
	Environment string

	// ClientName is the name that match the clientID
	// (Optional) The required information is the clientID
	ClientName string

	// ClientID is the application defined client on Silkroad.
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
	TokenExpirationTime uint16

	// UserAgent defines the UserAgent to send in the Headers for every request to the platform.
	UserAgent string

	// IAM endpoint struct
	IAM *IAMEndpoint
}

// URLFor returns the formated url of the API using the actual url scheme
func (c *Client) URLFor(endpoint, uri string) (url string) {
	switch c.Environment {
	case "production":
		url = fmt.Sprintf("https://%s.bqws.io%s", endpoint, uri)
	default:
		url = fmt.Sprintf("https://%s-%s.bqws.io%s", endpoint, c.Environment, uri)
	}
	return
}

// NewClient returns a new Silkroad API client.
// If a nil httpClient is provided, it will return a http.DefaultClient.
// If a empty environment is provided, it will use production as environment.
func NewClient(httpClient *http.Client, environment, clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod string, tokenExpirationTime uint16) (*Client, error) {

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

	//

	thisClient = &Client{
		client:                 httpClient,
		Environment:            environment,
		ClientName:             clientName,
		ClientID:               clientID,
		ClientSecret:           clientSecret,
		ClientDomain:           clientDomain,
		ClientScopes:           clientScopes,
		ClientJWTSigningMethod: clientJWTSigningMethod,
		TokenExpirationTime:    tokenExpirationTime,
		UserAgent:              userAgent,
	}

	thisClient.IAM = &IAMEndpoint{client: thisClient}

	return thisClient, nil
}
