package cli

import (
	color "github.com/fatih/color"

	scanner "github.com/Eiyarrr/Warren/scanner"
)

var (
	colorOpen    = color.New(color.FgGreen).SprintFunc()
	colorTimeout = color.New(color.FgYellow).SprintFunc()
	colorClosed  = color.New(color.FgHiBlack).SprintFunc()
	colorErr     = color.New(color.FgRed).SprintFunc()
	colorDefault = color.New(color.FgWhite).SprintFunc()
)

func getStatusColor(status scanner.Status) func(a ...any) string {
	switch status {
	case scanner.Open:
		return colorOpen
	case scanner.Timeout:
		return colorTimeout
	case scanner.Closed:
		return colorClosed
	case scanner.Err:
		return colorErr
	default:
		return colorDefault
	}
}
