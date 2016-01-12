package corbel

import "errors"

var (
	errMissingClientParams        = errors.New("Client: Missing parameters for the Client. client ID or Secret cannot be empty.")
	errInvalidEnvironment         = errors.New("Client: Environment is not valid.")
	errInvalidJWTSigningMethod    = errors.New("Client: Invalid JWT Signing Method.")
	errIdentifierEmpty            = errors.New("Client: Identifier can't be empty.")
	errUserNotFound               = errors.New("Client: User not found.")
	errInvalidTokenExpirationTime = errors.New("Client: Invalid TokenExpirationTime. Allowed range: 1-3600 seconds.")
	errHTTPNotAuthorized          = errors.New("HTTP: 401 Not authorized")
	errHTTPConflict               = errors.New("HTTP: 409 Conflict")
	errHTTPInvalidEntity          = errors.New("HTTP: 422 Invalid Entity")
	errJWTEncodingError           = errors.New("JWT: Encoding Error")
	errResponseError              = errors.New("HTTP: Response error")
	errURLParse                   = errors.New("HTTP: URL Parse Error")
	errJSONUnmarshalError         = errors.New("Encoding: JSON Unmarshal error")
	errInvalidLogLevel            = errors.New("Invalid log level")
)
