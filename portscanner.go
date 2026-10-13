package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
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
func portscanner(startPort int, endPort int, scribe ledger.Ledger) error {
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

	results := beginScanning(startPort, endPort, scribe)
	fmt.Println(results)

	return nil
}

// Looping & semaphore & randomization logic
func beginScanning(startPort int, endPort int, scribe ledger.Ledger) []Result {
	fmt.Println("Beginning port scan")
	_ = scribe.Record(ledger.Info, "Beginning port scan")

	for port := startPort; port <= endPort; port++ {
		fmt.Println(port)
	}

	return nil
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

		// converts err to a net.Error type then checks if it's a timeout
		if errors.As(err, &netErr) && netErr.Timeout() {
			result.status = Timeout
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
