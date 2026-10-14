package portscanner

import (
	"errors"
	"net"
	"strconv"
	"syscall"
	"time"
)

// Scans a single port
func scan(host string, port int, results chan <- Result) {
	result := Result{
		host: host,
		port: port,
	}

	// create network address
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// attempt to make a TCP connection onto addr
	// note: must have conn var to allow for closing of the connection even
	//       with _, the connection is still made, just not set to a var
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		// errors.As() turns err into a netErr
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			result.status = Timeout
		} else if errors.Is(err, syscall.ECONNREFUSED) {
			result.status = Closed
		} else {
			result.status = Err
		}
		result.err = err

		results <- result
		return
	}
	defer conn.Close()

	result.status = Open
	result.err = nil
	results <- result
}
