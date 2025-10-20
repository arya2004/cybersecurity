package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strings"
)

// Common weak passwords
var commonPasswords = map[string]bool{
	"password": true, "123456": true, "123456789": true, "12345678": true,
	"12345": true, "1234567": true, "password1": true, "123123": true,
	"qwerty": true, "abc123": true, "111111": true, "admin": true,
	"letmein": true, "welcome": true, "monkey": true, "dragon": true,
}

// Keyboard patterns
var keyboardPatterns = []string{
	"qwerty", "asdf", "zxcv", "qazwsx", "123456", "098765",
	"qwertyuiop", "asdfgh", "zxcvbnm",
}

// PasswordAnalyzer analyzes password strength
type PasswordAnalyzer struct {
	Password string
	Length   int
}

// NewPasswordAnalyzer creates a new password analyzer
func NewPasswordAnalyzer(password string) *PasswordAnalyzer {
	return &PasswordAnalyzer{
		Password: password,
		Length:   len(password),
	}
}

// CalculateEntropy calculates Shannon entropy
func (pa *PasswordAnalyzer) CalculateEntropy() float64 {
	if pa.Length == 0 {
		return 0.0
	}

	freq := make(map[rune]int)
	for _, char := range pa.Password {
		freq[char]++
	}

	entropy := 0.0
	for _, count := range freq {
		probability := float64(count) / float64(pa.Length)
		entropy -= probability * math.Log2(probability)
	}

	return entropy * float64(pa.Length)
}

// GetCharacterSets detects character sets used
func (pa *PasswordAnalyzer) GetCharacterSets() map[string]bool {
	return map[string]bool{
		"lowercase": regexp.MustCompile(`[a-z]`).MatchString(pa.Password),
		"uppercase": regexp.MustCompile(`[A-Z]`).MatchString(pa.Password),
		"numbers":   regexp.MustCompile(`\d`).MatchString(pa.Password),
		"special":   regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};:'",.<>?/\\|` + "`~]").MatchString(pa.Password),
	}
}

// CalculateCharsetSize calculates character set size
func (pa *PasswordAnalyzer) CalculateCharsetSize() int {
	charset := pa.GetCharacterSets()
	size := 0
	if charset["lowercase"] {
		size += 26
	}
	if charset["uppercase"] {
		size += 26
	}
	if charset["numbers"] {
		size += 10
	}
	if charset["special"] {
		size += 32
	}
	return size
}

// DetectSequentialPatterns detects sequential patterns
func (pa *PasswordAnalyzer) DetectSequentialPatterns() []string {
	patterns := []string{}

	// Sequential numbers
	if regexp.MustCompile(`(012|123|234|345|456|567|678|789|987|876|765|654|543|432|321|210)`).MatchString(pa.Password) {
		patterns = append(patterns, "Sequential numbers")
	}

	// Sequential letters
	lower := strings.ToLower(pa.Password)
	if regexp.MustCompile(`(abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz)`).MatchString(lower) {
		patterns = append(patterns, "Sequential letters")
	}

	return patterns
}

// DetectKeyboardPatterns detects keyboard patterns
func (pa *PasswordAnalyzer) DetectKeyboardPatterns() []string {
	patterns := []string{}
	lower := strings.ToLower(pa.Password)

	for _, pattern := range keyboardPatterns {
		if strings.Contains(lower, pattern) {
			patterns = append(patterns, fmt.Sprintf("Keyboard pattern: '%s'", pattern))
		}
	}

	return patterns
}

// CheckCommonPassword checks if password is common
func (pa *PasswordAnalyzer) CheckCommonPassword() bool {
	return commonPasswords[strings.ToLower(pa.Password)]
}

// CalculateStrengthScore calculates password strength (0-100)
func (pa *PasswordAnalyzer) CalculateStrengthScore() (int, string) {
	score := 0

	// Length scoring (0-30 points)
	if pa.Length >= 16 {
		score += 30
	} else if pa.Length >= 12 {
		score += 25
	} else if pa.Length >= 8 {
		score += 15
	} else if pa.Length >= 6 {
		score += 5
	}

	// Character variety (0-25 points)
	charset := pa.GetCharacterSets()
	for _, present := range charset {
		if present {
			score += 6
		}
	}

	// Entropy bonus (0-25 points)
	entropy := pa.CalculateEntropy()
	if entropy >= 80 {
		score += 25
	} else if entropy >= 60 {
		score += 20
	} else if entropy >= 40 {
		score += 10
	} else if entropy >= 20 {
		score += 5
	}

	// Deduct for vulnerabilities
	vulnerabilities := 0
	if pa.CheckCommonPassword() {
		vulnerabilities += 10
	}
	vulnerabilities += len(pa.DetectSequentialPatterns()) * 2
	vulnerabilities += len(pa.DetectKeyboardPatterns()) * 3

	score -= vulnerabilities
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// Determine rating
	rating := ""
	if score >= 80 {
		rating = "VERY STRONG"
	} else if score >= 60 {
		rating = "STRONG"
	} else if score >= 40 {
		rating = "MODERATE"
	} else if score >= 20 {
		rating = "WEAK"
	} else {
		rating = "VERY WEAK"
	}

	return score, rating
}

