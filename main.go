package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	concurrencyLevel = 1000                   // Number of concurrent goroutines
	dialTimeout      = 500 * time.Millisecond // Timeout for each connection attempt
	pauseDuration    = 15 * time.Second       // Time to wait between scans
	totalPorts       = 65535                  // Total number of ports to scan
)

func main() {
	for {
		start := time.Now()

		var wg sync.WaitGroup
		ports := make(chan int, 1000)

		// Start goroutines to scan ports concurrently
		for i := 0; i < concurrencyLevel; i++ {
			wg.Add(1)
			go func() {
				scanPorts(ports, &wg)
			}()
		}

		// Send port numbers to the channel
		go func() {
			for i := 1; i <= totalPorts; i++ {
				ports <- i
			}
			close(ports)
		}()

		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("Scanned %d ports in %s\n", totalPorts, elapsed)

		fmt.Println("Scan complete. Waiting 10 seconds before next scan...")
		time.Sleep(pauseDuration)
	}
}

func scanPorts(ports <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		address := fmt.Sprintf("localhost:%d", port)
		conn, err := net.DialTimeout("tcp", address, dialTimeout)
		if err == nil {
			conn.Close()
			fmt.Printf("Port %d is open\n", port)

			// Identify the process running on this port
			process, err := getProcessRunningOnPort(port)
			if err != nil {
				fmt.Printf("Error getting process for port %d: %v\n", port, err)
			} else {
				fmt.Printf("Process running on port %d: %s\n", port, process)
			}
		}
	}
}

func getProcessRunningOnPort(port int) (string, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("netstat -ano | Select-String \":%d\" | ForEach-Object { Get-Process -PID ($_.Split()[4]) }", port))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Parse the output to extract the process name
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ProcessName") {
			fields := strings.Fields(line)
			processName := fields[1]
			return processName, nil
		}
	}

	return "", fmt.Errorf("no process found running on port %d", port)
}

func getProcessName(pid string) (string, error) {
	// Use the `ps` command to get the process name from the process ID
	cmd := exec.Command("ps", "-p", pid, "-o", "comm=")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
