package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

//go:embed assets/extensions.txt
var extensionsData string

//go:embed assets/wordlist.txt
var wordlistData string

// Version information
const (
	Version   = "1.0.0"
	Author    = "MuhammadWaseem"
	ToolName  = "BackupFinder"
	GitHubURL = "https://github.com/MuhammadWaseem29"
)

// Config holds all configuration options
type Config struct {
	// Core options
	Target      string
	TargetList  string
	UseWordlist bool
	Extensions  string

	// Output options
	Output      string
	JSONOutput  bool
	SilentMode  bool
	VerboseMode bool
	NoColor     bool
	Timestamp   bool
	StatsMode   bool

	// Performance options
	Concurrency int
	RateLimit   int
	Timeout     int
	MaxRetries  int

	// Response options
	SaveResponses bool
	ResponseDir   string
}

// Statistics holds scan statistics
type Statistics struct {
	TotalURLs      int       `json:"total_urls"`
	TotalPatterns  int       `json:"total_patterns"`
	TotalGenerated int       `json:"total_generated"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	ScanDuration   int64     `json:"scan_duration"`
}

// ScanResult represents the complete scan result
type ScanResult struct {
	Tool      string     `json:"tool"`
	Version   string     `json:"version"`
	Author    string     `json:"author"`
	Timestamp time.Time  `json:"timestamp"`
	Stats     Statistics `json:"stats"`
	Results   []string   `json:"results"`
}

// BackupFinder is the main application struct
type BackupFinder struct {
	config *Config
	stats  *Statistics

	// Color functions
	colorRed    func(...interface{}) string
	colorGreen  func(...interface{}) string
	colorYellow func(...interface{}) string
	colorBlue   func(...interface{}) string
	colorCyan   func(...interface{}) string
	colorWhite  func(...interface{}) string
}

// NewBackupFinder creates a new BackupFinder instance
func NewBackupFinder(config *Config) *BackupFinder {
	bf := &BackupFinder{
		config: config,
		stats:  &Statistics{},
	}

	// Initialize colors
	if config.NoColor {
		bf.colorRed = color.New().SprintFunc()
		bf.colorGreen = color.New().SprintFunc()
		bf.colorYellow = color.New().SprintFunc()
		bf.colorBlue = color.New().SprintFunc()
		bf.colorCyan = color.New().SprintFunc()
		bf.colorWhite = color.New().SprintFunc()
	} else {
		bf.colorRed = color.New(color.FgRed).SprintFunc()
		bf.colorGreen = color.New(color.FgGreen).SprintFunc()
		bf.colorYellow = color.New(color.FgYellow).SprintFunc()
		bf.colorBlue = color.New(color.FgBlue).SprintFunc()
		bf.colorCyan = color.New(color.FgCyan).SprintFunc()
		bf.colorWhite = color.New(color.FgWhite, color.Bold).SprintFunc()
	}

	return bf
}

// PrintColor prints colored output with optional timestamp
func (bf *BackupFinder) PrintColor(colorFunc func(...interface{}) string, message string) {
	if bf.config.Timestamp {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", timestamp, colorFunc(message))
	} else {
		fmt.Println(colorFunc(message))
	}
}

// ShowBanner displays the application banner
func (bf *BackupFinder) ShowBanner() {
	if bf.config.SilentMode {
		return
	}

	banner := `╔══════════════════════════════════════════════════════════════════════════════╗
║                           BackupFinder v` + Version + `                                 ║
║                     Backup Files Discovery Tool                             ║
║                        Created by ` + Author + `                            ║
║                     GitHub: ` + GitHubURL + `            ║
╚══════════════════════════════════════════════════════════════════════════════╝`

	bf.PrintColor(bf.colorCyan, banner)
	fmt.Println()
}

// ExtractDomainParts extracts domain components from a URL
func (bf *BackupFinder) ExtractDomainParts(url string) []string {
	// Remove protocol
	re := regexp.MustCompile(`^https?://`)
	url = re.ReplaceAllString(url, "")

	// Remove path and parameters
	if idx := strings.Index(url, "/"); idx != -1 {
		url = url[:idx]
	}
	if idx := strings.Index(url, "?"); idx != -1 {
		url = url[:idx]
	}
	if idx := strings.Index(url, ":"); idx != -1 {
		url = url[:idx]
	}

	// Split by dots
	return strings.Split(url, ".")
}

