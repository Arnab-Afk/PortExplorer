package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	for {
		var wg sync.WaitGroup
		ports := make(chan int, 100)

		// Start 1000 goroutines to scan ports concurrently
		for i := 0; i < 1000; i++ {
			wg.Add(1)
			go func() {
				scanPorts(ports, &wg)
			}()
		}

		// Send port numbers to the channel
		go func() {
			for i := 1; i <= 65535; i++ {
				ports <- i
			}
			close(ports)
		}()

		wg.Wait()

		fmt.Println("Scan complete. Waiting 10 seconds before next scan...")
		time.Sleep(10 * time.Second)
	}
}

func scanPorts(ports <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err == nil {
			conn.Close()
			fmt.Printf("Port %d is open\n", port)
		}
	}
}
