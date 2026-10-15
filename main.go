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
			os.Exit(1)
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

}

func runScan(args []string) error {
	scanFlags := flag.NewFlagSet("scan", flag.ContinueOnError)

	host := scanFlags.String("host", "127.0.0.1", "Target host")
	startPort := scanFlags.Int("start-port", 1, "First port in the range")
	endPort := scanFlags.Int("end-port", 1024, "Last port in the range")
	workers := scanFlags.Int("workers", 50, "Number of Go workers to use")
	timeout := scanFlags.Duration("timeout", 3*time.Second, "Time before scanner assumes a timeout")

	err := scanFlags.Parse(args)
	if err != nil {
		return err
	}

	err = verifyFlags(*startPort, *endPort, *workers, *timeout)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scribe, err := ledger.New(file)
	if err != nil {
		return err
	}

	err = scribe.Record(ledger.Info, "Ledger loaded!")
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Host: %s, Start port: %d, End port: %d, Worker count: %d, Timeout: %d\n", *host, *startPort, *endPort, *workers, *timeout)
	err = scribe.Record(ledger.Info, message)
	if err != nil {
		return err
	}

	results, err := portscanner.ScanRange(*startPort, *endPort, *host, *workers, *timeout)
	if err != nil {
		return err
	}

	for _, result := range results {
		fmt.Printf("Host: \"%s\" | Port: %d | Status: %s | Error: %v\n", result.Host, result.Port, result.Status.String(), result.Err)
	}

	return nil
}

// verifyFlags verifies that all user given flags are allowed by Warren.
func verifyFlags(startPort int, endPort int, workers int, timeout time.Duration) error {
	if startPort > endPort || startPort < 0 || endPort > 65535 {
		return portscanner.ErrInvalidPortRange
	} else if workers <= 0 {
		return portscanner.ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return portscanner.ErrInvalidTimeout
	}

	return nil
}
