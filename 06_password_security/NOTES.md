# Password Security Tool - Technical Notes

**Contributor:** vatsalgupta2004

## Cryptographic Security Principles

### 1. Entropy (Information Theory)

**Shannon Entropy Formula:**
```
H(X) = -Σ p(xi) × log₂(p(xi))
```

Where:
- H(X) = entropy in bits
- p(xi) = probability of character xi
- Σ = sum over all unique characters

**Example:**
- Password: "aaaaaa" → Low entropy (~0 bits)
- Password: "aB3!xY" → High entropy (~35 bits)

### 2. Password Strength Calculation

**Theoretical Keyspace:**
```
N = L^C

Where:
  N = total possible combinations
  L = character set size
  C = password length
```

**Character Set Sizes:**
- Lowercase only: 26
- + Uppercase: 52
- + Numbers: 62
- + Special chars: 94

**Example Calculations:**

| Length | Charset | Combinations | Bits |
|--------|---------|--------------|------|
| 8 | 94 | 6.1 × 10¹⁵ | 52.6 |
| 12 | 94 | 4.8 × 10²³ | 79.0 |
| 16 | 94 | 3.9 × 10³¹ | 105.3 |

### 3. Brute Force Attack Resistance

**Attack Speed Assumptions:**
- Modern GPU: ~10 billion passwords/second
- Specialized hardware: ~100 billion passwords/second
- Distributed attack: ~1 trillion passwords/second

**Time to Crack:**
```
Time = Total Combinations / Attempts per Second
```

### 4. Common Attack Vectors

**Dictionary Attack:**
- Uses list of common passwords
- Very fast (milliseconds to minutes)
- Success rate: 10-30% of passwords

**Brute Force:**
- Tries all possible combinations
- Time increases exponentially with length
- Most effective against short passwords

**Hybrid Attack:**
- Dictionary words + common substitutions
- Example: "password" → "p@ssw0rd"
- More effective than pure dictionary

**Rainbow Table:**
- Pre-computed hash tables
- Fast for unsalted hashes
- Defeated by salt + strong hashing

### 5. NIST Guidelines (SP 800-63B)

**Requirements:**
- ✅ Minimum 8 characters (longer recommended)
- ✅ Support all ASCII characters
- ✅ Check against breach databases
- ✅ No composition rules (but encourage variety)
- ❌ No periodic password changes
- ❌ No password hints
- ❌ No knowledge-based authentication

### 6. Password Hashing Best Practices

**Recommended Algorithms:**
- **Argon2id** (best choice)
- **bcrypt** (widely supported)
- **scrypt** (memory-hard)
- **PBKDF2** (acceptable minimum)

**Never Use:**
- MD5 (broken)
- SHA-1 (broken)
- SHA-256 without salt (vulnerable)
- Plain text (catastrophic)

### 7. Vulnerability Detection

**Sequential Patterns:**
```python
# Examples:
"abc123"     → Sequential letters + numbers
"qwerty"     → Keyboard pattern
"password1"  → Dictionary + sequential
```

**Pattern Detection Regex:**
```regex
Sequential numbers: (012|123|234|345|456|567|678|789)
Sequential letters: (abc|bcd|cde|def|efg|fgh|ghi|hij)
Repeated chars: (.)\1{2,}
```

### 8. Secure Random Generation

**Python (CORRECT):**
```python
import secrets
password = ''.join(secrets.choice(charset) for _ in range(16))
```

