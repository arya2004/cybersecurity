# Certificate Security Analyzer

## Overview
A comprehensive SSL/TLS certificate security analysis tool written in Go. This educational tool helps security professionals and students understand certificate security by performing detailed analysis of SSL/TLS certificates and identifying potential vulnerabilities.

## Features

### Security Checks
- ✅ **Expiration Validation** - Detects expired or expiring certificates
- 🔑 **Key Strength Analysis** - Evaluates cryptographic key sizes
- 🔐 **Signature Algorithm Check** - Identifies weak signature algorithms (MD5, SHA1)
- 📛 **Common Name Validation** - Verifies CN matches domain
- 🎯 **Subject Alternative Names (SAN)** - Validates SANs coverage
- ⏰ **Validity Period Assessment** - Checks for excessive validity periods
- 🔒 **Self-Signed Detection** - Identifies untrusted self-signed certificates
- 📋 **Key Usage Verification** - Validates key usage extensions

### Security Scoring
- Automated security score (0-100)
- Letter grade assignment (A-F)
- Severity-based issue categorization
- Detailed remediation guidance

### Report Generation
- Comprehensive security assessment reports
- Issue prioritization by severity
- Detailed certificate information display
- Professional formatting

## Installation

```bash
cd 11_certificate_analyzer
go build -o certanalyzer main.go
```

## Usage

### Interactive Mode
```bash
go run main.go
```

Then enter the domain name when prompted:
```
Enter domain name (e.g., google.com): example.com
```

### Example Output
```
==============================================================
  SSL/TLS CERTIFICATE SECURITY ANALYZER
  Educational Tool for Certificate Assessment
==============================================================

[*] Connecting to google.com:443...
[✓] Successfully retrieved 3 certificate(s)

============================================================
CERTIFICATE SECURITY ANALYSIS
============================================================

[1] Expiration Check:
    Not Before: 2024-01-15 08:21:58
    Not After:  2024-04-08 08:21:57
    [✓] Valid (45 days remaining)

[2] Key Strength Analysis:
    Algorithm: RSA
    Key Size: 2048 bits
    [✓] Strong key size

🔐 SECURITY SCORE: 95/100 [Grade A - EXCELLENT]
```

## Educational Value

### Learning Objectives
1. **PKI Fundamentals** - Understanding Public Key Infrastructure
2. **Certificate Validation** - Learning certificate verification process
3. **Cryptographic Security** - Recognizing weak vs strong algorithms
4. **TLS/SSL Best Practices** - Industry standard compliance
5. **Vulnerability Assessment** - Identifying security weaknesses

### Security Concepts Covered
- X.509 certificate structure
- Certificate chain validation
- Signature algorithm security
- Key size requirements
- Subject Alternative Names (SANs)
- Certificate Authorities (CAs)
- Certificate lifecycle management

## Security Issues Detected

### Critical Issues
- **Expired Certificates** - Certificate has passed NotAfter date
- **Weak Key Size** - RSA keys < 2048 bits or ECC keys < 256 bits
- **Self-Signed Certificates** - Not trusted by browsers

### High Severity Issues
- **Weak Signature Algorithms** - MD5, SHA1 algorithms
- **Missing SANs** - No Subject Alternative Names defined
- **Certificate Expiring Soon** - Less than 30 days remaining

### Medium Severity Issues
- **CN Mismatch** - Common Name doesn't match domain
- **Excessive Validity Period** - Validity > 825 days
- **Empty Common Name** - No CN set in certificate

### Low Severity Issues
- **No Key Usage** - Missing key usage constraints
- **Long Validity Period** - Between 398-825 days

## Remediation Guide

### For Expired Certificates
```
1. Generate new Certificate Signing Request (CSR)
2. Submit to Certificate Authority
3. Install new certificate immediately
4. Update all dependent systems
```

### For Weak Algorithms
```
1. Request certificate with SHA-256 or stronger
2. Use 2048-bit RSA minimum (4096-bit recommended)
3. Or use 256-bit ECC keys
4. Update certificate before expiration
```

### For Self-Signed Certificates
```
1. Obtain certificate from trusted CA (Let's Encrypt, DigiCert, etc.)
2. Properly configure certificate chain
3. Ensure root and intermediate certificates are included
4. Test in multiple browsers
```

## Technical Details

### Analyzed Certificate Properties
- Version number
- Serial number
- Subject Distinguished Name (DN)
- Issuer DN
- Validity period (NotBefore/NotAfter)
- Public key algorithm and size
- Signature algorithm
- Key usage extensions
- Extended key usage
- Subject Alternative Names
- Basic constraints

