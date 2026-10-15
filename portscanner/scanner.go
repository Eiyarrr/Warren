package portscanner

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// ScanRange scans the specified host in the inclusive port range using the
// given number of workers and timeout.
// Each result is printed to the console and logged using Ledger.
func ScanRange(startPort int, endPort int, host string, workers int, timeout time.Duration) ([]Result, error) {
	if startPort > endPort || startPort < 0 || endPort > 65535 {
		return nil, ErrInvalidPortRange
	} else if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return nil, ErrInvalidTimeout
	}

	fmt.Printf("Port range: %d-%d\n", startPort, endPort)
	fmt.Println("Hostname/IP: \"" + host + "\"")
	fmt.Println("Worker count: " + strconv.Itoa(workers))
	fmt.Println("Timeout: " + timeout.String())

	results := beginScanning(startPort, endPort, host, workers, timeout)

	return results, nil
}

// beginScanning distributes the specified port range across a worker pool and
// returns one Result for each port scanned.
func beginScanning(startPort int, endPort int, host string, workers int, timeout time.Duration) []Result {
	jobs := make(chan int)
	results := make(chan Result)

	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for job := range jobs {
				Scan(host, job, timeout, results)
			}
		}()
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			jobs <- port
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	retResults := []Result{}
	for result := range results {
		retResults = append(retResults, result)
	}

	return retResults
}
