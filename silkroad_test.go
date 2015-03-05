package silkroad

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	var client *Client

	client = NewClient(nil, "")

	if got, want := client.Environment, "production"; got != want {
		t.Errorf("NewClient Environment is %v, but want %v", got, want)
	}

	if got, want := client.client, http.DefaultClient; got != want {
		t.Errorf("NewClient HTTPClient is %v, but want %v", got, want)
	}
}

func TestClienturlFor(t *testing.T) {
	var client *Client

	client = NewClient(nil, "")

	if got, want := client.URLFor("iam", "/v1.0/oauth/token"), "https://iam.bqws.io/v1.0/oauth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}

	client = NewClient(nil, "qa")
	if got, want := client.URLFor("iam", "/v1.0/oauth/token"), "https://iam-qa.bqws.io/v1.0/oauth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}
}
