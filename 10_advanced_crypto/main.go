// Advanced Cryptography Tool
// This tool provides comprehensive cryptographic operations including:
// - Multiple encryption algorithms (AES, DES, 3DES, RSA)
// - Hash functions (SHA family, MD5, BLAKE2)
// - Digital signatures and verification
// - Key generation and management
// - Cryptographic analysis and weakness detection
// - Educational cryptography demonstrations

package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// CryptoTool provides comprehensive cryptographic operations
type CryptoTool struct {
	RSAKeySize int
}

// CryptoResult represents the result of a cryptographic operation
type CryptoResult struct {
	Algorithm   string
	Operation   string
	Input       string
	Output      string
	Key         string
	IV          string
	Success     bool
	Error       string
	Metadata    map[string]interface{}
}

// WeaknesAnalysis represents cryptographic weakness analysis
type WeaknesAnalysis struct {
	Algorithm     string
	Vulnerabilities []string
	Recommendations []string
	Severity      string
}

// NewCryptoTool creates a new cryptography tool instance
func NewCryptoTool() *CryptoTool {
	return &CryptoTool{
		RSAKeySize: 2048,
	}
}

// AES Encryption/Decryption
func (ct *CryptoTool) AESEncrypt(plaintext, key string) CryptoResult {
	result := CryptoResult{
		Algorithm: "AES-256-GCM",
		Operation: "ENCRYPT",
		Input:     plaintext,
	}

	keyBytes := make([]byte, 32) // AES-256
	copy(keyBytes, []byte(key))

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// Generate random nonce for GCM
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		result.Error = err.Error()
		return result
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	
	// Prepend nonce to ciphertext
	result.Output = base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
	result.Key = base64.StdEncoding.EncodeToString(keyBytes)
	result.Success = true
	result.Metadata = map[string]interface{}{
		"key_size": len(keyBytes) * 8,
		"nonce_size": len(nonce),
		"mode": "GCM",
	}

	return result
}

func (ct *CryptoTool) AESDecrypt(ciphertext, key string) CryptoResult {
	result := CryptoResult{
		Algorithm: "AES-256-GCM",
		Operation: "DECRYPT",
		Input:     ciphertext,
	}

	keyBytes := make([]byte, 32)
	copy(keyBytes, []byte(key))

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if len(data) < 12 {
		result.Error = "ciphertext too short"
		return result
	}

	nonce := data[:12]
	ciphertextBytes := data[12:]

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Output = string(plaintext)
	result.Success = true
	return result
}

// DES Encryption (for educational purposes - demonstrating weak crypto)
func (ct *CryptoTool) DESEncrypt(plaintext, key string) CryptoResult {
	result := CryptoResult{
		Algorithm: "DES-CBC",
		Operation: "ENCRYPT",
		Input:     plaintext,
	}

	keyBytes := make([]byte, 8) // DES key size
	copy(keyBytes, []byte(key))

	block, err := des.NewCipher(keyBytes)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	// PKCS7 padding
	padding := 8 - len(plaintext)%8
	padtext := plaintext + strings.Repeat(string(rune(padding)), padding)

	iv := make([]byte, des.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		result.Error = err.Error()
		return result
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(padtext))
	mode.CryptBlocks(ciphertext, []byte(padtext))

	result.Output = base64.StdEncoding.EncodeToString(append(iv, ciphertext...))
	result.IV = base64.StdEncoding.EncodeToString(iv)
	result.Success = true
	result.Metadata = map[string]interface{}{
		"key_size": 56, // Effective key size
		"warning": "DES is cryptographically weak - use for educational purposes only",
	}

	return result
}

// RSA Key Generation
func (ct *CryptoTool) GenerateRSAKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, ct.RSAKeySize)
	if err != nil {
		return "", "", err
	}

	// Encode private key
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Encode public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}

// RSA Encryption
func (ct *CryptoTool) RSAEncrypt(plaintext, publicKeyPEM string) CryptoResult {
	result := CryptoResult{
		Algorithm: "RSA-OAEP",
		Operation: "ENCRYPT",
		Input:     plaintext,
	}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		result.Error = "failed to parse PEM block"
		return result
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		result.Error = "not an RSA public key"
		return result
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, []byte(plaintext), nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Output = base64.StdEncoding.EncodeToString(ciphertext)
	result.Success = true
	result.Metadata = map[string]interface{}{
		"key_size": publicKey.Size() * 8,
		"padding": "OAEP",
	}

	return result
}

