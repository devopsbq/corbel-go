package corbel

// IAMService handles communication with the IAM service of Corbel.
// It takes care of all Identity and Authorization Management
//
// Full API info: http://docs.silkroadiam.apiary.io/
type IAMService struct {
	client *Client
}
