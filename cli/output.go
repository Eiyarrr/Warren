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

	arrOut := make([]Output, len(results))
	for i, result := range results {
		arrOut[i].FillFromScan(result)
	}

	for _, out := range arrOut {
		if openOnly && out.Status != scanner.Open {
			continue
		}

		var banner string
		if out.Status == scanner.Open {
			var err error
			banner, _, err = grabber.GrabBanner(out.Host, out.Port, 3*time.Second, 4000)
			if err != nil {
				fmt.Println(err)
			}
		}

		errMsg := "-"
		if out.Error != nil {
			errMsg = out.Error.Error()
		}

		coloredStatus := getStatusColor(out.Status)(out.Status.String())
		fmt.Fprintf(w, "%s\t%d\t%s\t%s%s\t%s\n", out.Host, out.Port, coloredStatus, "         ", errMsg, banner)
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
