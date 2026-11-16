package scanner

import (
	"sync"
	"time"
)

func runScanPool(host string, workers int, timeout time.Duration, jobs <-chan int) ([]Result, error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	}

	results := make(chan Result)
	var wg sync.WaitGroup

	for range workers {
		wg.Go(func() {
			for job := range jobs {
				results <- Scan(host, job, timeout)
			}
		})
	}

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