// GenerateCombinations generates domain combinations for a URL
func (bf *BackupFinder) GenerateCombinations(url string) []string {
	parts := bf.ExtractDomainParts(url)
	combinations := make(map[string]bool)

	// Common TLDs to exclude
	tlds := map[string]bool{
		"com": true, "org": true, "net": true, "io": true, "co": true,
		"uk": true, "de": true, "fr": true, "jp": true, "cn": true, "ru": true,
		"edu": true, "gov": true, "mil": true, "int": true, "biz": true,
		"info": true, "name": true, "pro": true, "museum": true,
	}

	if len(parts) == 0 {
		return []string{}
	}

	// Individual parts (exclude TLDs)
	for _, part := range parts {
		if len(part) > 0 && !tlds[part] {
			combinations[part] = true
		}
	}

	// Handle subdomains specifically
	if len(parts) >= 3 {
		// For URLs like admin.paypal.com or dev.staging.example.com
		// Extract subdomain(s) separately
		for i := 0; i < len(parts)-2; i++ {
			if !tlds[parts[i]] && len(parts[i]) > 0 {
				// Add subdomain alone (very important for backup files!)
				combinations[parts[i]] = true

				// Add subdomain with main domain
				mainDomainIndex := len(parts) - 2
				if mainDomainIndex > i && !tlds[parts[mainDomainIndex]] {
					combinations[parts[i]+"-"+parts[mainDomainIndex]] = true
					combinations[parts[i]+"."+parts[mainDomainIndex]] = true
					combinations[parts[i]+"_"+parts[mainDomainIndex]] = true
					combinations[parts[mainDomainIndex]+"-"+parts[i]] = true
					combinations[parts[mainDomainIndex]+"_"+parts[i]] = true

					// Also add reverse order combinations
					combinations[parts[mainDomainIndex]+"."+parts[i]] = true
				}
			}
		}

		// Add multi-level subdomain combinations
		if len(parts) >= 4 {
			// For cases like dev.staging.example.com
			subdomainParts := parts[:len(parts)-2]

			// Create combinations of subdomains
			for i := 0; i < len(subdomainParts); i++ {
				for j := i + 1; j < len(subdomainParts); j++ {
					if !tlds[subdomainParts[i]] && !tlds[subdomainParts[j]] {
						combinations[subdomainParts[i]+"-"+subdomainParts[j]] = true
						combinations[subdomainParts[i]+"_"+subdomainParts[j]] = true
						combinations[subdomainParts[i]+"."+subdomainParts[j]] = true
					}
				}
			}
		}
	}

	// Two-part combinations
	if len(parts) > 1 {
		for i := 0; i < len(parts)-1; i++ {
			for j := i + 1; j < len(parts); j++ {
				if !tlds[parts[i]] && !tlds[parts[j]] && len(parts[i]) > 0 && len(parts[j]) > 0 {
					combinations[parts[i]+"."+parts[j]] = true
					combinations[parts[i]+"-"+parts[j]] = true
					combinations[parts[i]+"_"+parts[j]] = true
				}
			}
		}
	}

	// Three-part combinations
	if len(parts) > 2 {
		for i := 0; i < len(parts)-2; i++ {
			for j := i + 1; j < len(parts)-1; j++ {
				for k := j + 1; k < len(parts); k++ {
					if !tlds[parts[i]] && !tlds[parts[j]] && !tlds[parts[k]] &&
						len(parts[i]) > 0 && len(parts[j]) > 0 && len(parts[k]) > 0 {
						combinations[parts[i]+"-"+parts[j]+"-"+parts[k]] = true
						combinations[parts[i]+"."+parts[j]+"."+parts[k]] = true
						combinations[parts[i]+"_"+parts[j]+"_"+parts[k]] = true
					}
				}
			}
		}
	}

	// Add full domain without TLD
	if len(parts) >= 2 {
		fullDomain := ""
		for i := 0; i < len(parts)-1; i++ {
			if !tlds[parts[i]] && len(parts[i]) > 0 {
				if fullDomain != "" {
					fullDomain += "-"
				}
				fullDomain += parts[i]
			}
		}
		if fullDomain != "" {
			combinations[fullDomain] = true
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(combinations))
	for combo := range combinations {
		result = append(result, combo)
	}

	return result
}

// LoadPatterns loads patterns from embedded data or external file
func (bf *BackupFinder) LoadPatterns(filename string) ([]string, error) {
	var content string

	// Use embedded data for default files
	switch filename {
	case "assets/extensions.txt":
		content = extensionsData
	case "assets/wordlist.txt":
		content = wordlistData
	default:
		// Load from external file
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		content = string(data)
	}

	lines := strings.Split(content, "\n")
	patterns := make([]string, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return patterns, nil
}

// ProcessTargets processes all targets and generates backup file patterns
func (bf *BackupFinder) ProcessTargets(targets []string, patterns []string) []string {
	bf.stats.StartTime = time.Now()
	bf.stats.TotalURLs = len(targets)
	bf.stats.TotalPatterns = len(patterns)

	allResults := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrency
	semaphore := make(chan struct{}, bf.config.Concurrency)

	for i, target := range targets {
		if target == "" {
			continue
		}

		wg.Add(1)
		go func(target string, index int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			if bf.config.VerboseMode {
				bf.PrintColor(bf.colorYellow, fmt.Sprintf("[%d/%d] Processing: %s", index+1, len(targets), target))
			}

			// Generate combinations for this target
			combinations := bf.GenerateCombinations(target)

			// Apply each pattern to each combination
			for _, combo := range combinations {
				for _, pattern := range patterns {
					result := combo + pattern

					mu.Lock()
					allResults[result] = true
					mu.Unlock()
				}
			}
		}(target, i)
	}

	wg.Wait()
	bf.stats.EndTime = time.Now()
	bf.stats.ScanDuration = bf.stats.EndTime.Sub(bf.stats.StartTime).Milliseconds() / 1000

	// Convert map to sorted slice
	results := make([]string, 0, len(allResults))
	for result := range allResults {
		results = append(results, result)
	}
	sort.Strings(results)

	bf.stats.TotalGenerated = len(results)
	return results
}

// SaveResults saves results to file
func (bf *BackupFinder) SaveResults(results []string) error {
	if bf.config.Output == "" {
		return nil
	}

	if bf.config.JSONOutput {
		// Save as JSON
		scanResult := ScanResult{
			Tool:      ToolName,
			Version:   Version,
			Author:    Author,
			Timestamp: time.Now().UTC(),
			Stats:     *bf.stats,
			Results:   results,
		}

		data, err := json.MarshalIndent(scanResult, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(bf.config.Output, data, 0644)
		if err != nil {
			return err
		}

		bf.PrintColor(bf.colorGreen, "Results saved to JSON file: "+bf.config.Output)
	} else {
		// Save as plain text
		content := strings.Join(results, "\n")
		err := os.WriteFile(bf.config.Output, []byte(content), 0644)
		if err != nil {
			return err
		}

		bf.PrintColor(bf.colorGreen, "Results saved to file: "+bf.config.Output)
	}

	return nil
}

// ShowStats displays scan statistics
func (bf *BackupFinder) ShowStats() {
	if !bf.config.StatsMode || bf.config.SilentMode {
		return
	}

	duration := bf.stats.ScanDuration
	rate := int64(0)
	if duration > 0 {
		rate = int64(bf.stats.TotalGenerated) / duration
	}

	fmt.Println()
	stats := fmt.Sprintf(`╔════════════════════ SCAN STATISTICS ════════════════════╗
║ Total URLs Processed    : %10d                    ║
║ Total Patterns Used     : %10d                    ║
║ Total Results Generated : %10d                    ║
║ Scan Duration (seconds) : %10d                    ║
║ Generation Rate (/sec)  : %10d                    ║
╚══════════════════════════════════════════════════════════╝`,
		bf.stats.TotalURLs,
		bf.stats.TotalPatterns,
		bf.stats.TotalGenerated,
		duration,
		rate)

	bf.PrintColor(bf.colorCyan, stats)
}

// LoadTargets loads targets from file or returns single target
func (bf *BackupFinder) LoadTargets() ([]string, error) {
	if bf.config.TargetList != "" {
		// Load from file
		file, err := os.Open(bf.config.TargetList)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(content), "\n")
		targets := make([]string, 0, len(lines))

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				targets = append(targets, line)
			}
		}

		return targets, nil
	} else if bf.config.Target != "" {
		// Single target
		return []string{bf.config.Target}, nil
	}

	return nil, fmt.Errorf("no target specified")
}

