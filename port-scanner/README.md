# Port Scanner

A fast, concurrent TCP port scanner written in Go for cybersecurity labs and network reconnaissance.

## 📋 Description

This port scanner is a network security tool that checks which ports are open on a target system. It uses concurrent goroutines for fast scanning and provides detailed information about open ports and their associated services.

## ⚠️ Legal Disclaimer

**IMPORTANT**: This tool is for educational and authorized security testing purposes only. 

- ✅ Only scan systems you own or have explicit written permission to test
- ❌ Unauthorized port scanning may be illegal in your jurisdiction
- ⚖️ The author is not responsible for misuse of this tool

Always comply with:
- Computer Fraud and Abuse Act (CFAA) in the US
- Computer Misuse Act in the UK
- Similar laws in your country

## ✨ Features

### 🚀 Performance
- **Concurrent Scanning**: Uses goroutines for parallel port checking
- **Configurable Workers**: Adjust concurrency level (default: 100)
- **Fast Results**: Scan 1000 ports in seconds
- **Timeout Control**: Configurable connection timeout

### 🎯 Functionality
- **TCP Port Scanning**: Check if ports accept connections
- **Service Detection**: Identifies common services (HTTP, SSH, FTP, etc.)
- **Custom Port Ranges**: Scan any port range (1-65535)
- **Real-time Output**: See open ports as they're discovered
- **Detailed Summary**: Comprehensive scan results report

### 📊 Output
- Open port detection with service names
- Total ports scanned
- Scan duration
- Formatted summary table

## 🚀 Installation

### Prerequisites
- Go 1.16 or higher
- Network access to target system

### Setup

1. Clone the repository:
```bash
git clone https://github.com/arya2004/cybersecurity.git
cd cybersecurity/port-scanner
```

2. No additional dependencies required (uses Go standard library)

## 💻 Usage

### Basic Syntax
```bash
go run main.go <hostname> <start_port> <end_port>
```

### Examples

#### Scan common ports on localhost
```bash
go run main.go localhost 1 1000
```

#### Scan web ports
```bash
go run main.go example.com 80 443
```

#### Scan all ports (takes longer)
```bash
go run main.go 192.168.1.1 1 65535
```

#### Scan specific service ports
```bash
go run main.go scanme.nmap.org 20 25
```

## 📸 Sample Output

```
╔═══════════════════════════════════════╗
║     Go Port Scanner v1.0              ║
║     Cybersecurity Lab Tool            ║
╚═══════════════════════════════════════╝

⚠️  WARNING: Only scan systems you have explicit permission to test!

Starting port scan on localhost (Ports 1-1000)
This may take a few moments...

[+] Port 22 is OPEN (SSH)
[+] Port 80 is OPEN (HTTP)
[+] Port 443 is OPEN (HTTPS)

══════════════════════════════════════════════════
Scan Summary for localhost
══════════════════════════════════════════════════
Port Range: 1-1000
Total Ports Scanned: 1000
Open Ports Found: 3
Scan Duration: 2.34 seconds

Open Ports Details:
──────────────────────────────────────────────────
PORT       STATE      SERVICE             
──────────────────────────────────────────────────
22         Open       SSH                 
80         Open       HTTP                
443        Open       HTTPS               
══════════════════════════════════════════════════
```

## 🔧 How It Works

### Architecture

1. **Input Validation**: Validates hostname and port range
2. **Worker Pool**: Creates concurrent workers (goroutines)
3. **TCP Connection**: Attempts connection to each port
4. **Service Detection**: Maps ports to known services
5. **Result Collection**: Thread-safe result aggregation
6. **Summary Generation**: Sorted and formatted output

### Technical Details

```go
// Key components:
- net.DialTimeout()   // TCP connection attempt
- sync.WaitGroup      // Goroutine synchronization  
- sync.Mutex          // Thread-safe result collection
- Channel-based work distribution
```

## 📚 Recognized Services

The scanner recognizes 18+ common services:

| Port | Service | Description |
|------|---------|-------------|
| 21 | FTP | File Transfer Protocol |
| 22 | SSH | Secure Shell |
| 23 | Telnet | Telnet Protocol |
| 25 | SMTP | Email (Send) |
| 53 | DNS | Domain Name System |
| 80 | HTTP | Web Server |
| 110 | POP3 | Email (Receive) |
| 143 | IMAP | Email (Receive) |
| 443 | HTTPS | Secure Web |
| 445 | SMB | File Sharing |
| 3306 | MySQL | Database |
| 3389 | RDP | Remote Desktop |
| 5432 | PostgreSQL | Database |
| 5900 | VNC | Remote Desktop |
| 6379 | Redis | Database |
| 8080 | HTTP-Alt | Alternative Web |
| 8443 | HTTPS-Alt | Alternative Secure Web |
| 27017 | MongoDB | Database |

## ⚙️ Configuration

### Adjustable Parameters (in code)

```go
timeout := 1 * time.Second  // Connection timeout
workers := 100              // Concurrent goroutines
```

Modify these in `main()` function for your needs:
- **Timeout**: Lower for faster scans, higher for slow networks
- **Workers**: Higher for faster scans (be careful with resource usage)

## 🎓 Educational Use Cases

### Cybersecurity Labs
- Network reconnaissance exercises
- Security assessment training
- Vulnerability identification
- Network topology mapping

### Learning Objectives
- Understanding network protocols
- TCP handshake mechanics
- Concurrent programming in Go
- Security tool development

## 🔒 Security Considerations

### Best Practices
- Always get written permission before scanning
- Scan your own systems for testing
- Use rate limiting for production networks
- Monitor scan impact on target systems

### Ethical Hacking
- Part of reconnaissance phase (OSINT)
- Helps identify unnecessary open ports
- Assists in security hardening
- Compliance with security policies

## 🐛 Troubleshooting

### "Connection refused" errors
- Normal for closed ports
- Target may have firewall blocking

### Slow scanning
- Increase timeout value
- Reduce number of workers
- Network latency issues

### No ports found
- Firewall blocking all ports
- Incorrect hostname/IP
- Target system is down

## 🔮 Future Enhancements

Potential improvements:
- UDP port scanning
- OS fingerprinting
- Banner grabbing
- XML/JSON output formats
- Stealth scanning techniques
- Port knock detection
- Service version detection

## 👨‍💻 Author

**Ashvin**
- GitHub: [@ashvin2005](https://github.com/ashvin2005)
- LinkedIn: [ashvin-tiwari](https://linkedin.com/in/ashvin-tiwari)

## 🎃 Hacktoberfest 2025

Created as part of Hacktoberfest 2025 contributions to the Cybersecurity Lab Codes repository.

## 📄 License

MIT License - Same as parent repository

## 🙏 Acknowledgments

- Go standard library `net` package
- NMAP for inspiration
- Cybersecurity community

## 📚 Additional Resources

- [RFC 793 - TCP](https://tools.ietf.org/html/rfc793)
- [IANA Port Numbers](https://www.iana.org/assignments/service-names-port-numbers)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)

---

**Remember: With great power comes great responsibility. Use wisely!** 🔐