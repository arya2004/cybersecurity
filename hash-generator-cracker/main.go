package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashResult stores hash generation results
type HashResult struct {
	Algorithm string
	Hash      string
	Input     string
	Time      time.Duration
}

// CrackResult stores hash cracking results
type CrackResult struct {
	Success   bool
	Password  string
	Attempts  int
	Time      time.Duration
	Algorithm string
}

// Common passwords dictionary for cracking
var commonPasswordsDictionary = []string{
	"password", "123456", "12345678", "qwerty", "abc123", "monkey",
	"1234567", "letmein", "trustno1", "dragon", "baseball", "iloveyou",
	"master", "sunshine", "ashley", "bailey", "passw0rd", "shadow",
	"123123", "654321", "superman", "qazwsx", "michael", "football",
	"admin", "welcome", "login", "password1", "password123", "root",
	"test", "guest", "user", "demo", "changeme", "default", "pass",
	"111111", "000000", "123456789", "qwertyuiop", "ninja", "mustang",
}

// GenerateHash creates hash using specified algorithm
func GenerateHash(input, algorithm string) (HashResult, error) {
	result := HashResult{
		Algorithm: algorithm,
		Input:     input,
	}

	start := time.Now()

	switch strings.ToLower(algorithm) {
	case "md5":
		hash := md5.Sum([]byte(input))
		result.Hash = hex.EncodeToString(hash[:])
	case "sha1":
		hash := sha1.Sum([]byte(input))
		result.Hash = hex.EncodeToString(hash[:])
	case "sha256":
		hash := sha256.Sum256([]byte(input))
		result.Hash = hex.EncodeToString(hash[:])
	case "sha512":
		hash := sha512.Sum512([]byte(input))
		result.Hash = hex.EncodeToString(hash[:])
	case "bcrypt":
		hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
		if err != nil {
			return result, err
		}
		result.Hash = string(hash)
	default:
		return result, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	result.Time = time.Since(start)
	return result, nil
}

// GenerateAllHashes generates hashes using all supported algorithms
func GenerateAllHashes(input string) []HashResult {
	algorithms := []string{"md5", "sha1", "sha256", "sha512", "bcrypt"}
	results := []HashResult{}

	for _, algo := range algorithms {
		result, err := GenerateHash(input, algo)
		if err != nil {
			fmt.Printf("Error generating %s hash: %v\n", algo, err)
			continue
		}
		results = append(results, result)
	}

	return results
}

// CrackHashDictionary attempts dictionary attack on hash
func CrackHashDictionary(targetHash, algorithm string, dictionary []string) CrackResult {
	result := CrackResult{
		Success:   false,
		Algorithm: algorithm,
		Attempts:  0,
	}

	start := time.Now()

	// Special handling for bcrypt (can't reverse, only compare)
	if strings.ToLower(algorithm) == "bcrypt" {
		for _, password := range dictionary {
			result.Attempts++
			err := bcrypt.CompareHashAndPassword([]byte(targetHash), []byte(password))
			if err == nil {
				result.Success = true
				result.Password = password
				result.Time = time.Since(start)
				return result
			}
		}
		result.Time = time.Since(start)
		return result
	}

	// For other algorithms, generate and compare hashes
	for _, password := range dictionary {
		result.Attempts++
		hashResult, err := GenerateHash(password, algorithm)
		if err != nil {
			continue
		}

		if hashResult.Hash == targetHash {
			result.Success = true
			result.Password = password
			result.Time = time.Since(start)
			return result
		}
	}

	result.Time = time.Since(start)
	return result
}

// CrackHashBruteForce attempts simple brute force (numeric only, for demo)
func CrackHashBruteForce(targetHash, algorithm string, maxLength int) CrackResult {
	result := CrackResult{
		Success:   false,
		Algorithm: algorithm,
		Attempts:  0,
	}

	start := time.Now()
	charset := "0123456789" // Numeric only for demo (add more for real use)

	// Try passwords of increasing length
	for length := 1; length <= maxLength; length++ {
		if tryBruteForceLength(targetHash, algorithm, charset, length, &result) {
			result.Time = time.Since(start)
			return result
		}
	}

	result.Time = time.Since(start)
	return result
}

// tryBruteForceLength helper function for brute force
func tryBruteForceLength(targetHash, algorithm, charset string, length int, result *CrackResult) bool {
	return tryBruteForceRecursive(targetHash, algorithm, charset, "", length, result)
}

// tryBruteForceRecursive recursive helper for brute force
func tryBruteForceRecursive(targetHash, algorithm, charset, current string, remaining int, result *CrackResult) bool {
	if remaining == 0 {
		result.Attempts++
		hashResult, err := GenerateHash(current, algorithm)
		if err != nil {
			return false
		}

		if hashResult.Hash == targetHash {
			result.Success = true
			result.Password = current
			return true
		}
		return false
	}

	for _, char := range charset {
		if tryBruteForceRecursive(targetHash, algorithm, charset, current+string(char), remaining-1, result) {
			return true
		}
	}

	return false
}

// PrintBanner displays the program banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║   Hash Generator & Cracker v1.0      ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintMenu displays the main menu
func PrintMenu() {
	fmt.Println("\n" + "═"*50)
	fmt.Println("MAIN MENU")
	fmt.Println("═"*50)
	fmt.Println("1. Generate Hash")
	fmt.Println("2. Generate All Hashes")
	fmt.Println("3. Crack Hash (Dictionary Attack)")
	fmt.Println("4. Crack Hash (Brute Force - Numeric)")
	fmt.Println("5. Compare Hash")
	fmt.Println("6. Exit")
	fmt.Println("═"*50)
	fmt.Print("Select option: ")
}

// PrintHashResult displays hash generation result
func PrintHashResult(result HashResult) {
	fmt.Println("\n" + "═"*50)
	fmt.Println("HASH GENERATION RESULT")
	fmt.Println("═"*50)
	fmt.Printf("Algorithm: %s\n", strings.ToUpper(result.Algorithm))
	fmt.Printf("Input: %s\n", result.Input)
	fmt.Printf("Hash: %s\n", result.Hash)
	fmt.Printf("Time: %v\n", result.Time)
	fmt.Println("═"*50)
}

// PrintAllHashResults displays multiple hash results
func PrintAllHashResults(results []HashResult) {
	fmt.Println("\n" + "═"*50)
	fmt.Println("ALL HASHES GENERATED")
	fmt.Println("═"*50)
	fmt.Printf("Input: %s\n", results[0].Input)
	fmt.Println("─"*50)

	for _, result := range results {
		fmt.Printf("\n%-10s: %s\n", strings.ToUpper(result.Algorithm), result.Hash)
		fmt.Printf("%-10s  Time: %v\n", "", result.Time)
	}
	fmt.Println("═"*50)
}

// PrintCrackResult displays hash cracking result
func PrintCrackResult(result CrackResult) {
	fmt.Println("\n" + "═"*50)
	fmt.Println("HASH CRACKING RESULT")
	fmt.Println("═"*50)
	fmt.Printf("Algorithm: %s\n", strings.ToUpper(result.Algorithm))
	fmt.Printf("Attempts: %d\n", result.Attempts)
	fmt.Printf("Time: %v\n", result.Time)
	fmt.Println("─"*50)

	if result.Success {
		fmt.Printf("✓ SUCCESS! Password found: %s\n", result.Password)
	} else {
		fmt.Println("✗ FAILED: Password not found")
	}
	fmt.Println("═"*50)
}

// HandleGenerateHash handles hash generation
func HandleGenerateHash(reader *bufio.Reader) {
	fmt.Print("\nEnter text to hash: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	fmt.Print("Select algorithm (md5/sha1/sha256/sha512/bcrypt): ")
	algorithm, _ := reader.ReadString('\n')
	algorithm = strings.TrimSpace(strings.ToLower(algorithm))

	result, err := GenerateHash(input, algorithm)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	PrintHashResult(result)
}

// HandleGenerateAllHashes handles generating all hashes
func HandleGenerateAllHashes(reader *bufio.Reader) {
	fmt.Print("\nEnter text to hash: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	results := GenerateAllHashes(input)
	PrintAllHashResults(results)
}

// HandleCrackDictionary handles dictionary attack
func HandleCrackDictionary(reader *bufio.Reader) {
	fmt.Println("\n⚠️  Educational purposes only! Only crack hashes you own.")
	fmt.Print("\nEnter hash to crack: ")
	targetHash, _ := reader.ReadString('\n')
	targetHash = strings.TrimSpace(targetHash)

	fmt.Print("Select algorithm (md5/sha1/sha256/sha512/bcrypt): ")
	algorithm, _ := reader.ReadString('\n')
	algorithm = strings.TrimSpace(strings.ToLower(algorithm))

	fmt.Println("\nStarting dictionary attack...")
	fmt.Printf("Dictionary size: %d passwords\n", len(commonPasswordsDictionary))

	result := CrackHashDictionary(targetHash, algorithm, commonPasswordsDictionary)
	PrintCrackResult(result)
}

// HandleCrackBruteForce handles brute force attack
func HandleCrackBruteForce(reader *bufio.Reader) {
	fmt.Println("\n⚠️  Warning: Brute force is slow! Only for educational purposes.")
	fmt.Println("Note: This demo only tries numeric passwords (0-9)")
	fmt.Print("\nEnter hash to crack: ")
	targetHash, _ := reader.ReadString('\n')
	targetHash = strings.TrimSpace(targetHash)

	fmt.Print("Select algorithm (md5/sha1/sha256): ")
	algorithm, _ := reader.ReadString('\n')
	algorithm = strings.TrimSpace(strings.ToLower(algorithm))

	fmt.Print("Maximum password length to try (1-6 recommended): ")
	var maxLength int
	fmt.Scanln(&maxLength)

	if maxLength > 6 {
		fmt.Println("Warning: Length > 6 may take very long time!")
	}

	fmt.Println("\nStarting brute force attack (numeric only)...")
	result := CrackHashBruteForce(targetHash, algorithm, maxLength)
	PrintCrackResult(result)
}

// HandleCompareHash handles hash comparison
func HandleCompareHash(reader *bufio.Reader) {
	fmt.Print("\nEnter password to verify: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	fmt.Print("Enter hash to compare: ")
	targetHash, _ := reader.ReadString('\n')
	targetHash = strings.TrimSpace(targetHash)

	fmt.Print("Select algorithm (md5/sha1/sha256/sha512/bcrypt): ")
	algorithm, _ := reader.ReadString('\n')
	algorithm = strings.TrimSpace(strings.ToLower(algorithm))

	fmt.Println("\n" + "═"*50)
	fmt.Println("HASH COMPARISON")
	fmt.Println("═"*50)

	if strings.ToLower(algorithm) == "bcrypt" {
		err := bcrypt.CompareHashAndPassword([]byte(targetHash), []byte(password))
		if err == nil {
			fmt.Println("✓ MATCH: Password matches the hash!")
		} else {
			fmt.Println("✗ NO MATCH: Password does not match the hash")
		}
	} else {
		result, err := GenerateHash(password, algorithm)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Generated Hash: %s\n", result.Hash)
		fmt.Printf("Target Hash:    %s\n", targetHash)
		fmt.Println("─"*50)

		if result.Hash == targetHash {
			fmt.Println("✓ MATCH: Password matches the hash!")
		} else {
			fmt.Println("✗ NO MATCH: Password does not match the hash")
		}
	}
	fmt.Println("═"*50)
}

func main() {
	PrintBanner()

	fmt.Println("\n⚠️  DISCLAIMER:")
	fmt.Println("This tool is for educational and authorized testing only.")
	fmt.Println("Only test systems and hashes you own or have permission to test.")

	reader := bufio.NewReader(os.Stdin)

	for {
		PrintMenu()

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			reader.ReadString('\n') // Clear buffer
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			HandleGenerateHash(reader)
		case 2:
			HandleGenerateAllHashes(reader)
		case 3:
			HandleCrackDictionary(reader)
		case 4:
			HandleCrackBruteForce(reader)
		case 5:
			HandleCompareHash(reader)
		case 6:
			fmt.Println("\nThank you for using Hash Generator & Cracker!")
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please select 1-6.")
		}
	}
}