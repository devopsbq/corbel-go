package silkroad

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	userAgent = fmt.Sprintf("go-silkroad/%s", Version)
	allowedEnvironments = []string{"production", "staging", "current", "next", "qa", "integration"}
	allowedEndpoints = []string{"iam", "oauth", "assets", "resources"}
	allowedJTWSigningMethods = []string{"HS256", "RSA"}
}

// Client is the struct that manages communication with the Silkroad APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	httpClient *http.Client

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
	TokenExpirationTime time.Duration

	// UserAgent defines the UserAgent to send in the Headers for every request to the platform.
	UserAgent string

	// Token is the actual token to send as Authentication Bearer
	CurrentToken string

	// CurrentTokenExpirationTime is the unix time where the token will expire
	CurrentTokenExpirationTime int64

	// CurrentRefreshToken is the current refresh token received from the IAM service
	CurrentRefreshToken string

	// IAM endpoint struct
	IAM       *IAMService
	Resources *ResourcesService
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

// NewRequest creates an API request.
// method is the HTTP method to use
// endpoint is the endpoint of SR to speak with
// url is the url to query. it must be preceded by a slash.
// body is, if specified, the value JSON encoded to be used as request body.
// mediaType is the desired media type to set in the request.
func (c *Client) NewRequest(method, endpoint, urlStr, mediaType string, body interface{}) (*http.Request, error) {
	u, _ := url.Parse(c.URLFor(endpoint, urlStr))

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	if c.CurrentToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.CurrentToken))
	}
	return req, nil
}

// Token returns the token to use as bearer. If the token has already expired
// it refresh it.
func (c *Client) Token() string {
	if c.CurrentTokenExpirationTime < time.Now().Unix() {
		return c.CurrentToken
	}
	return ""
}

// ReturnErrorByHTTPStatusCode returns the http error code or nil if it returns the
//   desired error
func ReturnErrorByHTTPStatusCode(res *http.Response, desiredStatusCode int) error {
	if res.StatusCode == desiredStatusCode {
		return nil
	}
	return errors.New(http.StatusText(res.StatusCode))
}

// NewClient returns a new Silkroad API client.
// If a nil httpClient is provided, it will return a http.DefaultClient.
// If a empty environment is provided, it will use production as environment.
func NewClient(httpClient *http.Client, environment, clientID, clientName, clientSecret, clientScopes, clientDomain, clientJWTSigningMethod string, tokenExpirationTime time.Duration) (*Client, error) {

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
		TokenExpirationTime:    tokenExpirationTime,
		UserAgent:              userAgent,
	}

	thisClient.IAM = &IAMService{client: thisClient}
	thisClient.Resources = &ResourcesService{client: thisClient}

	return thisClient, nil
}
