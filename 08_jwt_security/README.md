# JWT Security Analysis Tool

A comprehensive JWT (JSON Web Token) security analysis and testing tool written in Go for educational purposes.

## Features

- **JWT Parsing**: Complete JWT token parsing and structure analysis
- **Security Analysis**: Comprehensive vulnerability detection and assessment
- **Algorithm Testing**: Tests for algorithm confusion and bypass techniques
- **Signature Verification**: Weak secret detection and brute force testing
- **Token Creation**: Create custom JWT tokens for testing
- **Vulnerability Reporting**: Detailed security reports with recommendations

## Supported Vulnerabilities

### Critical Vulnerabilities
- **Algorithm None**: Tokens using 'none' algorithm (CVE-2015-9235)
- **Empty Signature**: Tokens with missing signatures
- **Expired Tokens**: Token expiration validation

### High-Risk Issues
- **Weak Secrets**: Detection of common/weak HMAC secrets
- **No Expiration**: Tokens without expiration claims
- **Sensitive Data**: Detection of sensitive information in payload

### Medium-Risk Issues
- **Algorithm Confusion**: RS256/HS256 key confusion (CVE-2016-10555)
- **Long Expiration**: Tokens with excessive validity periods
- **Case Sensitivity**: Algorithm case manipulation attempts

## Usage

```bash
go run main.go
```

### Menu Options

1. **Analyze JWT Token**: Complete security analysis of existing tokens
2. **Create JWT Token**: Generate custom tokens for testing
3. **Test Signature Bypass**: Attempt common bypass techniques
4. **Brute Force Weak Secrets**: Test for common weak secrets
5. **Exit**: Exit the application

## Example Analysis

```
=== JWT Analysis ===

--- Header ---
{
  "alg": "HS256",
  "typ": "JWT"
}

--- Payload ---
{
  "sub": "1234567890",
  "name": "John Doe",
  "iat": 1516239022,
  "exp": 1516242622
}

--- Signature ---
Signature: SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
Length: 43 characters
Expires: 2018-01-18T01:30:22Z
Status: EXPIRED

=== Security Analysis Report ===
Overall Severity: HIGH

--- Vulnerabilities Found (2) ---

1. EXPIRED_TOKEN (HIGH)
   Description: JWT token has expired
   Impact: Token should not be accepted

2. WEAK_SECRET (HIGH)
   Description: JWT signed with weak or common secret
   Impact: Attackers can forge tokens with common secrets

--- Security Recommendations ---
1. Always include 'exp' claim in JWT
2. Use strong, randomly generated secrets for HMAC signing
3. Implement proper token expiration validation
```

## Common Weak Secrets Tested

The tool tests against common weak secrets including:
- "secret"
- "password" 
- "123456"
- "admin"
- "test"
- "key"
- "jwt"
- "token"
- Empty strings
- And more...

## Security Testing Features

### Algorithm Confusion Testing
Tests for RS256/HS256 algorithm confusion where:
- RSA public key is used as HMAC secret
- Algorithm manipulation bypasses

### None Algorithm Testing
Generates tokens with 'none' algorithm to test for:
- Missing signature verification
- Algorithm validation bypass

### Signature Bypass Techniques
- Empty signature testing
- Algorithm case sensitivity
- Malformed token structure

## Educational Value

This tool demonstrates:
- JWT structure and security components
- Common JWT vulnerabilities and attacks
- Proper JWT validation techniques
- Security best practices for JWT implementation
- Vulnerability assessment methodologies

## Security Considerations

⚠️ **Important**: This tool is for educational and authorized testing purposes only. Use only on applications you own or have explicit permission to test.

## Technical Implementation

- **No External Dependencies**: Pure Go implementation
- **HMAC Verification**: SHA-256 HMAC signature verification
- **Base64 Handling**: Proper URL-safe base64 encoding/decoding
- **JSON Processing**: Robust JSON parsing and validation
- **Error Handling**: Comprehensive error handling and reporting

## Output Formats

The tool provides:
- Structured vulnerability reports
- Severity classifications (CRITICAL, HIGH, MEDIUM, LOW)
- Detailed impact assessments
- Actionable security recommendations
- CVE references where applicable
