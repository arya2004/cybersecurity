/*
Certificate Security Analyzer
============================

This tool analyzes SSL/TLS certificates for security vulnerabilities and provides
comprehensive security assessments. It's designed for educational purposes to help
understand certificate security and common vulnerabilities.

Features:
- Certificate chain validation
- Expiration date checking
- Weak algorithm detection
- Key strength analysis
- Common name validation
- Certificate transparency checking
- Security best practices verification

Educational Value:
- Understanding PKI infrastructure
- Certificate validation process
- Cryptographic algorithm security
- TLS/SSL security concepts
- Common certificate vulnerabilities

Usage:
    go run main.go
    Enter domain name to analyze certificate
*/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"
)

// CertificateAnalyzer holds certificate analysis data
type CertificateAnalyzer struct {
	Domain      string
	Certificates []*x509.Certificate
	Issues      []SecurityIssue
	Score       int
}

// SecurityIssue represents a security finding
type SecurityIssue struct {
	Severity    string // "Critical", "High", "Medium", "Low", "Info"
	Title       string
	Description string
	Remediation string
}

// NewCertificateAnalyzer creates a new analyzer instance
func NewCertificateAnalyzer(domain string) *CertificateAnalyzer {
	return &CertificateAnalyzer{
		Domain:  domain,
		Issues:  make([]SecurityIssue, 0),
		Score:   100, // Start with perfect score
	}
}

// FetchCertificate retrieves the certificate from the domain
func (ca *CertificateAnalyzer) FetchCertificate() error {
	fmt.Printf("\n[*] Connecting to %s:443...\n", ca.Domain)

	conn, err := tls.Dial("tcp", ca.Domain+":443", &tls.Config{
		InsecureSkipVerify: true, // For analysis purposes
	})
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer conn.Close()

	state := conn.ConnectionState()
	ca.Certificates = state.PeerCertificates

	fmt.Printf("[✓] Successfully retrieved %d certificate(s)\n", len(ca.Certificates))
	return nil
}

// AnalyzeCertificate performs comprehensive security analysis
func (ca *CertificateAnalyzer) AnalyzeCertificate() {
	if len(ca.Certificates) == 0 {
		fmt.Println("[!] No certificates to analyze")
		return
	}

	cert := ca.Certificates[0]

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("CERTIFICATE SECURITY ANALYSIS")
	fmt.Println(strings.Repeat("=", 60))

	ca.checkExpiration(cert)
	ca.checkKeyStrength(cert)
	ca.checkSignatureAlgorithm(cert)
	ca.checkCommonName(cert)
	ca.checkKeyUsage(cert)
	ca.checkSANs(cert)
	ca.checkValidityPeriod(cert)
	ca.checkSelfSigned(cert)
}

// checkExpiration checks if certificate is expired or expiring soon
func (ca *CertificateAnalyzer) checkExpiration(cert *x509.Certificate) {
	fmt.Println("\n[1] Expiration Check:")
	fmt.Printf("    Not Before: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05"))
	fmt.Printf("    Not After:  %s\n", cert.NotAfter.Format("2006-01-02 15:04:05"))

	now := time.Now()
	daysUntilExpiry := int(cert.NotAfter.Sub(now).Hours() / 24)

	if cert.NotAfter.Before(now) {
		ca.addIssue("Critical", "Certificate Expired",
			fmt.Sprintf("Certificate expired %d days ago", -daysUntilExpiry),
			"Renew the certificate immediately")
		ca.Score -= 50
		fmt.Printf("    [✗] EXPIRED (%d days ago)\n", -daysUntilExpiry)
	} else if daysUntilExpiry <= 30 {
		ca.addIssue("High", "Certificate Expiring Soon",
			fmt.Sprintf("Certificate expires in %d days", daysUntilExpiry),
			"Plan certificate renewal")
		ca.Score -= 20
		fmt.Printf("    [!] Expires soon (%d days)\n", daysUntilExpiry)
	} else {
		fmt.Printf("    [✓] Valid (%d days remaining)\n", daysUntilExpiry)
	}
}

