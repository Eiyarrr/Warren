package scanner_test

import (
	"net"
	"strconv"
	"testing"
	"time"
	"fmt"

	scanner "github.com/Eiyarrr/Warren/scanner"
)

func TestOpenPort(t *testing.T) {
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

	result := scanner.Scan(host, port, 3*time.Second)
	if result.Status != scanner.Open {
		t.Fatalf("expected status %q, got %q", scanner.Open, result.Status)
	}
}

func TestClosedPort(t *testing.T) {
	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	addr := listener.Addr().String()
	host, port, err := splitAddr(addr)
	if err != nil {
		t.Fatal(err)
	}

	err = listener.Close()
	if err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}
	result := scanner.Scan(host, port, 3*time.Second)
	if result.Status != scanner.Closed {
		t.Fatalf("expected status %q, got %q", scanner.Closed, result.Status)
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
