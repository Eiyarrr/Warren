// Warren is a network service mapper.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	ledger "github.com/Eiyarrr/Ledger"
	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

const LOGFILE = "warren_log"

func main() {
	startPort := flag.Int("sPort", 1, "Port to begin scanning with")
	endPort := flag.Int("ePort", 65535, "Port to end scanning with")
	hostname := flag.String("host", "", "Host to scan")
	workers := flag.Int("workers", 50, "Number of Go workers to use")
	timeoutSeconds := flag.Int("timeout", 3, "Seconds before a timeout")

	flag.Parse()

	err := verifyFlags(*startPort, *endPort, *workers, *timeoutSeconds)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer file.Close()

	scribe, err := ledger.New(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = scribe.Record(ledger.Info, "Ledger loaded!")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = portscanner.Portscanner(*startPort, *endPort, *hostname, *workers, time.Duration(*timeoutSeconds)*time.Second, *scribe)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// verifyFlags verifies that all user given flags are allowed by Warren.
func verifyFlags(startPort int, endPort int, workers int, timeoutSeconds int) error {
	if startPort > endPort {
		return errors.New("start port cannot be greater than end port")
	} else if startPort < 0 || endPort > 65535 {
		return errors.New("ports must be between 0 and 65535")
	}
	if workers <= 0 {
		return errors.New("worker count must be greater than 0")
	}
	if timeoutSeconds <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	return nil
}
