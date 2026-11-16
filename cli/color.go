package cli

import (
	color "github.com/fatih/color"

	portscanner "github.com/Eiyarrr/Warren/portscanner"
)

var (
	colorOpen    = color.New(color.FgGreen).SprintFunc()
	colorTimeout = color.New(color.FgYellow).SprintFunc()
	colorClosed  = color.New(color.FgHiBlack).SprintFunc()
	colorErr     = color.New(color.FgRed).SprintFunc()
	colorDefault = color.New(color.FgWhite).SprintFunc()
)

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
