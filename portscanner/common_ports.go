package portscanner

import (
	"fmt"
	"strconv"
	"time"
)

func ScanCommon(host string, workers int, timeout time.Duration) ([]Result, error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return nil, ErrInvalidTimeout
	}

	fmt.Println("Hostname/IP: \"" + host + "\"")
	fmt.Println("Worker count: " + strconv.Itoa(workers))
	fmt.Println("Timeout: " + timeout.String())

	return nil, nil
}
