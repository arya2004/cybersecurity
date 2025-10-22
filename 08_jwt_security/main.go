// JWT Security Tool - Advanced JWT Analysis and Security Testing
// This tool provides comprehensive JWT (JSON Web Token) security analysis including:
// - JWT parsing and validation
// - Signature verification bypass testing
// - Common JWT vulnerabilities detection
// - JWT manipulation and crafting
// - Security recommendations

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"bufio"
)

// JWT represents the structure of a JSON Web Token
type JWT struct {
	Header    map[string]interface{} `json:"header"`
	Payload   map[string]interface{} `json:"payload"`
	Signature string                 `json:"signature"`
	Raw       string                 `json:"raw"`
}

// JWTHeader represents the JWT header structure
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	KeyID     string `json:"kid,omitempty"`
}

// JWTPayload represents the JWT payload structure
type JWTPayload struct {
	Issuer         string      `json:"iss,omitempty"`
	Subject        string      `json:"sub,omitempty"`
	Audience       interface{} `json:"aud,omitempty"`
	ExpirationTime int64       `json:"exp,omitempty"`
	NotBefore      int64       `json:"nbf,omitempty"`
	IssuedAt       int64       `json:"iat,omitempty"`
	JWTID          string      `json:"jti,omitempty"`
	Custom         map[string]interface{} `json:"-"`
}

// VulnerabilityReport represents security issues found in JWT
type VulnerabilityReport struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Recommendations []string        `json:"recommendations"`
	Severity        string          `json:"severity"`
}

// Vulnerability represents a specific security issue
type Vulnerability struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Severity    string `json:"severity"`
	CVE         string `json:"cve,omitempty"`
}

// JWTSecurityTool provides JWT security analysis capabilities
type JWTSecurityTool struct {
	CommonSecrets []string
}

// NewJWTSecurityTool creates a new JWT security tool instance
func NewJWTSecurityTool() *JWTSecurityTool {
	return &JWTSecurityTool{
		CommonSecrets: []string{
			"secret",
			"password",
			"123456",
			"admin",
			"test",
			"key",
			"jwt",
			"token",
			"your-256-bit-secret",
			"supersecret",
			"",
		},
	}
}

// ParseJWT parses a JWT token and extracts its components
func (jst *JWTSecurityTool) ParseJWT(token string) (*JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	// Decode header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %v", err)
	}

	var header map[string]interface{}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %v", err)
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %v", err)
	}

	return &JWT{
		Header:    header,
		Payload:   payload,
		Signature: parts[2],
		Raw:       token,
	}, nil
}

// AnalyzeSecurity performs comprehensive security analysis of JWT
func (jst *JWTSecurityTool) AnalyzeSecurity(jwt *JWT) VulnerabilityReport {
	var vulnerabilities []Vulnerability
	var recommendations []string
	severity := "LOW"

	// Check for algorithm vulnerabilities
	if alg, ok := jwt.Header["alg"].(string); ok {
		vulns, recs := jst.checkAlgorithmVulnerabilities(alg)
		vulnerabilities = append(vulnerabilities, vulns...)
		recommendations = append(recommendations, recs...)
		
		if alg == "none" {
			severity = "CRITICAL"
		} else if alg == "HS256" {
			severity = "MEDIUM"
		}
	}

	// Check payload vulnerabilities
	vulns, recs := jst.checkPayloadVulnerabilities(jwt.Payload)
	vulnerabilities = append(vulnerabilities, vulns...)
	recommendations = append(recommendations, recs...)

	// Check for weak secrets (if HMAC)
	if alg, ok := jwt.Header["alg"].(string); ok && strings.HasPrefix(alg, "HS") {
		if jst.isWeakSecret(jwt) {
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Type:        "WEAK_SECRET",
				Description: "JWT signed with weak or common secret",
				Impact:      "Attackers can forge tokens with common secrets",
				Severity:    "HIGH",
			})
			recommendations = append(recommendations, "Use strong, randomly generated secrets for HMAC signing")
			severity = "HIGH"
		}
	}

	// Check signature bypass vulnerabilities
	vulns, recs = jst.checkSignatureBypass(jwt)
	vulnerabilities = append(vulnerabilities, vulns...)
	recommendations = append(recommendations, recs...)

	return VulnerabilityReport{
		Vulnerabilities: vulnerabilities,
		Recommendations: recommendations,
		Severity:        severity,
	}
}