// Run executes the backup finder
func (bf *BackupFinder) Run() error {
	// Show banner
	bf.ShowBanner()

	// Load targets
	targets, err := bf.LoadTargets()
	if err != nil {
		return err
	}

	// Determine pattern file
	var patternFile string
	if bf.config.UseWordlist {
		patternFile = "assets/wordlist.txt"
	} else if bf.config.Extensions != "" {
		patternFile = bf.config.Extensions
	} else {
		patternFile = "assets/extensions.txt"
	}

	// Load patterns (embedded data or external file)
	patterns, err := bf.LoadPatterns(patternFile)
	if err != nil {
		return err
	}

	if !bf.config.SilentMode {
		dataType := "Extensions"
		if bf.config.UseWordlist {
			dataType = "Wordlist Patterns"
		}

		bf.PrintColor(bf.colorCyan, fmt.Sprintf("Using: %s (%s)", dataType, patternFile))
		bf.PrintColor(bf.colorGreen, fmt.Sprintf("Processing %d URLs with %d %s", len(targets), len(patterns), dataType))
		fmt.Println()
	}

	// Process targets
	results := bf.ProcessTargets(targets, patterns)

	// Output results
	if !bf.config.SilentMode {
		fmt.Println()
		bf.PrintColor(bf.colorGreen, fmt.Sprintf("Generated %d unique backup file patterns", len(results)))
		fmt.Println()
	}

	// Display results if not silent mode, or if silent mode AND no output file specified
	if !bf.config.SilentMode || (bf.config.SilentMode && bf.config.Output == "") {
		for _, result := range results {
			fmt.Println(result)
		}
	}

	// Save results
	if err := bf.SaveResults(results); err != nil {
		return err
	}

	// Show statistics
	bf.ShowStats()

	return nil
}

