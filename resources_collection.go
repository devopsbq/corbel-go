package silkroad

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AddToCollection add the required struct formated as json to the desired collection
// resource must have exported variables and optionally its representation as JSON.
func (r *ResourcesService) AddToCollection(collectionName string, resource interface{}) error {

	var (
		resourceByte []byte
		req          *http.Request
		res          *http.Response
		err          error
	)

	resourceByte, err = json.Marshal(resource)
	if err != nil {
		return err
	}

	req, err = r.client.NewRequest("POST", "resources", fmt.Sprintf("/v1.0/resource/%s", collectionName), "application/json", resourceByte)
	if err != nil {
		return err
	}

	res, err = r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	return ReturnErrorByHTTPStatusCode(res, 201)
}
