package portscanner

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	ledger "github.com/Eiyarrr/Ledger"
)

// Portscanner scans the specified host in the inclusive port range using the
// given number of workers and timeout.
// Each result is printed to the console and logged using Ledger.
func Portscanner(startPort int, endPort int, host string, workers int, timeout time.Duration, scribe ledger.Ledger) error {
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
	message = "Worker count: " + strconv.Itoa(workers)
	err = scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Information printed and logged
	results := beginScanning(startPort, endPort, host, workers, timeout)
	for _, result := range results {
		message := fmt.Sprintf("Host: \"%s\" | Port: %d | Status: %s | Error: %v", result.Host, result.Port, result.Status.String(), result.Err)
		fmt.Println(message)
		err := scribe.Record(ledger.Info, message)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	return nil
}

// beginScanning distributes the specified port range across a worker pool and
// returns one Result for each port scanned.
func beginScanning(startPort int, endPort int, host string, workers int, timeout time.Duration) []Result {
	jobs := make(chan int)
	results := make(chan Result)

	var wg sync.WaitGroup

	// spawn the persistent workers immediately
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobs {
				Scan(host, job, timeout, results)
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

	retResults := []Result{}
	for result := range results {
		retResults = append(retResults, result)
	}

	return retResults
}
