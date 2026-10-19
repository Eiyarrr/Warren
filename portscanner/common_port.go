package portscanner

// Ports that may be revealed from commonly used services.
var commonPorts = map[int]string{
    21:   "FTP",
    22:   "SSH",
    25:   "SMTP",
    53:   "DNS",
    80:   "HTTP",
    110:  "POP3",
    143:  "IMAP",
    443:  "HTTPS",
    445:  "SMB",
    3306: "MySQL",
    5432: "PostgreSQL",
    8080: "Alternative HTTP",
}
