package portscanner

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Portscanner scans the specified host in the inclusive port range using the
// given number of workers and timeout.
// Each result is printed to the console and logged using Ledger.
func Portscanner(startPort int, endPort int, host string, workers int, timeout time.Duration) ([]Result, error) {
	if startPort > endPort {
		return nil, errors.New("start port cannot be greater than end port")
	} else if startPort < 0 || endPort > 65535 {
		return nil, errors.New("ports must be between 0 and 65535")
	} else if workers <= 0 {
		return nil, errors.New("worker count must be greater than 0")
	} else if timeout <= 0 {
		return nil, errors.New("timeout must be greater than 0")
	}

	fmt.Printf("Port range: %d-%d\n", startPort, endPort)
	fmt.Println("Hostname/IP: \"" + host + "\"")
	fmt.Println("Worker count: " + strconv.Itoa(workers))

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
