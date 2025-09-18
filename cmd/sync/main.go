package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coredds/GoHoliday/config"
	"github.com/coredds/GoHoliday/updater"
)

func main() {
	var (
		country   = flag.String("country", "", "Specific country to sync (e.g., US, GB, CA, AU)")
		dryRun    = flag.Bool("dry-run", false, "Show what would be synced without making changes")
		verbose   = flag.Bool("verbose", false, "Enable verbose output")
		timeout   = flag.Duration("timeout", 5*time.Minute, "Timeout for sync operation")
		outputDir = flag.String("output", "./sync_output", "Directory to save synced data")
		listOnly  = flag.Bool("list", false, "Only list available countries")
		validate  = flag.Bool("validate", false, "Validate existing data against Python source")
		force     = flag.Bool("force", false, "Force sync even if data appears up-to-date")
		token     = flag.String("token", "", "GitHub Personal Access Token for authentication (optional)")
	)
	flag.Parse()

	fmt.Println("GoHolidays Python Sync Tool")
	fmt.Println("===========================")

	// Get GitHub token from flag, config file, or environment variable
	githubToken := *token
	if githubToken == "" {
		githubToken = config.LoadGitHubToken()
	}

	// Create GitHub syncer with optional token
	var syncer updater.Syncer
	if githubToken != "" {
		syncer = updater.NewGitHubSyncerWithToken(githubToken)
		if *verbose {
			fmt.Println("Using authenticated GitHub API access")
		}
	} else {
		syncer = updater.NewGitHubSyncer()
		if *verbose {
			fmt.Println("Using unauthenticated GitHub API access (rate limited)")
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Validate token if provided
	if githubToken != "" {
		if githubSyncer, ok := syncer.(*updater.GitHubSyncer); ok {
			if err := githubSyncer.ValidateToken(ctx); err != nil {
				log.Fatalf("GitHub token validation failed: %v", err)
			}
			if *verbose {
				fmt.Println("✓ GitHub token validated successfully")
			}
		}
	}

	if *listOnly {
		if err := listCountries(ctx, syncer); err != nil {
			log.Fatalf("Failed to list countries: %v", err)
		}
		return
	}

	if *country != "" {
		if err := syncSingleCountry(ctx, syncer, *country, *outputDir, *dryRun, *verbose); err != nil {
			log.Fatalf("Failed to sync %s: %v", *country, err)
		}
		return
	}

	if *validate {
		if err := validateData(ctx, syncer, *outputDir, *verbose); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		return
	}

	// Default: sync all countries
	if err := syncAllCountries(ctx, syncer, *outputDir, *dryRun, *verbose, *force); err != nil {
		log.Fatalf("Failed to sync: %v", err)
	}
}

func listCountries(ctx context.Context, syncer updater.Syncer) error {
	fmt.Println("Fetching available countries from Python holidays repository...")

	countries, err := syncer.FetchCountryList(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("\nFound %d countries:\n", len(countries))
	fmt.Println("Code | Status")
	fmt.Println("-----|-------")

	for _, country := range countries {
		status := "Available"
		fmt.Printf("%-4s | %s\n", country, status)
	}

	return nil
}

func syncSingleCountry(ctx context.Context, syncer updater.Syncer, countryCode, outputDir string, dryRun, verbose bool) error {
	fmt.Printf("Syncing country: %s\n", countryCode)

	if dryRun {
		fmt.Println("DRY RUN MODE - No files will be modified")
	}

	// Create output directory
	if !dryRun {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Fetch Python source
	if verbose {
		fmt.Printf("Fetching Python source for %s...\n", countryCode)
	}

	pythonSource, err := syncer.FetchCountryFile(ctx, countryCode)
	if err != nil {
		return fmt.Errorf("failed to fetch source: %w", err)
	}

	if verbose {
		fmt.Printf("Source length: %d characters\n", len(pythonSource))
	}

	// Validate content
	if err := syncer.ValidatePythonContent(pythonSource); err != nil {
		return fmt.Errorf("invalid Python content: %w", err)
	}

	// Parse holiday definitions
	if verbose {
		fmt.Println("Parsing holiday definitions...")
	}

	countryData, err := syncer.ParseHolidayDefinitions(pythonSource)
	if err != nil {
		return fmt.Errorf("failed to parse definitions: %w", err)
	}

	// Display results
	fmt.Printf("\nCountry: %s (%s)\n", countryData.Name, countryData.CountryCode)
	fmt.Printf("Subdivisions: %d\n", len(countryData.Subdivisions))
	fmt.Printf("Holidays: %d\n", len(countryData.Holidays))
	fmt.Printf("Categories: %v\n", countryData.Categories)
	fmt.Printf("Languages: %v\n", countryData.Languages)

	if verbose {
		fmt.Println("\nHolidays found:")
		for name, holiday := range countryData.Holidays {
			fmt.Printf("  - %s: %s (%s)\n", name, holiday.Name, holiday.Calculation)
		}

		if len(countryData.Subdivisions) > 0 {
			fmt.Println("\nSubdivisions:")
			for code, name := range countryData.Subdivisions {
				fmt.Printf("  - %s: %s\n", code, name)
			}
		}
	}

	if !dryRun {
		// Save to file
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.json", strings.ToUpper(countryCode)))

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return fmt.Errorf("failed to create directory for country file: %w", err)
		}

		if err := saveCountryData(countryData, outputFile); err != nil {
			return fmt.Errorf("failed to save data to %s: %w", outputFile, err)
		}
		fmt.Printf("Data saved to: %s\n", outputFile)
	}

	return nil
}

func syncAllCountries(ctx context.Context, syncer updater.Syncer, outputDir string, dryRun, verbose, force bool) error {
	fmt.Println("Syncing all available countries...")

	if dryRun {
		fmt.Println("DRY RUN MODE - No files will be modified")
	}

	// Get country list
	countries, err := syncer.FetchCountryList(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d countries to sync\n", len(countries))

	// Create output directory
	if !dryRun {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	successful := 0
	failed := 0

	for i, country := range countries {
		fmt.Printf("\n[%d/%d] Syncing %s...", i+1, len(countries), country)

		if err := syncSingleCountry(ctx, syncer, country, outputDir, dryRun, verbose); err != nil {
			fmt.Printf(" FAILED: %v\n", err)
			failed++
			continue
		}

		fmt.Printf(" SUCCESS\n")
		successful++

		// Rate limiting between countries
		if i < len(countries)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Printf("\nSync completed: %d successful, %d failed\n", successful, failed)
	return nil
}

func validateData(ctx context.Context, syncer updater.Syncer, dataDir string, verbose bool) error {
	fmt.Println("Validating existing data against Python source...")

	// Get list of existing JSON files
	files, err := filepath.Glob(filepath.Join(dataDir, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to list data files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No data files found to validate")
		return nil
	}

	var validationErrors []string
	validatedCount := 0

	for _, file := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Extract country code from filename
		basename := filepath.Base(file)
		countryCode := strings.TrimSuffix(basename, ".json")
		countryCode = strings.ToUpper(countryCode)

		if verbose {
			fmt.Printf("Validating %s...\n", countryCode)
		}

		// Load existing data
		existingData, err := loadExistingData(file)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: failed to load existing data: %v", countryCode, err))
			continue
		}

		// Fetch fresh data from source
		sourceContent, err := syncer.FetchCountryFile(ctx, countryCode)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: failed to fetch source: %v", countryCode, err))
			continue
		}

		freshData, err := syncer.ParseHolidayDefinitions(sourceContent)
		if err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: failed to parse fresh data: %v", countryCode, err))
			continue
		}

		// Compare data
		if err := compareCountryData(existingData, freshData, countryCode, verbose); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("%s: %v", countryCode, err))
		} else {
			validatedCount++
		}
	}

	// Report results
	fmt.Printf("\nValidation Results:\n")
	fmt.Printf("- Files validated: %d\n", validatedCount)
	fmt.Printf("- Validation errors: %d\n", len(validationErrors))

	if len(validationErrors) > 0 {
		fmt.Printf("\nValidation Errors:\n")
		for _, err := range validationErrors {
			fmt.Printf("- %s\n", err)
		}
		return fmt.Errorf("validation failed with %d errors", len(validationErrors))
	}

	fmt.Println("All data files validated successfully")
	return nil
}

func loadExistingData(filename string) (*updater.CountryData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data updater.CountryData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func compareCountryData(existing, fresh *updater.CountryData, countryCode string, verbose bool) error {
	var differences []string

	// Compare basic metadata
	if existing.Name != fresh.Name {
		differences = append(differences, fmt.Sprintf("name changed: %s -> %s", existing.Name, fresh.Name))
	}

	// Compare holiday counts
	if len(existing.Holidays) != len(fresh.Holidays) {
		differences = append(differences, fmt.Sprintf("holiday count changed: %d -> %d", len(existing.Holidays), len(fresh.Holidays)))
	}

	// Compare individual holidays
	for key, existingHoliday := range existing.Holidays {
		freshHoliday, exists := fresh.Holidays[key]
		if !exists {
			differences = append(differences, fmt.Sprintf("holiday removed: %s", key))
			continue
		}

		if existingHoliday.Name != freshHoliday.Name {
			differences = append(differences, fmt.Sprintf("holiday %s name changed: %s -> %s", key, existingHoliday.Name, freshHoliday.Name))
		}

		if existingHoliday.Category != freshHoliday.Category {
			differences = append(differences, fmt.Sprintf("holiday %s category changed: %s -> %s", key, existingHoliday.Category, freshHoliday.Category))
		}
	}

	// Check for new holidays
	for key := range fresh.Holidays {
		if _, exists := existing.Holidays[key]; !exists {
			differences = append(differences, fmt.Sprintf("new holiday added: %s", key))
		}
	}

	if len(differences) > 0 {
		if verbose {
			fmt.Printf("  Differences found:\n")
			for _, diff := range differences {
				fmt.Printf("    - %s\n", diff)
			}
		}
		return fmt.Errorf("found %d differences", len(differences))
	}

	if verbose {
		fmt.Printf("  No differences found\n")
	}
	return nil
}

func saveCountryData(data *updater.CountryData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
