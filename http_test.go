package corbel

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestClientNewRequest(t *testing.T) {
	// IAMVersion is a helper struct for the test
	type structIAMVersion struct {
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

	endpoints := map[string]string{"iam": "https://iam-int.bqws.io", "resources": "https://resources-int.bqws.io"}
	client, err = NewClient(
		nil,
		endpoints,
		"a9fb0e79",
		"test-client",
		"90f6ed907ce7e2426e51aa52a18470195f4eb04725beb41569db3f796a018dbd",
		"",
		"silkroad-qa",
		"HS256",
		3000)

	req, err = client.NewRequest("GET", "iam", "/version", nil)
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
	if got, want := iamVersion.BuildGroupID, "io.corbel"; got != want {
		t.Errorf("/version unmarshaled json build.groupId is %v, but want %v", got, want)
	}
}
