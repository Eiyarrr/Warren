package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func RunInteractive(version string, logFile string) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("warren> ")

		if !scanner.Scan() {
			err := scanner.Err()
			if err != nil {
				return err
			}
			return nil
		}

		input := strings.ToLower(strings.TrimSpace(scanner.Text()))

		switch input {
		case "":
			continue
		case "exit":
			return nil
		case "version":
			fmt.Println(version)
		case "help":
			fmt.Println(HelpMessage)
		default:
			fields := strings.Fields(input)
			if fields[0] == "scan" {
				err := RunScan(fields[1:], logFile)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			} else {
				fmt.Println("Unknown command.")
			}
		}
	}
}
