package portscanner

type Result struct {
	Host   string
	Port   int
	Status Status
	Err    error
}
