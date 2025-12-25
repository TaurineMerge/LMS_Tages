package testing

import "errors"

var (
	// ErrTestNotFound is returned when the testing service responds with a "not_found" status.
	ErrTestNotFound = errors.New("test not found")

	// ErrServiceUnavailable is returned when the testing service is not reachable.
	ErrServiceUnavailable = errors.New("testing service is unavailable")

	// ErrInvalidResponse is returned when the response from the testing service fails schema validation.
	ErrInvalidResponse = errors.New("invalid response from testing service")
)
