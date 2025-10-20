# Password Strength Checker

A comprehensive password strength analyzer written in Go that evaluates password security using entropy calculation, pattern detection, and crack time estimation.

## 📋 Description

This tool performs detailed analysis of password strength using multiple security metrics including character composition, entropy calculation, common password detection, dictionary word checking, and pattern recognition. It provides actionable suggestions for creating stronger passwords.

## ✨ Features

### 🔍 Comprehensive Analysis
- **Character Composition Check**: Lowercase, uppercase, digits, special characters
- **Entropy Calculation**: Shannon entropy measurement in bits
- **Length Analysis**: Optimal length recommendations
- **Pattern Detection**: Identifies common patterns and sequences
- **Dictionary Check**: Detects common words
- **Common Password Detection**: Compares against known weak passwords

### 📊 Strength Metrics
- **Score System**: 0-100 point scoring
- **Strength Levels**: Very Weak, Weak, Medium, Strong, Very Strong
- **Crack Time Estimation**: Estimates time to crack with modern hardware
- **Visual Progress Bar**: ASCII-based strength visualization

### 💡 Security Recommendations
- Personalized improvement suggestions
- Best practice guidance
- Pattern avoidance tips
- Character variety recommendations

### 🎨 User Experience
- Color-coded strength indicators
- Masked password display for security
- Interactive analysis mode
- Detailed formatted reports
- Clear visual feedback

## 🚀 Installation

### Prerequisites
- Go 1.16 or higher

### Setup

1. Clone the repository:
```bash
git clone https://github.com/arya2004/cybersecurity.git
cd cybersecurity/password-strength-checker
```

2. No additional dependencies required (uses Go standard library)

## 💻 Usage

### Run the analyzer
```bash
go run main.go
```

### Interactive Mode
```bash
$ go run main.go
Enter a password to analyze (input is hidden):
Password: YourP@ssw0rd123

[Displays detailed analysis]

Analyze another password? (y/n): y
```

## 📸 Sample Output

```
╔═══════════════════════════════════════╗
║   Password Strength Checker v1.0     ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝

══════════════════════════════════════════════════
PASSWORD ANALYSIS REPORT
══════════════════════════════════════════════════
Password: *************** (Length: 15)
──────────────────────────────────────────────────
Character Composition:
  Lowercase Letters: ✓ Present
  Uppercase Letters: ✓ Present
  Digits: ✓ Present
  Special Characters: ✓ Present
──────────────────────────────────────────────────
Strength Metrics:
  Entropy: 89.54 bits
  Score: 90/100
  Strength: Very Strong
  Estimated Crack Time: 284 centuries
──────────────────────────────────────────────────
Strength: [██████████████████░░] 90%
══════════════════════════════════════════════════
```

### Weak Password Example

```
══════════════════════════════════════════════════
Password: ******** (Length: 8)
──────────────────────────────────────────────────
⚠️  WARNING: This is a commonly used password!
⚠️  WARNING: Contains dictionary words!
⚠️  Patterns Detected:
     - Sequential numbers detected
──────────────────────────────────────────────────
💡 Suggestions for Improvement:
  1. Increase length to at least 12 characters
  2. Add special characters (!@#$%^&*)
  3. Avoid common passwords
  4. Avoid dictionary words
  5. Avoid predictable patterns
──────────────────────────────────────────────────
Strength: [███░░░░░░░░░░░░░░░░░] 15%
══════════════════════════════════════════════════
```

## 🔐 How It Works

