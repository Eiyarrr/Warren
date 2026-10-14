package portscanner

type Result struct {
	host   string
	port   int
	status Status
	err    error
}
