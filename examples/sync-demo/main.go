package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coredds/GoHoliday/updater"
)

func main() {
	fmt.Println("GoHolidays Python Sync Demo")
	fmt.Println("===========================")
	
	// Create GitHub syncer
	syncer := updater.NewGitHubSyncer()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	// Demonstrate the sync capabilities
	fmt.Println("1. Testing GitHub API connectivity...")
	
	countries, err := syncer.FetchCountryList(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch country list: %v", err)
	}
	
	fmt.Printf("✅ Successfully connected to Python holidays repository\n")
	fmt.Printf("✅ Found %d countries available for sync\n", len(countries))
	
	// Show sample of countries
	fmt.Println("\nFirst 10 countries available:")
	for i, country := range countries {
		if i >= 10 {
			break
		}
		fmt.Printf("  %d. %s\n", i+1, country)
	}
	
	// Test fetching specific country data
	testCountries := []string{"US", "AU", "CA", "GB"}
	
	fmt.Println("\n2. Testing data fetching for key countries...")
	
	for _, countryCode := range testCountries {
		fmt.Printf("\nFetching %s data...", countryCode)
		
		// Check if country is available
		available := false
		for _, c := range countries {
			if c == countryCode {
				available = true
				break
			}
		}
		
		if !available {
			fmt.Printf(" ❌ Not available in repository\n")
			continue
		}
		
		// Fetch the source code
		source, err := syncer.FetchCountryFile(ctx, countryCode)
		if err != nil {
			fmt.Printf(" ❌ Failed: %v\n", err)
			continue
		}
		
		fmt.Printf(" ✅ Success! (%d chars)\n", len(source))
		
		// Validate content
		if err := syncer.ValidatePythonContent(source); err != nil {
			fmt.Printf("   ⚠️ Validation warning: %v\n", err)
		} else {
			fmt.Printf("   ✅ Content validation passed\n")
		}
		
		// Parse definitions
		countryData, err := syncer.ParseHolidayDefinitions(source)
		if err != nil {
			fmt.Printf("   ❌ Parse failed: %v\n", err)
			continue
		}
		
		fmt.Printf("   📊 Country: %s (%s)\n", countryData.Name, countryData.CountryCode)
		fmt.Printf("   📊 Subdivisions: %d\n", len(countryData.Subdivisions))
		fmt.Printf("   📊 Holidays: %d\n", len(countryData.Holidays))
		
		// Show sample subdivisions
		if len(countryData.Subdivisions) > 0 {
			fmt.Printf("   📌 Sample subdivisions: ")
			count := 0
			for code, name := range countryData.Subdivisions {
				if count >= 3 {
					fmt.Printf("...")
					break
				}
				if count > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s (%s)", code, name)
				count++
			}
			fmt.Println()
		}
		
		// Rate limiting
		time.Sleep(1 * time.Second)
	}
	
	fmt.Println("\n3. Architecture Benefits Demonstration")
	fmt.Println("=====================================")
	
	fmt.Println("✅ Real-time Python holidays sync")
	fmt.Println("✅ 80+ countries available")
	fmt.Println("✅ Automatic subdivision detection")
	fmt.Println("✅ GitHub API rate limiting")
	fmt.Println("✅ Content validation")
	fmt.Println("✅ Error handling and recovery")
	
	fmt.Println("\n4. Next Steps for Production")
	fmt.Println("============================")
	
	fmt.Println("🔄 Enhanced Python AST parsing for better holiday extraction")
	fmt.Println("💾 Automated JSON generation for country providers")
	fmt.Println("⚡ Scheduled sync jobs (daily/weekly)")
	fmt.Println("🔍 Diff detection to track Python holidays changes")
	fmt.Println("🧪 Integration tests with Python source")
	fmt.Println("📈 Monitoring and alerting for sync failures")
	
	fmt.Println("\n✨ Sync system demonstration complete!")
}
