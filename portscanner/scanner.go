package portscanner

import (
	"errors"
	"fmt"
	"os"
	"sync"

	ledger "github.com/Eiyarrr/Ledger"
)

// Needs to get args from func instead
func Portscanner(startPort int, endPort int, host string, workers int, scribe ledger.Ledger) error {
	// User arg verification
	if startPort > endPort {
		return errors.New("start port cannot be greater than end port")
	} else if startPort < 0 || endPort > 65535 {
		return errors.New("ports must be between 0 and 65535")
	}
	if workers <= 0 {
		return errors.New("worker count must be greater than 0")
	}

	// Record information attained from user
	message := fmt.Sprintf("Port range: %d-%d", startPort, endPort)
	err := scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	message = "Hostname/IP: \"" + host + "\""
	err = scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Information printed and logged
	beginScanning(startPort, endPort, host, workers, scribe)

	return nil
}

// Looping & semaphore & randomization logic
func beginScanning(startPort int, endPort int, host string, workers int, scribe ledger.Ledger) {
	jobs := make(chan int)
	results := make(chan Result)

	var wg sync.WaitGroup

	// spawn the persistent workers immediately
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobs {
				scan(host, job, results)
			}
		}()
	}

	// Inside a separate goroutine to prevent blocking from others
	go func() {
		for port := startPort; port <= endPort; port++ {
			jobs <- port
		}
		close(jobs)
	}()

	// also must be goroutine to prevent blocking
	go func() {
		// Wait() waits for all workers to call Done()
		wg.Wait()
		close(results)
	}()

	for result := range results {
		message := fmt.Sprintf("Host: \"%s\" | Port: %d | Status: %s | Error: %v", result.host, result.port, result.status.String(), result.err)
		fmt.Println(message)
		err := scribe.Record(ledger.Info, message)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
