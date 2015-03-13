package corbel

import (
	"fmt"
	"net/http"
)

// AddRelation adds the required relation to the resource in the collection
// with the _related_ resource. Additionally arbitrary information can be passed
// to as relation data or nil.
func (r *ResourcesService) AddRelation(collectionName, resourceID, relationName, relatedCollectionName, relatedID string, relationInfo interface{}) error {
	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = r.client.NewRequest("PUT", "resources", fmt.Sprintf("/v1.0/resource/%s/%s/%s;r=%s/%s", collectionName, resourceID, relationName, relatedCollectionName, relatedID), relationInfo)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	return ReturnErrorByHTTPStatusCode(res, 201)
}

// DeleteRelation deletes the desired relation between the origin and the related
// resource
func (r *ResourcesService) DeleteRelation(collectionName, resourceID, relationName, relatedCollectionName, relatedID string) error {

	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = r.client.NewRequest("DELETE", "resources", fmt.Sprintf("/v1.0/resource/%s/%s/%s;r=%s/%s", collectionName, resourceID, relationName, relatedCollectionName, relatedID), nil)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	return ReturnErrorByHTTPStatusCode(res, 204)
}

// DeleteAllRelations deletes all the relations by relationName of the desired resource
func (r *ResourcesService) DeleteAllRelations(collectionName, resourceID, relationName string) error {

	var (
		req *http.Request
		res *http.Response
		err error
	)

	req, err = r.client.NewRequest("DELETE", "resources", fmt.Sprintf("/v1.0/resource/%s/%s/%s", collectionName, resourceID, relationName), nil)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}

	return ReturnErrorByHTTPStatusCode(res, 204)
}
