package silkroad

import "errors"

var (
	errMissingClientParams        = errors.New("Missing parameters for the Client. client ID or Secret cannot be empty.")
	errInvalidEnvironment         = errors.New("Environment is not valid.")
	errInvalidJWTSigningMethod    = errors.New("Invalid JWT Signing Method.")
	errInvalidTokenExpirationTime = errors.New("Invalid TokenExpirationTime. Must be 1-3600 seconds.")
)
