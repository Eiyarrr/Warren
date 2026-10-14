package portscanner_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/Eiyarrr/Warren/portscanner"
)

func TestOpenPort(t *testing.T) {
	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	// get the network address
	addr := listener.Addr().String()

	// get host and port as strings
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split listener address: %v", err)
	}

	// convert port to int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	result := make(chan portscanner.Result, 1)
	portscanner.Scan(host, port, 3, result)
	close(result)

	r := <-result
	if r.Status != portscanner.Open {
		t.Fatalf("expected status %q, got %q", portscanner.Open, r.Status)
	}
}

func TestClosedPort(t *testing.T) {
	// :0 makes os decide
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	// get the network address
	addr := listener.Addr().String()
	// close the listener to make a closed port
	err = listener.Close()
	if err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	// get host and port as strings
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split listener address: %v", err)
	}

	// convert port to int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	result := make(chan portscanner.Result, 1)
	portscanner.Scan(host, port, 3, result)
	close(result)

	r := <-result
	if r.Status != portscanner.Closed {
		t.Fatalf("expected status %q, got %q", portscanner.Closed, r.Status)
	}
}
