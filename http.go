package corbel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// NewRequestContentType creates an API request.
// method is the HTTP method to use
// endpoint is the endpoint of SR to speak with
// urlStr is the url to query. it must be preceded by a slash.
// headerContentType is the header['Content-Type'] of the request.
// headerAccept is the header['Accept'] of the request.
// body is, if specified, the value JSON encoded to be used as request body.
func (c *Client) NewRequestContentType(method, endpoint, urlStr, headerContentType, headerAccept string, body interface{}) (*http.Request, error) {
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

	req.Header.Add("Content-Type", headerContentType)
	req.Header.Add("Accept", headerAccept)
	req.Header.Add("User-Agent", c.UserAgent)
	if token := c.Token(); token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	c.logger.Debugf("Request values -> Headers: %v, Host: %s, Method: %s, Form: %v, URL: %v, Body: %v",
		req.Header, req.Host, req.Method, req.Form, req.URL, req.Body)
	return req, nil
}

// NewRequest creates an API request using 'application/json', most common api query.
// method is the HTTP method to use
// endpoint is the endpoint of SR to speak with
// url is the url to query. it must be preceded by a slash.
// body is, if specified, the value JSON encoded to be used as request body.
func (c *Client) NewRequest(method, endpoint, urlStr string, body interface{}) (*http.Request, error) {
	return c.NewRequestContentType(method, endpoint, urlStr, "application/json", "application/json", body)
}

func returnErrorHTTPInterface(client *Client, req *http.Request, err error, object interface{}, desiredStatusCode int) (string, error) {
	var (
		res        *http.Response
		objectByte []byte
	)
	if err != nil {
		return "", err
	}

	res, err = client.httpClient.Do(req)
	client.logger.Debugf("Response values -> Header: %v, Code: %d, Status: %s, Body: %v",
		res.Header, res.StatusCode, res.Status, res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	objectByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errResponseError
	}

	err = json.Unmarshal(objectByte, &object)
	if err != nil {
		return "", errJSONUnmarshalError
	}

	// fmt.Println(string(objectByte)) // for debug

	return returnErrorByHTTPStatusCode(res, desiredStatusCode)
}

func returnErrorHTTPSimple(client *Client, req *http.Request, err error, desiredStatusCode int) (string, error) {
	var (
		res *http.Response
	)
	if err != nil {
		return "", err
	}

	res, err = client.httpClient.Do(req)
	client.logger.Debugf("Response values -> Header: %v, Code: %d, Status: %s, Body: %v",
		res.Header, res.StatusCode, res.Status, res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	// objectByte, _ := ioutil.ReadAll(res.Body) // for debug
	// fmt.Println(string(objectByte)) // for debug

	return returnErrorByHTTPStatusCode(res, desiredStatusCode)
}

// returnErrorByHTTPStatusCode returns the http error code or nil if it returns the
// desired error
func returnErrorByHTTPStatusCode(res *http.Response, desiredStatusCode int) (string, error) {
	var (
		location       *url.URL
		locationString string
	)
	location, _ = res.Location()
	if location == nil {
		locationString = ""
	} else {
		locationString = location.String()
	}

	if res.StatusCode == desiredStatusCode {
		return locationString, nil
	}
	if http.StatusText(res.StatusCode) == "" {
		return "", fmt.Errorf("HTTP Error %d", res.StatusCode)
	}
	return locationString, fmt.Errorf("%d %s", res.StatusCode, http.StatusText(res.StatusCode))
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
