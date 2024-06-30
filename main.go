package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	concurrencyLevel = 1000                   // Number of concurrent goroutines
	dialTimeout      = 500 * time.Millisecond // Timeout for each connection attempt
	totalPorts       = 65535                  // Total number of ports to scan
)

var (
	results      = make(map[int]string)
	resultsMutex sync.Mutex
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/api/results", resultsHandler)
	http.HandleFunc("/scan", scanHandler)

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Port Scanner</title>
			<link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" rel="stylesheet">
		</head>
		<body>
			<div class="container">
				<h1 class="mt-5">Port Scanner</h1>
<form action="/scan" method="post">
    <button type="submit" class="btn btn-primary mt-3">Start Scan</button>
</form>
<div class="mt-5">
    <h2>Scan Results</h2>
    <table class="table">
        <thead>
            <tr>
                <th>Port</th>
                <th>Status</th>
                <th>Process Details</th>
            </tr>
        </thead>
        <tbody id="results">
        </tbody>
    </table>
</div>
<script>
    function fetchResults() {
        fetch('/api/results')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                const resultsElement = document.getElementById('results');
                resultsElement.innerHTML = '';
                data.forEach(result => {
                    const row = document.createElement('tr');
                    row.innerHTML = `<td>${result.port}</td><td>${result.status}</td><td>${result.details}</td>`;
                    resultsElement.appendChild(row);
                });
            })
            .catch(error => console.error('There has been a problem with your fetch operation:', error));
    }
    setInterval(fetchResults, 5000); // Fetch results every 5 seconds
</script>
		</body>
		</html>
	`)
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		go startScan()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	resultsMutex.Lock()
	defer resultsMutex.Unlock()

	type Result struct {
		Port    int    `json:"port"`
		Status  string `json:"status"`
		Details string `json:"details"`
	}
	var resultList []Result
	for port, details := range results {
		resultList = append(resultList, Result{Port: port, Status: "Open", Details: details})
	}
	jsonResponse(w, resultList)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func startScan() {
	resultsMutex.Lock()
	results = make(map[int]string)
	resultsMutex.Unlock()

	var wg sync.WaitGroup
	ports := make(chan int, concurrencyLevel)

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
}

func scanPorts(ports <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		address := fmt.Sprintf("localhost:%d", port)
		conn, err := net.DialTimeout("tcp", address, dialTimeout)
		if err == nil {
			conn.Close()
			pid := getPID(port)
			var details string
			if pid != "" {
				details = getProcessDetails(pid)
			}
			resultsMutex.Lock()
			results[port] = details
			resultsMutex.Unlock()
		}
	}
}
