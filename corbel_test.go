package corbel

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClientNewClient(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClient(nil, nil, "", "", "", "", "", "", 3000)
	if err == nil {
		t.Error("NewClient must fail if JWT clientJWTSigningMethod is not an allowed method.")
	}

	client, err = NewClient(nil, nil, "", "", "", "", "", "HS256", 0)
	if err == nil {
		t.Error("NewClient must fail if Token expiration time is 0")
	}

	client, err = NewClient(nil, nil, "", "", "", "", "", "HS256", 3601)
	if err == nil {
		t.Error("NewClient must fail if Token expiration time is over 3600")
	}

	client, err = NewClient(nil, nil, "", "", "", "", "", "HS256", 3000)
	if err == nil {
		t.Error("NewClient must fail if client Id or Secret are not passed, but it did't raise an error")
	}

	endpoints := map[string]string{"iam": "https://iam.bqws.io", "resources": "https://resources.bqws.io"}
	client, err = NewClient(nil, nil, "someID", "", "someSecret", "", "", "HS256", 3000)
	if err != nil {
		t.Error("NewClient must not fail if client Id or Secret are provided, but it raised an error")
	}

	if got, want := client.Endpoints["iam"], endpoints["iam"]; got != want {
		t.Errorf("NewClient Environment is %v, but want %v", got, want)
	}

	if got, want := client.httpClient, http.DefaultClient; got != want {
		t.Errorf("NewClient HTTPClient is %v, but want %v", got, want)
	}

	if got, want := client.UserAgent, fmt.Sprintf("corbel-go/%s", Version); got != want {
		t.Errorf("NewClient HTTPClient is %v, but want %v", got, want)
	}

	if got, want := client.CurrentToken, ""; got != want {
		t.Errorf("NewClient Token is %v, but want %v", got, want)
	}

	if got, want := client.CurrentToken, client.Token(); got != want {
		t.Errorf("NewClient Token is %v, but want %v", got, want)
	}
}

func TestClientURLFor(t *testing.T) {
	var (
		client *Client
	)

	client, _ = NewClient(nil, nil, "someID", "", "someSecret", "", "", "HS256", 3000)
	if got, want := client.URLFor("iam", "/v1.0/auth/token"), "https://iam.bqws.io/v1.0/auth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}

	endpoints := map[string]string{"iam": "https://iam-qa.bqws.io", "resources": "https://resources-qa.bqws.io"}
	client, _ = NewClient(nil, endpoints, "someID", "", "someSecret", "", "", "HS256", 3000)
	if got, want := client.URLFor("iam", "/v1.0/auth/token"), "https://iam-qa.bqws.io/v1.0/auth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}
}
