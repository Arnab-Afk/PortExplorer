# Go Port Scanner with Process Details
=====================================

This is a simple port scanner written in Go that scans all 65535 ports on localhost and provides information about open ports and the processes using them. It uses a concurrency level of 1000 goroutines to speed up the scanning process.

## Usage

To run the port scanner, simply execute the `main.go` file:

```sh
go run main.go

The scanner will print the status of each port as it is scanned, along with the process ID (PID) and details of the process using the port (if applicable).

Functions
The following functions are used in the port scanner:

main(): The main function that runs the scanner in an infinite loop with a 15-second pause between scans.
scanPorts(ports <-chan int, wg *sync.WaitGroup): A goroutine that scans ports concurrently. It takes a channel of ports and a wait group as arguments.
getPID(port int) string: A function that returns the PID of the process using the given port. It uses the netstat command to find the PID.
getProcessDetails(pid string) string: A function that returns the details of the process with the given PID. It uses the Get-WmiObject PowerShell command to retrieve the details.
Configuration
The following constants can be found at the beginning of the main.go file:

concurrencyLevel: The number of concurrent goroutines used for scanning ports.
dialTimeout: The timeout for each connection attempt.
pauseDuration: The time to wait between scans.
totalPorts: The total number of ports to scan.
Feel free to modify these constants to adjust the scanner's behavior.

Dependencies
This port scanner relies on the following external packages:

net: For handling network connections.
os/exec: For executing external commands.
sync: For synchronizing goroutines.
time: For handling time-related operations.
These packages are part of the standard Go library, so no additional installation is required.

License
This port scanner is released under the MIT License. See the LICENSE file for more information.
