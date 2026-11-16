package bannergrabber_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	bannergrabber "github.com/Eiyarrr/Warren/bannergrabber"
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

	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("failed to split listener address: %v", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("failed to parse port: %v", err)
	}

	_, _, err = bannergrabber.GrabBanner(host, port, timeout, maxBytes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

/*
needs:
timeout
truncation
empty banner
immediate close
*/
