package grabber

import "errors"

// Errors returned by Warren banner grabbing operations.
var (
	// ErrBadMaxBytes indicates that the max bytes given was 0 or less
	ErrBadMaxBytes = errors.New("max bytes must be greater than 0")
)
