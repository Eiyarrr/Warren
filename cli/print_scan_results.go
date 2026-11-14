package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	color "github.com/fatih/color"

	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

func printResults(results []portscanner.Result, openOnly bool) {
	colorOpen := color.New(color.FgGreen).SprintFunc()
	colorTimeout := color.New(color.FgYellow).SprintFunc()
	colorClosed := color.New(color.FgHiBlack).SprintFunc()
	colorErr := color.New(color.FgRed).SprintFunc()

	statusColor := color.New().SprintFunc()
	for _, result := range results {
		if openOnly && !(result.Status == portscanner.Open) {
			continue
		}

		switch result.Status {
		case portscanner.Open:
			statusColor = colorOpen
		case portscanner.Timeout:
			statusColor = colorTimeout
		case portscanner.Closed:
			statusColor = colorClosed
		case portscanner.Err:
			statusColor = colorErr
		default:
			statusColor = color.New(color.FgWhite).SprintFunc()
		}

		message := fmt.Sprintf("Host: %q | Port: %d | Status: %s | Error: %v", result.Host, result.Port, statusColor(result.Status.String()), result.Err)
		fmt.Println(message)
	}
}

func printHeader(target string, workers int, timeout time.Duration) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

	fmt.Fprintln(w, "TARGET \t", target)
	fmt.Fprintln(w, "WORKERS \t", workers)
	fmt.Fprintln(w, "TIMEOUT \t", timeout)

	w.Flush()
}
