package portscanner

import (
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

	jobs := make(chan int)

	results, err := runScanPool(host, workers, timeout, jobs)
	if err != nil {
		return results, err
	}

	go func() {
		for port := startPort; port <= endPort; port++ {
			jobs <- port
		}
		close(jobs)
	}()

	return results, nil
}
