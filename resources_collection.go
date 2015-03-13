package corbel

import (
	"fmt"
	"net/http"
)

// AddToCollection add the required struct formated as json to the desired collection
// resource must have exported variables and optionally its representation as JSON.
func (r *ResourcesService) AddToCollection(collectionName string, resource interface{}) error {
	var (
		req *http.Request
		err error
	)

	req, err = r.client.NewRequest("POST", "resources", fmt.Sprintf("/v1.0/resource/%s", collectionName), resource)
	return returnErrorHTTPSimple(r.client, req, err, 201)
}

// UpdateInCollection updates the required struct formated as json to the desired collection
// resource must have exported variables and optionally its representation as JSON.
func (r *ResourcesService) UpdateInCollection(collectionName, id string, resource interface{}) error {
	var (
		req *http.Request
		err error
	)

	req, err = r.client.NewRequest("PUT", "resources", fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id), resource)
	return returnErrorHTTPSimple(r.client, req, err, 204)
}

// SearchCollection gets the desired objects in base of a search query
func (r *ResourcesService) SearchCollection(collectionName string) *Search {
	return NewSearch(r.client, "resources", fmt.Sprintf("/v1.0/resource/%s", collectionName))
}

// GetFromCollection gets the desired object from the collection by id
func (r *ResourcesService) GetFromCollection(collectionName, id string, resource interface{}) error {
	var (
		req *http.Request
		err error
	)

	req, err = r.client.NewRequest("GET", "resources", fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id), nil)
	return returnErrorHTTPInterface(r.client, req, err, resource, 200)
}

// DeleteFromCollection deletes the desired resource from the platform by id
func (r *ResourcesService) DeleteFromCollection(collectionName, id string) error {

	var (
		req *http.Request
		err error
	)

	req, err = r.client.NewRequest("DELETE", "resources", fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id), nil)
	return returnErrorHTTPSimple(r.client, req, err, 204)
}
