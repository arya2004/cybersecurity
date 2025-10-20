package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// PasswordAnalysis contains detailed password analysis results
type PasswordAnalysis struct {
	Password           string
	Length             int
	HasLowercase       bool
	HasUppercase       bool
	HasDigits          bool
	HasSpecialChars    bool
	EntropyBits        float64
	Strength           string
	Score              int
	CrackTimeSeconds   float64
	CrackTimeReadable  string
	IsCommon           bool
	ContainsDictionary bool
	Patterns           []string
	Suggestions        []string
}

// Common weak passwords list
var commonPasswords = []string{
	"password", "123456", "12345678", "qwerty", "abc123", "monkey",
	"1234567", "letmein", "trustno1", "dragon", "baseball", "iloveyou",
	"master", "sunshine", "ashley", "bailey", "passw0rd", "shadow",
	"123123", "654321", "superman", "qazwsx", "michael", "football",
}

// Common dictionary words to check
var commonWords = []string{
	"love", "admin", "user", "welcome", "hello", "test", "demo",
	"pass", "access", "login", "root", "system", "guest", "temp",
}

// AnalyzePassword performs comprehensive password strength analysis
func AnalyzePassword(password string) PasswordAnalysis {
	analysis := PasswordAnalysis{
		Password:  password,
		Length:    len(password),
		Patterns:  []string{},
		Suggestions: []string{},
	}

	// Character type checks
	analysis.HasLowercase = containsLowercase(password)
	analysis.HasUppercase = containsUppercase(password)
	analysis.HasDigits = containsDigits(password)
	analysis.HasSpecialChars = containsSpecialChars(password)

	// Calculate entropy
	analysis.EntropyBits = calculateEntropy(password)

	// Check for common passwords
	analysis.IsCommon = isCommonPassword(password)

	// Check for dictionary words
	analysis.ContainsDictionary = containsDictionaryWord(password)

	// Detect patterns
	analysis.Patterns = detectPatterns(password)

	// Calculate score
	analysis.Score = calculateScore(analysis)

	// Determine strength level
	analysis.Strength = determineStrength(analysis.Score)

	// Estimate crack time
	analysis.CrackTimeSeconds = estimateCrackTime(analysis.EntropyBits)
	analysis.CrackTimeReadable = formatCrackTime(analysis.CrackTimeSeconds)

	// Generate suggestions
	analysis.Suggestions = generateSuggestions(analysis)

	return analysis
}

