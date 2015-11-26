package corbel

import "fmt"

// IAMGroup is the representation of a Group object used by IAM
type IAMGroup struct {
	ID     string   `json:"id,omitempty"`
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"` // must be omitempty
}

//UserAddGroups add groups to user's list of groups
func (i *IAMService) UserAddGroups(userID string, groupIDs []string) error {
	if userID == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("PUT", "iam", fmt.Sprintf("/v1.0/user/%s/groups", userID), groupIDs)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

//UserDeleteGroup deletes a groups from the user group list
func (i *IAMService) UserDeleteGroup(userID, groupID string) error {
	if userID == "" || groupID == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/user/%s/groups/%s", userID, groupID), nil)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

// GroupAdd adds a new group to iam into the current domain
func (i *IAMService) GroupAdd(group *IAMGroup) (string, error) {
	req, err := i.client.NewRequest("POST", "iam", "/v1.0/group", group)
	return returnErrorHTTPSimple(i.client, req, err, 201)
}

// GroupGetAll gets all Groups of the client current domain
func (i *IAMService) GroupGetAll(groups []*IAMGroup) error {
	req, err := i.client.NewRequest("GET", "iam", "/v1.0/group", nil)
	_, err = returnErrorHTTPInterface(i.client, req, err, groups, 200)
	return err
}

// GroupGet gets the desired IAMGroup from the domain by id
func (i *IAMService) GroupGet(id string, group *IAMGroup) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("GET", "iam", fmt.Sprintf("/v1.0/group/%s", id), nil)
	_, err = returnErrorHTTPInterface(i.client, req, err, group, 200)
	return err
}

// GroupDelete deletes the desired group from IAM by id
func (i *IAMService) GroupDelete(id string) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/group/%s", id), nil)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

// GroupSetScopes set the scopes of a Group
func (i *IAMService) GroupSetScopes(id string, scopes []string) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("PUT", "iam", fmt.Sprintf("/v1.0/group/%s/scopes", id), scopes)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

// GroupDeleteScope deletes a scope of a group
func (i *IAMService) GroupDeleteScope(id string, scope string) error {
	if id == "" {
		return errIdentifierEmpty
	}
	req, err := i.client.NewRequest("DELETE", "iam", fmt.Sprintf("/v1.0/group/%s/scopes/%s", id, scope), nil)
	_, err = returnErrorHTTPSimple(i.client, req, err, 204)
	return err
}

// TODO: Search groups
