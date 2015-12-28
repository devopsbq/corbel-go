package corbel

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

// Resource te
type Resource struct {
	ACL map[string]string
}