// checkAlgorithmVulnerabilities checks for algorithm-related vulnerabilities
func (jst *JWTSecurityTool) checkAlgorithmVulnerabilities(algorithm string) ([]Vulnerability, []string) {
	var vulnerabilities []Vulnerability
	var recommendations []string

	switch algorithm {
	case "none":
		vulnerabilities = append(vulnerabilities, Vulnerability{
			Type:        "ALGORITHM_NONE",
			Description: "JWT uses 'none' algorithm - no signature verification",
			Impact:      "Anyone can create valid tokens without signature",
			Severity:    "CRITICAL",
			CVE:         "CVE-2015-9235",
		})
		recommendations = append(recommendations, "Never use 'none' algorithm in production")
		recommendations = append(recommendations, "Implement proper signature verification")

	case "HS256":
		vulnerabilities = append(vulnerabilities, Vulnerability{
			Type:        "HMAC_ALGORITHM",
			Description: "JWT uses HMAC algorithm which requires shared secret",
			Impact:      "Vulnerable to secret brute force attacks",
			Severity:    "MEDIUM",
		})
		recommendations = append(recommendations, "Consider using RS256 (RSA) for better security")
		recommendations = append(recommendations, "Ensure HMAC secret is sufficiently long and random")

	case "RS256":
		// RS256 is generally secure, but check for key confusion
		vulnerabilities = append(vulnerabilities, Vulnerability{
			Type:        "KEY_CONFUSION",
			Description: "Potential RS256/HS256 key confusion vulnerability",
			Impact:      "Public key might be used as HMAC secret",
			Severity:    "MEDIUM",
			CVE:         "CVE-2016-10555",
		})
		recommendations = append(recommendations, "Validate algorithm in JWT header matches expected algorithm")
	}

	return vulnerabilities, recommendations
}

// checkPayloadVulnerabilities checks for payload-related security issues
func (jst *JWTSecurityTool) checkPayloadVulnerabilities(payload map[string]interface{}) ([]Vulnerability, []string) {
	var vulnerabilities []Vulnerability
	var recommendations []string

	// Check expiration time
	if exp, ok := payload["exp"]; ok {
		if expFloat, ok := exp.(float64); ok {
			expirationTime := time.Unix(int64(expFloat), 0)
			if time.Now().After(expirationTime) {
				vulnerabilities = append(vulnerabilities, Vulnerability{
					Type:        "EXPIRED_TOKEN",
					Description: "JWT token has expired",
					Impact:      "Token should not be accepted",
					Severity:    "HIGH",
				})
			} else if expirationTime.Sub(time.Now()) > 24*time.Hour {
				vulnerabilities = append(vulnerabilities, Vulnerability{
					Type:        "LONG_EXPIRATION",
					Description: "JWT has very long expiration time",
					Impact:      "Increases token hijacking window",
					Severity:    "MEDIUM",
				})
				recommendations = append(recommendations, "Use shorter token expiration times (recommended: 15-60 minutes)")
			}
		}
	} else {
		vulnerabilities = append(vulnerabilities, Vulnerability{
			Type:        "NO_EXPIRATION",
			Description: "JWT does not have expiration time",
			Impact:      "Token valid indefinitely if compromised",
			Severity:    "HIGH",
		})
		recommendations = append(recommendations, "Always include 'exp' claim in JWT")
	}

	// Check for sensitive information in payload
	sensitiveFields := []string{"password", "secret", "key", "token", "ssn", "credit_card"}
	for field := range payload {
		for _, sensitive := range sensitiveFields {
			if strings.Contains(strings.ToLower(field), sensitive) {
				vulnerabilities = append(vulnerabilities, Vulnerability{
					Type:        "SENSITIVE_DATA",
					Description: fmt.Sprintf("Potentially sensitive field '%s' in JWT payload", field),
					Impact:      "Sensitive information exposed in JWT",
					Severity:    "MEDIUM",
				})
				recommendations = append(recommendations, "Avoid storing sensitive information in JWT payload")
				break
			}
		}
	}

	return vulnerabilities, recommendations
}

