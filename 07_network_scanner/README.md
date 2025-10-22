# Network Security Scanner

A comprehensive network security scanning tool written in Go for educational purposes.

## Features

- **Port Scanning**: TCP port scanning with configurable concurrency
- **Service Detection**: Identifies common services running on open ports
- **Banner Grabbing**: Attempts to retrieve service banners for fingerprinting
- **Host Discovery**: Discovers active hosts in network ranges
- **OS Fingerprinting**: Basic operating system detection
- **Vulnerability Assessment**: Identifies potential security issues
- **Interactive Interface**: User-friendly menu system

## Usage

```bash
go run main.go
```

### Scan Options

1. **Single Host Common Ports**: Scans the most common ports on a target
2. **Port Range Scan**: Scans a custom range of ports
3. **Network Discovery**: Finds active hosts in a network range
4. **Comprehensive Scan**: Performs detailed scanning with additional checks

## Example Commands

```bash
# Scan common ports on a host
Target: 192.168.1.1

# Scan port range
Target: 192.168.1.1
Start Port: 1
End Port: 1000

# Network discovery
Network: 192.168.1.0/24
```

## Security Considerations

⚠️ **Important**: This tool is for educational purposes only. Only use on networks you own or have explicit permission to test.

## Technical Details

- **Concurrency**: Configurable concurrent scanning for performance
- **Timeout**: 5-second timeout for connection attempts
- **Service Database**: Built-in database of common services
- **Error Handling**: Comprehensive error handling and reporting

## Common Ports Scanned

- 21 (FTP)
- 22 (SSH)
- 23 (Telnet)
- 25 (SMTP)
- 53 (DNS)
- 80 (HTTP)
- 443 (HTTPS)
- 3389 (RDP)
- And more...

## Output Format

```
=== Scan Results for 192.168.1.1 ===
Port     State        Service         Banner
------------------------------------------------------------
22       open         SSH             OpenSSH 7.4
80       open         HTTP            Apache/2.4.6
443      open         HTTPS           

Summary: 3 open ports found
OS Detection: Linux/Unix (SSH detected)

Potential Vulnerabilities:
- HTTP service - check for web vulnerabilities
- HTTPS service - verify SSL/TLS configuration
```

## Educational Value

This scanner demonstrates:
- Network programming in Go
- Concurrent programming patterns
- TCP connection handling
- Service fingerprinting techniques
- Basic vulnerability assessment
- Network security concepts
