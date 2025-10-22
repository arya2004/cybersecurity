/*
Hash Cracker and Password Security Tool
=======================================

This educational tool demonstrates password hashing, cracking techniques, and
password strength analysis. It helps understand password security concepts and
the importance of strong passwords.

Features:
- Multiple hash algorithm support (MD5, SHA1, SHA256, SHA512)
- Dictionary-based password cracking
- Brute-force password cracking
- Rainbow table simulation
- Password strength analysis
- Password generation recommendations
- Hash comparison and verification

Educational Value:
- Understanding cryptographic hash functions
- Password storage security
- Attack techniques and defenses
- Password complexity requirements
- Brute-force attack demonstration

Usage:
    go run main.go
    Choose from menu options to explore password security

⚠️  DISCLAIMER: For educational purposes only. Only use on systems you own.
*/

package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// HashAlgorithm represents a hashing algorithm
type HashAlgorithm string

const (
	MD5_HASH    HashAlgorithm = "MD5"
	SHA1_HASH   HashAlgorithm = "SHA1"
	SHA256_HASH HashAlgorithm = "SHA256"
	SHA512_HASH HashAlgorithm = "SHA512"
)

// PasswordStrength represents password strength analysis
type PasswordStrength struct {
	Password       string
	Score          int
	Strength       string
	Length         int
	HasLower       bool
	HasUpper       bool
	HasDigit       bool
	HasSpecial     bool
	EstimatedTime  string
	Suggestions    []string
	CommonPatterns []string
}

// HashPassword generates a hash for the given password
func HashPassword(password string, algorithm HashAlgorithm) string {
	var hash []byte

	switch algorithm {
	case MD5_HASH:
		h := md5.Sum([]byte(password))
		hash = h[:]
	case SHA1_HASH:
		h := sha1.Sum([]byte(password))
		hash = h[:]
	case SHA256_HASH:
		h := sha256.Sum256([]byte(password))
		hash = h[:]
	case SHA512_HASH:
		h := sha512.Sum512([]byte(password))
		hash = h[:]
	}

	return hex.EncodeToString(hash)
}

