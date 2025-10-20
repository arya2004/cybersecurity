package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// Connection represents a network connection
type Connection struct {
	LocalAddr  string
	LocalPort  string
	RemoteAddr string
	RemotePort string
	Protocol   string
	State      string
	PID        string
	Process    string
}

// TrafficStats stores traffic statistics
type TrafficStats struct {
	TotalConnections   int
	ActiveConnections  int
	TCPConnections     int
	UDPConnections     int
	ListeningPorts     int
	EstablishedConns   int
	TimeWaitConns      int
	UniqueRemoteIPs    int
	ConnectionsByPort  map[string]int
	ConnectionsByIP    map[string]int
	ProtocolDistribution map[string]int
}

// PortInfo stores information about well-known ports
var wellKnownPorts = map[string]string{
	"20":    "FTP Data",
	"21":    "FTP Control",
	"22":    "SSH",
	"23":    "Telnet",
	"25":    "SMTP",
	"53":    "DNS",
	"80":    "HTTP",
	"110":   "POP3",
	"143":   "IMAP",
	"443":   "HTTPS",
	"445":   "SMB",
	"3306":  "MySQL",
	"3389":  "RDP",
	"5432":  "PostgreSQL",
	"5900":  "VNC",
	"6379":  "Redis",
	"8080":  "HTTP-Alt",
	"8443":  "HTTPS-Alt",
	"27017": "MongoDB",
}

// GetActiveConnections retrieves current network connections
func GetActiveConnections() ([]Connection, error) {
	connections := []Connection{}
	
	// Get TCP connections
	tcpConns, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error getting network interfaces: %v", err)
	}

	// Simulate network connections for demonstration
	// In a real implementation, this would parse /proc/net/tcp, netstat, or use system calls
	connections = append(connections, Connection{
		LocalAddr:  "127.0.0.1",
		LocalPort:  "8080",
		RemoteAddr: "0.0.0.0",
		RemotePort: "*",
		Protocol:   "TCP",
		State:      "LISTENING",
		PID:        "1234",
		Process:    "web-server",
	})

	connections = append(connections, Connection{
		LocalAddr:  "192.168.1.100",
		LocalPort:  "54321",
		RemoteAddr: "93.184.216.34",
		RemotePort: "443",
		Protocol:   "TCP",
		State:      "ESTABLISHED",
		PID:        "5678",
		Process:    "chrome",
	})

	connections = append(connections, Connection{
		LocalAddr:  "0.0.0.0",
		LocalPort:  "22",
		RemoteAddr: "0.0.0.0",
		RemotePort: "*",
		Protocol:   "TCP",
		State:      "LISTENING",
		PID:        "999",
		Process:    "sshd",
	})

	connections = append(connections, Connection{
		LocalAddr:  "192.168.1.100",
		LocalPort:  "12345",
		RemoteAddr: "142.250.185.46",
		RemotePort: "443",
		Protocol:   "TCP",
		State:      "ESTABLISHED",
		PID:        "5678",
		Process:    "chrome",
	})

	// Note: This is a simplified demo. Real implementation would:
	// - Parse /proc/net/tcp and /proc/net/udp on Linux
	// - Use netstat or ss command
	// - Use Windows API on Windows
	// - Use system calls for cross-platform support

	fmt.Println("📡 Note: This is a demonstration with sample data.")
	fmt.Println("   In production, this would read actual system network connections.")
	
	return connections, nil
}

