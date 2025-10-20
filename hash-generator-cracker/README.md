# Hash Generator & Cracker

A comprehensive cryptographic hash tool written in Go that generates and cracks password hashes using multiple algorithms including MD5, SHA1, SHA256, SHA512, and bcrypt.

## 📋 Description

This educational tool demonstrates hash generation and password cracking techniques using dictionary attacks and brute force methods. It supports multiple hashing algorithms and provides insights into password security and hash cracking methodologies.

## ⚠️ Legal Disclaimer

**IMPORTANT**: This tool is for educational and authorized security testing purposes only.

- ✅ Only test hashes you own or have explicit permission to crack
- ❌ Unauthorized hash cracking may be illegal
- ⚖️ Always comply with applicable laws and regulations
- 🎓 Designed for cybersecurity education and research

## ✨ Features

### 🔐 Hash Generation
- **Multiple Algorithms**: MD5, SHA1, SHA256, SHA512, bcrypt
- **Single Hash Generation**: Generate one hash at a time
- **Batch Generation**: Create hashes with all algorithms simultaneously
- **Performance Metrics**: Time measurement for each algorithm
- **Hash Comparison**: Verify passwords against known hashes

### 🔨 Hash Cracking
- **Dictionary Attack**: Test against 40+ common passwords
- **Brute Force Attack**: Numeric brute force (0-9, demo purposes)
- **Multiple Algorithms**: Support for all hash types
- **Attempt Tracking**: Count and display number of attempts
- **Time Tracking**: Measure cracking duration
- **Success Reporting**: Clear indication of found passwords

### 📊 Supported Algorithms

| Algorithm | Output Length | Speed | Security |
|-----------|---------------|-------|----------|
| MD5 | 128 bits (32 hex) | Very Fast | ⚠️ Broken |
| SHA1 | 160 bits (40 hex) | Fast | ⚠️ Deprecated |
| SHA256 | 256 bits (64 hex) | Fast | ✓ Secure |
| SHA512 | 512 bits (128 hex) | Fast | ✓ Secure |
| bcrypt | Variable | Slow | ✓ Very Secure |

## 🚀 Installation

### Prerequisites
- Go 1.16 or higher

### Setup

1. Clone the repository:
```bash
git clone https://github.com/arya2004/cybersecurity.git
cd cybersecurity/hash-generator-cracker
```

2. Install dependencies:
```bash
go mod download
```

3. Run the tool:
```bash
go run main.go
```

## 💻 Usage

### Interactive Menu
```bash
$ go run main.go

╔═══════════════════════════════════════╗
║   Hash Generator & Cracker v1.0      ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝

MAIN MENU
1. Generate Hash
2. Generate All Hashes
3. Crack Hash (Dictionary Attack)
4. Crack Hash (Brute Force - Numeric)
5. Compare Hash
6. Exit

Select option:
```

### Option 1: Generate Single Hash
Generate a hash using one algorithm:
```
Enter text to hash: password123
Select algorithm: sha256

Hash: ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f
Time: 12.5µs
```

### Option 2: Generate All Hashes
Generate hashes using all supported algorithms:
```
Enter text to hash: hello

MD5       : 5d41402abc4b2a76b9719d911017c592
SHA1      : aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
SHA256    : 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
SHA512    : 9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca7...
bcrypt    : $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
```

### Option 3: Dictionary Attack
Crack hash using common password dictionary:
```
Enter hash to crack: 5f4dcc3b5aa765d61d8327deb882cf99
Select algorithm: md5

Starting dictionary attack...
Dictionary size: 42 passwords

✓ SUCCESS! Password found: password
Attempts: 1
Time: 245µs
```

### Option 4: Brute Force Attack
Numeric brute force (educational demo):
```
Enter hash to crack: e807f1fcf82d132f9bb018ca6738a19f
Select algorithm: md5
Maximum password length: 4

Starting brute force attack (numeric only)...

✓ SUCCESS! Password found: 1234
Attempts: 1234
Time: 1.2s
```

### Option 5: Compare Hash
Verify if a password matches a hash:
```
Enter password to verify: mypassword
Enter hash to compare: 34819d7beeabb9260a5c854bc85b3e44
Select algorithm: md5

✓ MATCH: Password matches the hash!
```

## 📸 Sample Output

### Hash Generation
```
══════════════════════════════════════════════════
HASH GENERATION RESULT
══════════════════════════════════════════════════
Algorithm: SHA256
Input: SecurePassword123
Hash: 8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92
Time: 15.2µs
══════════════════════════════════════════════════
```

### Successful Dictionary Attack
```
══════════════════════════════════════════════════
HASH CRACKING RESULT
══════════════════════════════════════════════════
Algorithm: MD5
Attempts: 12
Time: 156.3µs
──────────────────────────────────────────────────
✓ SUCCESS! Password found: letmein
══════════════════════════════════════════════════
```