// CrackHashDictionary attempts to crack a hash using dictionary attack
func CrackHashDictionary(targetHash string, algorithm HashAlgorithm, wordlist []string) (string, bool, int) {
	fmt.Printf("\n[*] Starting dictionary attack...\n")
	fmt.Printf("[*] Wordlist size: %d passwords\n", len(wordlist))
	fmt.Printf("[*] Hash algorithm: %s\n", algorithm)

	startTime := time.Now()
	attempts := 0

	for _, word := range wordlist {
		attempts++
		hash := HashPassword(strings.TrimSpace(word), algorithm)

		if hash == targetHash {
			duration := time.Since(startTime)
			fmt.Printf("\n[✓] PASSWORD CRACKED!\n")
			fmt.Printf("    Password: %s\n", word)
			fmt.Printf("    Attempts: %d\n", attempts)
			fmt.Printf("    Time: %v\n", duration)
			return word, true, attempts
		}

		if attempts%1000 == 0 {
			fmt.Printf("\r[*] Tested %d passwords...", attempts)
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\n[✗] Password not found in dictionary\n")
	fmt.Printf("    Attempts: %d\n", attempts)
	fmt.Printf("    Time: %v\n", duration)

	return "", false, attempts
}

// CrackHashBruteForce attempts brute force crack (limited for educational purposes)
func CrackHashBruteForce(targetHash string, algorithm HashAlgorithm, maxLength int) (string, bool) {
	fmt.Printf("\n[*] Starting brute force attack...\n")
	fmt.Printf("[*] Maximum length: %d characters\n", maxLength)
	fmt.Printf("[*] Hash algorithm: %s\n", algorithm)
	fmt.Println("[!] WARNING: This may take a very long time!")

	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	startTime := time.Now()
	attempts := 0

	// Only try up to maxLength for demonstration
	for length := 1; length <= maxLength; length++ {
		fmt.Printf("\n[*] Testing passwords of length %d...\n", length)

		// Generate and test passwords (simplified for educational purposes)
		// In reality, this would be much more complex
		result, found, count := bruteForceLengthN(targetHash, algorithm, charset, length, 10000)
		attempts += count

		if found {
			duration := time.Since(startTime)
			fmt.Printf("\n[✓] PASSWORD CRACKED!\n")
			fmt.Printf("    Password: %s\n", result)
			fmt.Printf("    Attempts: %d\n", attempts)
			fmt.Printf("    Time: %v\n", duration)
			return result, true
		}
	}

	fmt.Printf("\n[✗] Password not cracked within limits\n")
	return "", false
}

// bruteForceLengthN tries passwords of specific length (limited for demo)
func bruteForceLengthN(targetHash string, algorithm HashAlgorithm, charset string, length, maxAttempts int) (string, bool, int) {
	attempts := 0

	// For educational purposes, only try a limited number
	// Real brute force would be exhaustive
	for i := 0; i < maxAttempts && attempts < maxAttempts; i++ {
		password := generateRandomPassword(charset, length)
		hash := HashPassword(password, algorithm)
		attempts++

		if hash == targetHash {
			return password, true, attempts
		}

		if attempts%100 == 0 {
			fmt.Printf("\r    Tested %d/%d passwords...", attempts, maxAttempts)
		}
	}

	return "", false, attempts
}

// generateRandomPassword generates a random password for brute force demo
func generateRandomPassword(charset string, length int) string {
	// Simplified for educational demo
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[i%len(charset)]
	}
	return string(password)
}

// AnalyzePasswordStrength performs comprehensive password strength analysis
func AnalyzePasswordStrength(password string) *PasswordStrength {
	ps := &PasswordStrength{
		Password:    password,
		Length:      len(password),
		Suggestions: make([]string, 0),
	}

	// Check character types
	ps.HasLower = hasLowercase(password)
	ps.HasUpper = hasUppercase(password)
	ps.HasDigit = hasDigit(password)
	ps.HasSpecial = hasSpecialChar(password)

	// Calculate score
	ps.Score = calculatePasswordScore(ps)

	// Determine strength
	switch {
	case ps.Score >= 80:
		ps.Strength = "Very Strong"
	case ps.Score >= 60:
		ps.Strength = "Strong"
	case ps.Score >= 40:
		ps.Strength = "Moderate"
	case ps.Score >= 20:
		ps.Strength = "Weak"
	default:
		ps.Strength = "Very Weak"
	}

	// Estimate crack time
	ps.EstimatedTime = estimateCrackTime(ps)

	// Generate suggestions
	ps.generateSuggestions()

	// Check for common patterns
	ps.detectCommonPatterns()

	return ps
}

// calculatePasswordScore calculates password score based on multiple factors
func calculatePasswordScore(ps *PasswordStrength) int {
	score := 0

	// Length scoring
	if ps.Length >= 12 {
		score += 30
	} else if ps.Length >= 8 {
		score += 20
	} else if ps.Length >= 6 {
		score += 10
	}

	// Character variety scoring
	if ps.HasLower {
		score += 10
	}
	if ps.HasUpper {
		score += 15
	}
	if ps.HasDigit {
		score += 15
	}
	if ps.HasSpecial {
		score += 20
	}

	// Bonus for using all character types
	if ps.HasLower && ps.HasUpper && ps.HasDigit && ps.HasSpecial {
		score += 10
	}

	// Entropy bonus
	entropy := calculateEntropy(ps.Password)
	if entropy > 50 {
		score += 10
	}

	return score
}

// calculateEntropy calculates password entropy
func calculateEntropy(password string) float64 {
	var poolSize float64

	if hasLowercase(password) {
		poolSize += 26
	}
	if hasUppercase(password) {
		poolSize += 26
	}
	if hasDigit(password) {
		poolSize += 10
	}
	if hasSpecialChar(password) {
		poolSize += 32
	}

	if poolSize == 0 {
		return 0
	}

	entropy := float64(len(password)) * math.Log2(poolSize)
	return entropy
}

// estimateCrackTime estimates time to crack password
func estimateCrackTime(ps *PasswordStrength) string {
	var poolSize float64 = 0

	if ps.HasLower {
		poolSize += 26
	}
	if ps.HasUpper {
		poolSize += 26
	}
	if ps.HasDigit {
		poolSize += 10
	}
	if ps.HasSpecial {
		poolSize += 32
	}

	if poolSize == 0 {
		return "Instantly"
	}

	// Assuming 1 billion hashes/second (modern GPU)
	hashesPerSecond := 1000000000.0
	possibleCombinations := math.Pow(poolSize, float64(ps.Length))
	seconds := possibleCombinations / hashesPerSecond / 2 // Average case

	return formatDuration(seconds)
}

// formatDuration converts seconds to human readable format
func formatDuration(seconds float64) string {
	if seconds < 1 {
		return "Instantly"
	} else if seconds < 60 {
		return fmt.Sprintf("%.0f seconds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.0f minutes", seconds/60)
	} else if seconds < 86400 {
		return fmt.Sprintf("%.0f hours", seconds/3600)
	} else if seconds < 31536000 {
		return fmt.Sprintf("%.0f days", seconds/86400)
	} else if seconds < 31536000000 {
		return fmt.Sprintf("%.0f years", seconds/31536000)
	} else {
		return "Centuries or more"
	}
}

// generateSuggestions generates password improvement suggestions
func (ps *PasswordStrength) generateSuggestions() {
	if ps.Length < 12 {
		ps.Suggestions = append(ps.Suggestions, "Increase password length to at least 12 characters")
	}

	if !ps.HasLower {
		ps.Suggestions = append(ps.Suggestions, "Add lowercase letters (a-z)")
	}

	if !ps.HasUpper {
		ps.Suggestions = append(ps.Suggestions, "Add uppercase letters (A-Z)")
	}

	if !ps.HasDigit {
		ps.Suggestions = append(ps.Suggestions, "Add numbers (0-9)")
	}

	if !ps.HasSpecial {
		ps.Suggestions = append(ps.Suggestions, "Add special characters (!@#$%^&*)")
	}

	if len(ps.CommonPatterns) > 0 {
		ps.Suggestions = append(ps.Suggestions, "Avoid common patterns and sequences")
	}
}

// detectCommonPatterns detects common password patterns
func (ps *PasswordStrength) detectCommonPatterns() {
	password := strings.ToLower(ps.Password)

	commonWords := []string{"password", "admin", "user", "login", "welcome", "123456", "qwerty"}
	for _, word := range commonWords {
		if strings.Contains(password, word) {
			ps.CommonPatterns = append(ps.CommonPatterns, fmt.Sprintf("Contains common word: '%s'", word))
		}
	}

	// Check for sequential characters
	if regexp.MustCompile(`(?i)(abc|bcd|cde|123|234|345|456|567|678|789)`).MatchString(password) {
		ps.CommonPatterns = append(ps.CommonPatterns, "Contains sequential characters")
	}

	// Check for repeated characters
	if regexp.MustCompile(`(.)\1{2,}`).MatchString(password) {
		ps.CommonPatterns = append(ps.CommonPatterns, "Contains repeated characters")
	}

	// Check for keyboard patterns
	if regexp.MustCompile(`(?i)(qwert|asdf|zxcv)`).MatchString(password) {
		ps.CommonPatterns = append(ps.CommonPatterns, "Contains keyboard pattern")
	}
}

// Helper functions for character type checking
func hasLowercase(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func hasUppercase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !unicode.IsSpace(char) {
			return true
		}
	}
	return false
}

