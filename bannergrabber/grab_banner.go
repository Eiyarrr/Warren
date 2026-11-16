package bannergrabber

import (
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

func GrabBanner(host string, port int, timeout time.Duration, maxBytes int) (string, bool, error) {
	if maxBytes <= 0 {
		return "", false, ErrBadMaxBytes
	}

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

	buf := make([]byte, maxBytes)
	totalWritten := 0
	for {
		// truncation is not an error, needs a new way to find
		n, err := conn.Read(buf[totalWritten:])
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return "", false, err
		}

		totalWritten += n

		if totalWritten >= maxBytes {
			break
		}
	}

	banner := string(buf[:totalWritten])

	return banner, false, nil
}
