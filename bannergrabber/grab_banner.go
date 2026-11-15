package bannergrabber

import (
	"net"
	"strconv"
	"time"
)

// takes in info, makes conn, reads, returns banner, if truncated, err
func GrabBanner(host string, port int, timeout time.Duration, maxBytes int) (string, bool, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return "", false, err
	}

	truncate := false

	buf := make([]byte, maxBytes)
	n, err := conn.Read(buf)
	if n > maxBytes {
		truncate = true
	}

	banner := string(buf)

	return banner, truncate, nil
}
