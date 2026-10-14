package portscanner

import (
	"testing"
)

func TestWorkers(t *testing.T) {
	host := "127.0.0.1"
	startPort := 50000
	endPort := 50009
	workers := 3
	timeoutSeconds := 3

	results := beginScanning(startPort, endPort, host, workers, timeoutSeconds)

	expectedCount := endPort - startPort + 1
	seenPorts := make(map[int]bool)

	if len(results) != expectedCount {
		t.Fatalf("expected %d results, got %d", expectedCount, len(results))
	}

	for _, result := range results {
		port := result.Port

		if port < startPort || port > endPort {
			t.Fatalf("port %d is outside the requested range", port)
		}

		if seenPorts[port] {
			t.Fatalf("expected unique ports, got duplicate at %d", port)
		}
		seenPorts[port] = true
	}
}