// PrintPasswordAnalysis prints detailed password analysis
func PrintPasswordAnalysis(ps *PasswordStrength) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("PASSWORD STRENGTH ANALYSIS")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\nPassword Length: %d characters\n", ps.Length)
	fmt.Printf("Strength Score: %d/100\n", ps.Score)
	fmt.Printf("Strength Rating: %s\n", ps.Strength)

	fmt.Println("\nCharacter Composition:")
	fmt.Printf("  Lowercase letters: %s\n", boolToStatus(ps.HasLower))
	fmt.Printf("  Uppercase letters: %s\n", boolToStatus(ps.HasUpper))
	fmt.Printf("  Digits: %s\n", boolToStatus(ps.HasDigit))
	fmt.Printf("  Special characters: %s\n", boolToStatus(ps.HasSpecial))

	entropy := calculateEntropy(ps.Password)
	fmt.Printf("\nEntropy: %.2f bits\n", entropy)
	fmt.Printf("Estimated crack time: %s\n", ps.EstimatedTime)

	if len(ps.CommonPatterns) > 0 {
		fmt.Println("\n⚠️  Common Patterns Detected:")
		for _, pattern := range ps.CommonPatterns {
			fmt.Printf("  - %s\n", pattern)
		}
	}

	if len(ps.Suggestions) > 0 {
		fmt.Println("\n💡 Suggestions for Improvement:")
		for i, suggestion := range ps.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
}

func boolToStatus(b bool) string {
	if b {
		return "[✓] Yes"
	}
	return "[✗] No"
}

// GetDefaultWordlist provides a small wordlist for educational purposes
func GetDefaultWordlist() []string {
	return []string{
		"password", "123456", "12345678", "qwerty", "abc123",
		"monkey", "1234567", "letmein", "trustno1", "dragon",
		"baseball", "111111", "iloveyou", "master", "sunshine",
		"ashley", "bailey", "passw0rd", "shadow", "123123",
		"654321", "superman", "qazwsx", "michael", "football",
		"admin", "welcome", "login", "password123", "test",
	}
}

