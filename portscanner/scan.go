package portscanner

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// ScanRange scans the specified host in the inclusive port range using the
// given number of workers and timeout.
// Each result is printed to the console and logged using Ledger.
func ScanRange(startPort int, endPort int, common bool, host string, workers int, timeout time.Duration) ([]Result, error) {
	if common {
		return ScanCommon(host, workers, timeout)
	}

	if startPort > endPort || startPort < 0 || endPort > 65535 {
		return nil, ErrInvalidPortRange
	} else if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return nil, ErrInvalidTimeout
	}

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

// Scan scans a single port, returning the result into a channel.
// Scan attempts to establish a TCP connection to the specified host and port
// using the provided timeout.
// It returns exactly one Result to results, indicating whether the port is
// open, closed, timed out, or encountered some other error.
//
// The timeout must be positive, the results channel must be writeable.
func Scan(host string, port int, timeout time.Duration, results chan<- Result) {
	result := Result{
		Host: host,
		Port: port,
	}

	if timeout <= 0 {
		result.Status = Err
		result.Err = ErrInvalidTimeout

		results <- result
		return
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			result.Status = Timeout
		} else if errors.Is(err, syscall.ECONNREFUSED) {
			result.Status = Closed
		} else {
			result.Status = Err
		}
		result.Err = err

		results <- result
		return
	}
	defer conn.Close()

	result.Status = Open
	result.Err = nil
	results <- result
}
