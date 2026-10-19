## Warren
Warren is a network service mapper in progress. I created it because I am interested in networks and how information can be unintentionally leaked over the internet. This is still a work in progress and I intend to add more features over time.

*Warren is under the MIT License, for more information, see [LICENSE](LICENSE)*

## Features
- **Port scanning**: ScanRange scans a range of ports and returns `([]Result, error)`. Scan attempts to connect to a single port and sends its `Result` through a channel. ScanCommon scans common ports (as seen in `portscanner/common_ports.go`) and returns `([]Result, error)` Connection errors are reported through the `Result.Err` field.
- **Configurable Scanning**: The caller of any scan function specifies the `host`, `port range` (or `common ports`), `worker count`, and `connection timeout`.
- **Result Filtering**: Using the flag `--open-only` displays only ports that are marked as `Open`.
- **Error Handling**: Warren provides exported sentinel error values for specific port-scanning configuration failures.
- **REPL-style CLI**: Warren can run individual commands directly or launch an interactive REPL for executing multiple commands in a single session.

## Usage
Warren is a work in progress, but if you wish to use it, clone the repository, enter its directory, and run `go run .`.
If you wish to see all allowed parameters, run `go run . --help`.

## Testing
To run the test suite for this project, clone the repository, enter its directory, and run `go test ./...`
