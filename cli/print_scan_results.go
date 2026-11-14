package cli

import (
	"fmt"
	"time"

	color "github.com/fatih/color"

	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

func printResults(results []portscanner.Result, openOnly bool) {
	colorOpen := color.New(color.FgGreen).SprintFunc()
	colorTimeout := color.New(color.FgYellow).SprintFunc()
	colorClosed := color.New(color.FgBlue).SprintFunc()
	colorErr := color.New(color.FgRed).SprintFunc()

	statusColor := color.New().SprintFunc()
	for _, result := range results {
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

		message := statusColor("Host: \"%s\" | Port: %s | Status: %s | Error: %v\n", result.Host, result.Port, result.Status.String(), result.Err)
		fmt.Println(message)
	}

	/*
		if openOnly {
			for _, result := range results {
				port := strconv.Itoa(result.Port)
				commonPort, ok := portscanner.CommonPorts[result.Port]
				if ok {
					port = fmt.Sprint(port, " / ", commonPort)
				}

				if result.Status == portscanner.Open {
					fmt.Printf("Host: \"%s\" | Port: %s | Status: %s | Error: %v\n", result.Host, port, result.Status.String(), result.Err)
				}
			}
		} else {
			for _, result := range results {
				port := strconv.Itoa(result.Port)
				commonPort, ok := portscanner.CommonPorts[result.Port]
				if ok {
					port = fmt.Sprint(port, " / ", commonPort)
				}
				fmt.Printf("Host: \"%s\" | Port: %s | Status: %s | Error: %v\n", result.Host, port, result.Status.String(), result.Err)
			}
		}
	*/
}

func printHeader(target string, workers int, timeout time.Duration) {
	fmt.Println(target)
	fmt.Println(workers)
	fmt.Println(timeout)
}