**Python (WRONG - Don't use):**
```python
import random  # NOT cryptographically secure!
password = ''.join(random.choice(charset) for _ in range(16))
```

**Go (CORRECT):**
```go
import "crypto/rand"
num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
```

**Go (WRONG - Don't use):**
```go
import "math/rand"  // NOT cryptographically secure!
rand.Intn(len(chars))
```

### 9. Password Storage Security

**Secure Storage Flow:**
```
1. Generate random salt (16+ bytes)
2. Combine password + salt
3. Hash with strong algorithm (Argon2id)
4. Store: salt + hash (never plain password)
```

**Verification:**
```
1. Retrieve stored salt + hash
2. Hash provided password + salt
3. Compare hashes (constant-time comparison)
```

### 10. Real-World Statistics

**Most Common Passwords (2024):**
1. "123456" (used by millions)
2. "password"
3. "123456789"
4. "12345678"
5. "qwerty"

**Average Password Statistics:**
- Average length: 9.6 characters
- 52% use only lowercase + numbers
- 34% reuse passwords across sites
- 23% use dictionary words

### 11. Strength Scoring Algorithm

**Our Implementation:**
```
Total Score = Length Score + Variety Score + Entropy Score - Vulnerabilities

Length Score (0-30):
  16+ chars: 30 points
  12-15 chars: 25 points
  8-11 chars: 15 points
  6-7 chars: 5 points

Variety Score (0-25):
  Each charset (lower/upper/numbers/special): +6.25 points

Entropy Score (0-25):
  80+ bits: 25 points
  60-79 bits: 20 points
  40-59 bits: 10 points
  20-39 bits: 5 points

Vulnerability Penalty (0-20):
  Common password: -10 points
  Sequential patterns: -2 points each
  Keyboard patterns: -3 points each
  Repeated chars: -2 points each
  Common words: -2 points each
```

### 12. Time to Crack Examples

**8-character passwords:**
| Type | Time |
|------|------|
| Numbers only | ~2.5 hours |
| Lowercase only | ~2 days |
| Mixed case | ~1 month |
| + Numbers | ~2 months |
| + Special | ~7 months |

**12-character passwords:**
| Type | Time |
|------|------|
| Numbers only | ~316 years |
| Lowercase only | ~16,000 years |
| Mixed case | ~300,000 years |
| + Numbers | ~2 million years |
| + Special | ~200 million years |

### 13. Implementation Notes

**Python Advantages:**
- Rich string manipulation
- Easy regex patterns
- Interactive CLI
- Extensive standard library

**Go Advantages:**
- Better performance
- Concurrent generation
- Static typing
- Easy distribution

**Both Implementations:**
- Use cryptographically secure RNG
- Follow same scoring algorithm
- Detect same vulnerabilities
- Provide similar analysis

### 14. Testing Recommendations

**Test Cases:**
```python
# Weak passwords
"password"     → Should score < 20
"123456"       → Should score < 10
"qwerty"       → Should score < 15

# Medium passwords
"MyPass123"    → Should score 40-60
"SecureP@ss"   → Should score 50-70

# Strong passwords
"K7#mQ9@vL2$pR5!x"  → Should score 80+
"correct-horse-battery-staple"  → Should score 70+
```

### 15. Security Considerations

**DO:**
✅ Use this tool for analysis and generation
✅ Use password managers for storage
✅ Enable 2FA/MFA everywhere
✅ Use unique passwords per site
✅ Use passphrases for memorable passwords

**DON'T:**
❌ Store passwords in plain text
❌ Transmit passwords over HTTP
❌ Email passwords
❌ Share passwords
❌ Write passwords on paper/sticky notes

### 16. Future Research Areas

- Machine learning for pattern detection
- Integration with breach databases (HIBP API)
- Behavioral analysis (typing patterns)
- Biometric augmentation
- Quantum-resistant algorithms

### 17. References

1. **NIST SP 800-63B**: Digital Identity Guidelines
2. **OWASP ASVS**: Application Security Verification Standard
3. **Shannon, C.E.**: "A Mathematical Theory of Communication" (1948)
4. **Bonneau, J.**: "The Science of Guessing" (2012)
5. **Argon2 RFC**: RFC 9106 (2021)

---

**Author:** vatsalgupta2004  
**License:** MIT  
**Created for:** Hacktoberfest 2025  
**Repository:** [Cybersecurity Lab Codes](https://github.com/arya2004/cybersecurity)