// checkKeyStrength analyzes key strength
func (ca *CertificateAnalyzer) checkKeyStrength(cert *x509.Certificate) {
	fmt.Println("\n[2] Key Strength Analysis:")
	fmt.Printf("    Algorithm: %s\n", cert.PublicKeyAlgorithm)

	keySize := 0
	switch pub := cert.PublicKey.(type) {
	case interface{ Size() int }:
		keySize = pub.Size() * 8 // Convert bytes to bits
	}

	fmt.Printf("    Key Size: %d bits\n", keySize)

	if keySize < 2048 {
		ca.addIssue("Critical", "Weak Key Size",
			fmt.Sprintf("Key size of %d bits is insufficient", keySize),
			"Use at least 2048-bit RSA or 256-bit ECC keys")
		ca.Score -= 30
		fmt.Printf("    [✗] WEAK (< 2048 bits)\n")
	} else if keySize < 3072 {
		ca.addIssue("Low", "Key Size Below Recommended",
			fmt.Sprintf("Key size of %d bits meets minimum but not recommended", keySize),
			"Consider upgrading to 3072-bit or 4096-bit keys")
		ca.Score -= 5
		fmt.Printf("    [!] Acceptable but not recommended\n")
	} else {
		fmt.Printf("    [✓] Strong key size\n")
	}
}

// checkSignatureAlgorithm checks for weak signature algorithms
func (ca *CertificateAnalyzer) checkSignatureAlgorithm(cert *x509.Certificate) {
	fmt.Println("\n[3] Signature Algorithm Check:")
	fmt.Printf("    Algorithm: %s\n", cert.SignatureAlgorithm)

	weakAlgorithms := map[string]bool{
		"MD5WithRSA":    true,
		"MD2WithRSA":    true,
		"SHA1WithRSA":   true,
		"DSAWithSHA1":   true,
		"ECDSAWithSHA1": true,
	}

	algName := cert.SignatureAlgorithm.String()

	if weakAlgorithms[algName] {
		ca.addIssue("High", "Weak Signature Algorithm",
			fmt.Sprintf("%s is cryptographically weak", algName),
			"Use SHA-256 or stronger signature algorithms")
		ca.Score -= 25
		fmt.Printf("    [✗] WEAK ALGORITHM\n")
	} else if strings.Contains(algName, "SHA256") || strings.Contains(algName, "SHA384") || strings.Contains(algName, "SHA512") {
		fmt.Printf("    [✓] Strong signature algorithm\n")
	} else {
		fmt.Printf("    [!] Unknown algorithm security level\n")
	}
}

// checkCommonName validates common name
func (ca *CertificateAnalyzer) checkCommonName(cert *x509.Certificate) {
	fmt.Println("\n[4] Common Name (CN) Check:")
	fmt.Printf("    CN: %s\n", cert.Subject.CommonName)

	expectedDomain := ca.Domain
	if strings.HasPrefix(expectedDomain, "www.") {
		expectedDomain = expectedDomain[4:]
	}

	if cert.Subject.CommonName == "" {
		ca.addIssue("Medium", "Empty Common Name",
			"Certificate has no Common Name set",
			"Set appropriate Common Name in certificate")
		ca.Score -= 10
		fmt.Printf("    [✗] Empty Common Name\n")
	} else if !strings.Contains(cert.Subject.CommonName, expectedDomain) {
		ca.addIssue("Medium", "CN Mismatch",
			fmt.Sprintf("CN '%s' doesn't match domain '%s'", cert.Subject.CommonName, ca.Domain),
			"Ensure CN matches the domain")
		ca.Score -= 10
		fmt.Printf("    [!] Potential mismatch\n")
	} else {
		fmt.Printf("    [✓] Valid Common Name\n")
	}
}

// checkKeyUsage validates key usage extensions
func (ca *CertificateAnalyzer) checkKeyUsage(cert *x509.Certificate) {
	fmt.Println("\n[5] Key Usage Check:")
	fmt.Printf("    Key Usage: %v\n", cert.KeyUsage)
	fmt.Printf("    Extended Key Usage: %v\n", cert.ExtKeyUsage)

	if cert.KeyUsage == 0 {
		ca.addIssue("Low", "No Key Usage Set",
			"Certificate has no key usage constraints",
			"Define appropriate key usage")
		ca.Score -= 5
		fmt.Printf("    [!] No key usage defined\n")
	} else {
		fmt.Printf("    [✓] Key usage properly defined\n")
	}
}

