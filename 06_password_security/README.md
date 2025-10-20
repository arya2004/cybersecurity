# Password Security Tool

**Contributor:** vatsalgupta2004

A comprehensive password security toolkit that analyzes password strength, detects vulnerabilities, and generates cryptographically secure passwords.

## Features

### 🔐 Password Strength Analyzer
- **Entropy Calculation**: Measures password randomness and unpredictability
- **Length Analysis**: Checks if password meets minimum security requirements
- **Character Variety**: Detects use of uppercase, lowercase, numbers, and special characters
- **Common Pattern Detection**: Identifies sequential characters, repeated patterns, and keyboard walks
- **Dictionary Attack Detection**: Checks against common weak passwords
- **Strength Score**: Provides overall security rating (0-100)

### 🎲 Secure Password Generator
- **Cryptographically Secure**: Uses `secrets` module for true randomness
- **Customizable Length**: Generate passwords from 8 to 128 characters
- **Character Set Options**: Include/exclude uppercase, lowercase, numbers, special characters
- **Memorable Passwords**: Generate passphrases using random words
- **Multiple Passwords**: Batch generation capability
- **Excludes Ambiguous Characters**: Optional exclusion of similar-looking characters (O/0, l/1)

### 🛡️ Vulnerability Detection
- **Breach Database Check**: Warns about commonly compromised passwords
- **Keyboard Pattern Detection**: Identifies patterns like "qwerty", "asdf"
- **Sequential Characters**: Detects "abc", "123", etc.
- **Repeated Characters**: Flags excessive character repetition
- **Personal Information**: Warns against using common names, dates

## Implementation

### Python Version (`password_security.py`)
- Full-featured implementation with interactive CLI
- Color-coded output for better readability
- Comprehensive vulnerability checks
- Secure random generation using `secrets` module
- Entropy calculation using Shannon's formula

### Go Version (`main.go`)
- High-performance implementation
- Concurrent password generation
- Compatible with Go's `crypto/rand` for security
- CLI interface with flags
- Efficient string operations

## Usage

### Python
```bash
python password_security.py
```

**Options:**
1. Analyze password strength
2. Generate secure password
3. Generate passphrase
4. Batch generate passwords
5. Check password against common breaches

### Go
```bash
go run main.go
```

**Command-line flags:**
```bash
# Analyze password
go run main.go -analyze "YourPassword123!"

# Generate password
go run main.go -generate -length 16

# Generate with specific characters
go run main.go -generate -length 20 -uppercase -lowercase -numbers -special

# Generate passphrase
go run main.go -passphrase -words 5
```

## Security Concepts Demonstrated

### 1. **Password Entropy**
```
Entropy = L × log₂(N)
where:
  L = password length
  N = size of character set
```

Higher entropy = harder to crack via brute force.

### 2. **Character Sets**
- Lowercase (a-z): 26 characters
- Uppercase (A-Z): 26 characters
- Numbers (0-9): 10 characters
- Special (!@#$%^&*): ~32 characters

### 3. **Strength Criteria**
- **Weak** (0-40): Easily crackable within minutes
- **Moderate** (41-60): Crackable with significant resources
- **Strong** (61-80): Difficult to crack with current technology
- **Very Strong** (81-100): Practically uncrackable with brute force

### 4. **Common Vulnerabilities**
- Dictionary words
- Keyboard patterns (qwerty, asdf)
- Sequential patterns (abc, 123)
- Repeated characters (aaa, 111)
- Common substitutions (@ for a, 3 for e)

## Example Output

### Password Analysis
```
🔍 Analyzing Password: MyP@ssw0rd123!
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Strength Score: 65/100 - STRONG

✅ Length: 14 characters (Good)
✅ Uppercase letters: Present
✅ Lowercase letters: Present
✅ Numbers: Present
✅ Special characters: Present
⚠️  Contains common word: 'password'
⚠️  Sequential characters detected
🔐 Entropy: 82.3 bits

Estimated Time to Crack:
  • Brute Force: ~3.2 years
  • Dictionary Attack: ~5 minutes
```

### Password Generation
```
🎲 Generated Secure Password:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Password: K7#mQ9@vL2$pR5!x
Length: 16 characters
Entropy: 95.4 bits
Strength: VERY STRONG (93/100)
Time to Crack: ~12,000 years
```

## Educational Value

This tool demonstrates:
1. **Cryptographic Security**: Using secure random number generators
2. **Entropy Theory**: Measuring password unpredictability
3. **Pattern Recognition**: Detecting common security weaknesses
4. **String Manipulation**: Efficient text processing
5. **Security Best Practices**: NIST password guidelines compliance

## Requirements

### Python
```bash
pip install colorama  # For colored output (optional)
```

No external dependencies required for core functionality.

### Go
```bash
go mod init password_security
```

Standard library only - no external dependencies.

## Testing

### Test Cases Included
- Weak passwords: "password", "123456", "qwerty"
- Medium passwords: "MyPassword123", "SecureP@ss"
- Strong passwords: "K7#mQ9@vL2$pR5!x", "Tr0ub4dor&3"
- Passphrases: "correct horse battery staple"

## Security Considerations

⚠️ **Important Notes:**
- Generated passwords are cryptographically secure
- Tool uses `secrets` (Python) and `crypto/rand` (Go) modules
- Never transmit passwords over unencrypted connections
- Never store passwords in plain text
- Use password managers for storage
- This tool is for educational purposes

## Real-World Applications

1. **User Registration**: Enforce strong password policies
2. **Security Audits**: Analyze existing passwords
3. **Password Managers**: Generate secure credentials
4. **Security Training**: Teach password best practices
5. **Compliance**: Meet security standard requirements

## Time to Crack Estimates

| Password Type | Length | Character Set | Brute Force Time |
|--------------|--------|---------------|------------------|
| Numbers only | 8 | 10 | ~2.5 hours |
| Lowercase | 8 | 26 | ~2 days |
| Mixed case | 8 | 52 | ~1 month |
| + Numbers | 8 | 62 | ~2 months |
| + Special | 8 | 94 | ~7 months |
| + Special | 12 | 94 | ~6,000 years |
| + Special | 16 | 94 | ~1.5 million years |

*Assumes 10 billion attempts per second

## NIST Guidelines Compliance

This tool follows NIST SP 800-63B recommendations:
- ✅ Minimum length of 8 characters
- ✅ No complexity requirements (but encouraged)
- ✅ Check against breach databases
- ✅ No periodic password changes required
- ✅ Allow all printable ASCII characters

## Future Enhancements

- [ ] Integration with Have I Been Pwned API
- [ ] Password policy configuration
- [ ] Multi-language support
- [ ] Web interface
- [ ] Mobile app version
- [ ] Password strength meter visualization
- [ ] Custom dictionary support
- [ ] Password history tracking

## References

- [NIST SP 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html)
- [OWASP Password Guidelines](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [Shannon Entropy](https://en.wikipedia.org/wiki/Entropy_(information_theory))
- [Password Strength](https://en.wikipedia.org/wiki/Password_strength)

## License

MIT License - See main repository LICENSE file

## Author

**Contributor:** vatsalgupta2004  
**Repository:** [Cybersecurity Lab Codes](https://github.com/arya2004/cybersecurity)  
**Created for:** Hacktoberfest 2025
