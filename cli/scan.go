package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	ledger "github.com/Eiyarrr/Ledger"
	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

func runScan(args []string, logFile string) error {
	scanFlags := flag.NewFlagSet("scan", flag.ContinueOnError)

	host := scanFlags.String("host", "127.0.0.1", "Target host")
	startPort := scanFlags.Int("start-port", 1, "First port in the range")
	endPort := scanFlags.Int("end-port", 1024, "Last port in the range")
	workers := scanFlags.Int("workers", 50, "Number of Go workers to use")
	timeout := scanFlags.Duration("timeout", 3*time.Second, "Time before scanner assumes a timeout")
	common := scanFlags.Bool("common", false, "Whether only the common ports should be scanned")
	openOnly := scanFlags.Bool("open-only", false, "Whether to only print the open ports")

	err := scanFlags.Parse(args)
	if errors.Is(err, flag.ErrHelp) {
		return nil
	}
	if err != nil {
		return err
	}

	err = verifyScanFlags(*startPort, *endPort, *workers, *timeout)
	if err != nil {
		return err
	}

	if logFile != "" {
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

		fmt.Printf("Port range: %d-%d\n", startPort, endPort)
		fmt.Printf("Hostname/IP: %q\n", *host)
		fmt.Printf("Worker count: %d\n", *workers)
		fmt.Printf("Timeout: %s\n", timeout.String())
		fmt.Print("Common?: ", *common)
		fmt.Print("OpenOnly?: ", *openOnly)

		message := fmt.Sprintf("Host: %s, Start port: %d, End port: %d, Worker count: %d, Timeout: %s, Common?: %v, OpenOnly?: %v\n", *host, *startPort, *endPort, *workers, *timeout, *common, *openOnly)
		err = scribe.Record(ledger.Info, message)
		if err != nil {
			return err
		}
	}

	results, err := portscanner.ScanRange(*startPort, *endPort, *common, *host, *workers, *timeout)
	if err != nil {
		return err
	}

	if *openOnly {
		for _, result := range results {
			if result.Status == portscanner.Open {
				fmt.Printf("Host: \"%s\" | Port: %d | Status: %s | Error: %v\n", result.Host, result.Port, result.Status.String(), result.Err)
			}
		}
	} else {
		for _, result := range results {
			fmt.Printf("Host: \"%s\" | Port: %d | Status: %s | Error: %v\n", result.Host, result.Port, result.Status.String(), result.Err)
		}
	}

	return nil
}

func verifyScanFlags(startPort int, endPort int, workers int, timeout time.Duration) error {
	if startPort > endPort || startPort < 0 || endPort > 65535 || startPort > 65535 || endPort < 0 {
		return portscanner.ErrInvalidPortRange
	} else if workers <= 0 {
		return portscanner.ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return portscanner.ErrInvalidTimeout
	}

	return nil
}
