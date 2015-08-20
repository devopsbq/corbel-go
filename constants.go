package corbel

const (
	grantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	// AclAdmin is used to set admin privileges to something in a resource
	AclAdmin = "ADMIN"
	// AclWrite is used to set write privileges to something in a resource
	AclWrite = "WRITE"
	// AclRead is used to set read privileges to something in a resource
	AclRead = "READ"
	aclAll  = "ALL"
)
