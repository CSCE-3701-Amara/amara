package http

import "errors"

var ErrMissingRegistrationHandler = errors.New("registration handler must not be nil")
