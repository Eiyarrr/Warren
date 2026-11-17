package cli

import (
	scanner "github.com/Eiyarrr/Warren/scanner"
)

type Output struct {
	Host   string
	Port   int
	Status scanner.Status
	Error  error
	Banner string
}

func (op *Output) FillFromScan(result scanner.Result) {
	op.Host = result.Host
	op.Port = result.Port
	op.Status = result.Status
	op.Error = result.Error
}
