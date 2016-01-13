package corbel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/Sirupsen/logrus"
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
	url, _ := url.Parse(c.URLFor(endpoint, urlStr))
	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			c.logger.Debugf("failed to encode body: %v", err)
			return nil, errJSONMarshalError
		}
	}

	c.logger.WithFields(logrus.Fields{
		"method": method, "accept": headerAccept, "url": url.String(), "body": buf.String(),
	}).Debug("new request")
	req, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		c.logger.Debugf("failed to create request: %v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", headerContentType)
	req.Header.Add("Accept", headerAccept)
	req.Header.Add("User-Agent", c.UserAgent)
	if token := c.Token(); token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	c.logger.Debugf("request headers: %v", req.Header)
	return req, nil
}

// NewRequest creates an API request using 'application/json', most common api query.
// method is the HTTP method to use
// endpoint is the endpoint of SR to speak with
// url is the url to query. it must be preceded by a slash.
// body is, if specified, the value JSON encoded to be used as request body.
func (c *Client) NewRequest(method, endpoint, url string, body interface{}) (*http.Request, error) {
	return c.NewRequestContentType(method, endpoint, url, "application/json", "application/json", body)
}

func returnErrorHTTPInterface(client *Client, req *http.Request, errr error, object interface{}, desiredStatusCode int) (string, error) {
	if errr != nil {
		return "", errr
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		client.logger.Debugf("failed to make request: %v", err)
		return "", err
	}
	objectByte, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if object != nil {
		if err != nil {
			return "", errResponseError
		}
		if err = json.Unmarshal(objectByte, &object); err != nil {
			return "", errJSONUnmarshalError
		}
	}
	client.logger.WithFields(logrus.Fields{
		"code": res.StatusCode, "status": res.Status, "body": string(objectByte),
	}).Debug("response received")
	return returnErrorByHTTPStatusCode(res, desiredStatusCode)
}

func returnErrorHTTPSimple(client *Client, req *http.Request, err error, desiredStatusCode int) (string, error) {
	return returnErrorHTTPInterface(client, req, err, nil, desiredStatusCode)
}

// returnErrorByHTTPStatusCode returns the http error code or nil if it returns the
// desired error
func returnErrorByHTTPStatusCode(res *http.Response, desiredStatusCode int) (string, error) {
	location, _ := res.Location()
	locationString := ""
	if location != nil {
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
