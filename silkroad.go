package silkroad

import (
	"fmt"
	"net/http"
)

// Client is the struct that manages communication with the Silkroad APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	client *http.Client

	// Environment is used to define the target environment to speak with.
	Environment string
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
func NewClient(httpClient *http.Client, environment string) (client *Client) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if environment == "" {
		environment = "production"
	}

	client = &Client{
		client:      httpClient,
		Environment: environment,
	}

	return
}
