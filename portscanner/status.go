package portscanner

// Status represents the possible states of a port.
type Status int

// Predefined statuses of ports.
const (
	// Open indicates that a successful TCP connection was established
	Open Status = iota

	// Closed indicates that the connection attempt was refused.
	Closed

	// Timeout indicates that the connection attempt exceeded the allotted timespan.
	// This may indicate filtering, but does not prove the port is closed.
	Timeout

	// Err indicates that there was some kind of other error that occurred.
	Err
)

// String returns a human readable representation of the status.
func (s Status) String() string {
	switch s {
	case Open:
		return "Open"
	case Timeout:
		return "Timeout"
	case Closed:
		return "Closed"
	case Err:
		return "Err"
	default:
		return "Unknown"
	}
}
