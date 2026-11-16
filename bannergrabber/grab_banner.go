package bannergrabber

import (
	"context"
	"errors"
	"io"
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
	// truncation is not an error, needs a new way to find
	n, err := conn.Read(buf)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "", false, nil
		}
		if errors.Is(err, io.EOF) {
			return "", false, nil
		}
		return "", false, err
	}

	banner := string(buf[:n])

	return banner, truncate, nil
}