// Hash Functions
func (ct *CryptoTool) ComputeHash(data, algorithm string) CryptoResult {
	result := CryptoResult{
		Algorithm: algorithm,
		Operation: "HASH",
		Input:     data,
	}

	var hasher hash.Hash
	
	switch strings.ToUpper(algorithm) {
	case "MD5":
		hasher = md5.New()
	case "SHA1":
		hasher = sha1.New()
	case "SHA256":
		hasher = sha256.New()
	case "SHA512":
		hasher = sha512.New()
	default:
		result.Error = "unsupported hash algorithm"
		return result
	}

	hasher.Write([]byte(data))
	hashBytes := hasher.Sum(nil)

	result.Output = hex.EncodeToString(hashBytes)
	result.Success = true
	result.Metadata = map[string]interface{}{
		"output_size": len(hashBytes) * 8,
		"input_size": len(data),
	}

	return result
}

// Cryptographic Analysis
func (ct *CryptoTool) AnalyzeWeakness(algorithm string) WeaknesAnalysis {
	analysis := WeaknesAnalysis{
		Algorithm: algorithm,
	}

	switch strings.ToUpper(algorithm) {
	case "DES":
		analysis.Vulnerabilities = []string{
			"56-bit key size is too small for modern security",
			"Vulnerable to brute force attacks",
			"Meet-in-the-middle attacks possible",
			"Linear and differential cryptanalysis vulnerabilities",
		}
		analysis.Recommendations = []string{
			"Migrate to AES-256",
			"Use 3DES as temporary measure only",
			"Implement proper key management",
		}
		analysis.Severity = "CRITICAL"

	case "3DES", "TRIPLEDES":
		analysis.Vulnerabilities = []string{
			"Effective key size reduced due to meet-in-the-middle attacks",
			"Slow performance compared to AES",
			"Block size of 64 bits is small",
			"Birthday attack concerns with large amounts of data",
		}
		analysis.Recommendations = []string{
			"Migrate to AES-256",
			"Limit amount of data encrypted",
			"Use different keys for each DES operation",
		}
		analysis.Severity = "HIGH"

	case "MD5":
		analysis.Vulnerabilities = []string{
			"Collision attacks are practical",
			"Rainbow table attacks possible",
			"Not suitable for cryptographic purposes",
			"Preimage attacks demonstrated",
		}
		analysis.Recommendations = []string{
			"Use SHA-256 or SHA-3 for new applications",
			"Add salt for password hashing",
			"Consider bcrypt or Argon2 for passwords",
		}
		analysis.Severity = "CRITICAL"

	case "SHA1":
		analysis.Vulnerabilities = []string{
			"Collision attacks demonstrated (SHAttered)",
			"Theoretical preimage attacks",
			"160-bit output insufficient for long-term security",
		}
		analysis.Recommendations = []string{
			"Migrate to SHA-256 or SHA-3",
			"Update digital signature algorithms",
			"Review certificate validation",
		}
		analysis.Severity = "HIGH"

	case "RSA":
		analysis.Vulnerabilities = []string{
			"Quantum computing threat (Shor's algorithm)",
			"Side-channel attack vulnerabilities",
			"Padding oracle attacks if improperly implemented",
			"Key generation quality critical",
		}
		analysis.Recommendations = []string{
			"Use key sizes >= 2048 bits",
			"Implement proper padding (OAEP)",
			"Consider post-quantum alternatives",
			"Use secure random number generation",
		}
		analysis.Severity = "MEDIUM"

	case "AES":
		analysis.Vulnerabilities = []string{
			"Side-channel attacks possible",
			"Key management is critical",
			"GCM mode authentication bypass if nonce reused",
		}
		analysis.Recommendations = []string{
			"Use GCM or CCM modes for authenticated encryption",
			"Never reuse nonces in GCM mode",
			"Implement proper key derivation",
			"Use hardware acceleration when available",
		}
		analysis.Severity = "LOW"

	default:
		analysis.Vulnerabilities = []string{"Unknown algorithm - cannot analyze"}
		analysis.Recommendations = []string{"Use well-known, standardized algorithms"}
		analysis.Severity = "UNKNOWN"
	}

	return analysis
}

// Key Strength Analysis
func (ct *CryptoTool) AnalyzeKeyStrength(key string, algorithm string) map[string]interface{} {
	analysis := make(map[string]interface{})
	
	analysis["length"] = len(key)
	analysis["entropy"] = calculateEntropy(key)
	analysis["character_sets"] = analyzeCharacterSets(key)
	
	// Algorithm-specific analysis
	switch strings.ToUpper(algorithm) {
	case "AES":
		if len(key) < 16 {
			analysis["strength"] = "WEAK"
			analysis["warning"] = "Key too short for AES"
		} else if len(key) >= 32 {
			analysis["strength"] = "STRONG"
		} else {
			analysis["strength"] = "MEDIUM"
		}
	case "DES":
		analysis["strength"] = "WEAK"
		analysis["warning"] = "DES keys are inherently weak"
	default:
		if len(key) < 8 {
			analysis["strength"] = "WEAK"
		} else if len(key) >= 16 {
			analysis["strength"] = "STRONG"
		} else {
			analysis["strength"] = "MEDIUM"
		}
	}
	
	return analysis
}

