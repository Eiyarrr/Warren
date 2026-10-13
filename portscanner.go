package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"syscall"
	"time"

	ledger "github.com/Eiyarrr/Ledger"
)

type Status int

const (
	Open    Status = iota //0
	Closed                //1
	Timeout               //2
	Err                   //3
)

type Result struct {
	host   string
	port   int
	status Status
	err    error
}

// Needs to get args from func instead
func portscanner(startPort int, endPort int, host string, scribe ledger.Ledger) error {
	// User arg verification
	if startPort > endPort {
		return errors.New("start port cannot be greater than end port")
	}
	if startPort < 0 || endPort > 65535 {
		return errors.New("ports must be between 0 and 65535")
	}

	// Record information attained from user
	message := fmt.Sprintf("Port range: %d-%d", startPort, endPort)
	_ = scribe.Record(ledger.Info, message)

	results := beginScanning(startPort, endPort, host)
	fmt.Println(results)
	// add scribe recording

	return nil
}

// Looping & semaphore & randomization logic
func beginScanning(startPort int, endPort int, host string) []Result {
	results := []Result{}
	for port := startPort; port <= endPort; port++ {
		results = append(results, scan(host, port))
	}

	return results
}

// Scans a single port
func scan(host string, port int) Result {
	result := Result{
		host: host,
		port: port,
	}

	// create network address
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// attempt to make a TCP connection onto addr
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		var netErr net.Error

		// errors.As() turns err into a netErr
		if errors.As(err, &netErr) && netErr.Timeout() {
			result.status = Timeout
		} else if errors.Is(err, syscall.ECONNREFUSED) {
			result.status = Closed
		} else {
			result.status = Err
		}

		result.err = err

		return result
	}
	defer conn.Close()

	result.status = Open
	result.err = nil
	return result
}
