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
	err := cli.Run(os.Args[1:], version, logFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
