package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// FileRecord stores file integrity information
type FileRecord struct {
	Path         string    `json:"path"`
	Hash         string    `json:"hash"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	LastChecked  time.Time `json:"last_checked"`
	Status       string    `json:"status"`
}

// IntegrityDatabase stores baseline and current state
type IntegrityDatabase struct {
	BaselineDate time.Time              `json:"baseline_date"`
	Files        map[string]FileRecord  `json:"files"`
}

const (
	dbFileName = "integrity_baseline.json"
)

var (
	db = &IntegrityDatabase{
		Files: make(map[string]FileRecord),
	}
)

// CalculateFileHash generates SHA256 hash of file
func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

// GetFileInfo retrieves file metadata
func GetFileInfo(filePath string) (int64, time.Time, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, time.Time{}, err
	}
	return info.Size(), info.ModTime(), nil
}

// CreateFileRecord generates a complete file record
func CreateFileRecord(filePath string) (FileRecord, error) {
	record := FileRecord{
		Path:        filePath,
		LastChecked: time.Now(),
		Status:      "OK",
	}

	// Calculate hash
	hash, err := CalculateFileHash(filePath)
	if err != nil {
		return record, err
	}
	record.Hash = hash

	// Get file info
	size, modTime, err := GetFileInfo(filePath)
	if err != nil {
		return record, err
	}
	record.Size = size
	record.ModTime = modTime

	return record, nil
}

// CreateBaseline scans directory and creates integrity baseline
func CreateBaseline(directory string) error {
	fmt.Printf("\n🔍 Creating baseline for: %s\n", directory)
	fmt.Println("─"*50)

	db.BaselineDate = time.Now()
	db.Files = make(map[string]FileRecord)

	fileCount := 0
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("⚠️  Error accessing %s: %v\n", path, err)
			return nil // Continue walking
		}

		// Skip directories and the database file itself
		if info.IsDir() || filepath.Base(path) == dbFileName {
			return nil
		}

		fmt.Printf("📄 Processing: %s\n", path)
		record, err := CreateFileRecord(path)
		if err != nil {
			fmt.Printf("⚠️  Error processing %s: %v\n", path, err)
			return nil
		}

		db.Files[path] = record
		fileCount++
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Println("─"*50)
	fmt.Printf("✓ Baseline created: %d files processed\n", fileCount)
	
	return SaveDatabase()
}

// VerifyIntegrity checks files against baseline
func VerifyIntegrity(directory string) ([]FileRecord, error) {
	fmt.Printf("\n🔍 Verifying integrity for: %s\n", directory)
	fmt.Println("─"*50)

	var changes []FileRecord
	checkedFiles := make(map[string]bool)

	// Check existing files
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() || filepath.Base(path) == dbFileName {
			return nil
		}

		checkedFiles[path] = true
		baseline, exists := db.Files[path]

		if !exists {
			// New file detected
			record, err := CreateFileRecord(path)
			if err != nil {
				return nil
			}
			record.Status = "NEW"
			changes = append(changes, record)
			fmt.Printf("🆕 NEW FILE: %s\n", path)
			return nil
		}

		// Verify existing file
		current, err := CreateFileRecord(path)
		if err != nil {
			return nil
		}

		if current.Hash != baseline.Hash {
			current.Status = "MODIFIED"
			changes = append(changes, current)
			fmt.Printf("✏️  MODIFIED: %s\n", path)
		} else if current.Size != baseline.Size {
			current.Status = "SIZE_CHANGED"
			changes = append(changes, current)
			fmt.Printf("📏 SIZE CHANGED: %s\n", path)
		}

		return nil
	})

	if err != nil {
		return changes, err
	}

	// Check for deleted files
	for path := range db.Files {
		if !checkedFiles[path] {
			record := db.Files[path]
			record.Status = "DELETED"
			record.LastChecked = time.Now()
			changes = append(changes, record)
			fmt.Printf("🗑️  DELETED: %s\n", path)
		}
	}

	fmt.Println("─"*50)
	if len(changes) == 0 {
		fmt.Println("✓ No changes detected - All files intact!")
	} else {
		fmt.Printf("⚠️  %d change(s) detected\n", len(changes))
	}

	return changes, nil
}

// GenerateReport creates detailed integrity report
func GenerateReport(changes []FileRecord) {
	fmt.Println("\n" + "═"*50)
	fmt.Println("INTEGRITY VERIFICATION REPORT")
	fmt.Println("═"*50)
	fmt.Printf("Report Date: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("Baseline Date: %s\n", db.BaselineDate.Format("2006-01-02 15:04:05"))
	fmt.Println("─"*50)

	if len(changes) == 0 {
		fmt.Println("✓ SYSTEM INTEGRITY: INTACT")
		fmt.Println("No unauthorized changes detected.")
		fmt.Println("═"*50)
		return
	}

	// Categorize changes
	newFiles := []FileRecord{}
	modifiedFiles := []FileRecord{}
	deletedFiles := []FileRecord{}
	sizeChanged := []FileRecord{}

	for _, change := range changes {
		switch change.Status {
		case "NEW":
			newFiles = append(newFiles, change)
		case "MODIFIED":
			modifiedFiles = append(modifiedFiles, change)
		case "DELETED":
			deletedFiles = append(deletedFiles, change)
		case "SIZE_CHANGED":
			sizeChanged = append(sizeChanged, change)
		}
	}

	// Summary
	fmt.Println("SUMMARY:")
	fmt.Printf("  New Files: %d\n", len(newFiles))
	fmt.Printf("  Modified Files: %d\n", len(modifiedFiles))
	fmt.Printf("  Deleted Files: %d\n", len(deletedFiles))
	fmt.Printf("  Size Changed: %d\n", len(sizeChanged))
	fmt.Println("─"*50)

	// Details
	if len(newFiles) > 0 {
		fmt.Println("\n🆕 NEW FILES:")
		for _, file := range newFiles {
			fmt.Printf("  • %s\n", file.Path)
			fmt.Printf("    Hash: %s\n", file.Hash)
			fmt.Printf("    Size: %d bytes\n", file.Size)
		}
	}

	if len(modifiedFiles) > 0 {
		fmt.Println("\n✏️  MODIFIED FILES:")
		for _, file := range modifiedFiles {
			baseline := db.Files[file.Path]
			fmt.Printf("  • %s\n", file.Path)
			fmt.Printf("    Old Hash: %s\n", baseline.Hash)
			fmt.Printf("    New Hash: %s\n", file.Hash)
			fmt.Printf("    Modified: %s\n", file.ModTime.Format("2006-01-02 15:04:05"))
		}
	}

	if len(sizeChanged) > 0 {
		fmt.Println("\n📏 SIZE CHANGED FILES:")
		for _, file := range sizeChanged {
			baseline := db.Files[file.Path]
			fmt.Printf("  • %s\n", file.Path)
			fmt.Printf("    Old Size: %d bytes\n", baseline.Size)
			fmt.Printf("    New Size: %d bytes\n", file.Size)
		}
	}

	if len(deletedFiles) > 0 {
		fmt.Println("\n🗑️  DELETED FILES:")
		for _, file := range deletedFiles {
			fmt.Printf("  • %s\n", file.Path)
			fmt.Printf("    Last Hash: %s\n", file.Hash)
		}
	}

	fmt.Println("═"*50)
}

// SaveDatabase persists database to JSON file
func SaveDatabase() error {
	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dbFileName, data, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("💾 Database saved: %s\n", dbFileName)
	return nil
}

// LoadDatabase loads database from JSON file
func LoadDatabase() error {
	data, err := ioutil.ReadFile(dbFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no baseline found (use create command first)")
		}
		return err
	}

	err = json.Unmarshal(data, db)
	if err != nil {
		return err
	}

	fmt.Printf("📂 Loaded baseline: %d files from %s\n", 
		len(db.Files), db.BaselineDate.Format("2006-01-02 15:04:05"))
	return nil
}

// ShowStatus displays current baseline status
func ShowStatus() {
	if len(db.Files) == 0 {
		fmt.Println("\n⚠️  No baseline created yet")
		fmt.Println("Use 'create' command to create a baseline")
		return
	}

	fmt.Println("\n" + "═"*50)
	fmt.Println("BASELINE STATUS")
	fmt.Println("═"*50)
	fmt.Printf("Baseline Date: %s\n", db.BaselineDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total Files: %d\n", len(db.Files))
	fmt.Println("─"*50)

	// Calculate total size
	var totalSize int64
	for _, record := range db.Files {
		totalSize += record.Size
	}

	fmt.Printf("Total Size: %.2f MB\n", float64(totalSize)/(1024*1024))
	fmt.Println("═"*50)
}

// PrintBanner displays program banner
func PrintBanner() {
	banner := `
╔═══════════════════════════════════════╗
║   File Integrity Monitor v1.0        ║
║   Cybersecurity Lab Tool              ║
╚═══════════════════════════════════════╝
`
	fmt.Println(banner)
}

// PrintUsage displays usage information
func PrintUsage() {
	fmt.Println("\nUsage: go run main.go <command> <directory>")
	fmt.Println("\nCommands:")
	fmt.Println("  create <dir>   - Create integrity baseline for directory")
	fmt.Println("  verify <dir>   - Verify directory against baseline")
	fmt.Println("  status         - Show baseline information")
	fmt.Println("\nExamples:")
	fmt.Println("  go run main.go create ./test-files")
	fmt.Println("  go run main.go verify ./test-files")
	fmt.Println("  go run main.go status")
	fmt.Println("\nNote: Baseline is saved to 'integrity_baseline.json'")
}

func main() {
	PrintBanner()

	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("Error: Directory path required")
			PrintUsage()
			os.Exit(1)
		}
		directory := os.Args[2]

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			fmt.Printf("Error: Directory '%s' does not exist\n", directory)
			os.Exit(1)
		}

		err := CreateBaseline(directory)
		if err != nil {
			fmt.Printf("Error creating baseline: %v\n", err)
			os.Exit(1)
		}

	case "verify":
		if len(os.Args) < 3 {
			fmt.Println("Error: Directory path required")
			PrintUsage()
			os.Exit(1)
		}
		directory := os.Args[2]

		err := LoadDatabase()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		changes, err := VerifyIntegrity(directory)
		if err != nil {
			fmt.Printf("Error verifying integrity: %v\n", err)
			os.Exit(1)
		}

		GenerateReport(changes)

	case "status":
		err := LoadDatabase()
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error loading database: %v\n", err)
			os.Exit(1)
		}
		ShowStatus()

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		PrintUsage()
		os.Exit(1)
	}
}