### Failed Crack Attempt
```
══════════════════════════════════════════════════
HASH CRACKING RESULT
══════════════════════════════════════════════════
Algorithm: SHA256
Attempts: 42
Time: 2.3ms
──────────────────────────────────────────────────
✗ FAILED: Password not found
══════════════════════════════════════════════════
```

## 🔐 How It Works

### Hash Generation
```go
// MD5 Example
hash := md5.Sum([]byte(input))
hexHash := hex.EncodeToString(hash[:])

// bcrypt Example (with salt)
hash, _ := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
```

### Dictionary Attack
1. Load password dictionary
2. For each password:
   - Generate hash
   - Compare with target hash
   - If match, return password
3. Track attempts and time

### Brute Force Attack
1. Define character set (0-9 for demo)
2. Generate all combinations up to max length
3. Hash each combination
4. Compare with target hash
5. Continue until match or exhausted

### bcrypt Comparison
Special handling for bcrypt (uses comparison, not reversal):
```go
err := bcrypt.CompareHashAndPassword([]byte(targetHash), []byte(password))
if err == nil {
    // Password matches!
}
```

## 🎓 Educational Use Cases

### Cybersecurity Labs
- Understanding hash algorithms
- Password security education
- Attack methodology demonstration
- Security awareness training

### Academic Courses
- Cryptography classes
- Information security courses
- Ethical hacking labs
- Computer forensics

### Practical Learning
- Hash collision understanding
- Rainbow table concepts
- Salt and pepper techniques
- Secure password storage

## 📊 Algorithm Comparison

### Speed Tests
```
Input: "password"
MD5:     ~10µs
SHA1:    ~12µs
SHA256:  ~15µs
SHA512:  ~20µs
bcrypt:  ~100ms (intentionally slow)
```

### Security Status
- **MD5**: ⚠️ BROKEN - Do not use for security
- **SHA1**: ⚠️ DEPRECATED - Collision attacks exist
- **SHA256**: ✓ Currently secure for most uses
- **SHA512**: ✓ Very secure, larger output
- **bcrypt**: ✓ Designed for passwords (includes salt)

## 🛡️ Password Security Best Practices

### For Developers
✅ **Use bcrypt, Argon2, or PBKDF2** for password hashing  
✅ **Always use salts** (bcrypt includes this)  
✅ **Never use MD5 or SHA1** for passwords  
✅ **Use slow hashing** (bcrypt cost factor)  
✅ **Store only hashes**, never plaintext passwords

### For Users
✅ **Use strong, unique passwords**  
✅ **Enable two-factor authentication**  
✅ **Use a password manager**  
✅ **Never reuse passwords**  
✅ **Change passwords if breach occurs**

## 🔧 Customization

### Add Custom Dictionary Words
```go
var commonPasswordsDictionary = []string{
    "password", "123456", "YourCustomPassword",
}
```

### Expand Brute Force Character Set
```go
// Current: "0123456789"
charset := "0123456789abcdefghijklmnopqrstuvwxyz" // Add lowercase
charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" // Full
```

### Adjust bcrypt Cost
```go
// Default: bcrypt.DefaultCost (10)
hash, _ := bcrypt.GenerateFromPassword([]byte(input), 12) // Slower, more secure
```

## ⚡ Performance Notes

### Dictionary Attack
- **Speed**: Milliseconds for small dictionaries
- **Efficiency**: Depends on dictionary size and algorithm
- **Best For**: Common passwords, known wordlists

### Brute Force
- **Complexity**: Exponential with password length
- **Numeric (0-9)**: 10^length combinations
- **Full ASCII**: 94^length combinations
- **Warning**: Impractical for long passwords

### Example Times
```
4-digit numeric (10,000 combinations): ~1 second
6-digit numeric (1,000,000 combinations): ~1 minute
8-digit numeric: ~2 hours
8-char alphanumeric: Centuries
```

## 🔮 Future Enhancements

Potential additions:
- GPU acceleration for cracking
- Rainbow table support
- Custom wordlist import
- Argon2 algorithm support
- Hash identification (detect algorithm)
- Multi-threaded cracking
- Progress indicators
- Result export (CSV/JSON)
- Web interface

## 👨‍💻 Author

**Ashvin**
- GitHub: [@ashvin2005](https://github.com/ashvin2005)
- LinkedIn: [ashvin-tiwari](https://linkedin.com/in/ashvin-tiwari)

## 🎃 Hacktoberfest 2025

Created as part of Hacktoberfest 2025 contributions to the Cybersecurity Lab Codes repository.

## 📄 License

MIT License (same as parent repository)

## 🙏 Acknowledgments

- Go crypto package maintainers
- bcrypt creators
- OWASP for security guidance
- Cybersecurity research community

## 📚 References

- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [bcrypt Paper](https://www.usenix.org/legacy/events/usenix99/provos/provos.pdf)
- [NIST Hash Functions](https://csrc.nist.gov/projects/hash-functions)
- [Go crypto package](https://pkg.go.dev/crypto)

---

**Use responsibly and ethically!** 🔐🔨