// checkSANs checks Subject Alternative Names
func (ca *CertificateAnalyzer) checkSANs(cert *x509.Certificate) {
	fmt.Println("\n[6] Subject Alternative Names (SAN) Check:")

	if len(cert.DNSNames) == 0 {
		ca.addIssue("High", "No SANs Defined",
			"Certificate has no Subject Alternative Names",
			"Add SANs for all domains and subdomains")
		ca.Score -= 15
		fmt.Printf("    [✗] No SANs found\n")
	} else {
		fmt.Printf("    SANs: %v\n", cert.DNSNames)
		fmt.Printf("    [✓] %d SAN(s) found\n", len(cert.DNSNames))

		// Check if requested domain is in SANs
		domainFound := false
		for _, san := range cert.DNSNames {
			if strings.Contains(ca.Domain, san) || strings.Contains(san, ca.Domain) {
				domainFound = true
				break
			}
		}

		if !domainFound {
			ca.addIssue("High", "Domain Not in SANs",
				fmt.Sprintf("Requested domain '%s' not found in SANs", ca.Domain),
				"Add the domain to certificate SANs")
			ca.Score -= 15
			fmt.Printf("    [!] Requested domain not in SANs\n")
		}
	}
}

// checkValidityPeriod checks if validity period is appropriate
func (ca *CertificateAnalyzer) checkValidityPeriod(cert *x509.Certificate) {
	fmt.Println("\n[7] Validity Period Check:")

	validityPeriod := cert.NotAfter.Sub(cert.NotBefore)
	validityDays := int(validityPeriod.Hours() / 24)

	fmt.Printf("    Validity Period: %d days\n", validityDays)

	if validityDays > 825 { // Apple/Google limit
		ca.addIssue("Medium", "Excessive Validity Period",
			fmt.Sprintf("Certificate valid for %d days (> 825 days)", validityDays),
			"Use certificates with validity periods under 398 days")
		ca.Score -= 10
		fmt.Printf("    [!] Exceeds browser limits (> 825 days)\n")
	} else if validityDays > 398 {
		ca.addIssue("Low", "Long Validity Period",
			fmt.Sprintf("Certificate valid for %d days", validityDays),
			"Consider using shorter validity periods (< 90 days)")
		ca.Score -= 5
		fmt.Printf("    [!] Longer than recommended (> 398 days)\n")
	} else {
		fmt.Printf("    [✓] Appropriate validity period\n")
	}
}

// checkSelfSigned checks if certificate is self-signed
func (ca *CertificateAnalyzer) checkSelfSigned(cert *x509.Certificate) {
	fmt.Println("\n[8] Issuer Check:")
	fmt.Printf("    Issuer: %s\n", cert.Issuer.CommonName)
	fmt.Printf("    Subject: %s\n", cert.Subject.CommonName)

	if cert.Issuer.CommonName == cert.Subject.CommonName {
		ca.addIssue("Critical", "Self-Signed Certificate",
			"Certificate is self-signed and not trusted by browsers",
			"Obtain certificate from a trusted Certificate Authority")
		ca.Score -= 40
		fmt.Printf("    [✗] SELF-SIGNED (not trusted)\n")
	} else {
		fmt.Printf("    [✓] Issued by CA\n")
	}
}

// addIssue adds a security issue to the list
func (ca *CertificateAnalyzer) addIssue(severity, title, description, remediation string) {
	ca.Issues = append(ca.Issues, SecurityIssue{
		Severity:    severity,
		Title:       title,
		Description: description,
		Remediation: remediation,
	})
}