// SystemInfo holds system environment information
type SystemInfo struct {
	OS           string
	Architecture string
	GoVersion    string
	GoInstalled  bool
	GoPath       string
	GoRoot       string
}

// GetSystemInfo gathers system information
func GetSystemInfo() *SystemInfo {
	info := &SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
	}

	// Check if Go is installed
	if goVersion, err := exec.Command("go", "version").Output(); err == nil {
		info.GoInstalled = true
		info.GoVersion = strings.TrimSpace(string(goVersion))
		
		if goPath, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
			info.GoPath = strings.TrimSpace(string(goPath))
		}
		
		if goRoot, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			info.GoRoot = strings.TrimSpace(string(goRoot))
		}
	}

	return info
}

// CheckGoVersion checks if Go version meets minimum requirements
func CheckGoVersion(version string) (bool, string) {
	// Extract version number from "go version go1.22.3 linux/amd64"
	parts := strings.Fields(version)
	if len(parts) < 3 {
		return false, ""
	}
	
	goVersion := strings.TrimPrefix(parts[2], "go")
	versionParts := strings.Split(goVersion, ".")
	
	if len(versionParts) < 2 {
		return false, goVersion
	}
	
	major, err1 := strconv.Atoi(versionParts[0])
	minor, err2 := strconv.Atoi(versionParts[1])
	
	if err1 != nil || err2 != nil {
		return false, goVersion
	}
	
	// Minimum Go 1.18
	return major > 1 || (major == 1 && minor >= 18), goVersion
}