// EstimateCrackTime estimates brute force crack time
func (pa *PasswordAnalyzer) EstimateCrackTime() string {
	charsetSize := pa.CalculateCharsetSize()
	if charsetSize == 0 {
		return "Instant"
	}

	// Assume 10 billion attempts per second
	attemptsPerSecond := 10_000_000_000.0
	totalCombinations := math.Pow(float64(charsetSize), float64(pa.Length))
	seconds := totalCombinations / attemptsPerSecond

	if seconds < 1 {
		return "Instant"
	} else if seconds < 60 {
		return fmt.Sprintf("%.1f seconds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.1f minutes", seconds/60)
	} else if seconds < 86400 {
		return fmt.Sprintf("%.1f hours", seconds/3600)
	} else if seconds < 31536000 {
		return fmt.Sprintf("%.1f days", seconds/86400)
	} else if seconds < 31536000*100 {
		return fmt.Sprintf("%.1f years", seconds/31536000)
	} else {
		return fmt.Sprintf("%.0f+ years", seconds/31536000)
	}
}

// PrintAnalysis prints password analysis
func (pa *PasswordAnalyzer) PrintAnalysis() {
	score, rating := pa.CalculateStrengthScore()
	charset := pa.GetCharacterSets()
	entropy := pa.CalculateEntropy()
	crackTime := pa.EstimateCrackTime()

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🔍 Password Analysis")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("\n📊 Strength Score: %d/100 - %s\n", score, rating)
	fmt.Printf("🔐 Entropy: %.2f bits\n", entropy)
	fmt.Printf("📏 Length: %d characters\n", pa.Length)
	fmt.Printf("🎲 Character Set Size: %d\n", pa.CalculateCharsetSize())
	fmt.Printf("⏱️  Estimated Crack Time: %s\n", crackTime)

	fmt.Println("\n✓ Character Types:")
	for name, present := range charset {
		symbol := "❌"
		if present {
			symbol = "✅"
		}
		fmt.Printf("  %s %s\n", symbol, strings.Title(name))
	}

	if pa.CheckCommonPassword() {
		fmt.Println("\n⚠️  WARNING: This is a commonly used password!")
	}

	vulnerabilities := append(pa.DetectSequentialPatterns(), pa.DetectKeyboardPatterns()...)
	if len(vulnerabilities) > 0 {
		fmt.Println("\n⚠️  Vulnerabilities Found:")
		for _, vuln := range vulnerabilities {
			fmt.Printf("  ⚠️  %s\n", vuln)
		}
	} else {
		fmt.Println("\n✅ No common vulnerabilities detected!")
	}

	fmt.Println(strings.Repeat("=", 60))
}

// GeneratePassword generates a cryptographically secure password
func GeneratePassword(length int, useUpper, useLower, useNumbers, useSpecial bool) (string, error) {
	if length < 4 {
		return "", fmt.Errorf("password length must be at least 4")
	}

	chars := ""
	if useLower {
		chars += "abcdefghijklmnopqrstuvwxyz"
	}
	if useUpper {
		chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if useNumbers {
		chars += "0123456789"
	}
	if useSpecial {
		chars += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	if chars == "" {
		return "", fmt.Errorf("at least one character set must be selected")
	}

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		password[i] = chars[num.Int64()]
	}

	return string(password), nil
}

func main() {
	// Command-line flags
	analyzePtr := flag.String("analyze", "", "Password to analyze")
	generatePtr := flag.Bool("generate", false, "Generate a new password")
	lengthPtr := flag.Int("length", 16, "Password length")
	upperPtr := flag.Bool("uppercase", true, "Include uppercase letters")
	lowerPtr := flag.Bool("lowercase", true, "Include lowercase letters")
	numbersPtr := flag.Bool("numbers", true, "Include numbers")
	specialPtr := flag.Bool("special", true, "Include special characters")

	flag.Parse()

	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🔐 Password Security Tool (Go)")
	fmt.Println(strings.Repeat("=", 60))

	if *analyzePtr != "" {
		// Analyze mode
		analyzer := NewPasswordAnalyzer(*analyzePtr)
		analyzer.PrintAnalysis()
	} else if *generatePtr {
		// Generate mode
		password, err := GeneratePassword(*lengthPtr, *upperPtr, *lowerPtr, *numbersPtr, *specialPtr)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		fmt.Printf("\n🎲 Generated Password: %s\n", password)

		// Analyze generated password
		analyzer := NewPasswordAnalyzer(password)
		score, rating := analyzer.CalculateStrengthScore()
		fmt.Printf("Strength: %d/100 - %s\n", score, rating)
		fmt.Println(strings.Repeat("=", 60))
	} else {
		// Interactive mode
		fmt.Println("\nUsage:")
		fmt.Println("  Analyze: go run main.go -analyze \"YourPassword123!\"")
		fmt.Println("  Generate: go run main.go -generate -length 16")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}
}