// CalculateStatistics computes traffic statistics
func CalculateStatistics(connections []Connection) TrafficStats {
	stats := TrafficStats{
		ConnectionsByPort:    make(map[string]int),
		ConnectionsByIP:      make(map[string]int),
		ProtocolDistribution: make(map[string]int),
	}

	uniqueIPs := make(map[string]bool)

	for _, conn := range connections {
		stats.TotalConnections++

		// Protocol distribution
		stats.ProtocolDistribution[conn.Protocol]++

		if conn.Protocol == "TCP" {
			stats.TCPConnections++
		} else if conn.Protocol == "UDP" {
			stats.UDPConnections++
		}

		// State tracking
		switch conn.State {
		case "ESTABLISHED":
			stats.EstablishedConns++
			stats.ActiveConnections++
		case "LISTENING":
			stats.ListeningPorts++
		case "TIME_WAIT":
			stats.TimeWaitConns++
		}

		// Port tracking
		if conn.LocalPort != "*" {
			stats.ConnectionsByPort[conn.LocalPort]++
		}
		if conn.RemotePort != "*" {
			stats.ConnectionsByPort[conn.RemotePort]++
		}

		// IP tracking
		if conn.RemoteAddr != "0.0.0.0" && conn.RemoteAddr != "*" {
			stats.ConnectionsByIP[conn.RemoteAddr]++
			uniqueIPs[conn.RemoteAddr] = true
		}
	}

	stats.UniqueRemoteIPs = len(uniqueIPs)

	return stats
}

// PrintBanner displays program banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║   Network Traffic Analyzer v1.0      ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintStatistics displays traffic statistics
func PrintStatistics(stats TrafficStats) {
	fmt.Println("\n" + "═"*60)
	fmt.Println("NETWORK TRAFFIC STATISTICS")
	fmt.Println("═"*60)
	
	// Overall stats
	fmt.Println("Overall:")
	fmt.Printf("  Total Connections: %d\n", stats.TotalConnections)
	fmt.Printf("  Active Connections: %d\n", stats.ActiveConnections)
	fmt.Printf("  Listening Ports: %d\n", stats.ListeningPorts)
	fmt.Printf("  Unique Remote IPs: %d\n", stats.UniqueRemoteIPs)
	fmt.Println("─"*60)

	// Protocol distribution
	fmt.Println("Protocol Distribution:")
	for protocol, count := range stats.ProtocolDistribution {
		percentage := float64(count) / float64(stats.TotalConnections) * 100
		fmt.Printf("  %s: %d (%.1f%%)\n", protocol, count, percentage)
	}
	fmt.Println("─"*60)

	// Connection states
	fmt.Println("Connection States:")
	fmt.Printf("  ESTABLISHED: %d\n", stats.EstablishedConns)
	fmt.Printf("  LISTENING: %d\n", stats.ListeningPorts)
	fmt.Printf("  TIME_WAIT: %d\n", stats.TimeWaitConns)
	fmt.Println("─"*60)

	// Top ports
	fmt.Println("Top Active Ports:")
	type portCount struct {
		Port  string
		Count int
	}
	var ports []portCount
	for port, count := range stats.ConnectionsByPort {
		ports = append(ports, portCount{port, count})
	}
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Count > ports[j].Count
	})

	for i, pc := range ports {
		if i >= 10 {
			break
		}
		service := wellKnownPorts[pc.Port]
		if service == "" {
			service = "Unknown"
		}
		fmt.Printf("  Port %s (%s): %d connections\n", pc.Port, service, pc.Count)
	}
	fmt.Println("─"*60)

	// Top remote IPs
	fmt.Println("Top Remote IPs:")
	type ipCount struct {
		IP    string
		Count int
	}
	var ips []ipCount
	for ip, count := range stats.ConnectionsByIP {
		ips = append(ips, ipCount{ip, count})
	}
	sort.Slice(ips, func(i, j int) bool {
		return ips[i].Count > ips[j].Count
	})

	for i, ic := range ips {
		if i >= 10 {
			break
		}
		fmt.Printf("  %s: %d connections\n", ic.IP, ic.Count)
	}
	fmt.Println("═"*60)
}

