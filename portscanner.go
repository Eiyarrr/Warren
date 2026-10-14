package main

import (
	"errors"
	"fmt"
	"net"
	"os"
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
func portscanner(startPort int, endPort int, host string, workers int, scribe ledger.Ledger) error {
	// User arg verification
	if startPort > endPort {
		return errors.New("start port cannot be greater than end port")
	} else if startPort < 0 || endPort > 65535 {
		return errors.New("ports must be between 0 and 65535")
	}
	if workers <= 0 {
		return errors.New("worker count must be greater than 0")
	}

	// Record information attained from user
	message := fmt.Sprintf("Port range: %d-%d", startPort, endPort)
	err := scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	message = "Hostname/IP: \"" + host + "\""
	err = scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Information printed and logged
	beginScanning(startPort, endPort, host, workers, scribe)

	return nil
}

// Looping & semaphore & randomization logic
func beginScanning(startPort int, endPort int, host string, workers int, scribe ledger.Ledger) {
	for port := startPort; port <= endPort; port++ {
		result := scan(host, port)

		message := fmt.Sprintf("Host: \"%s\" | Port: %d | Status: %s | Error: %s", result.host, result.port, result.status.String(), result.err)
		fmt.Println(message)
		err := scribe.Record(ledger.Info, message)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
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
	// note: must have conn var to allow for closing of the connection even
	//       with _, the connection is still made, just not set to a var
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		// errors.As() turns err into a netErr
		var netErr net.Error
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

func (s Status) String() string {
	switch s {
	case Open:
		return "Open"
	case Timeout:
		return "Timeout"
	case Closed:
		return "Closed"
	case Err:
		return "Err"
	default:
		return ""
	}
}
