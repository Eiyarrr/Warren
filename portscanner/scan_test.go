package portscanner_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/Eiyarrr/Warren/portscanner"
)

func TestOpenPort(t *testing.T) {
	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split listener address: %v", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	result := portscanner.Scan(host, port, 3*time.Second)

	if result.Status != portscanner.Open {
		t.Fatalf("expected status %q, got %q", portscanner.Open, result.Status)
	}
}

func TestClosedPort(t *testing.T) {
	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	addr := listener.Addr().String()
	err = listener.Close()
	if err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split listener address: %v", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	result := portscanner.Scan(host, port, 3*time.Second)
	if result.Status != portscanner.Closed {
		t.Fatalf("expected status %q, got %q", portscanner.Closed, result.Status)
	}
}
