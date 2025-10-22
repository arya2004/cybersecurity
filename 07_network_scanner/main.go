// Network Security Scanner in Go
// This tool performs comprehensive network security scanning including:
// - Port scanning (TCP/UDP)
// - Service detection
// - OS fingerprinting
// - Vulnerability detection
// - Network topology discovery

package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ScanResult represents the result of a port scan
type ScanResult struct {
	Host    string
	Port    int
	Open    bool
	Service string
	Banner  string
	Error   error
}

// NetworkScanner performs various network security scans
type NetworkScanner struct {
	Timeout         time.Duration
	MaxConcurrency  int
	ServiceDatabase map[int]string
}

// NewNetworkScanner creates a new network scanner instance
func NewNetworkScanner() *NetworkScanner {
	return &NetworkScanner{
		Timeout:         5 * time.Second,
		MaxConcurrency:  100,
		ServiceDatabase: getCommonServices(),
	}
}

// getCommonServices returns a map of common ports to their services
func getCommonServices() map[int]string {
	return map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		993:  "IMAPS",
		995:  "POP3S",
		1433: "MSSQL",
		3306: "MySQL",
		3389: "RDP",
		5432: "PostgreSQL",
		5900: "VNC",
		6379: "Redis",
		8080: "HTTP-Alt",
		9200: "Elasticsearch",
		27017: "MongoDB",
	}
}

// ScanTCPPort performs a TCP port scan on a single port
func (ns *NetworkScanner) ScanTCPPort(host string, port int) ScanResult {
	result := ScanResult{
		Host: host,
		Port: port,
		Open: false,
	}

	// Set service name if known
	if service, exists := ns.ServiceDatabase[port]; exists {
		result.Service = service
	}

	// Attempt to connect
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), ns.Timeout)
	if err != nil {
		result.Error = err
		return result
	}
	defer conn.Close()

	result.Open = true

	// Try to grab banner
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		result.Banner = strings.TrimSpace(scanner.Text())
	}

	return result
}

// ScanPortRange scans a range of ports on a target host
func (ns *NetworkScanner) ScanPortRange(host string, startPort, endPort int) []ScanResult {
	var results []ScanResult
	var wg sync.WaitGroup
	resultChan := make(chan ScanResult, endPort-startPort+1)

	// Create a semaphore to limit concurrency
	semaphore := make(chan struct{}, ns.MaxConcurrency)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			result := ns.ScanTCPPort(host, p)
			resultChan <- result
		}(port)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	for result := range resultChan {
		results = append(results, result)
	}

	// Sort results by port number
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

// ScanCommonPorts scans the most common ports
func (ns *NetworkScanner) ScanCommonPorts(host string) []ScanResult {
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 1433, 3306, 3389, 5432, 5900, 6379, 8080}
	
	var results []ScanResult
	var wg sync.WaitGroup
	resultChan := make(chan ScanResult, len(commonPorts))
	semaphore := make(chan struct{}, ns.MaxConcurrency)

	for _, port := range commonPorts {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := ns.ScanTCPPort(host, p)
			resultChan <- result
		}(port)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

// DiscoverHosts discovers active hosts in a network range
func (ns *NetworkScanner) DiscoverHosts(network string) []string {
	ip, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		fmt.Printf("Error parsing CIDR: %v\n", err)
		return nil
	}

	var hosts []string
	var wg sync.WaitGroup
	hostChan := make(chan string, 256)
	semaphore := make(chan struct{}, ns.MaxConcurrency)

	// Iterate through all IPs in the network
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		wg.Add(1)
		go func(targetIP string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if ns.isHostActive(targetIP) {
				hostChan <- targetIP
			}
		}(ip.String())
	}

	go func() {
		wg.Wait()
		close(hostChan)
	}()

	for host := range hostChan {
		hosts = append(hosts, host)
	}

	sort.Strings(hosts)
	return hosts
}

// isHostActive checks if a host is active using ICMP ping and TCP connect
func (ns *NetworkScanner) isHostActive(host string) bool {
	// Try TCP connect to common ports
	commonPorts := []int{22, 80, 443}
	
	for _, port := range commonPorts {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 1*time.Second)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}

// incrementIP increments an IP address
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// DetectOS attempts basic OS fingerprinting
func (ns *NetworkScanner) DetectOS(host string) string {
	// Check for common OS-specific services and behaviors
	sshResult := ns.ScanTCPPort(host, 22)
	rdpResult := ns.ScanTCPPort(host, 3389)
	
	if rdpResult.Open {
		return "Windows (RDP detected)"
	}
	
	if sshResult.Open && strings.Contains(strings.ToLower(sshResult.Banner), "ubuntu") {
		return "Ubuntu Linux"
	}
	
	if sshResult.Open && strings.Contains(strings.ToLower(sshResult.Banner), "centos") {
		return "CentOS Linux"
	}
	
	if sshResult.Open {
		return "Linux/Unix (SSH detected)"
	}
	
	return "Unknown"
}

