package main

import (
	"flag"
	"fmt"
	"os"

	ledger "github.com/Eiyarrr/Ledger"
)

const LOGFILE = "warren_log"

// Needs to ask which service to use when more are created
func main() {
	// Get args
	startPort := flag.Int("sPort", 0, "Port to begin scanning with")
	endPort := flag.Int("ePort", 65535, "Port to end scanning with")

	// Lock in args
	flag.Parse()

	// Get ledger up and running
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

	err = portscanner(*startPort, *endPort, *scribe)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
