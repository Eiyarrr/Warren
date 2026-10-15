package portscanner

// Result contains the target address, port, status, and any error
// that was encountered during a TCP port scan.
type Result struct {
	Host   string
	Port   int
	Status Status
	Err    error
}
