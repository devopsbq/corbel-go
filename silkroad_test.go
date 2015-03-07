package silkroad

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClient(nil, "", "", "", "", "", "", "", 3000)
	if err == nil {
		t.Error("NewClient must fail if JWT clientJWTSigningMethod is not an allowed method.")
	}

	client, err = NewClient(nil, "", "", "", "", "", "", "HS256", 0)
	if err == nil {
		t.Error("NewClient must fail if Token expiration time is 0")
	}

	client, err = NewClient(nil, "", "", "", "", "", "", "HS256", 3601)
	if err == nil {
		t.Error("NewClient must fail if Token expiration time is over 3600")
	}

	client, err = NewClient(nil, "", "", "", "", "", "", "HS256", 3000)
	if err == nil {
		t.Error("NewClient must fail if client Id or Secret are not passed, but it did't raise an error")
	}

	client, err = NewClient(nil, "wrongEnvironment", "someID", "", "someSecret", "", "", "HS256", 3000)
	if err == nil {
		t.Error("NewClient must fail if a wrong environment name is provided, but it did not raised an error")
	}

	client, err = NewClient(nil, "", "someID", "", "someSecret", "", "", "HS256", 3000)
	if err != nil {
		t.Error("NewClient must not fail if client Id or Secret are provided, but it raised an error")
	}

	if got, want := client.Environment, "production"; got != want {
		t.Errorf("NewClient Environment is %v, but want %v", got, want)
	}

	if got, want := client.httpClient, http.DefaultClient; got != want {
		t.Errorf("NewClient HTTPClient is %v, but want %v", got, want)
	}

	if got, want := client.UserAgent, fmt.Sprintf("go-silkroad/%s", Version); got != want {
		t.Errorf("NewClient HTTPClient is %v, but want %v", got, want)
	}
}

func TestClientURLFor(t *testing.T) {
	var (
		client *Client
	)

	client, _ = NewClient(nil, "", "someID", "", "someSecret", "", "", "HS256", 3000)

	if got, want := client.URLFor("iam", "/v1.0/auth/token"), "https://iam.bqws.io/v1.0/auth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}

	client, _ = NewClient(nil, "qa", "someID", "", "someSecret", "", "", "HS256", 3000)
	if got, want := client.URLFor("iam", "/v1.0/auth/token"), "https://iam-qa.bqws.io/v1.0/auth/token"; got != want {
		t.Errorf("urlFor url is %v, but want %v", got, want)
	}
}

func TestClientNewRequest(t *testing.T) {
	// IAMVersion is a helper struct for the test
	type structIAMVersion struct {
		BuildUser       string `json:"build.user,omitempty"`
		BuildGroupID    string `json:"build.groupId,omitempty"`
		BuildArtifactID string `json:"build.artifactId,omitempty"`
	}

	var (
		client     *Client
		err        error
		req        *http.Request
		iamVersion *structIAMVersion
		contents   []byte
	)

	client, err = NewClient(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		3000)

	req, err = client.NewRequest("GET", "iam", "/version", "application/json", nil)
	if err != nil {
		t.Errorf("Request failed: %s", err.Error())
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		t.Errorf("Request failed: %s", err.Error())
	}

	defer res.Body.Close()
	contents, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading the request body. %s", err.Error())
	}

	err = json.Unmarshal(contents, &iamVersion)
	if err != nil {
		t.Errorf("Error unmarshalling the response. %s", err.Error())
	}

	if got, want := iamVersion.BuildArtifactID, "iam"; got != want {
		t.Errorf("/version unmarshaled json build.artifactId is %v, but want %v", got, want)
	}
	if got, want := iamVersion.BuildGroupID, "com.bqreaders.silkroad"; got != want {
		t.Errorf("/version unmarshaled json build.groupId is %v, but want %v", got, want)
	}
	if got, want := iamVersion.BuildUser, "jenkins"; got != want {
		t.Errorf("/version unmarshaled json build.user is %v, but want %v", got, want)
	}
}

func TestClientIAMOauthToken(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClient(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	token, err := client.IAM.OauthToken()
	if got := err; got != nil {
		t.Errorf("GetToken must not fail. Got: %v  Want: nil", got)
	}

	if got, want := strings.Count(token, "."), 2; got != want {
		t.Errorf("GetToken must return a token with 2 dots. Got: %v  Want: %v", got, want)
	}
}

func TestClientIAMOauthTokenUpgrade(t *testing.T) {
	var (
		client *Client
		err    error
	)

	client, err = NewClient(
		nil,
		"qa",
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		10)

	err = client.IAM.OauthTokenUpgrade("aaaaaa")
	if err != errClientNotAuthorized {
		t.Errorf("OauthTokenUpgrade must fail since it got an invalid token. %s", err)
	}

	// TODO: correct tests with Assets workflow
}