// containsLowercase checks if password has lowercase letters
func containsLowercase(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

// containsUppercase checks if password has uppercase letters
func containsUppercase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

// containsDigits checks if password has digits
func containsDigits(s string) bool {
	for _, char := range s {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

// containsSpecialChars checks if password has special characters
func containsSpecialChars(s string) bool {
	specialChars := `!@#$%^&*()_+-=[]{}|;:,.<>?/~` + "`"
	for _, char := range s {
		if strings.ContainsRune(specialChars, char) {
			return true
		}
	}
	return false
}

// calculateEntropy calculates Shannon entropy of the password
func calculateEntropy(password string) float64 {
	if len(password) == 0 {
		return 0
	}

	// Count character types
	charsetSize := 0
	if containsLowercase(password) {
		charsetSize += 26
	}
	if containsUppercase(password) {
		charsetSize += 26
	}
	if containsDigits(password) {
		charsetSize += 10
	}
	if containsSpecialChars(password) {
		charsetSize += 32
	}

	if charsetSize == 0 {
		return 0
	}

	// Entropy = log2(charset^length)
	entropy := float64(len(password)) * math.Log2(float64(charsetSize))
	return entropy
}

// isCommonPassword checks if password is in common passwords list
func isCommonPassword(password string) bool {
	lowerPassword := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}
	return false
}

// containsDictionaryWord checks for common dictionary words
func containsDictionaryWord(password string) bool {
	lowerPassword := strings.ToLower(password)
	for _, word := range commonWords {
		if strings.Contains(lowerPassword, word) {
			return true
		}
	}
	return false
}

// detectPatterns identifies common password patterns
func detectPatterns(password string) []string {
	patterns := []string{}

	// Check for sequential numbers
	if matched, _ := regexp.MatchString(`\d{3,}`, password); matched {
		patterns = append(patterns, "Sequential numbers detected")
	}

	// Check for repeated characters
	if matched, _ := regexp.MatchString(`(.)\1{2,}`, password); matched {
		patterns = append(patterns, "Repeated characters detected")
	}

	// Check for keyboard patterns
	keyboardPatterns := []string{"qwerty", "asdf", "zxcv", "1234", "abcd"}
	lowerPassword := strings.ToLower(password)
	for _, pattern := range keyboardPatterns {
		if strings.Contains(lowerPassword, pattern) {
			patterns = append(patterns, "Keyboard pattern detected: "+pattern)
			break
		}
	}

	// Check for year patterns
	if matched, _ := regexp.MatchString(`(19|20)\d{2}`, password); matched {
		patterns = append(patterns, "Year pattern detected")
	}

	return patterns
}

// calculateScore assigns a score from 0-100
func calculateScore(analysis PasswordAnalysis) int {
	score := 0

	// Length score (max 30 points)
	if analysis.Length >= 12 {
		score += 30
	} else if analysis.Length >= 8 {
		score += 20
	} else if analysis.Length >= 6 {
		score += 10
	}

	// Character variety score (max 40 points)
	if analysis.HasLowercase {
		score += 10
	}
	if analysis.HasUppercase {
		score += 10
	}
	if analysis.HasDigits {
		score += 10
	}
	if analysis.HasSpecialChars {
		score += 10
	}

	// Entropy score (max 20 points)
	if analysis.EntropyBits >= 60 {
		score += 20
	} else if analysis.EntropyBits >= 40 {
		score += 15
	} else if analysis.EntropyBits >= 28 {
		score += 10
	}

	// Deductions
	if analysis.IsCommon {
		score -= 40
	}
	if analysis.ContainsDictionary {
		score -= 10
	}
	if len(analysis.Patterns) > 0 {
		score -= 10 * len(analysis.Patterns)
	}

	// Ensure score is within 0-100
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// determineStrength returns strength level based on score
func determineStrength(score int) string {
	if score >= 80 {
		return "Very Strong"
	} else if score >= 60 {
		return "Strong"
	} else if score >= 40 {
		return "Medium"
	} else if score >= 20 {
		return "Weak"
	}
	return "Very Weak"
}

// estimateCrackTime estimates time to crack password (in seconds)
func estimateCrackTime(entropyBits float64) float64 {
	// Assume 1 billion guesses per second (modern hardware)
	guessesPerSecond := 1e9
	totalCombinations := math.Pow(2, entropyBits)
	// Average case: half of all combinations
	return totalCombinations / (2 * guessesPerSecond)
}

// formatCrackTime converts seconds to readable format
func formatCrackTime(seconds float64) string {
	if seconds < 1 {
		return "Instant"
	} else if seconds < 60 {
		return fmt.Sprintf("%.0f seconds", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.0f minutes", seconds/60)
	} else if seconds < 86400 {
		return fmt.Sprintf("%.0f hours", seconds/3600)
	} else if seconds < 2592000 {
		return fmt.Sprintf("%.0f days", seconds/86400)
	} else if seconds < 31536000 {
		return fmt.Sprintf("%.0f months", seconds/2592000)
	} else if seconds < 315360000 {
		return fmt.Sprintf("%.0f years", seconds/31536000)
	}
	return "Centuries"
}

// generateSuggestions provides improvement recommendations
func generateSuggestions(analysis PasswordAnalysis) []string {
	suggestions := []string{}

	if analysis.Length < 12 {
		suggestions = append(suggestions, "Increase length to at least 12 characters")
	}
	if !analysis.HasLowercase {
		suggestions = append(suggestions, "Add lowercase letters (a-z)")
	}
	if !analysis.HasUppercase {
		suggestions = append(suggestions, "Add uppercase letters (A-Z)")
	}
	if !analysis.HasDigits {
		suggestions = append(suggestions, "Add numbers (0-9)")
	}
	if !analysis.HasSpecialChars {
		suggestions = append(suggestions, "Add special characters (!@#$%^&*)")
	}
	if analysis.IsCommon {
		suggestions = append(suggestions, "Avoid common passwords")
	}
	if analysis.ContainsDictionary {
		suggestions = append(suggestions, "Avoid dictionary words")
	}
	if len(analysis.Patterns) > 0 {
		suggestions = append(suggestions, "Avoid predictable patterns")
	}

	return suggestions
}

// PrintBanner displays the program banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║   Password Strength Checker v1.0     ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintAnalysis displays the detailed password analysis
func PrintAnalysis(analysis PasswordAnalysis) {
	// Determine color code for strength
	var strengthColor string
	switch analysis.Strength {
	case "Very Strong":
		strengthColor = "\033[32m" // Green
	case "Strong":
		strengthColor = "\033[36m" // Cyan
	case "Medium":
		strengthColor = "\033[33m" // Yellow
	case "Weak":
		strengthColor = "\033[31m" // Red
	case "Very Weak":
		strengthColor = "\033[35m" // Magenta
	}
	resetColor := "\033[0m"

	fmt.Println("\n" + "═"*50)
	fmt.Println("PASSWORD ANALYSIS REPORT")
	fmt.Println("═"*50)

	// Mask password for display
	maskedPassword := strings.Repeat("*", len(analysis.Password))
	fmt.Printf("Password: %s (Length: %d)\n", maskedPassword, analysis.Length)
	fmt.Println("─"*50)

	// Character composition
	fmt.Println("Character Composition:")
	fmt.Printf("  Lowercase Letters: %s\n", boolToStatus(analysis.HasLowercase))
	fmt.Printf("  Uppercase Letters: %s\n", boolToStatus(analysis.HasUppercase))
	fmt.Printf("  Digits: %s\n", boolToStatus(analysis.HasDigits))
	fmt.Printf("  Special Characters: %s\n", boolToStatus(analysis.HasSpecialChars))
	fmt.Println("─"*50)

	// Strength metrics
	fmt.Println("Strength Metrics:")
	fmt.Printf("  Entropy: %.2f bits\n", analysis.EntropyBits)
	fmt.Printf("  Score: %d/100\n", analysis.Score)
	fmt.Printf("  Strength: %s%s%s\n", strengthColor, analysis.Strength, resetColor)
	fmt.Printf("  Estimated Crack Time: %s\n", analysis.CrackTimeReadable)
	fmt.Println("─"*50)

	// Warnings
	if analysis.IsCommon {
		fmt.Println("⚠️  WARNING: This is a commonly used password!")
	}
	if analysis.ContainsDictionary {
		fmt.Println("⚠️  WARNING: Contains dictionary words!")
	}
	if len(analysis.Patterns) > 0 {
		fmt.Println("⚠️  Patterns Detected:")
		for _, pattern := range analysis.Patterns {
			fmt.Printf("     - %s\n", pattern)
		}
	}
	if len(analysis.Patterns) > 0 || analysis.IsCommon || analysis.ContainsDictionary {
		fmt.Println("─"*50)
	}

	// Suggestions
	if len(analysis.Suggestions) > 0 {
		fmt.Println("💡 Suggestions for Improvement:")
		for i, suggestion := range analysis.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
		fmt.Println("─"*50)
	}

	// Progress bar
	fmt.Printf("Strength: [")
	bars := analysis.Score / 5
	for i := 0; i < 20; i++ {
		if i < bars {
			fmt.Print("█")
		} else {
			fmt.Print("░")
		}
	}
	fmt.Printf("] %d%%\n", analysis.Score)
	fmt.Println("═"*50)
}

// boolToStatus converts boolean to readable status
func boolToStatus(b bool) string {
	if b {
		return "✓ Present"
	}
	return "✗ Missing"
}

func main() {
	PrintBanner()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter a password to analyze (input is hidden):")
	fmt.Println("Note: Your password will be masked for security.")
	fmt.Print("\nPassword: ")

	// Read password
	password, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}

	// Clean up input
	password = strings.TrimSpace(password)

	if password == "" {
		fmt.Println("Error: Password cannot be empty")
		os.Exit(1)
	}

	// Analyze password
	analysis := AnalyzePassword(password)

	// Print results
	PrintAnalysis(analysis)

	// Interactive mode option
	fmt.Print("\nAnalyze another password? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		main()
	}

	fmt.Println("\nThank you for using Password Strength Checker!")
}