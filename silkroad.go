package silkroad

import "net/http"

// Client is the struct that manages communication with the Silkroad APIs.
type Client struct {
	// client is the HTTP client to communicate with the API.
	client *http.Client

	// Environment is used to define the target environment to speak with.
	Environment string
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
