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

	client, err = NewClientForEnvironment(nil, "", "", "", "", "", "", "", 3000)
	if err == nil {
		t.Error("NewClient must fail if JWT clientJWTSigningMethod is not an allowed method.")
	}

	client, err = NewClientForEnvironment(nil, "", "", "", "", "", "", "HS256", 0)
	if err == nil {
		t.Error("NewClient must fail if Token expiration time is 0")
	}

	client, err = NewClientForEnvironment(nil, "", "", "", "", "", "", "HS256", 3601)
	if err == nil {
		t.Error("NewClient must fail if Token expiration time is over 3600")
	}

	client, err = NewClientForEnvironment(nil, "", "", "", "", "", "", "HS256", 3000)
	if err == nil {
		t.Error("NewClient must fail if client Id or Secret are not passed, but it did't raise an error")
	}

	client, err = NewClientForEnvironment(nil, "wrongEnvironment", "someID", "", "someSecret", "", "", "HS256", 3000)
	if err == nil {
		t.Error("NewClient must fail if a wrong environment name is provided, but it did not raised an error")
	}

	client, err = NewClientForEnvironment(nil, "", "someID", "", "someSecret", "", "", "HS256", 3000)
	if err != nil {
		t.Error("NewClient must not fail if client Id or Secret are provided, but it raised an error")
	}

	if got, want := client.Environment, "production"; got != want {
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
}

func TestClientURLFor(t *testing.T) {
	var (
		client *Client
	)

	client, _ = NewClient(nil, "someID", "", "someSecret", "", "", "HS256", 3000)

	if got, want := client.URLFor("iam", "/v1.0/auth/token"), "https://iam.bqws.io/v1.0/auth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}

	client, _ = NewClientForEnvironment(nil, "qa", "someID", "", "someSecret", "", "", "HS256", 3000)
	if got, want := client.URLFor("iam", "/v1.0/auth/token"), "https://iam-qa.bqws.io/v1.0/auth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}
}