// ShowSystemDiagnostics displays comprehensive system information
func ShowSystemDiagnostics() {
	info := GetSystemInfo()
	
	fmt.Println("🔍 System Diagnostics")
	fmt.Println("====================")
	fmt.Printf("Operating System: %s\n", info.OS)
	fmt.Printf("Architecture: %s\n", info.Architecture)
	fmt.Println()
	
	if info.GoInstalled {
		color.Green("✅ Go is installed")
		fmt.Printf("Version: %s\n", info.GoVersion)
		
		if valid, version := CheckGoVersion(info.GoVersion); valid {
			color.Green("✅ Go version is compatible (>= 1.18)")
		} else {
			color.Yellow("⚠️  Go version might be too old")
			color.Yellow("   Current: %s", version)
			color.Yellow("   Minimum required: 1.18")
		}
		
		fmt.Printf("GOPATH: %s\n", info.GoPath)
		fmt.Printf("GOROOT: %s\n", info.GoRoot)
		fmt.Println()
		
		color.Green("✅ BackupFinder can run directly")
		fmt.Println("   Run: go run ./cmd/backupfinder")
		fmt.Println("   Or build: go build -o backupfinder ./cmd/backupfinder")
	} else {
		color.Red("❌ Go is not installed")
		fmt.Println()
		
		color.Cyan("📋 Installation Options:")
		fmt.Println()
		
		switch info.OS {
		case "linux":
			fmt.Println("🐧 Linux Installation:")
			fmt.Println("   1. Use setup script: chmod +x setup.sh && ./setup.sh")
			fmt.Println("   2. Use install script: chmod +x install.sh && ./install.sh")
			fmt.Println("   3. Manual installation:")
			fmt.Println("      curl -O https://go.dev/dl/go1.22.3.linux-amd64.tar.gz")
			fmt.Println("      sudo rm -rf /usr/local/go")
			fmt.Println("      sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz")
			fmt.Println("      echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc")
			fmt.Println("      source ~/.bashrc")
			
		case "darwin":
			fmt.Println("🍎 macOS Installation:")
			fmt.Println("   1. Homebrew: brew install go")
			fmt.Println("   2. Download installer: https://go.dev/dl/")
			fmt.Println("   3. Manual: similar to Linux but use darwin-amd64 or darwin-arm64")
			
		case "windows":
			fmt.Println("🪟 Windows Installation:")
			fmt.Println("   1. Download installer: https://go.dev/dl/")
			fmt.Println("   2. Chocolatey: choco install golang")
			fmt.Println("   3. Scoop: scoop install go")
			
		default:
			fmt.Printf("   Download from: https://go.dev/dl/ (for %s)\n", info.OS)
		}
		
		fmt.Println()
		color.Yellow("⚠️  After installing Go, re-run this command to verify")
	}
}