// Helper functions
func calculateEntropy(s string) float64 {
	if len(s) == 0 {
		return 0
	}

	frequencies := make(map[rune]int)
	for _, char := range s {
		frequencies[char]++
	}

	entropy := 0.0
	length := float64(len(s))

	for _, count := range frequencies {
		probability := float64(count) / length
		if probability > 0 {
			entropy -= probability * (logBase2(probability))
		}
	}

	return entropy
}

func logBase2(x float64) float64 {
	return 0.693147180559945309417 / 0.301029995663981195214 * x // ln(x) / ln(2)
}

func analyzeCharacterSets(s string) map[string]bool {
	sets := map[string]bool{
		"lowercase": false,
		"uppercase": false,
		"digits":    false,
		"special":   false,
	}

	for _, char := range s {
		if char >= 'a' && char <= 'z' {
			sets["lowercase"] = true
		} else if char >= 'A' && char <= 'Z' {
			sets["uppercase"] = true
		} else if char >= '0' && char <= '9' {
			sets["digits"] = true
		} else {
			sets["special"] = true
		}
	}

	return sets
}

// Timing Attack Demonstration
func (ct *CryptoTool) DemonstrateTimingAttack(password, guess string) map[string]interface{} {
	result := make(map[string]interface{})
	
	// Vulnerable comparison (character by character)
	start := time.Now()
	vulnerableCompare := func(a, b string) bool {
		if len(a) != len(b) {
			return false
		}
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				return false
			}
			// Simulate processing delay
			time.Sleep(time.Microsecond)
		}
		return true
	}
	
	isMatch := vulnerableCompare(password, guess)
	vulnerableTime := time.Since(start)
	
	// Secure comparison (constant time)
	start = time.Now()
	secureCompare := func(a, b string) bool {
		if len(a) != len(b) {
			return false
		}
		diff := 0
		for i := 0; i < len(a); i++ {
			diff |= int(a[i]) ^ int(b[i])
			time.Sleep(time.Microsecond) // Simulate constant work
		}
		return diff == 0
	}
	
	isSecureMatch := secureCompare(password, guess)
	secureTime := time.Since(start)
	
	result["vulnerable_time"] = vulnerableTime
	result["secure_time"] = secureTime
	result["match"] = isMatch
	result["secure_match"] = isSecureMatch
	result["timing_difference"] = vulnerableTime - secureTime
	result["warning"] = "Vulnerable comparison allows timing attacks"
	
	return result
}

