package corbel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// NewRequest creates an API request.
// method is the HTTP method to use
// endpoint is the endpoint of SR to speak with
// url is the url to query. it must be preceded by a slash.
// body is, if specified, the value JSON encoded to be used as request body.
func (c *Client) NewRequest(method, endpoint, urlStr string, body interface{}) (*http.Request, error) {
	u, _ := url.Parse(c.URLFor(endpoint, urlStr))

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)
	if c.CurrentToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.CurrentToken))
	}
	return req, nil
}

func returnErrorHTTPInterface(client *Client, req *http.Request, err error, object interface{}, desiredStatusCode int) error {
	var (
		res        *http.Response
		objectByte []byte
	)
	if err != nil {
		return err
	}

	res, err = client.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	objectByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errResponseError
	}

	err = json.Unmarshal(objectByte, &object)
	if err != nil {
		return errJSONUnmarshalError
	}

	return returnErrorByHTTPStatusCode(res, desiredStatusCode)
}

func returnErrorHTTPSimple(client *Client, req *http.Request, err error, desiredStatusCode int) error {
	var (
		res *http.Response
	)
	if err != nil {
		return err
	}

	res, err = client.httpClient.Do(req)
	if err != nil {
		return err
	}

	return returnErrorByHTTPStatusCode(res, desiredStatusCode)
}

// returnErrorByHTTPStatusCode returns the http error code or nil if it returns the
// desired error
func returnErrorByHTTPStatusCode(res *http.Response, desiredStatusCode int) error {
	if res.StatusCode == desiredStatusCode {
		return nil
	}
	if http.StatusText(res.StatusCode) == "" {
		return fmt.Errorf("HTTP Error %d", res.StatusCode)
	}
	return errors.New(http.StatusText(res.StatusCode))
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qv, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qv.Encode()
	return u.String(), nil
}