// PrintConnections displays detailed connection list
func PrintConnections(connections []Connection, filter string) {
	fmt.Println("\n" + "═"*60)
	fmt.Println("ACTIVE NETWORK CONNECTIONS")
	if filter != "" {
		fmt.Printf("Filter: %s\n", filter)
	}
	fmt.Println("═"*60)
	fmt.Printf("%-6s %-22s %-22s %-12s %-10s\n", 
		"PROTO", "LOCAL ADDRESS", "REMOTE ADDRESS", "STATE", "PID/PROCESS")
	fmt.Println("─"*60)

	displayed := 0
	for _, conn := range connections {
		// Apply filter
		if filter != "" {
			switch filter {
			case "established":
				if conn.State != "ESTABLISHED" {
					continue
				}
			case "listening":
				if conn.State != "LISTENING" {
					continue
				}
			case "tcp":
				if conn.Protocol != "TCP" {
					continue
				}
			case "udp":
				if conn.Protocol != "UDP" {
					continue
				}
			}
		}

		localAddr := fmt.Sprintf("%s:%s", conn.LocalAddr, conn.LocalPort)
		remoteAddr := fmt.Sprintf("%s:%s", conn.RemoteAddr, conn.RemotePort)
		pidProcess := fmt.Sprintf("%s/%s", conn.PID, conn.Process)

		// Truncate if too long
		if len(localAddr) > 22 {
			localAddr = localAddr[:19] + "..."
		}
		if len(remoteAddr) > 22 {
			remoteAddr = remoteAddr[:19] + "..."
		}

		fmt.Printf("%-6s %-22s %-22s %-12s %-10s\n",
			conn.Protocol, localAddr, remoteAddr, conn.State, pidProcess)
		displayed++
	}

	fmt.Println("─"*60)
	fmt.Printf("Total: %d connections displayed\n", displayed)
	fmt.Println("═"*60)
}

// MonitorTraffic continuously monitors network traffic
func MonitorTraffic(duration time.Duration, interval time.Duration) {
	fmt.Printf("\n🔍 Monitoring network traffic for %v (updating every %v)\n", duration, interval)
	fmt.Println("Press Ctrl+C to stop monitoring")
	fmt.Println("─"*60)

	var mutex sync.Mutex
	ticker := time.NewTicker(interval)
	timeout := time.After(duration)

	connectionHistory := make(map[string]int)

	for {
		select {
		case <-timeout:
			ticker.Stop()
			fmt.Println("\n✓ Monitoring complete")
			
			// Print summary
			fmt.Println("\nMonitoring Summary:")
			fmt.Printf("  Duration: %v\n", duration)
			fmt.Printf("  Unique connections observed: %d\n", len(connectionHistory))
			return

		case t := <-ticker.C:
			connections, err := GetActiveConnections()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			mutex.Lock()
			for _, conn := range connections {
				key := fmt.Sprintf("%s:%s->%s:%s", 
					conn.LocalAddr, conn.LocalPort, conn.RemoteAddr, conn.RemotePort)
				connectionHistory[key]++
			}
			mutex.Unlock()

			stats := CalculateStatistics(connections)
			
			// Clear screen (simplified)
			fmt.Printf("\n[%s] Active: %d | Established: %d | Listening: %d | Unique IPs: %d\n",
				t.Format("15:04:05"), 
				stats.ActiveConnections,
				stats.EstablishedConns,
				stats.ListeningPorts,
				stats.UniqueRemoteIPs)
		}
	}
}

