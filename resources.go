package corbel

// ResourcesService handles the interface for retrival resource's representation
// on Corbel.
//
// Full API info: http://docs.corbelresources.apiary.io/
type ResourcesService struct {
	client *Client
}

// ResourceBasic models a resource item, any specific resource should use this as anonymous field
type ResourceBasic struct {
	service        *ResourcesService
	collectionName string
	ID             string            `json:"id,omitempty"`
	ACL            map[string]string `json:"_acl,omitempty"`
}

//Resource defines any resource
type Resource interface {
	setID(ID string)
	setCollection(collectionName string)
	setService(service *ResourcesService)
	setPermissionTo(ID, permission string) error
	remokePermissionTo(ID string) error

	AllowUser(userID, permission string) error
	RevokeUserPermission(userID string) error
	/*
		AllowGroup(groupID, permission string) error
		RevokeGroupPermission(userID string) error
	*/
	AllowAll(permission string) error
	RevokeAllPermission() error
}

func (r *ResourceBasic) setID(ID string) {
	r.ID = ID
}

func (r *ResourceBasic) setCollection(collectionName string) {
	r.collectionName = collectionName
}

func (r *ResourceBasic) setService(service *ResourcesService) {
	r.service = service
}

func (r *ResourceBasic) setPermissionTo(ID, permission string) error {
	if r.ACL == nil {
		r.ACL = make(map[string]string)
	}
	r.ACL[ID] = permission
	return r.service.UpdateResourceACL(r.collectionName, r.ID, r.ACL)
}

func (r *ResourceBasic) remokePermissionTo(ID string) error {
	if r.ACL == nil {
		return nil
	}
	delete(r.ACL, ID)
	return r.service.UpdateResourceACL(r.collectionName, r.ID, r.ACL)
}

//AllowUser grants the setted permission to the user with userID
func (r *ResourceBasic) AllowUser(userID, permission string) error {
	return r.setPermissionTo(userID, permission) // "userId:" + userID
}

//RevokeUserPermission revokes the setted permission to the user with userID
func (r *ResourceBasic) RevokeUserPermission(userID string) error {
	return r.remokePermissionTo(userID) // "userId:" + userID
}

/*
//AllowGroup grants the setted permission to the group with groupID
func (r *ResourceBasic) AllowGroup(groupID, permission string) error {
	return r.setPermissionTo(groupID, permission) // "groupId:" + groupID
}

//RevokeGroupPermission revokes the setted permission to the group with groupID
func (r *ResourceBasic) RevokeGroupPermission(groupID string) error {
	return r.remokePermissionTo(groupID) // "groupId:" + groupID
}
*/

//AllowAll grants the setted permission to all
func (r *ResourceBasic) AllowAll(permission string) error {
	return r.setPermissionTo(aclAll, permission)
}

//RevokeAllPermission revokes the setted permission to all
func (r *ResourceBasic) RevokeAllPermission() error {
	return r.remokePermissionTo(aclAll)
}
