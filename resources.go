package corbel

// ResourcesService handles the interface for retrival resource's representation
// on Corbel.
//
// Full API info: http://docs.silkroadresources.apiary.io/
type ResourcesService struct {
	client *Client
}

// Resource te
type Resource struct {
	ACL map[string]string
}
