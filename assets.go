package corbel

// AssetsService handles communication with the Assets service of Corbel.
// It takes care of user assets that gives special scopes to the resources.
//
// Full API info: http://docs.silkroadassets.apiary.io/
type AssetsService struct {
	client *Client
}
