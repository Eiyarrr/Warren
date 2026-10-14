package portscanner

import (
	"errors"
	"net"
	"strconv"
	"syscall"
	"time"
)

// Scans a single port
func Scan(host string, port int, timeoutSeconds time.Duration, results chan<- Result) {
	result := Result{
		Host: host,
		Port: port,
	}

	// create network address
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// attempt to make a TCP connection onto addr
	// note: must have conn var to allow for closing of the connection even
	//       with _, the connection is still made, just not set to a var
	conn, err := net.DialTimeout("tcp", addr, timeoutSeconds)
	if err != nil {
		// errors.As() turns err into a netErr
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			result.Status = Timeout
		} else if errors.Is(err, syscall.ECONNREFUSED) {
			result.Status = Closed
		} else {
			result.Status = Err
		}
		result.Err = err

		results <- result
		return
	}
	defer conn.Close()

	result.Status = Open
	result.Err = nil
	results <- result
}
