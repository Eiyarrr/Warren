package portscanner

import (
	"time"
	"sync"
)

func ScanCommon(host string, workers int, timeout time.Duration) ([]Result, error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return nil, ErrInvalidTimeout
	}

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
		for port := range commonPorts {
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

	return retResults, nil
}