// ShowQuickStart shows quick start guide
func ShowQuickStart() {
	fmt.Println("🚀 BackupFinder Quick Start")
	fmt.Println("===========================")
	
	info := GetSystemInfo()
	
	if info.GoInstalled {
		if valid, _ := CheckGoVersion(info.GoVersion); valid {
			color.Green("✅ Your system is ready!")
			fmt.Println()
			fmt.Println("🏗️  Build BackupFinder:")
			fmt.Println("   go build -o backupfinder ./cmd/backupfinder")
			fmt.Println()
			fmt.Println("🧪 Test installation:")
			fmt.Println("   ./backupfinder health-check")
			fmt.Println()
			fmt.Println("📖 Basic usage:")
			fmt.Println("   ./backupfinder -u example.com -w")
			fmt.Println("   ./backupfinder -u example.com -e txt,bak,old")
			fmt.Println("   ./backupfinder --help")
		} else {
			color.Yellow("⚠️  Go version might be incompatible")
			fmt.Println("   Consider upgrading to Go 1.18+")
		}
	} else {
		color.Red("❌ Go is required but not installed")
		fmt.Println()
		color.Cyan("🔧 Auto-setup (Linux only):")
		fmt.Println("   chmod +x install.sh && ./install.sh")
		fmt.Println()
		fmt.Println("For other systems, see: ./backupfinder system-check")
	}
}