### Scoring System
```
Starting Score: 100 points

Deductions:
- Expired certificate: -50 points
- Self-signed certificate: -40 points
- Weak key size (< 2048 bits): -30 points
- Weak signature algorithm: -25 points
- Expiring soon (< 30 days): -20 points
- No SANs: -15 points
- Domain not in SANs: -15 points
- CN mismatch: -10 points
- Excessive validity (> 825 days): -10 points
- No key usage: -5 points
- Long validity (398-825 days): -5 points
```

### Grade Assignment
```
A: 90-100 points (Excellent)
B: 80-89 points (Good)
C: 70-79 points (Fair)
D: 60-69 points (Poor)
F: 0-59 points (Critical)
```

## Code Structure

```
main.go
├── CertificateAnalyzer (struct)
│   ├── FetchCertificate() - Retrieves certificate via TLS
│   ├── AnalyzeCertificate() - Performs all security checks
│   ├── checkExpiration() - Validates expiration dates
│   ├── checkKeyStrength() - Analyzes key cryptographic strength
│   ├── checkSignatureAlgorithm() - Detects weak algorithms
│   ├── checkCommonName() - Validates CN against domain
│   ├── checkKeyUsage() - Verifies key usage extensions
│   ├── checkSANs() - Validates Subject Alternative Names
│   ├── checkValidityPeriod() - Checks validity duration
│   ├── checkSelfSigned() - Detects self-signed certs
│   └── PrintReport() - Generates security report
└── SecurityIssue (struct) - Represents a security finding
```

## Use Cases

### 1. Security Auditing
```bash
# Audit your organization's certificates
go run main.go
Enter domain: company.com
```

### 2. Certificate Monitoring
```bash
# Check certificate expiration
# Integrate into monitoring systems
# Set up alerts for expiring certificates
```

### 3. Compliance Checking
```bash
# Verify compliance with security standards
# Ensure proper key sizes
# Validate signature algorithms
```

### 4. Educational Labs
```bash
# Teach certificate security concepts
# Demonstrate vulnerability detection
# Practice security assessment
```

## Limitations

- Requires active network connection
- Only analyzes server certificates (not client certificates)
- Does not verify certificate revocation status (OCSP/CRL)
- Does not perform certificate transparency log checking
- Basic chain validation only

## Future Enhancements

- [ ] OCSP and CRL revocation checking
- [ ] Certificate Transparency log verification
- [ ] Support for client certificate analysis
- [ ] Batch processing of multiple domains
- [ ] JSON/CSV report export
- [ ] Certificate chain visualization
- [ ] Historical certificate tracking
- [ ] Integration with vulnerability databases

## Security Notice

⚠️ **IMPORTANT**: This tool is designed for educational purposes and authorized security testing only.

### Legal and Ethical Use
- Only analyze certificates of systems you own or have explicit authorization to test
- Respect privacy and security policies
- Comply with applicable laws and regulations
- Do not use for malicious purposes
- Use responsibly in production environments

## References

### Standards and Guidelines
- [RFC 5280](https://tools.ietf.org/html/rfc5280) - X.509 Certificate Standard
- [RFC 6125](https://tools.ietf.org/html/rfc6125) - Certificate Validation
- [CA/Browser Forum Baseline Requirements](https://cabforum.org/baseline-requirements-documents/)
- [NIST SP 800-52 Rev. 2](https://csrc.nist.gov/publications/detail/sp/800-52/rev-2/final) - TLS Guidelines

### Resources
- [SSL Labs SSL/TLS Best Practices](https://github.com/ssllabs/research/wiki/SSL-and-TLS-Deployment-Best-Practices)
- [Mozilla SSL Configuration Generator](https://ssl-config.mozilla.org/)
- [OWASP Transport Layer Protection Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Transport_Layer_Protection_Cheat_Sheet.html)

## Contributing

Contributions are welcome! Please ensure:
- Code follows Go best practices
- Security checks are accurate and documented
- Educational value is maintained
- Comprehensive comments are included
- Test coverage is provided

## License

This tool is part of the Cybersecurity Lab Codes repository and is licensed under the MIT License.

## Author

Developed for educational purposes to help security professionals and students understand SSL/TLS certificate security.

---

**Disclaimer**: This tool is for educational and authorized testing purposes only. The authors are not responsible for any misuse of this tool.
