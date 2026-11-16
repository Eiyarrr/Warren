package cli_test

import (
	"errors"
	"testing"

	cli "github.com/Eiyarrr/Warren/cli"
	scanner "github.com/Eiyarrr/Warren/scanner"
)

func TestRunScan(t *testing.T) {
	args := []string{"scan", "-start-port", "80", "-end-port", "80"}

	err := cli.Run(args, "")
	if err != nil {
		t.Fatalf("failed to scan: %v", err)
	}
}

func TestInvalidPorts(t *testing.T) {
	args := []string{"scan", "-start-port", "100", "-end-port", "90"}

	err := cli.Run(args, "")
	if !errors.Is(err, scanner.ErrInvalidPortRange) {
		t.Fatalf("expected %v, got %v", scanner.ErrInvalidPortRange, err)
	}
}

func TestInvalidWorkerCount(t *testing.T) {
	args := []string{"scan", "-workers", "-10"}

	err := cli.Run(args, "")
	if !errors.Is(err, scanner.ErrInvalidWorkerCount) {
		t.Fatalf("expected %v, got %v", scanner.ErrInvalidWorkerCount, err)
	}
}

func TestInvalidTimeout(t *testing.T) {
	args := []string{"scan", "-timeout", "-1s"}

	err := cli.Run(args, "")
	if !errors.Is(err, scanner.ErrInvalidTimeout) {
		t.Fatalf("expected %v, got %v", scanner.ErrInvalidTimeout, err)
	}
}

func TestUnknownCommand(t *testing.T) {
	args := []string{"Unknown Command"}

	err := cli.Run(args, "")
	if err == nil {
		t.Fatalf("expected error for unknown command")
	}
}
