package portscanner

import "errors"

// Errors returned by Warren portscanner operations.
var (
	// ErrInvalidPortRange indicates that the specified port range was invalid.
	ErrInvalidPortRange = errors.New("invalid port range")

	// ErrInvalidWorkerCount indicates that the worker count is less than 1.
	ErrInvalidWorkerCount = errors.New("worker count must be greater than 0")

	// ErrInvalidTimeout indicates that the timeout is 0 or negative.
	ErrInvalidTimeout = errors.New("timeout must be greater than 0")
)
