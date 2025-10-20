package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// PortScanResult represents the result of scanning a single port
type PortScanResult struct {
	Port    int
	State   string
	Service string
}

// ScanPort attempts to connect to a specific port on the target host
func ScanPort(protocol, hostname string, port int, timeout time.Duration) PortScanResult {
	result := PortScanResult{Port: port, State: "Closed"}
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, timeout)

	if err != nil {
		return result
	}

	defer conn.Close()
	result.State = "Open"
	result.Service = getServiceName(port)
	return result
}

// getServiceName returns common service names for well-known ports
func getServiceName(port int) string {
	services := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		445:  "SMB",
		3306: "MySQL",
		3389: "RDP",
		5432: "PostgreSQL",
		5900: "VNC",
		6379: "Redis",
		8080: "HTTP-Alt",
		8443: "HTTPS-Alt",
		27017: "MongoDB",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "Unknown"
}

// ScanPorts performs concurrent port scanning on the target host
func ScanPorts(hostname string, startPort, endPort int, timeout time.Duration, workers int) []PortScanResult {
	var results []PortScanResult
	var mutex sync.Mutex
	var wg sync.WaitGroup

	ports := make(chan int, workers)

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				result := ScanPort("tcp", hostname, port, timeout)
				if result.State == "Open" {
					mutex.Lock()
					results = append(results, result)
					mutex.Unlock()
					fmt.Printf("[+] Port %d is OPEN (%s)\n", port, result.Service)
				}
			}
		}()
	}

	// Send ports to workers
	for port := startPort; port <= endPort; port++ {
		ports <- port
	}
	close(ports)

	wg.Wait()

	// Sort results by port number
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

// PrintBanner displays the program banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║     Go Port Scanner v1.0              ║
║     Cybersecurity Lab Tool            ║
╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintUsage displays usage information
func PrintUsage() {
	fmt.Println("Usage: go run main.go <hostname> <start_port> <end_port>")
	fmt.Println("\nExample:")
	fmt.Println("  go run main.go localhost 1 1000")
	fmt.Println("  go run main.go scanme.nmap.org 20 80")
	fmt.Println("\nNOTE: Only scan systems you have permission to test!")
}

// ValidateInput validates command line arguments
func ValidateInput(args []string) (string, int, int, error) {
	if len(args) != 4 {
		return "", 0, 0, fmt.Errorf("invalid number of arguments")
	}

	hostname := args[1]
	startPort, err := strconv.Atoi(args[2])
	if err != nil || startPort < 1 || startPort > 65535 {
		return "", 0, 0, fmt.Errorf("invalid start port")
	}

	endPort, err := strconv.Atoi(args[3])
	if err != nil || endPort < 1 || endPort > 65535 {
		return "", 0, 0, fmt.Errorf("invalid end port")
	}

	if startPort > endPort {
		return "", 0, 0, fmt.Errorf("start port must be less than or equal to end port")
	}

	return hostname, startPort, endPort, nil
}

// PrintSummary displays the scan summary
func PrintSummary(hostname string, startPort, endPort int, results []PortScanResult, duration time.Duration) {
	fmt.Println("\n" + "═"*50)
	fmt.Printf("Scan Summary for %s\n", hostname)
	fmt.Println("═"*50)
	fmt.Printf("Port Range: %d-%d\n", startPort, endPort)
	fmt.Printf("Total Ports Scanned: %d\n", endPort-startPort+1)
	fmt.Printf("Open Ports Found: %d\n", len(results))
	fmt.Printf("Scan Duration: %.2f seconds\n", duration.Seconds())
	
	if len(results) > 0 {
		fmt.Println("\nOpen Ports Details:")
		fmt.Println("─"*50)
		fmt.Printf("%-10s %-10s %-20s\n", "PORT", "STATE", "SERVICE")
		fmt.Println("─"*50)
		for _, result := range results {
			fmt.Printf("%-10d %-10s %-20s\n", result.Port, result.State, result.Service)
		}
	}
	fmt.Println("═"*50)
}

func main() {
	PrintBanner()

	// Validate input
	hostname, startPort, endPort, err := ValidateInput(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
		PrintUsage()
		os.Exit(1)
	}

	// Configuration
	timeout := 1 * time.Second
	workers := 100 // Number of concurrent workers

	// Disclaimer
	fmt.Println("⚠️  WARNING: Only scan systems you have explicit permission to test!")
	fmt.Printf("\nStarting port scan on %s (Ports %d-%d)\n", hostname, startPort, endPort)
	fmt.Println("This may take a few moments...\n")

	// Perform scan
	startTime := time.Now()
	results := ScanPorts(hostname, startPort, endPort, timeout, workers)
	duration := time.Since(startTime)

	// Print summary
	PrintSummary(hostname, startPort, endPort, results, duration)

	// Exit with appropriate status
	if len(results) > 0 {
		os.Exit(0)
	} else {
		fmt.Println("\nNo open ports found.")
		os.Exit(0)
	}
}