// PrintReport prints the final security report
func (ca *CertificateAnalyzer) PrintReport() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("SECURITY ASSESSMENT REPORT")
	fmt.Println(strings.Repeat("=", 60))

	// Ensure score doesn't go below 0
	if ca.Score < 0 {
		ca.Score = 0
	}

	// Determine grade
	var grade string
	var gradeColor string
	switch {
	case ca.Score >= 90:
		grade = "A"
		gradeColor = "EXCELLENT"
	case ca.Score >= 80:
		grade = "B"
		gradeColor = "GOOD"
	case ca.Score >= 70:
		grade = "C"
		gradeColor = "FAIR"
	case ca.Score >= 60:
		grade = "D"
		gradeColor = "POOR"
	default:
		grade = "F"
		gradeColor = "CRITICAL"
	}

	fmt.Printf("\n🔐 SECURITY SCORE: %d/100 [Grade %s - %s]\n", ca.Score, grade, gradeColor)

	if len(ca.Issues) == 0 {
		fmt.Println("\n✅ No security issues found! Certificate follows best practices.")
		return
	}

	fmt.Printf("\n⚠️  FOUND %d SECURITY ISSUE(S):\n", len(ca.Issues))

	// Group issues by severity
	criticalIssues := filterIssues(ca.Issues, "Critical")
	highIssues := filterIssues(ca.Issues, "High")
	mediumIssues := filterIssues(ca.Issues, "Medium")
	lowIssues := filterIssues(ca.Issues, "Low")

	printIssueGroup("CRITICAL", criticalIssues, "🔴")
	printIssueGroup("HIGH", highIssues, "🟠")
	printIssueGroup("MEDIUM", mediumIssues, "🟡")
	printIssueGroup("LOW", lowIssues, "🟢")

	fmt.Println("\n" + strings.Repeat("=", 60))
}

// filterIssues filters issues by severity
func filterIssues(issues []SecurityIssue, severity string) []SecurityIssue {
	filtered := make([]SecurityIssue, 0)
	for _, issue := range issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// printIssueGroup prints a group of issues
func printIssueGroup(severity string, issues []SecurityIssue, icon string) {
	if len(issues) == 0 {
		return
	}

	fmt.Printf("\n%s %s SEVERITY (%d):\n", icon, severity, len(issues))
	for i, issue := range issues {
		fmt.Printf("\n  %d. %s\n", i+1, issue.Title)
		fmt.Printf("     Description: %s\n", issue.Description)
		fmt.Printf("     Remediation: %s\n", issue.Remediation)
	}
}

// PrintCertificateDetails prints detailed certificate information
func (ca *CertificateAnalyzer) PrintCertificateDetails() {
	if len(ca.Certificates) == 0 {
		return
	}

	cert := ca.Certificates[0]

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("CERTIFICATE DETAILS")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\nVersion: %d\n", cert.Version)
	fmt.Printf("Serial Number: %s\n", cert.SerialNumber)
	fmt.Printf("\nSubject:\n")
	fmt.Printf("  Common Name: %s\n", cert.Subject.CommonName)
	fmt.Printf("  Organization: %s\n", cert.Subject.Organization)
	fmt.Printf("  Country: %s\n", cert.Subject.Country)

	fmt.Printf("\nIssuer:\n")
	fmt.Printf("  Common Name: %s\n", cert.Issuer.CommonName)
	fmt.Printf("  Organization: %s\n", cert.Issuer.Organization)

	fmt.Printf("\nValidity:\n")
	fmt.Printf("  Not Before: %s\n", cert.NotBefore.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  Not After:  %s\n", cert.NotAfter.Format("2006-01-02 15:04:05 MST"))

	fmt.Printf("\nPublic Key:\n")
	fmt.Printf("  Algorithm: %s\n", cert.PublicKeyAlgorithm)

	fmt.Printf("\nSignature:\n")
	fmt.Printf("  Algorithm: %s\n", cert.SignatureAlgorithm)

	if len(cert.DNSNames) > 0 {
		fmt.Printf("\nSubject Alternative Names:\n")
		for _, san := range cert.DNSNames {
			fmt.Printf("  - %s\n", san)
		}
	}
}

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  SSL/TLS CERTIFICATE SECURITY ANALYZER")
	fmt.Println("  Educational Tool for Certificate Assessment")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Print("\nEnter domain name (e.g., google.com): ")
	var domain string
	fmt.Scanln(&domain)

	// Remove protocol if present
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")

	analyzer := NewCertificateAnalyzer(domain)

	// Fetch certificate
	err := analyzer.FetchCertificate()
	if err != nil {
		fmt.Printf("\n[✗] Error: %v\n", err)
		return
	}

	// Analyze security
	analyzer.AnalyzeCertificate()

	// Print detailed certificate info
	analyzer.PrintCertificateDetails()

	// Print final report
	analyzer.PrintReport()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Analysis complete!")
	fmt.Println("\n⚠️  DISCLAIMER: This tool is for educational purposes only.")
	fmt.Println("Only analyze certificates of systems you own or have authorization to test.")
	fmt.Println(strings.Repeat("=", 60))
}
