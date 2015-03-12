package corbel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IAMUser is the representation of an User object used by IAM
type IAMUser struct {
	ID          string                 `json:"id,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	Username    string                 `json:"username,omitempty"`
	Email       string                 `json:"email,omitempty"`
	FirstName   string                 `json:"firstName,omitempty"`
	LastName    string                 `json:"lastName,omitempty"`
	ProfileURL  string                 `json:"profileUrl,omitempty"`
	PhoneNumber string                 `json:"phoneNumber,omitempty"`
	Scopes      []string               `json:"scopes,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Country     string                 `json:"country,omitempty"`
	CreatedDate int                    `json:"createdDate,omitempty"`
	CreatedBy   string                 `json:"createdBy,omitempty"`
}

// Add adds an IAMUser defined struct to the domain of the client
func (i *IAMService) Add(user *IAMUser) error {
	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = i.client.NewRequest("POST", "iam", "/v1.0/user/", user)
	if err != nil {
		return err
	}

	res, err = i.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	return ReturnErrorByHTTPStatusCode(res, 201)
}

// Update updates an user by using IAMUser
func (i *IAMService) Update(id string, user *IAMUser) error {
	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = i.client.NewRequest("PUT", "iam", fmt.Sprintf("/v1.0/user/%s", id), user)
	if err != nil {
		return err
	}

	res, err = i.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	return ReturnErrorByHTTPStatusCode(res, 204)
}

// Get gets the desired IAMUuser from the domain by id
func (i *IAMService) Get(id string, user *IAMUser) error {

	var (
		req      *http.Request
		res      *http.Response
		userByte []byte
		err      error
	)

	req, err = i.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/user/%s", id), nil)
	if err != nil {
		return err
	}

	res, err = i.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	userByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errResponseError
	}

	err = json.Unmarshal(userByte, &user)
	if err != nil {
		return errJSONUnmarshalError
	}

	return ReturnErrorByHTTPStatusCode(res, 200)
}

// Delete deletes the desired user from IAM by id
func (i *IAMService) Delete(id string) error {

	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/user/%s", id), nil)
	if err != nil {
		return err
	}

	res, err = i.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	return ReturnErrorByHTTPStatusCode(res, 204)
}

// Search gets the desired objects in base of a search query
func (i *IAMService) Search() *Search {
	return NewSearch(i.client, "iam", "/v1.0/user")
}
