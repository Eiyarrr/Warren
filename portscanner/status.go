package portscanner

type Status int

const (
	Open    Status = iota //0
	Closed                //1
	Timeout               //2
	Err                   //3
)

func (s Status) String() string {
	switch s {
	case Open:
		return "Open"
	case Timeout:
		return "Timeout"
	case Closed:
		return "Closed"
	case Err:
		return "Err"
	default:
		return ""
	}
}
