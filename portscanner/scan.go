package portscanner

import (
	"errors"
	"net"
	"strconv"
	"syscall"
	"time"
)

// Scan scans a single port, returning the result into a channel.
// Scan attempts to establish a TCP connection to the specified host and port
// using the provided timeout.
// It returns exactly one Result to results, indicating whether the port is
// open, closed, timed out, or encountered some other error.
//
// The timeout must be positive, the results channel must be writeable.
func Scan(host string, port int, timeout time.Duration, results chan<- Result) {
	result := Result{
		Host: host,
		Port: port,
	}

	if timeout <= 0 {
		result.Status = Err
		result.Err = errors.New("timeout must be greater than 0")

		results <- result
		return
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
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
