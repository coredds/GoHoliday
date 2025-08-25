package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/your-username/goholidays/updater"
)

func main() {
	var (
		country     = flag.String("country", "", "Specific country to sync (e.g., US, GB, CA, AU)")
		dryRun      = flag.Bool("dry-run", false, "Show what would be synced without making changes")
		verbose     = flag.Bool("verbose", false, "Enable verbose output")
		timeout     = flag.Duration("timeout", 5*time.Minute, "Timeout for sync operation")
		outputDir   = flag.String("output", "./sync_output", "Directory to save synced data")
		listOnly    = flag.Bool("list", false, "Only list available countries")
		validate    = flag.Bool("validate", false, "Validate existing data against Python source")
		force       = flag.Bool("force", false, "Force sync even if data appears up-to-date")
	)
	flag.Parse()

	fmt.Println("GoHolidays Python Sync Tool")
	fmt.Println("===========================")

	// Create GitHub syncer
	syncer := updater.NewGitHubSyncer()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

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

func listCountries(ctx context.Context, syncer *updater.GitHubSyncer) error {
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

func syncSingleCountry(ctx context.Context, syncer *updater.GitHubSyncer, countryCode, outputDir string, dryRun, verbose bool) error {
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
		outputFile := fmt.Sprintf("%s/%s.json", outputDir, countryCode)
		if err := saveCountryData(countryData, outputFile); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}
		fmt.Printf("Data saved to: %s\n", outputFile)
	}
	
	return nil
}

func syncAllCountries(ctx context.Context, syncer *updater.GitHubSyncer, outputDir string, dryRun, verbose, force bool) error {
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

func validateData(ctx context.Context, syncer *updater.GitHubSyncer, dataDir string, verbose bool) error {
	fmt.Println("Validating existing data against Python source...")
	
	// This would compare existing JSON files with fresh Python source
	// Implementation would check for discrepancies and report them
	
	fmt.Println("Validation complete - feature not yet implemented")
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
