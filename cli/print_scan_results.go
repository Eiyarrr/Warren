package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	color "github.com/fatih/color"

	"github.com/Eiyarrr/Warren/bannergrabber"
	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

var (
	colorOpen    = color.New(color.FgGreen).SprintFunc()
	colorTimeout = color.New(color.FgYellow).SprintFunc()
	colorClosed  = color.New(color.FgHiBlack).SprintFunc()
	colorErr     = color.New(color.FgRed).SprintFunc()
	colorDefault = color.New(color.FgWhite).SprintFunc()
)

func printResults(results []portscanner.Result, openOnly bool) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)

	fmt.Println()
	fmt.Fprintln(w, "HOST\tPORT\tSTATUS\tERROR")

	for _, result := range results {
		if openOnly && result.Status != portscanner.Open {
			continue
		}

		errMsg := "-"
		if result.Err != nil {
			errMsg = result.Err.Error()
		}

		banner, _, err := bannergrabber.GrabBanner(result.Host, result.Port, 3*time.Second, 4000)
		if err != nil {
			fmt.Println(err)
		}

		coloredStatus := getStatusColor(result.Status)(result.Status.String())
		fmt.Fprintf(w, "%s\t%d\t%s\t%s%s\t%s\n", result.Host, result.Port, coloredStatus, "         ", errMsg, banner)
	}

	w.Flush()
}

func getStatusColor(status portscanner.Status) func(a ...any) string {
	switch status {
	case portscanner.Open:
		return colorOpen
	case portscanner.Timeout:
		return colorTimeout
	case portscanner.Closed:
		return colorClosed
	case portscanner.Err:
		return colorErr
	default:
		return colorDefault
	}
}

func printHeader(target string, workers int, timeout time.Duration) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "TARGET \t", target)
	fmt.Fprintln(w, "WORKERS \t", workers)
	fmt.Fprintln(w, "TIMEOUT \t", timeout)

	w.Flush()
}
