package silkroad

import "errors"

var (
	errMissingClientParams = errors.New("Missing parameters for the Client. client ID or Secret cannot be empty.")
  errInvalidEnvironment = errors.New("Environment is not valid.")
)
