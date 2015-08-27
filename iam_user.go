package corbel

import (
	"fmt"
	"net/http"
)

// IAMUser is the representation of an User object used by IAM
type IAMUser struct {
	ID          string                 `json:"id,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	Username    string                 `json:"username,omitempty"`
	Email       string                 `json:"email,omitempty"`
	Password    string                 `json:"password,omitempty"`
	FirstName   string                 `json:"firstName,omitempty"`
	LastName    string                 `json:"lastName,omitempty"`
	ProfileURL  string                 `json:"profileUrl,omitempty"`
	PhoneNumber string                 `json:"phoneNumber,omitempty"`
	Scopes      []string               `json:"scopes"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Country     string                 `json:"country,omitempty"`
	CreatedDate int                    `json:"createdDate,omitempty"`
	CreatedBy   string                 `json:"createdBy,omitempty"`
	Groups      []string               `json:"groups,omitempty"`
}

// UserAdd adds an IAMUser defined struct to the domain of the client
func (i *IAMService) UserAdd(user *IAMUser) (string, error) {
	req, err := i.client.NewRequest("POST", "iam", "/v1.0/user", user)
	return returnErrorHTTPSimple(i.client, req, err, 201)
}

// UserExists checks if an user exists in the domain of the client
func (i *IAMService) UserExists(username string) bool {
	var (
		req *http.Request
		err error
	)

	req, err = i.client.NewRequest("HEAD", "iam", fmt.Sprintf("/v1.0/username/%s", username), nil)
	_, err = returnErrorHTTPSimple(i.client, req, err, 200)
	if err != nil {
		return false
	}
	return true
}

// UserUpdate updates an user by using IAMUser
func (i *IAMService) UserUpdate(id string, user *IAMUser) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("PUT", "iam", fmt.Sprintf("/v1.0/user/%s", id), user)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

//UserUpdateMe updates the user authenticated by the current token
func (i *IAMService) UserUpdateMe(user *IAMUser) error {
	return i.UserUpdate("me", user)
}

// UserGet gets the desired IAMUuser from the domain by id
func (i *IAMService) UserGet(id string, user *IAMUser) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/user/%s", id), nil)
	_, err = returnErrorHTTPInterface(i.client, req, err, user, 200)
	return err
}

// UserGetMe gets the user authenticated by the current token
func (i *IAMService) UserGetMe(user *IAMUser) error {
	return i.UserGet("me", user)
}

// UserDelete deletes the desired user from IAM by id
func (i *IAMService) UserDelete(id string) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/user/%s", id), nil)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

//UserDeleteMe deletes the user authenticated by the current token
func (i *IAMService) UserDeleteMe() error {
	return i.UserDelete("me")
}

// UserSearch gets the desired objects in base of a search query
func (i *IAMService) UserSearch() *Search {
	return NewSearch(i.client, "iam", "/v1.0/user")
}

// UserByUsername allow to find an user based on its username
func (i *IAMService) UserByUsername(username string) (*IAMUser, error) {
	search := i.UserSearch()
	search.Query.Eq["username"] = username

	var arrUsers []*IAMUser
	err := search.Page(0, &arrUsers)
	if err != nil {
		return nil, err
	}
	if len(arrUsers) == 0 {
		return nil, errUserNotFound
	}
	return arrUsers[0], nil
}
