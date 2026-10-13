package main

import (
	"errors"
	"fmt"

	ledger "github.com/Eiyarrr/Ledger"
)

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
	err := scribe.Record(ledger.Info, message)
	if err != nil {
		return err
	}


	portStatuses := beginScanning(startPort, endPort, scribe)
	fmt.Println(portStatuses)

	return nil
}

// Looping & semaphore & randomization logic
func beginScanning(startPort int, endPort int, scribe ledger.Ledger) []bool {
	fmt.Println("Beginning port scan")
	_ = scribe.Record(ledger.Info, "Beginning port scan")

	for port := startPort; port <= endPort; port++ {
		fmt.Println(port)
	}

	return nil
}

// Scans a single port
func scan(port int) bool {
	return false
}
