package portscanner

import (
	"time"
)

func ScanCommon(host string, workers int, timeout time.Duration) ([]Result, error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return nil, ErrInvalidTimeout
	}

	jobs := make(chan int)
	go func() {
		for port := range CommonPorts {
			jobs <- port
		}
		close(jobs)
	}()

	results, err := runScanPool(host, workers, timeout, jobs)
	if err != nil {
		return results, err
	}

	return results, nil
}
