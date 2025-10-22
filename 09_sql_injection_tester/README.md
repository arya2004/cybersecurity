# SQL Injection Testing Tool

A comprehensive SQL injection vulnerability testing tool written in Go for educational and authorized security testing purposes.

## Features

- **Multiple Injection Types**: Error-based, Boolean-based, Time-based, Union-based testing
- **Database Support**: MySQL, MSSQL, PostgreSQL, Oracle, SQLite detection
- **Advanced Payloads**: 50+ specialized payloads for different scenarios
- **Automated Testing**: Comprehensive parameter testing with confidence scoring
- **Evidence Collection**: Detailed vulnerability evidence and proof-of-concept
- **Educational Focus**: Learning tool for understanding SQL injection mechanics

## Supported Injection Types

### Error-Based SQL Injection
- Database error message extraction
- Version fingerprinting through errors
- Information schema exploitation
- Custom error pattern detection

### Boolean-Based SQL Injection
- True/false condition testing
- Blind SQL injection detection
- Response comparison analysis
- Logic-based vulnerability assessment

### Time-Based SQL Injection
- Delay function testing (SLEEP, WAITFOR)
- Response time analysis
- Blind injection through timing
- Database-specific time functions

### Union-Based SQL Injection
- Column enumeration
- Data extraction through UNION
- NULL value testing
- Information disclosure

## Database Coverage

- **MySQL**: Error patterns, functions, and syntax
- **Microsoft SQL Server**: WAITFOR, error messages, functions
- **PostgreSQL**: pg_sleep, error patterns, syntax
- **Oracle**: ORA error codes, functions
- **SQLite**: Error patterns and syntax
- **Generic**: Common SQL patterns and errors

## Payload Categories

### Authentication Bypass
```sql
admin'--
admin'/*
' OR '1'='1
```

### Information Gathering
```sql
' AND @@version IS NOT NULL--
' AND version() IS NOT NULL--
' AND user() IS NOT NULL--
```

### Error Generation
```sql
'
"
' AND EXTRACTVALUE(1, CONCAT(0x7e, (SELECT version()), 0x7e))--
```

### Time Delays
```sql
'; WAITFOR DELAY '00:00:05'--
'; SELECT SLEEP(5)--
'; SELECT pg_sleep(5)--
```

## Usage

```bash
go run main.go
```

### Testing Options

1. **Single URL Test**: Test one parameter on a target URL
2. **Multiple Parameters**: Test multiple parameters simultaneously
3. **Payload Database**: View all available payloads and statistics
4. **Custom Payload**: Test custom SQL injection payloads

### Example Usage

```
Enter target URL: http://example.com/login.php?user=admin&pass=password
Enter parameter name to test: user

Starting SQL injection test on http://example.com/login.php?user=admin&pass=password (parameter: user)
Progress: 1/50 - Testing: ERROR_BASED
Progress: 2/50 - Testing: BOOLEAN_BASED
...
```

## Vulnerability Detection

### High Confidence Indicators
- Database-specific error messages
- Successful time-based delays
- Boolean logic manipulation
- Union query success

### Medium Confidence Indicators
- Generic SQL error patterns
- Response size variations
- Timeout behaviors
- Content differences

### Evidence Collection
- Error message extraction
- Response time measurements
- Content comparison analysis
- Pattern matching results

## Example Output

```
🚨 FOUND 2 POTENTIAL SQL INJECTION VULNERABILITIES:
============================================================

[1] VULNERABILITY DETECTED
URL: http://example.com/search.php?q=test
Parameter: q
Payload: ' OR '1'='1
Type: BOOLEAN_BASED
Database: GENERIC
Risk Level: HIGH
Confidence: HIGH
Error Type: BOOLEAN_DIFFERENCE
Evidence:
  - True condition response differs from baseline
  - False condition response differs from true condition
Description: Classic boolean-based SQL injection

[2] VULNERABILITY DETECTED
URL: http://example.com/search.php?q=test
Parameter: q
Payload: '
Type: ERROR_BASED
Database: MYSQL
Risk Level: HIGH
Confidence: HIGH
Error Type: MYSQL_ERROR
Evidence:
  - Found MYSQL error pattern: You have an error in your SQL syntax
Description: Single quote to trigger SQL error
```

## Security Recommendations

The tool provides comprehensive security recommendations:

1. **Use Parameterized Queries**: Prevent injection through prepared statements
2. **Input Validation**: Implement strict input validation and sanitization
3. **Stored Procedures**: Use stored procedures where appropriate
4. **Least Privilege**: Apply minimal database permissions
5. **Error Handling**: Implement custom error pages
6. **WAF Protection**: Deploy web application firewalls
7. **Regular Testing**: Conduct ongoing security assessments

## Educational Value

This tool demonstrates:
- SQL injection attack vectors and techniques
- Database-specific vulnerability patterns
- Proper vulnerability assessment methodology
- Security testing best practices
- Defensive programming concepts

## Technical Implementation

- **HTTP Client**: Robust HTTP request handling with timeouts
- **Pattern Matching**: Regular expression-based error detection
- **Payload Management**: Structured payload database with metadata
- **Response Analysis**: Advanced response comparison algorithms
- **Evidence Collection**: Comprehensive vulnerability proof collection

## Legal and Ethical Use

⚠️ **IMPORTANT**: This tool is designed for:
- Educational purposes
- Authorized security testing
- Applications you own
- Penetration testing with explicit permission

**DO NOT** use this tool for:
- Unauthorized testing
- Malicious attacks
- Systems you don't own
- Illegal activities

## Performance Features

- **Concurrent Testing**: Efficient parallel request processing
- **Timeout Management**: Configurable request timeouts
- **Progress Tracking**: Real-time testing progress indication
- **Result Caching**: Baseline response caching for comparison