// checkSignatureBypass checks for signature bypass vulnerabilities
func (jst *JWTSecurityTool) checkSignatureBypass(jwt *JWT) ([]Vulnerability, []string) {
	var vulnerabilities []Vulnerability
	var recommendations []string

	// Check for algorithm confusion
	if alg, ok := jwt.Header["alg"].(string); ok {
		if alg != strings.ToUpper(alg) && alg != strings.ToLower(alg) {
			vulnerabilities = append(vulnerabilities, Vulnerability{
				Type:        "ALGORITHM_CASE_SENSITIVITY",
				Description: "JWT algorithm has mixed case which might bypass validation",
				Impact:      "Potential signature bypass",
				Severity:    "MEDIUM",
			})
			recommendations = append(recommendations, "Implement case-sensitive algorithm validation")
		}
	}

	// Check for empty signature
	if jwt.Signature == "" {
		vulnerabilities = append(vulnerabilities, Vulnerability{
			Type:        "EMPTY_SIGNATURE",
			Description: "JWT has empty signature",
			Impact:      "Token can be forged without signature",
			Severity:    "CRITICAL",
		})
		recommendations = append(recommendations, "Reject tokens with empty signatures")
	}

	return vulnerabilities, recommendations
}

// isWeakSecret checks if JWT is signed with a weak/common secret
func (jst *JWTSecurityTool) isWeakSecret(jwt *JWT) bool {
	parts := strings.Split(jwt.Raw, ".")
	if len(parts) != 3 {
		return false
	}

	headerPayload := parts[0] + "." + parts[1]
	expectedSignature := parts[2]

	for _, secret := range jst.CommonSecrets {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(headerPayload))
		calculatedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		
		if calculatedSignature == expectedSignature {
			return true
		}
	}

	return false
}

// CreateJWT creates a new JWT with specified claims
func (jst *JWTSecurityTool) CreateJWT(header map[string]interface{}, payload map[string]interface{}, secret string) (string, error) {
	// Encode header
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)

	// Encode payload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Create signature
	headerPayload := encodedHeader + "." + encodedPayload
	
	var signature string
	if alg, ok := header["alg"].(string); ok {
		switch alg {
		case "none":
			signature = ""
		case "HS256":
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write([]byte(headerPayload))
			signature = base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
		default:
			return "", fmt.Errorf("unsupported algorithm: %s", alg)
		}
	}

	return headerPayload + "." + signature, nil
}

// PrintJWTInfo prints detailed information about a JWT
func (jst *JWTSecurityTool) PrintJWTInfo(jwt *JWT) {
	fmt.Println("\n=== JWT Analysis ===")
	
	// Header information
	fmt.Println("\n--- Header ---")
	headerBytes, _ := json.MarshalIndent(jwt.Header, "", "  ")
	fmt.Println(string(headerBytes))

	// Payload information
	fmt.Println("\n--- Payload ---")
	payloadBytes, _ := json.MarshalIndent(jwt.Payload, "", "  ")
	fmt.Println(string(payloadBytes))

	// Signature information
	fmt.Println("\n--- Signature ---")
	fmt.Printf("Signature: %s\n", jwt.Signature)
	fmt.Printf("Length: %d characters\n", len(jwt.Signature))

	// Timestamp analysis
	if exp, ok := jwt.Payload["exp"]; ok {
		if expFloat, ok := exp.(float64); ok {
			expirationTime := time.Unix(int64(expFloat), 0)
			fmt.Printf("Expires: %s\n", expirationTime.Format(time.RFC3339))
			if time.Now().After(expirationTime) {
				fmt.Println("Status: EXPIRED")
			} else {
				fmt.Printf("Valid for: %s\n", time.Until(expirationTime).String())
			}
		}
	}

	if iat, ok := jwt.Payload["iat"]; ok {
		if iatFloat, ok := iat.(float64); ok {
			issuedTime := time.Unix(int64(iatFloat), 0)
			fmt.Printf("Issued: %s\n", issuedTime.Format(time.RFC3339))
		}
	}
}