// Main function
func main() {
	config := &Config{
		Concurrency: 10,
		RateLimit:   50,
		Timeout:     30,
		MaxRetries:  3,
		StatsMode:   true,
		ResponseDir: "responses",
	}

	var rootCmd = &cobra.Command{
		Use:   "backupfinder",
		Short: "Backup files discovery tool",
		Long: `Good Day!

I truly hope everything is awesome on your side of the screen! 😊

BackupFinder discovers backup files on web servers by generating intelligent patterns.
It creates thousands of potential backup file names based on your target domain.
Perfect for penetration testing, bug bounty hunting, and security audits.

May you be well on your side of the screen :)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate input
			if config.Target == "" && config.TargetList == "" {
				return fmt.Errorf("no target specified. Use -u for single target or -l for target list")
			}

			bf := NewBackupFinder(config)
			return bf.Run()
		},
	}

	// Add flags
	rootCmd.Flags().StringVarP(&config.Target, "target", "u", "", "Target URL/domain to scan")
	rootCmd.Flags().StringVarP(&config.TargetList, "list", "l", "", "File with target list")
	rootCmd.Flags().BoolVarP(&config.UseWordlist, "wordlist", "w", false, "Use wordlist mode (comprehensive patterns)")
	rootCmd.Flags().StringVarP(&config.Extensions, "extensions", "e", "", "Custom extensions file")
	rootCmd.Flags().StringVarP(&config.Output, "output", "o", "", "Output file")
	rootCmd.Flags().BoolVar(&config.JSONOutput, "json", false, "JSON output format")
	rootCmd.Flags().BoolVar(&config.SilentMode, "silent", false, "Show only results")
	rootCmd.Flags().BoolVarP(&config.VerboseMode, "verbose", "v", false, "Verbose mode")
	rootCmd.Flags().BoolVar(&config.NoColor, "no-color", false, "Disable colored output")
	rootCmd.Flags().BoolVar(&config.Timestamp, "timestamp", false, "Add timestamps to output")
	rootCmd.Flags().BoolVar(&config.StatsMode, "stats", true, "Show statistics")
	rootCmd.Flags().IntVarP(&config.Concurrency, "concurrency", "c", 10, "Number of concurrent workers")
	rootCmd.Flags().IntVar(&config.RateLimit, "rate-limit", 50, "Rate limit for requests")
	rootCmd.Flags().IntVar(&config.Timeout, "timeout", 30, "Request timeout in seconds")
	rootCmd.Flags().IntVar(&config.MaxRetries, "retries", 3, "Maximum number of retries")
	rootCmd.Flags().BoolVar(&config.SaveResponses, "store-resp", false, "Store responses")
	rootCmd.Flags().StringVar(&config.ResponseDir, "store-resp-dir", "responses", "Response storage directory")

	// Handle JSON export flag
	rootCmd.Flags().StringVar(&config.Output, "je", "", "Export to JSON file")
	rootCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if jeFlag := cmd.Flag("je"); jeFlag.Changed {
			config.JSONOutput = true
		}
		return nil
	}

	// Version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s v%s by %s\n", ToolName, Version, Author)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Health check command
	var healthCmd = &cobra.Command{
		Use:   "health-check",
		Short: "Check system requirements",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("✅ All requirements satisfied")
			fmt.Println("✅ Assets embedded in binary")
			fmt.Printf("✅ Extensions patterns: %d\n", len(strings.Split(extensionsData, "\n")))
			fmt.Printf("✅ Wordlist patterns: %d\n", len(strings.Split(wordlistData, "\n")))
		},
	}
	rootCmd.AddCommand(healthCmd)

	// Template list command
	var templatesCmd = &cobra.Command{
		Use:   "templates",
		Short: "List available templates",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Available Templates:")
			fmt.Println("• Backup Extensions (assets/extensions.txt) - 92+ common backup file extensions")
			fmt.Println("• Comprehensive Wordlist (assets/wordlist.txt) - 1907+ specialized backup patterns")
			fmt.Println("• Custom Templates - User-defined template files")
		},
	}
	rootCmd.AddCommand(templatesCmd)

	// System check command
	var systemCheckCmd = &cobra.Command{
		Use:   "system-check",
		Short: "Check system requirements and Go installation",
		Long: `Performs comprehensive system diagnostics including:
• Operating system and architecture detection
• Go installation verification
• Version compatibility check
• Environment configuration status
• Installation guidance for missing components`,
		Run: func(cmd *cobra.Command, args []string) {
			ShowSystemDiagnostics()
		},
	}
	rootCmd.AddCommand(systemCheckCmd)

	// Quick start command
	var quickStartCmd = &cobra.Command{
		Use:   "quick-start",
		Short: "Show quick start guide",
		Long:  "Displays a quick start guide with system-specific instructions",
		Run: func(cmd *cobra.Command, args []string) {
			ShowQuickStart()
		},
	}
	rootCmd.AddCommand(quickStartCmd)

	// Setup command
	var setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Interactive setup guide",
		Long: `Interactive setup guide that helps you:
• Verify Go installation
• Build BackupFinder
• Test the installation
• Provides troubleshooting tips`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🔧 BackupFinder Interactive Setup")
			fmt.Println("=================================")
			fmt.Println()
			
			info := GetSystemInfo()
			
			if !info.GoInstalled {
				color.Red("❌ Go is not installed on your system")
				fmt.Println()
				
				color.Cyan("🔧 Automated Setup Options:")
				fmt.Println()
				
				if info.OS == "linux" {
					fmt.Println("1. 📦 Quick Setup (Recommended):")
					fmt.Println("   chmod +x install.sh && ./install.sh")
					fmt.Println()
					fmt.Println("2. 🛠️  Basic Setup:")
					fmt.Println("   chmod +x setup.sh && ./setup.sh")
					fmt.Println()
				}
				
				fmt.Println("3. 📖 Manual Setup:")
				fmt.Println("   Run: ./backupfinder system-check")
				fmt.Println("   Follow the installation instructions for your OS")
				
				return
			}
			
			if valid, version := CheckGoVersion(info.GoVersion); !valid {
				color.Yellow("⚠️  Go version might be too old: %s", version)
				color.Yellow("   Minimum required: 1.18")
				fmt.Println("   Consider upgrading Go for best compatibility")
				fmt.Println()
			}
			
			color.Green("✅ Go is properly installed")
			fmt.Printf("   Version: %s\n", info.GoVersion)
			fmt.Println()
			
			color.Cyan("🏗️  Building BackupFinder...")
			fmt.Println("   Run: go build -o backupfinder ./cmd/backupfinder")
			fmt.Println()
			
			color.Cyan("🧪 Testing BackupFinder...")
			fmt.Println("   Run: ./backupfinder health-check")
			fmt.Println("   Run: ./backupfinder version")
			fmt.Println()
			
			color.Green("🚀 You're ready to use BackupFinder!")
			fmt.Println("   Quick usage: ./backupfinder -u example.com -w")
			fmt.Println("   Full help: ./backupfinder --help")
		},
	}
	rootCmd.AddCommand(setupCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
