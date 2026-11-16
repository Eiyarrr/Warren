package cli

import (
	"errors"
	"flag"
	"time"

	scanner "github.com/Eiyarrr/Warren/scanner"
)

func runScan(args []string) error {
	scanFlags := flag.NewFlagSet("scan", flag.ContinueOnError)

	host := scanFlags.String("host", "127.0.0.1", "Host to scan")
	startPort := scanFlags.Int("start-port", 1, "First port to scan")
	endPort := scanFlags.Int("end-port", 1024, "Last port to scan")
	workers := scanFlags.Int("workers", 50, "Number of concurrent workers")
	timeout := scanFlags.Duration("timeout", 3*time.Second, "Maximum time to wait for a connection attempt")
	common := scanFlags.Bool("common", false, "Scan only common ports")
	openOnly := scanFlags.Bool("open-only", false, "Display only open ports")

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

	results := []scanner.Result{}
	if *common {
		results, err = scanner.ScanCommon(*host, *workers, *timeout)
	} else {
		results, err = scanner.ScanRange(*startPort, *endPort, *host, *workers, *timeout)
		if err != nil {
			return err
		}
	}

	printHeader(*host, *workers, *timeout)
	printResults(results, *openOnly)

	return nil
}

func verifyScanFlags(startPort int, endPort int, workers int, timeout time.Duration) error {
	if startPort > endPort || startPort < 0 || endPort > 65535 || startPort > 65535 || endPort < 0 {
		return scanner.ErrInvalidPortRange
	} else if workers <= 0 {
		return scanner.ErrInvalidWorkerCount
	} else if timeout <= 0 {
		return scanner.ErrInvalidTimeout
	}

	return nil
}
