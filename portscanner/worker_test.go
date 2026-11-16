package portscanner

import (
	"testing"
	"time"
)

func TestWorkers(t *testing.T) {
	host := "127.0.0.1"
	startPort := 50000
	endPort := 50009
	workers := 3
	timeout := 3 * time.Second

	jobs := make(chan int)

	results, err := runScanPool(host, workers, timeout, jobs)
	if err != nil {
		t.Fatalf("unexpected error running scan pool: %v", err)
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			jobs <- port
		}
		close(jobs)
	}()

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
