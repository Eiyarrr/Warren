// Warren is a network service mapper.
package main

import (
	"fmt"
	"os"

	"github.com/Eiyarrr/Warren/cli"
)

const logFile = "warren_log"
const version = "v0.1.0 - DEV"

func main() {
	if len(os.Args) == 1 {
		err := cli.RunInteractive(version, logFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	switch os.Args[1] {
	case "scan":
		err := cli.RunScan(os.Args[2:], logFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return

	case "help":
		fmt.Println(cli.HelpMessage)
		return

	case "version":
		fmt.Println(version)
		return

	default:
		fmt.Println("Error, unknown command")
		os.Exit(1)
	}
}
