package main

import (
	"fmt"
	"log"
	"time"

	"github.com/your-username/goholidays/config"
)

func main() {
	fmt.Println("GoHolidays Configuration System Demo")
	fmt.Println("====================================")

	// Initialize the holiday manager with configuration
	manager := config.NewHolidayManager()

	// Show supported countries
	fmt.Println("\n1. Supported Countries:")
	countries := manager.GetSupportedCountries()
	for _, country := range countries {
		fmt.Printf("   - %s\n", country)
	}

	// Demonstrate holidays for a specific country
	year := 2024
	countryCode := "US"
	
	fmt.Printf("\n2. Holidays for %s in %d:\n", countryCode, year)
	holidays, err := manager.GetHolidays(countryCode, year)
	if err != nil {
		log.Printf("Error getting holidays: %v", err)
		return
	}

	// Sort and display holidays
	var dates []time.Time
	for date := range holidays {
		dates = append(dates, date)
	}
	
	// Simple sort by date
	for i := 0; i < len(dates); i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i].After(dates[j]) {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	for _, date := range dates {
		holiday := holidays[date]
		fmt.Printf("   %s: %s (%s)\n", 
			date.Format("2006-01-02"), 
			holiday.Name, 
			holiday.Category)
	}

	// Demonstrate country configuration
	fmt.Printf("\n3. Configuration for %s:\n", countryCode)
	info, err := manager.GetCountryInfo(countryCode)
	if err != nil {
		log.Printf("Error getting country info: %v", err)
		return
	}

	fmt.Printf("   Enabled: %v\n", info["enabled"])
	if subdivisions, ok := info["supported_subdivisions"].([]string); ok && len(subdivisions) > 0 {
		fmt.Printf("   Supported Subdivisions: %v\n", subdivisions[:min(3, len(subdivisions))])
		if len(subdivisions) > 3 {
			fmt.Printf("   ... and %d more\n", len(subdivisions)-3)
		}
	}

	// Demonstrate regional holidays (if supported)
	fmt.Printf("\n4. Regional Holidays for %s (California) in %d:\n", countryCode, year)
	regionalHolidays, err := manager.GetHolidaysWithSubdivisions(countryCode, year, []string{"CA"})
	if err != nil {
		log.Printf("Error getting regional holidays: %v", err)
	} else {
		regionalCount := len(regionalHolidays) - len(holidays)
		if regionalCount > 0 {
			fmt.Printf("   Found %d additional regional holidays\n", regionalCount)
		} else {
			fmt.Printf("   No additional regional holidays found\n")
		}
	}

	// Try another country
	countryCode = "GB"
	fmt.Printf("\n5. Sample holidays for %s in %d:\n", countryCode, year)
	gbHolidays, err := manager.GetHolidays(countryCode, year)
	if err != nil {
		log.Printf("Error getting holidays: %v", err)
	} else {
		count := 0
		for date, holiday := range gbHolidays {
			if count < 5 { // Show first 5
				fmt.Printf("   %s: %s\n", 
					date.Format("2006-01-02"), 
					holiday.Name)
				count++
			}
		}
		if len(gbHolidays) > 5 {
			fmt.Printf("   ... and %d more holidays\n", len(gbHolidays)-5)
		}
	}

	fmt.Println("\n6. Configuration System Features:")
	fmt.Println("   ✓ Environment-based configuration (dev/staging/prod)")
	fmt.Println("   ✓ Country-specific overrides and exclusions")
	fmt.Println("   ✓ Custom holiday definitions")
	fmt.Println("   ✓ Subdivision filtering")
	fmt.Println("   ✓ Category-based filtering")
	fmt.Println("   ✓ Holiday name localization")
	fmt.Println("   ✓ Output formatting and timezone support")
	fmt.Println("   ✓ Performance optimization settings")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