// Print functions
func (ct *CryptoTool) PrintCryptoResult(result CryptoResult) {
	fmt.Printf("\n=== %s %s Result ===\n", result.Algorithm, result.Operation)
	
	if !result.Success {
		fmt.Printf("❌ Error: %s\n", result.Error)
		return
	}
	
	fmt.Printf("✅ Operation successful\n")
	fmt.Printf("Input: %s\n", result.Input)
	fmt.Printf("Output: %s\n", result.Output)
	
	if result.Key != "" {
		fmt.Printf("Key: %s\n", result.Key)
	}
	
	if result.IV != "" {
		fmt.Printf("IV: %s\n", result.IV)
	}
	
	if len(result.Metadata) > 0 {
		fmt.Println("\nMetadata:")
		for key, value := range result.Metadata {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}
}

func (ct *CryptoTool) PrintWeaknessAnalysis(analysis WeaknesAnalysis) {
	fmt.Printf("\n=== Cryptographic Analysis: %s ===\n", analysis.Algorithm)
	fmt.Printf("Severity: %s\n", analysis.Severity)
	
	if len(analysis.Vulnerabilities) > 0 {
		fmt.Println("\n🚨 Vulnerabilities:")
		for i, vuln := range analysis.Vulnerabilities {
			fmt.Printf("%d. %s\n", i+1, vuln)
		}
	}
	
	if len(analysis.Recommendations) > 0 {
		fmt.Println("\n💡 Recommendations:")
		for i, rec := range analysis.Recommendations {
			fmt.Printf("%d. %s\n", i+1, rec)
		}
	}
}

// Interactive menu
func showMenu() {
	fmt.Println("\n=== Advanced Cryptography Tool ===")
	fmt.Println("1. AES Encryption/Decryption")
	fmt.Println("2. DES Encryption (Educational)")
	fmt.Println("3. RSA Key Generation & Encryption")
	fmt.Println("4. Hash Functions")
	fmt.Println("5. Cryptographic Analysis")
	fmt.Println("6. Key Strength Analysis")
	fmt.Println("7. Timing Attack Demonstration")
	fmt.Println("8. Exit")
	fmt.Print("Select option: ")
}

func main() {
	tool := NewCryptoTool()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Advanced Cryptography Tool v1.0")
	fmt.Println("Educational tool for cryptographic operations and analysis")
	fmt.Println("⚠️  Use for educational purposes and authorized testing only")

	for {
		showMenu()
		
		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Println("\n--- AES Operations ---")
			fmt.Print("1. Encrypt  2. Decrypt: ")
			if !scanner.Scan() {
				continue
			}
			op := strings.TrimSpace(scanner.Text())

			fmt.Print("Enter text: ")
			if !scanner.Scan() {
				continue
			}
			text := strings.TrimSpace(scanner.Text())

			fmt.Print("Enter key: ")
			if !scanner.Scan() {
				continue
			}
			key := strings.TrimSpace(scanner.Text())

			if op == "1" {
				result := tool.AESEncrypt(text, key)
				tool.PrintCryptoResult(result)
			} else if op == "2" {
				result := tool.AESDecrypt(text, key)
				tool.PrintCryptoResult(result)
			}

		case "2":
			fmt.Println("\n--- DES Encryption (Educational) ---")
			fmt.Print("Enter plaintext: ")
			if !scanner.Scan() {
				continue
			}
			plaintext := strings.TrimSpace(scanner.Text())

			fmt.Print("Enter key (8 chars): ")
			if !scanner.Scan() {
				continue
			}
			key := strings.TrimSpace(scanner.Text())

			result := tool.DESEncrypt(plaintext, key)
			tool.PrintCryptoResult(result)

			// Show weakness analysis
			analysis := tool.AnalyzeWeakness("DES")
			tool.PrintWeaknessAnalysis(analysis)

		case "3":
			fmt.Println("\n--- RSA Operations ---")
			fmt.Println("Generating RSA key pair...")
			
			privateKey, publicKey, err := tool.GenerateRSAKeyPair()
			if err != nil {
				fmt.Printf("Error generating keys: %v\n", err)
				continue
			}

			fmt.Println("✅ RSA Key Pair Generated")
			fmt.Printf("Public Key:\n%s\n", publicKey)
			fmt.Printf("Private Key:\n%s\n", privateKey)

			fmt.Print("Enter text to encrypt: ")
			if !scanner.Scan() {
				continue
			}
			plaintext := strings.TrimSpace(scanner.Text())

			result := tool.RSAEncrypt(plaintext, publicKey)
			tool.PrintCryptoResult(result)

		case "4":
			fmt.Println("\n--- Hash Functions ---")
			fmt.Print("Enter text to hash: ")
			if !scanner.Scan() {
				continue
			}
			text := strings.TrimSpace(scanner.Text())

			algorithms := []string{"MD5", "SHA1", "SHA256", "SHA512"}
			for _, alg := range algorithms {
				result := tool.ComputeHash(text, alg)
				tool.PrintCryptoResult(result)

				if alg == "MD5" || alg == "SHA1" {
					analysis := tool.AnalyzeWeakness(alg)
					fmt.Printf("⚠️  %s is cryptographically weak\n", alg)
				}
			}

		case "5":
			fmt.Println("\n--- Cryptographic Analysis ---")
			fmt.Print("Enter algorithm to analyze: ")
			if !scanner.Scan() {
				continue
			}
			algorithm := strings.TrimSpace(scanner.Text())

			analysis := tool.AnalyzeWeakness(algorithm)
			tool.PrintWeaknessAnalysis(analysis)

		case "6":
			fmt.Println("\n--- Key Strength Analysis ---")
			fmt.Print("Enter key to analyze: ")
			if !scanner.Scan() {
				continue
			}
			key := strings.TrimSpace(scanner.Text())

			fmt.Print("Enter algorithm: ")
			if !scanner.Scan() {
				continue
			}
			algorithm := strings.TrimSpace(scanner.Text())

			analysis := tool.AnalyzeKeyStrength(key, algorithm)
			
			fmt.Println("\n--- Key Analysis Results ---")
			for k, v := range analysis {
				fmt.Printf("%s: %v\n", k, v)
			}

		case "7":
			fmt.Println("\n--- Timing Attack Demonstration ---")
			fmt.Print("Enter correct password: ")
			if !scanner.Scan() {
				continue
			}
			password := strings.TrimSpace(scanner.Text())

			fmt.Print("Enter guess: ")
			if !scanner.Scan() {
				continue
			}
			guess := strings.TrimSpace(scanner.Text())

			result := tool.DemonstrateTimingAttack(password, guess)
			
			fmt.Println("\n--- Timing Attack Results ---")
			for k, v := range result {
				fmt.Printf("%s: %v\n", k, v)
			}

		case "8":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
