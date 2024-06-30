package main

import (
	"fmt"
	"net"
	"time"
)

// ScanPort checks if a port is open on a given host
func ScanPort(protocol, hostname string, port int) bool {
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func main() {
	hostname := "scanme.nmap.org"
	startPort := 20
	endPort := 1024

	fmt.Printf("Scanning ports %d-%d on %s\n", startPort, endPort, hostname)
	for port := startPort; port <= endPort; port++ {
		if ScanPort("tcp", hostname, port) {
			fmt.Printf("Port %d is open\n", port)
		}
	}
}
