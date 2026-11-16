package bannergrabber

import (
	"errors"
	"net"
	"strconv"
	"time"
)

func GrabBanner(host string, port int, timeout time.Duration, maxBytes int) (string, bool, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "", false, nil
		}
		return "", false, err
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return "", false, err
	}
	truncate := false

	buf := make([]byte, maxBytes)
	n, err := conn.Read(buf)
	if err != nil {
		return "", false, err
	}
	if n > maxBytes {
		truncate = true
	}

	banner := string(buf)

	return banner, truncate, nil
}