### 1. Character Composition Analysis
Checks for presence of:
- Lowercase letters (a-z)
- Uppercase letters (A-Z)
- Digits (0-9)
- Special characters (!@#$%^&*)

### 2. Entropy Calculation
```
Entropy = Length × log2(Character Set Size)
```

Character set sizes:
- Lowercase only: 26
- + Uppercase: 52
- + Digits: 62
- + Special chars: 94

### 3. Pattern Detection
Identifies:
- Sequential numbers (123, 456)
- Repeated characters (aaa, 111)
- Keyboard patterns (qwerty, asdf)
- Year patterns (1990, 2024)

### 4. Scoring Algorithm
Maximum 100 points from:
- **Length (30 points)**
  - 12+ chars: 30 pts
  - 8-11 chars: 20 pts
  - 6-7 chars: 10 pts
  
- **Character Variety (40 points)**
  - Lowercase: 10 pts
  - Uppercase: 10 pts
  - Digits: 10 pts
  - Special: 10 pts

- **Entropy (20 points)**
  - 60+ bits: 20 pts
  - 40-59 bits: 15 pts
  - 28-39 bits: 10 pts

- **Deductions**
  - Common password: -40 pts
  - Dictionary word: -10 pts
  - Each pattern: -10 pts

### 5. Crack Time Estimation
Assumes 1 billion guesses/second (modern GPU):
```
Time = 2^entropy / (2 × 10^9 guesses/sec)
```

## 📊 Strength Levels

| Score | Level | Color | Description |
|-------|-------|-------|-------------|
| 80-100 | Very Strong | Green | Excellent security |
| 60-79 | Strong | Cyan | Good security |
| 40-59 | Medium | Yellow | Acceptable |
| 20-39 | Weak | Red | Poor security |
| 0-19 | Very Weak | Magenta | Unacceptable |

## 🛡️ Security Best Practices

### Strong Password Guidelines
✅ **Minimum 12 characters** (16+ recommended)  
✅ **Mix character types** (upper, lower, digits, special)  
✅ **Avoid common words** or names  
✅ **Use unique passwords** for each account  
✅ **Consider passphrases** (e.g., "Coffee!Morning@2024#Smile")  
✅ **Use password managers** (LastPass, 1Password, Bitwarden)

### What to Avoid
❌ Personal information (birthdays, names)  
❌ Dictionary words  
❌ Sequential patterns (12345, abcde)  
❌ Keyboard patterns (qwerty, asdfgh)  
❌ Repeated characters (aaaa, 1111)  
❌ Common substitutions (p@ssw0rd)

## 🎓 Educational Use Cases

### Cybersecurity Training
- Password policy education
- Security awareness programs
- User training workshops
- Best practice demonstrations

### Academic Labs
- Cryptography courses
- Information security classes
- Computer science projects
- Security research

## 🔧 Customization

### Add Custom Common Passwords
```go
var commonPasswords = []string{
    "password", "123456", "YourCustomPassword",
}
```

### Add Custom Dictionary Words
```go
var commonWords = []string{
    "love", "admin", "YourCustomWord",
}
```

### Adjust Scoring Weights
Modify the `calculateScore` function to change scoring criteria.

## 📈 Algorithm Details

### Entropy Formula
```
H = L × log₂(R)

Where:
H = Entropy (bits)
L = Password length
R = Character set size
```

### Crack Time Formula
```
T = 2^H / (2 × G)

Where:
T = Time (seconds)
H = Entropy (bits)
G = Guesses per second (10^9)
```

## 💡 Example Passwords

### Very Weak
- `password` (Score: 0)
- `123456` (Score: 5)
- `qwerty` (Score: 10)

### Weak
- `Password1` (Score: 25)
- `welcome123` (Score: 30)

### Medium
- `MyP@ss2024` (Score: 50)
- `Hello!World9` (Score: 55)

### Strong
- `C0ff33&M0rning!` (Score: 70)
- `MyStr0ng#Pass24` (Score: 75)

### Very Strong
- `T!ger$Run@2024#Fast` (Score: 90)
- `Blu3Sky!Moon&Star$99` (Score: 95)

## 🔮 Future Enhancements

Potential additions:
- Password generation feature
- Passphrase generator
- Breach database integration (Have I Been Pwned API)
- Multi-language support
- Web interface
- Password history tracking
- Compliance checking (NIST, ISO standards)
- Zxcvbn algorithm integration

## 👨‍💻 Author

**Ashvin**
- GitHub: [@ashvin2005](https://github.com/ashvin2005)
- LinkedIn: [ashvin-tiwari](https://linkedin.com/in/ashvin-tiwari)

## 🎃 Hacktoberfest 2025

Created as part of Hacktoberfest 2025 contributions to the Cybersecurity Lab Codes repository.

## 📄 License

MIT License (same as parent repository)

## 🙏 Acknowledgments

- NIST Digital Identity Guidelines
- OWASP Password Storage Cheat Sheet
- zxcvbn password strength estimator

## 📚 References

- [NIST SP 800-63B](https://pages.nist.gov/800-63-3/sp800-63b.html)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [Have I Been Pwned](https://haveibeenpwned.com/)

---

**Secure your passwords, secure your life!** 🔐💪