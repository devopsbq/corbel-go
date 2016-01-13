package corbel

import (
	"fmt"
	"net/http"
)

// ResourcesService handles the interface for retrival resource's representation
// on Corbel.
//
// Full API info: http://docs.corbelresources.apiary.io/
type ResourcesService struct {
	client *Client
}

// UserACL defines the content of an ACL for a user
type UserACL struct {
	Permission string                 `json:"permission"`
	Properties map[string]interface{} `json:"properties"`
}

func (r *ResourcesService) createRequest(method, accept, uri string, body interface{}) (*http.Request, error) {
	return r.client.NewRequestContentType(method, "resources", uri, "application/json", accept, body)
}

// CollectionRequest perform a specific collection request on resources
func (r *ResourcesService) CollectionRequest(method, accept, collectionName string, send interface{}) (*http.Request, error) {
	uri := fmt.Sprintf("/v1.0/resource/%s", collectionName)
	return r.createRequest(method, accept, uri, send)
}

// ResourceRequest perform a specific resource request on resources
func (r *ResourcesService) ResourceRequest(method, accept, collectionName, id string, send interface{}) (*http.Request, error) {
	uri := fmt.Sprintf("/v1.0/resource/%s/%s", collectionName, id)
	return r.createRequest(method, accept, uri, send)
}

// RelationRequest perform a specific relation request on resources
func (r *ResourcesService) RelationRequest(method, accept, collectionName, resourceID, relationName, relatedCollectionName, relatedID string, send interface{}) (*http.Request, error) {
	uri := fmt.Sprintf("/v1.0/resource/%s/%s/%s", collectionName, resourceID, relationName)
	if relatedCollectionName != "" || relatedID != "" {
		uri = fmt.Sprintf("%s;r=%s", uri, relatedCollectionName)
	}
	if relatedID != "" {
		uri = fmt.Sprintf("%s/%s", uri, relatedID)
	}
	return r.createRequest(method, accept, uri, send)
}
