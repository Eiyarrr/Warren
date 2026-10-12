package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	ledger "github.com/Eiyarrr/Ledger"
)

const LOGFILE = "warren_log"

func main() {
	// Get args
	startPort := flag.Int("sPort", 0, "Port to begin scanning with")
	endPort := flag.Int("ePort", 65535, "Port to end scanning with")

	// Lock in args
	flag.Parse()

	// User arg verification
	if *startPort > *endPort {
		fmt.Fprintln(os.Stderr, errors.New("start port cannot be greater than end port"))
	}

	// Get ledger up and running
	file, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	scribe, err := ledger.New(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = scribe.Record(ledger.Info, "Ledger loaded!")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// Record information attained from user
	message := fmt.Sprintf("Port range: %d-%d", *startPort, *endPort)
	err = scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
