package silkroad

// IAMEndpoint handles communication with the IAM service of Silkroad.
// It takes care of all Identity and Authorization Management
//
// Full API info: http://docs.silkroadiam.apiary.io/
type IAMService struct {
	client *Client
}
