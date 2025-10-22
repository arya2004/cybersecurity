# Advanced Cryptography Tool

A comprehensive educational cryptography tool written in Go that demonstrates various cryptographic algorithms, their implementations, vulnerabilities, and security analysis.

## Features

- **Multiple Encryption Algorithms**: AES-256-GCM, DES, 3DES, RSA-OAEP
- **Hash Functions**: MD5, SHA-1, SHA-256, SHA-512
- **Key Generation**: RSA key pair generation with configurable key sizes
- **Cryptographic Analysis**: Weakness detection and security recommendations
- **Key Strength Analysis**: Entropy calculation and character set analysis
- **Timing Attack Demonstration**: Educational timing attack simulation
- **Interactive Interface**: User-friendly menu-driven operation

## Supported Algorithms

### Symmetric Encryption
- **AES-256-GCM**: Advanced Encryption Standard with Galois/Counter Mode
- **DES**: Data Encryption Standard (educational purposes - cryptographically weak)
- **3DES**: Triple DES (legacy algorithm demonstration)

### Asymmetric Encryption
- **RSA-OAEP**: RSA with Optimal Asymmetric Encryption Padding
- **Key Sizes**: Configurable (default 2048-bit)

### Hash Functions
- **SHA-256**: Secure Hash Algorithm 256-bit
- **SHA-512**: Secure Hash Algorithm 512-bit
- **SHA-1**: Secure Hash Algorithm 160-bit (deprecated)
- **MD5**: Message Digest 5 (cryptographically broken)

## Usage

```bash
go run main.go
```

### Menu Options

1. **AES Encryption/Decryption**: 
   - Secure symmetric encryption using AES-256-GCM
   - Automatic nonce generation for security
   - Base64 encoded output for transport

2. **DES Encryption (Educational)**:
   - Demonstrates weak encryption algorithm
   - Shows why DES should not be used in production
   - Includes comprehensive weakness analysis

3. **RSA Key Generation & Encryption**:
   - Generates RSA key pairs
   - RSA-OAEP encryption for security
   - PEM format key export

4. **Hash Functions**:
   - Multiple hash algorithm comparison
   - Identifies weak hash functions
   - Hex-encoded output

5. **Cryptographic Analysis**:
   - Algorithm weakness detection
   - Vulnerability identification
   - Security recommendations

6. **Key Strength Analysis**:
   - Entropy calculation
   - Character set analysis
   - Algorithm-specific strength assessment

7. **Timing Attack Demonstration**:
   - Shows vulnerable vs secure comparison
   - Educational timing attack simulation
   - Performance measurement

## Example Operations

### AES Encryption
```
Enter text: Hello, World!
Enter key: MySecretKey123456
```
Output:
```
✅ Operation successful
Input: Hello, World!
Output: SGVsbG8sIFdvcmxkIQ==...
Metadata:
  key_size: 256
  nonce_size: 12
  mode: GCM
```

### Cryptographic Analysis
```
Enter algorithm to analyze: DES
```
Output:
```
=== Cryptographic Analysis: DES ===
Severity: CRITICAL

🚨 Vulnerabilities:
1. 56-bit key size is too small for modern security
2. Vulnerable to brute force attacks
3. Meet-in-the-middle attacks possible
4. Linear and differential cryptanalysis vulnerabilities

💡 Recommendations:
1. Migrate to AES-256
2. Use 3DES as temporary measure only
3. Implement proper key management
```

### Key Strength Analysis
```
Enter key to analyze: password123
Enter algorithm: AES
```
Output:
```
--- Key Analysis Results ---
length: 11
entropy: 3.45
character_sets: map[digits:true lowercase:true special:false uppercase:false]
strength: WEAK
warning: Key too short for AES
```

## Security Analysis Features

### Algorithm Weakness Detection
- **DES**: 56-bit key vulnerability analysis
- **MD5**: Collision attack detection
- **SHA-1**: SHAttered attack awareness
- **RSA**: Quantum computing threat assessment

### Key Strength Metrics
- **Entropy Calculation**: Shannon entropy measurement
- **Character Set Analysis**: Uppercase, lowercase, digits, special characters
- **Length Assessment**: Algorithm-specific minimum requirements
- **Strength Classification**: WEAK, MEDIUM, STRONG ratings

### Timing Attack Education
- **Vulnerable Comparison**: Character-by-character comparison timing
- **Secure Comparison**: Constant-time comparison implementation
- **Performance Measurement**: Microsecond-level timing analysis
- **Educational Value**: Understanding side-channel attacks

## Educational Value

This tool demonstrates:

### Cryptographic Concepts
- Symmetric vs asymmetric encryption
- Hash function properties and weaknesses
- Padding schemes (OAEP, PKCS7)
- Authenticated encryption (GCM mode)
- Key derivation and management

### Security Principles
- Why certain algorithms are deprecated
- Importance of key strength and randomness
- Side-channel attack vulnerabilities
- Proper cryptographic implementation practices

### Real-World Applications
- TLS/SSL protocol understanding
- Password security principles
- Digital signature concepts
- Secure communication foundations

## Technical Implementation

### Security Features
- **Cryptographically Secure Random**: Uses crypto/rand for all random generation
- **Proper Padding**: OAEP for RSA, PKCS7 for block ciphers
- **Authenticated Encryption**: GCM mode for AES prevents tampering
- **Key Derivation**: Demonstrates proper key handling techniques

### Performance Considerations
- **Efficient Algorithms**: Uses Go's optimized crypto libraries
- **Memory Management**: Secure handling of sensitive data
- **Error Handling**: Comprehensive error reporting and validation
- **Resource Management**: Proper cleanup of cryptographic contexts

## Security Warnings

⚠️ **Important Security Notes**:

1. **Educational Purpose**: This tool is designed for learning and authorized testing only
2. **Weak Algorithms**: DES and MD5 are included only for educational demonstration
3. **Key Management**: In production, use proper key management systems
4. **Random Generation**: Always use cryptographically secure random number generators
5. **Side-Channel Protection**: Production code should implement side-channel protections

## Best Practices Demonstrated

### Secure Implementation
- Never reuse nonces in GCM mode
- Use authenticated encryption when possible
- Implement constant-time comparisons for secrets
- Validate all inputs and handle errors securely

### Algorithm Selection
- Use AES-256 for symmetric encryption
- Use RSA-2048+ or modern ECC for asymmetric encryption
- Use SHA-256+ for hash functions
- Avoid deprecated algorithms (DES, MD5, SHA-1)

### Key Management
- Generate keys using cryptographically secure random sources
- Use appropriate key sizes for each algorithm
- Implement proper key derivation functions
- Store keys securely and rotate regularly

## Dependencies

This tool uses only Go standard library packages:
- `crypto/*`: Cryptographic implementations
- `encoding/*`: Data encoding/decoding
- `fmt`, `os`, `bufio`: Basic I/O operations

No external dependencies required for enhanced security and reliability.