// PrintVulnerabilityReport prints the security analysis report
func (jst *JWTSecurityTool) PrintVulnerabilityReport(report VulnerabilityReport) {
	fmt.Println("\n=== Security Analysis Report ===")
	fmt.Printf("Overall Severity: %s\n", report.Severity)
	
	if len(report.Vulnerabilities) == 0 {
		fmt.Println("✓ No critical vulnerabilities found")
	} else {
		fmt.Printf("\n--- Vulnerabilities Found (%d) ---\n", len(report.Vulnerabilities))
		for i, vuln := range report.Vulnerabilities {
			fmt.Printf("\n%d. %s (%s)\n", i+1, vuln.Type, vuln.Severity)
			fmt.Printf("   Description: %s\n", vuln.Description)
			fmt.Printf("   Impact: %s\n", vuln.Impact)
			if vuln.CVE != "" {
				fmt.Printf("   CVE: %s\n", vuln.CVE)
			}
		}
	}

	if len(report.Recommendations) > 0 {
		fmt.Println("\n--- Security Recommendations ---")
		for i, rec := range report.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}
	}
}

// Interactive menu
func showMenu() {
	fmt.Println("\n=== JWT Security Tool ===")
	fmt.Println("1. Analyze JWT token")
	fmt.Println("2. Create JWT token")
	fmt.Println("3. Test signature bypass")
	fmt.Println("4. Brute force weak secrets")
	fmt.Println("5. Exit")
	fmt.Print("Select option: ")
}

func main() {
	tool := NewJWTSecurityTool()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("JWT Security Analysis Tool v1.0")
	fmt.Println("Educational tool for JWT security testing")
	fmt.Println("Use responsibly and only on applications you own or have permission to test")

	for {
		showMenu()
		
		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Enter JWT token: ")
			if !scanner.Scan() {
				continue
			}
			token := strings.TrimSpace(scanner.Text())

			jwt, err := tool.ParseJWT(token)
			if err != nil {
				fmt.Printf("Error parsing JWT: %v\n", err)
				continue
			}

			tool.PrintJWTInfo(jwt)
			report := tool.AnalyzeSecurity(jwt)
			tool.PrintVulnerabilityReport(report)

		case "2":
			header := map[string]interface{}{
				"alg": "HS256",
				"typ": "JWT",
			}

			payload := map[string]interface{}{
				"sub":  "1234567890",
				"name": "John Doe",
				"iat":  time.Now().Unix(),
				"exp":  time.Now().Add(time.Hour).Unix(),
			}

			fmt.Print("Enter secret key: ")
			if !scanner.Scan() {
				continue
			}
			secret := strings.TrimSpace(scanner.Text())

			token, err := tool.CreateJWT(header, payload, secret)
			if err != nil {
				fmt.Printf("Error creating JWT: %v\n", err)
				continue
			}

			fmt.Printf("Created JWT: %s\n", token)

		case "3":
			fmt.Print("Enter JWT token: ")
			if !scanner.Scan() {
				continue
			}
			token := strings.TrimSpace(scanner.Text())

			// Test algorithm confusion
			fmt.Println("\nTesting signature bypass techniques...")

			// Test with 'none' algorithm
			parts := strings.Split(token, ".")
			if len(parts) == 3 {
				noneHeader := map[string]interface{}{
					"alg": "none",
					"typ": "JWT",
				}
				headerBytes, _ := json.Marshal(noneHeader)
				encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
				noneToken := encodedHeader + "." + parts[1] + "."
				fmt.Printf("None algorithm bypass: %s\n", noneToken)
			}

		case "4":
			fmt.Print("Enter JWT token: ")
			if !scanner.Scan() {
				continue
			}
			token := strings.TrimSpace(scanner.Text())

			jwt, err := tool.ParseJWT(token)
			if err != nil {
				fmt.Printf("Error parsing JWT: %v\n", err)
				continue
			}

			fmt.Println("Testing common weak secrets...")
			if tool.isWeakSecret(jwt) {
				fmt.Println("⚠️  JWT signed with weak/common secret!")
			} else {
				fmt.Println("✓ No common weak secrets found")
			}

		case "5":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
