// Warren is a network service mapper.
package main

import (
	"fmt"
	"os"

	"github.com/Eiyarrr/Warren/cli"
)

const version = "v0.1.0 - DEV"

func main() {
	err := cli.Run(os.Args[1:], version)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
