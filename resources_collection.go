package corbel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// AddToCollection add the required struct formated as json to the desired collection
// resource must have exported variables and optionally its representation as JSON.
func (r *ResourcesService) AddToCollection(collectionName string, resource interface{}) error {
	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = r.client.NewRequest("POST", "resources", fmt.Sprintf("/v1.0/resource/%s", collectionName), resource)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	return ReturnErrorByHTTPStatusCode(res, 201)
}

// UpdateInCollection updates the required struct formated as json to the desired collection
// resource must have exported variables and optionally its representation as JSON.
func (r *ResourcesService) UpdateInCollection(collectionName, id string, resource interface{}) error {
	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = r.client.NewRequest("PUT", "resources", fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id), resource)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	return ReturnErrorByHTTPStatusCode(res, 204)
}

// SearchCollection gets the desired objects in base of a search query
func (r *ResourcesService) SearchCollection(collectionName string) *Search {
	return NewSearch(r.client, "resources", fmt.Sprintf("/v1.0/resource/%s", collectionName))
}

// GetFromCollection gets the desired object from the collection by id
func (r *ResourcesService) GetFromCollection(collectionName, id string, resource interface{}) error {

	var (
		req          *http.Request
		res          *http.Response
		resourceByte []byte
		err          error
	)

	req, err = r.client.NewRequest("GET", "resources", fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id), nil)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	resourceByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errResponseError
	}

	err = json.Unmarshal(resourceByte, &resource)
	if err != nil {
		return errJSONUnmarshalError
	}

	return ReturnErrorByHTTPStatusCode(res, 200)
}

// DeleteFromCollection deletes the desired resource from the platform by id
func (r *ResourcesService) DeleteFromCollection(collectionName, id string) error {

	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = r.client.NewRequest("DELETE", "resources", fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id), nil)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	return ReturnErrorByHTTPStatusCode(res, 204)
}
