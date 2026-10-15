// Warren is a network service mapper.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	ledger "github.com/Eiyarrr/Ledger"
	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

const logFile = "warren_log"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Interactive mode is not implemented yet.")
		return
	}

	switch os.Args[1] {
	case "scan":
		err := runScan(os.Args[2:])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	case "help":
		fmt.Println("Help!")
		return
	case "version":
		fmt.Println("v0.1.0 - DEV")
		return
	default:
		fmt.Println("Error, unknown command")
		os.Exit(1)
	}

	startPort := flag.Int("sPort", 1, "Port to begin scanning with")
	endPort := flag.Int("ePort", 65535, "Port to end scanning with")
	hostname := flag.String("host", "", "Host to scan")
	workers := flag.Int("workers", 50, "Number of Go workers to use")
	timeoutSeconds := flag.Int("timeout", 3, "Seconds before a timeout")
	common := flag.Bool("common", false, "Scan just the common ports?")

	flag.Parse()

	err := verifyFlags(*startPort, *endPort, *workers, *timeoutSeconds)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	var results []portscanner.Result
	err = nil
	if *common {
		results, err = portscanner.ScanCommon(*hostname, *workers, time.Duration(*timeoutSeconds)*time.Second)
	} else {
		results, err = portscanner.ScanRange(*startPort, *endPort, *hostname, *workers, time.Duration(*timeoutSeconds)*time.Second)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		message := fmt.Sprintf("Error occurred while scanning: %v", err)
		err = scribe.Record(ledger.Fatal, message)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}

	for _, result := range results {
		fmt.Printf("Host: \"%s\" | Port: %d | Status: %s | Error: %v\n", result.Host, result.Port, result.Status.String(), result.Err)
	}
}

// verifyFlags verifies that all user given flags are allowed by Warren.
func verifyFlags(startPort int, endPort int, workers int, timeoutSeconds int) error {
	if startPort > endPort || startPort < 0 || endPort > 65535 {
		return portscanner.ErrInvalidPortRange
	} else if workers <= 0 {
		return portscanner.ErrInvalidWorkerCount
	} else if timeoutSeconds <= 0 {
		return portscanner.ErrInvalidTimeout
	}

	return nil
}

func runScan(args []string) error {
	scanFlags := flag.NewFlagSet("scan", flag.ContinueOnError)

	host := scanFlags.String("host", "127.0.0.1", "Target host")
	startPort := scanFlags.Int("start-port", 1, "First port in the range")
	endPort := scanFlags.Int("end-port", 65535, "Last port in the range")

	err := scanFlags.Parse(args)
	if err != nil {
		return err
	}

	fmt.Printf("Host: %s, Start port: %d, End port: %d\n", *host, *startPort, *endPort)

	return nil
}
