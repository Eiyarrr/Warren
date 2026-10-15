package cli

import (
	"errors"
	"fmt"
)

func Run(args []string, version string, logFile string) error {
	if len(args) == 0 {
		return runInteractive(version, logFile)
	}

	switch args[0] {
	case "scan":
		err := runScan(args[1:], logFile)
		if err != nil {
			return err
		}
		return nil

	case "help":
		fmt.Println(HelpMessage)
		return nil

	case "version":
		fmt.Println(version)
		return nil

	default:
		return errors.New("unknown command")
	}
}
