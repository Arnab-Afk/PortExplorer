package main

import (
	"fmt"
	"net"
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
		}
	}
}
