package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coredds/goholiday/config"
	"github.com/coredds/goholiday/updater"
)

func main() {
	// Example 1: Using GitHub token from config loader
	fmt.Println("=== GitHub Token Authentication Example ===")

	// Get token from config loader (checks env var and config file)
	token := config.LoadGitHubToken()
	if token == "" {
		fmt.Println("No GitHub token found (checked GITHUB_TOKEN env var and config/github_token.txt)")
		fmt.Println("Using unauthenticated access (rate limited to 60 requests/hour)")
	} else {
		fmt.Println("Found GitHub token, using authenticated access (5000 requests/hour)")
		fmt.Printf("Token: %s...%s\n", token[:8], token[len(token)-4:])
	}

	// Create syncer with token (empty string = unauthenticated)
	syncer := updater.NewGitHubSyncerWithToken(token)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Validate token if provided
	if token != "" {
		fmt.Println("Validating GitHub token...")
		if err := syncer.ValidateToken(ctx); err != nil {
			log.Fatalf("Token validation failed: %v", err)
		}
		fmt.Println("✓ Token validated successfully!")
	}

	// Example: List available countries
	fmt.Println("\nFetching available countries...")
	countries, err := syncer.FetchCountryList(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch countries: %v", err)
	}

	fmt.Printf("Found %d countries: %v\n", len(countries), countries[:min(5, len(countries))])
	if len(countries) > 5 {
		fmt.Printf("... and %d more\n", len(countries)-5)
	}

	// Example: Fetch a specific country file
	if len(countries) > 0 {
		countryCode := countries[0]
		fmt.Printf("\nFetching Python source for %s...\n", countryCode)

		source, err := syncer.FetchCountryFile(ctx, countryCode)
		if err != nil {
			log.Fatalf("Failed to fetch country file: %v", err)
		}

		fmt.Printf("✓ Successfully fetched %d bytes of Python source\n", len(source))

		// Parse the source
		countryData, err := syncer.ParseHolidayDefinitions(source)
		if err != nil {
			log.Fatalf("Failed to parse holiday definitions: %v", err)
		}

		fmt.Printf("✓ Parsed %d holidays for %s (%s)\n",
			len(countryData.Holidays), countryData.Name, countryData.CountryCode)
	}

	fmt.Println("\n=== Example completed successfully! ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