// DetectSuspiciousActivity looks for potentially suspicious connections
func DetectSuspiciousActivity(connections []Connection) {
	fmt.Println("\n" + "═"*60)
	fmt.Println("SUSPICIOUS ACTIVITY DETECTION")
	fmt.Println("═"*60)

	suspicious := []string{}

	// Check for unusual ports
	for _, conn := range connections {
		if conn.State == "ESTABLISHED" {
			port := conn.RemotePort
			if port != "*" && wellKnownPorts[port] == "" {
				portNum := 0
				fmt.Sscanf(port, "%d", &portNum)
				if portNum > 49152 { // Ephemeral port range
					suspicious = append(suspicious, 
						fmt.Sprintf("⚠️  Unusual high port: %s:%s -> %s:%s (%s)",
							conn.LocalAddr, conn.LocalPort, conn.RemoteAddr, port, conn.Process))
				}
			}
		}

		// Check for multiple connections from same process
		// (Real implementation would track and count)
	}

	// Check for too many connections to same IP
	ipCounts := make(map[string]int)
	for _, conn := range connections {
		if conn.State == "ESTABLISHED" && conn.RemoteAddr != "0.0.0.0" {
			ipCounts[conn.RemoteAddr]++
		}
	}

	for ip, count := range ipCounts {
		if count > 10 {
			suspicious = append(suspicious, 
				fmt.Sprintf("⚠️  High connection count to %s: %d connections", ip, count))
		}
	}

	if len(suspicious) == 0 {
		fmt.Println("✓ No suspicious activity detected")
	} else {
		fmt.Println("Potential Issues Found:")
		for _, alert := range suspicious {
			fmt.Println("  " + alert)
		}
	}
	fmt.Println("═"*60)
}

// ExportConnections exports connections to file
func ExportConnections(connections []Connection, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	
	// Write header
	fmt.Fprintf(writer, "# Network Traffic Analysis Report\n")
	fmt.Fprintf(writer, "# Generated: %s\n\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(writer, "Protocol,LocalAddress,LocalPort,RemoteAddress,RemotePort,State,PID,Process\n")

	// Write connections
	for _, conn := range connections {
		fmt.Fprintf(writer, "%s,%s,%s,%s,%s,%s,%s,%s\n",
			conn.Protocol, conn.LocalAddr, conn.LocalPort,
			conn.RemoteAddr, conn.RemotePort, conn.State,
			conn.PID, conn.Process)
	}

	writer.Flush()
	fmt.Printf("📄 Report exported to: %s\n", filename)
	return nil
}

// PrintMenu displays interactive menu
func PrintMenu() {
	fmt.Println("\n" + "═"*60)
	fmt.Println("MAIN MENU")
	fmt.Println("═"*60)
	fmt.Println("1. View All Connections")
	fmt.Println("2. View Established Connections")
	fmt.Println("3. View Listening Ports")
	fmt.Println("4. View Statistics")
	fmt.Println("5. Detect Suspicious Activity")
	fmt.Println("6. Monitor Traffic (Real-time)")
	fmt.Println("7. Export Report")
	fmt.Println("8. Exit")
	fmt.Println("═"*60)
	fmt.Print("Select option: ")
}

func main() {
	PrintBanner()

	fmt.Println("\n📡 Network Traffic Analyzer")
	fmt.Println("Monitor and analyze network connections on your system\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		PrintMenu()

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			reader.ReadString('\n') // Clear buffer
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		connections, err := GetActiveConnections()
		if err != nil {
			fmt.Printf("Error getting connections: %v\n", err)
			continue
		}

		switch choice {
		case 1:
			PrintConnections(connections, "")
		case 2:
			PrintConnections(connections, "established")
		case 3:
			PrintConnections(connections, "listening")
		case 4:
			stats := CalculateStatistics(connections)
			PrintStatistics(stats)
		case 5:
			DetectSuspiciousActivity(connections)
		case 6:
			fmt.Print("\nMonitor duration (seconds) [60]: ")
			durationStr, _ := reader.ReadString('\n')
			durationStr = strings.TrimSpace(durationStr)
			duration := 60
			if durationStr != "" {
				fmt.Sscanf(durationStr, "%d", &duration)
			}
			MonitorTraffic(time.Duration(duration)*time.Second, 2*time.Second)
		case 7:
			filename := fmt.Sprintf("network_analysis_%s.csv", time.Now().Format("20060102_150405"))
			err := ExportConnections(connections, filename)
			if err != nil {
				fmt.Printf("Error exporting: %v\n", err)
			}
		case 8:
			fmt.Println("\nThank you for using Network Traffic Analyzer! 🌐")
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please select 1-8.")
		}

		fmt.Print("\nPress Enter to continue...")
		reader.ReadString('\n')
	}
}