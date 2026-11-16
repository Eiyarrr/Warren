package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	grabber "github.com/Eiyarrr/Warren/grabber"
	scanner "github.com/Eiyarrr/Warren/scanner"
)

func printResults(results []scanner.Result, openOnly bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)

	fmt.Println()
	fmt.Fprintln(w, "HOST\tPORT\tSTATUS\tERROR\tBANNER")

	for _, result := range results {
		if openOnly && result.Status != scanner.Open {
			continue
		}

		errMsg := "-"
		if result.Err != nil {
			errMsg = result.Err.Error()
		}

		banner := "=="
		banner, _, err := grabber.GrabBanner(result.Host, result.Port, 3*time.Second, 4000)
		if err != nil {
			fmt.Println(err)
		}

		coloredStatus := getStatusColor(result.Status)(result.Status.String())
		fmt.Fprintf(w, "%s\t%d\t%s\t%s%s\t%s\n", result.Host, result.Port, coloredStatus, "         ", errMsg, banner)
	}

	w.Flush()
}

func printHeader(target string, workers int, timeout time.Duration) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "TARGET \t", target)
	fmt.Fprintln(w, "WORKERS \t", workers)
	fmt.Fprintln(w, "TIMEOUT \t", timeout)

	w.Flush()
}
