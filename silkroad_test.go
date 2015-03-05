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
