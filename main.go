package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ports := make(chan int, 100)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			scanPorts(ports, &wg)
		}()
	}

	go func() {
		for i := 1; i <= 65535; i++ {
			ports <- i
		}
		close(ports)
	}()

	wg.Wait()
}

func scanPorts(ports <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			continue
		}
		defer conn.Close()
		fmt.Printf("Port %d is open\n", port)
	}
}
