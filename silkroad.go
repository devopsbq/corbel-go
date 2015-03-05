package silkroad

import (
	"fmt"
	"net/http"
)

var (
	userAgent           string
	allowedEnvironments []string
	allowedEndpoints    []string
)

// init defines constants that will be used later
func init() {
	userAgent = fmt.Sprintf("go-silkroad/%s", Version)
	allowedEnvironments = []string{"production", "staging", "current", "next", "qa", "integration"}
	allowedEndpoints = []string{"iam", "oauth", "assets", "resources"}
}

// Client is the struct that manages communication with the Silkroad APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	client *http.Client

	// Environment is used to define the target environment to speak with.
	Environment string

	// ClientID is the application defined client on Silkroad
	ClientID string

	// ClientSecret is the application secret hash that match with clientID
	ClientSecret string
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
func NewClient(httpClient *http.Client, environment, clientID, clientSecret string) (*Client, error) {

	var client *Client

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if environment == "" {
		environment = "production"
	}

	if StringInSlice(allowedEnvironments, environment) == false {
		return nil, errInvalidEnvironment
	}

	if clientID == "" || clientSecret == "" {
		return nil, errMissingClientParams
	}

	client = &Client{
		client:       httpClient,
		Environment:  environment,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	return client, nil
}
