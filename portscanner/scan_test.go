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

	// get host and port as a string
	host, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("failed to split listener address: %v", err)
	}

	// convert port to int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	results := make(chan portscanner.Result, 1)
	portscanner.Scan(host, port, results)
	close(results)

	r := <-results
	if r.Status.String() != "Open" {
		t.Fatalf("expected status \"Open\", got \"%s\"", r.Status.String())
	}
}
