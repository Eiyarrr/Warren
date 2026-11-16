package grabber_test

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	grabber "github.com/Eiyarrr/Warren/grabber"
)

func TestNoSend(t *testing.T) {
	timeout := 3 * time.Second
	maxBytes := 1024

	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	host, port, err := splitAddr(addr)
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = grabber.GrabBanner(host, port, timeout, maxBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormal(t *testing.T) {
	timeout := 3 * time.Second
	maxBytes := 1024

	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	host, port, err := splitAddr(addr)
	if err != nil {
		t.Fatal(err)
	}

	testBanner := "TEST BANNER"
	errChan := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- fmt.Errorf("unexpected error accepting connection: %v", err)
			return
		}
		defer conn.Close()

		err = conn.SetWriteDeadline(time.Now().Add(timeout))
		if err != nil {
			errChan <- fmt.Errorf("unexpected error setting write deadline: %v", err)
			return
		}

		bannerBytes := []byte(testBanner)
		n, err := conn.Write(bannerBytes)
		if err != nil {
			errChan <- fmt.Errorf("unexpected error when writing to connection: %v", err)
			return
		}
		if n != len(bannerBytes) {
			errChan <- fmt.Errorf("failed to write all bytes of banner")
			return
		}

		errChan <- nil
	}()

	banner, truncated, err := grabber.GrabBanner(host, port, timeout, maxBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if banner != testBanner {
		t.Fatalf("incorrect banner. got %q, expected %q", banner, testBanner)
	}
	if truncated {
		t.Fatalf("incorrect truncated status. got %v, expected %v", truncated, false)
	}

	err = <-errChan
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantClose(t *testing.T) {
	timeout := 3 * time.Second
	maxBytes := 1024

	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	host, port, err := splitAddr(addr)
	if err != nil {
		t.Fatal(err)
	}

	errChan := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- fmt.Errorf("unexpected error accepting connection: %v", err)
			return
		}
		err = conn.Close()
		if err != nil {
			errChan <- fmt.Errorf("unexpected error closing connection: %v", err)
			return
		}
		errChan <- nil
	}()

	banner, truncated, err := grabber.GrabBanner(host, port, timeout, maxBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
    if banner != "" {
        t.Fatalf("expected empty banner, got %q", banner)
    }
	if truncated {
		t.Fatalf("incorrect truncated status. got %v, expected %v", truncated, false)
	}

	err = <-errChan
	if err != nil {
		t.Fatal(err)
	}
}

func TestTruncation(t *testing.T) {
	timeout := 3 * time.Second
	maxBytes := 1024

	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()
	host, port, err := splitAddr(addr)
	if err != nil {
		t.Fatal(err)
	}

	testBanner := string(make([]byte, maxBytes+1))
	errChan := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			errChan <- fmt.Errorf("unexpected error accepting connection: %v", err)
			return
		}
		defer conn.Close()

		err = conn.SetWriteDeadline(time.Now().Add(timeout))
		if err != nil {
			errChan <- fmt.Errorf("unexpected error setting write deadline: %v", err)
			return
		}

		bannerBytes := []byte(testBanner)
		n, err := conn.Write(bannerBytes)
		if err != nil {
			errChan <- fmt.Errorf("unexpected error when writing to connection: %v", err)
			return
		}
		if n != len(bannerBytes) {
			errChan <- fmt.Errorf("failed to write all bytes of banner")
			return
		}

		errChan <- nil
	}()

	banner, truncated, err := grabber.GrabBanner(host, port, timeout, maxBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if banner != testBanner[:maxBytes] {
		t.Fatalf("incorrect banner. got %q, expected %q", banner, testBanner[:maxBytes])
	}
	if !truncated {
		t.Fatalf("incorrect truncated status. got %v, expected %v", truncated, true)
	}

	err = <-errChan
	if err != nil {
		t.Fatal(err)
	}
}

func splitAddr(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return host, 0, fmt.Errorf("failed to parse port: %v", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return host, port, fmt.Errorf("failed to parse port: %v", err)
	}

	return host, port, err
}
