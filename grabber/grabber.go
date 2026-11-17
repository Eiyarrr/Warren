package grabber

import (
	"errors"
	"io"
	"net"
	"strconv"
	"syscall"
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
		if errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) {
			return "", false, nil
		}
		return "", false, err
	}
	defer conn.Close()

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return "", false, err
	}

	buf := make([]byte, maxBytes+1)
	totalWritten := 0
	for {
		// truncation is not an error, needs a new way to find
		n, err := conn.Read(buf[totalWritten:])
		totalWritten += n

		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break
			}
			if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.EPIPE) {
				break
			}
			return "", false, err
		}

		if totalWritten > maxBytes {
			break
		}

		if n == 0 {
			break
		}
	}
	truncate := false
	if totalWritten > maxBytes {
		truncate = true
		totalWritten--
	}

	banner := string(buf[:totalWritten])

	return banner, truncate, nil
}
