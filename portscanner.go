package main

import (
	"errors"
	"fmt"
	"os"

	ledger "github.com/Eiyarrr/Ledger"
)

// Needs to get args from func instead
func portscanner(startPort int, endPort int, scribe ledger.Ledger) {
	// User arg verification
	if startPort > endPort {
		fmt.Fprintln(os.Stderr, errors.New("start port cannot be greater than end port"))
	}

	// Record information attained from user
	message := fmt.Sprintf("Port range: %d-%d", startPort, endPort)
	err := scribe.Record(ledger.Info, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