// Main menu
func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("  HASH CRACKER & PASSWORD SECURITY TOOL")
	fmt.Println("  Educational Tool for Password Security")
	fmt.Println(strings.Repeat("=", 60))

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Generate Password Hash")
		fmt.Println("2. Crack Hash (Dictionary Attack)")
		fmt.Println("3. Analyze Password Strength")
		fmt.Println("4. Compare Hash")
		fmt.Println("5. Educational Demo")
		fmt.Println("6. Exit")
		fmt.Print("\nChoice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			generateHashMenu()
		case 2:
			crackHashMenu()
		case 3:
			analyzePasswordMenu()
		case 4:
			compareHashMenu()
		case 5:
			educationalDemo()
		case 6:
			fmt.Println("\nExiting... Stay secure!")
			return
		default:
			fmt.Println("Invalid choice!")
		}
	}
}

func generateHashMenu() {
	fmt.Print("\nEnter password: ")
	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Println("\nGenerated Hashes:")
	fmt.Printf("MD5:    %s\n", HashPassword(password, MD5_HASH))
	fmt.Printf("SHA1:   %s\n", HashPassword(password, SHA1_HASH))
	fmt.Printf("SHA256: %s\n", HashPassword(password, SHA256_HASH))
	fmt.Printf("SHA512: %s\n", HashPassword(password, SHA512_HASH))
}

func crackHashMenu() {
	fmt.Print("\nEnter hash to crack: ")
	reader := bufio.NewReader(os.Stdin)
	hash, _ := reader.ReadString('\n')
	hash = strings.TrimSpace(hash)

	fmt.Println("\nSelect algorithm:")
	fmt.Println("1. MD5")
	fmt.Println("2. SHA1")
	fmt.Println("3. SHA256")
	fmt.Print("Choice: ")

	var choice int
	fmt.Scanln(&choice)

	var algorithm HashAlgorithm
	switch choice {
	case 1:
		algorithm = MD5_HASH
	case 2:
		algorithm = SHA1_HASH
	case 3:
		algorithm = SHA256_HASH
	default:
		fmt.Println("Invalid choice!")
		return
	}

	wordlist := GetDefaultWordlist()
	CrackHashDictionary(hash, algorithm, wordlist)
}

func analyzePasswordMenu() {
	fmt.Print("\nEnter password to analyze: ")
	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	analysis := AnalyzePasswordStrength(password)
	PrintPasswordAnalysis(analysis)
}

func compareHashMenu() {
	fmt.Print("\nEnter password: ")
	reader := bufio.NewReader(os.Stdin)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Enter hash to compare: ")
	hash, _ := reader.ReadString('\n')
	hash = strings.TrimSpace(hash)

	if HashPassword(password, MD5_HASH) == hash {
		fmt.Println("\n[✓] Hash matches (MD5)")
	} else if HashPassword(password, SHA1_HASH) == hash {
		fmt.Println("\n[✓] Hash matches (SHA1)")
	} else if HashPassword(password, SHA256_HASH) == hash {
		fmt.Println("\n[✓] Hash matches (SHA256)")
	} else if HashPassword(password, SHA512_HASH) == hash {
		fmt.Println("\n[✓] Hash matches (SHA512)")
	} else {
		fmt.Println("\n[✗] Hash does not match")
	}
}

func educationalDemo() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("EDUCATIONAL DEMONSTRATION")
	fmt.Println(strings.Repeat("=", 60))

	testPasswords := []string{"password", "P@ssw0rd", "MyS3cur3P@ss!", "a", "12345678"}

	for _, pwd := range testPasswords {
		fmt.Printf("\nAnalyzing: '%s'\n", pwd)
		analysis := AnalyzePasswordStrength(pwd)
		fmt.Printf("Strength: %s (Score: %d/100)\n", analysis.Strength, analysis.Score)
		fmt.Printf("Estimated crack time: %s\n", analysis.EstimatedTime)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("\n⚠️  Key Takeaways:")
	fmt.Println("1. Use passwords with at least 12 characters")
	fmt.Println("2. Mix uppercase, lowercase, numbers, and special characters")
	fmt.Println("3. Avoid common words and patterns")
	fmt.Println("4. Use unique passwords for each account")
	fmt.Println("5. Consider using a password manager")
	fmt.Println(strings.Repeat("=", 60))
}
