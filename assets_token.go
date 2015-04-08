package corbel

import (
	"fmt"
	"net/http"
	"net/url"
)

// UpgradeToken returns the assertion to upgrade the token on IAM to get the
// 'purchased' scopes
func (a *AssetsService) UpgradeToken() error {
	var (
		req      *http.Request
		res      *http.Response
		location *url.URL
		err      error
	)

	req, err = a.client.NewRequest("GET", "assets", "/v1.0/asset/access", nil)
	if err != nil {
		return err
	}

	req.Header.Add("No-Redirect", "true")
	res, err = a.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	location, err = res.Location()
	if err != nil {
		return err
	}

	req, err = a.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/oauth/token/upgrade?%s", location.RawQuery), nil)
	if err != nil {
		return err
	}

	return returnErrorHTTPSimple(a.client, req, err, 204)
}