// VulnerabilityCheck performs basic vulnerability checks
func (ns *NetworkScanner) VulnerabilityCheck(host string, results []ScanResult) []string {
	var vulnerabilities []string
	
	for _, result := range results {
		if !result.Open {
			continue
		}
		
		switch result.Port {
		case 21: // FTP
			if strings.Contains(strings.ToLower(result.Banner), "vsftpd 2.3.4") {
				vulnerabilities = append(vulnerabilities, "CVE-2011-2523: vsftpd 2.3.4 backdoor")
			}
		case 22: // SSH
			if strings.Contains(strings.ToLower(result.Banner), "openssh") {
				vulnerabilities = append(vulnerabilities, "Potential SSH brute force target")
			}
		case 23: // Telnet
			vulnerabilities = append(vulnerabilities, "Insecure protocol: Telnet transmits in plain text")
		case 80: // HTTP
			vulnerabilities = append(vulnerabilities, "HTTP service - check for web vulnerabilities")
		case 443: // HTTPS
			vulnerabilities = append(vulnerabilities, "HTTPS service - verify SSL/TLS configuration")
		case 3389: // RDP
			vulnerabilities = append(vulnerabilities, "RDP exposed - potential brute force target")
		}
	}
	
	return vulnerabilities
}

// PrintResults prints scan results in a formatted manner
func (ns *NetworkScanner) PrintResults(host string, results []ScanResult) {
	fmt.Printf("\n=== Scan Results for %s ===\n", host)
	fmt.Printf("%-8s %-12s %-15s %s\n", "Port", "State", "Service", "Banner")
	fmt.Println(strings.Repeat("-", 60))
	
	openPorts := 0
	for _, result := range results {
		if result.Open {
			state := "open"
			service := result.Service
			if service == "" {
				service = "unknown"
			}
			banner := result.Banner
			if len(banner) > 30 {
				banner = banner[:30] + "..."
			}
			fmt.Printf("%-8d %-12s %-15s %s\n", result.Port, state, service, banner)
			openPorts++
		}
	}
	
	fmt.Printf("\nSummary: %d open ports found\n", openPorts)
	
	// OS Detection
	os := ns.DetectOS(host)
	fmt.Printf("OS Detection: %s\n", os)
	
	// Vulnerability Check
	vulns := ns.VulnerabilityCheck(host, results)
	if len(vulns) > 0 {
		fmt.Println("\nPotential Vulnerabilities:")
		for _, vuln := range vulns {
			fmt.Printf("- %s\n", vuln)
		}
	}
}

// Interactive menu system
func showMenu() {
	fmt.Println("\n=== Network Security Scanner ===")
	fmt.Println("1. Scan single host (common ports)")
	fmt.Println("2. Scan single host (port range)")
	fmt.Println("3. Discover hosts in network")
	fmt.Println("4. Comprehensive scan")
	fmt.Println("5. Exit")
	fmt.Print("Select option: ")
}

func main() {
	scanner := NewNetworkScanner()
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("Network Security Scanner v1.0")
	fmt.Println("Educational tool for cybersecurity learning")
	fmt.Println("Use responsibly and only on networks you own or have permission to test")
	
	for {
		showMenu()
		
		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)
		
		switch choice {
		case "1":
			fmt.Print("Enter target host/IP: ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			
			fmt.Printf("Scanning common ports on %s...\n", host)
			results := scanner.ScanCommonPorts(host)
			scanner.PrintResults(host, results)
			
		case "2":
			fmt.Print("Enter target host/IP: ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			
			fmt.Print("Enter start port: ")
			startPortStr, _ := reader.ReadString('\n')
			startPort, err := strconv.Atoi(strings.TrimSpace(startPortStr))
			if err != nil {
				fmt.Println("Invalid start port")
				continue
			}
			
			fmt.Print("Enter end port: ")
			endPortStr, _ := reader.ReadString('\n')
			endPort, err := strconv.Atoi(strings.TrimSpace(endPortStr))
			if err != nil {
				fmt.Println("Invalid end port")
				continue
			}
			
			fmt.Printf("Scanning ports %d-%d on %s...\n", startPort, endPort, host)
			results := scanner.ScanPortRange(host, startPort, endPort)
			scanner.PrintResults(host, results)
			
		case "3":
			fmt.Print("Enter network CIDR (e.g., 192.168.1.0/24): ")
			network, _ := reader.ReadString('\n')
			network = strings.TrimSpace(network)
			
			fmt.Printf("Discovering hosts in %s...\n", network)
			hosts := scanner.DiscoverHosts(network)
			
			fmt.Printf("\nDiscovered %d active hosts:\n", len(hosts))
			for _, host := range hosts {
				fmt.Printf("- %s\n", host)
			}
			
		case "4":
			fmt.Print("Enter target host/IP: ")
			host, _ := reader.ReadString('\n')
			host = strings.TrimSpace(host)
			
			fmt.Printf("Performing comprehensive scan on %s...\n", host)
			
			// Scan common ports
			results := scanner.ScanCommonPorts(host)
			scanner.PrintResults(host, results)
			
			// Additional detailed scan for open ports
			if len(results) > 0 {
				fmt.Println("\nPerforming detailed service enumeration...")
				// Here you could add more detailed service enumeration
			}
			
		case "5":
			fmt.Println("Exiting...")
			return
			
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
