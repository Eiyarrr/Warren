## Warren
Warren is a network service mapper in progress. I created it because I am interested in networks and how information can be unintentionally leaked over the internet. This is still a work in progress and I intend to add more features over time.

*Warren is under the MIT License, for more information, see [LICENSE](LICENSE)*

## Features
- **Portscanning**: ScanRange scans a range of ports and returns ([]Result, error). Scan attempts to connect to a single port and sends its result through a channel. Connection errors are reported through the Result.Err field.
- **Error Handling**: Warren provides exported sentinel error values for specific port-scanning configuration failures.

## Testing
To run the test suite for this project, clone the repository, enter its directory, and run `go test ./...`
