package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHolidays Multi-Country Demo with Australia")
	fmt.Println("============================================")
	
	year := 2024
	countries := map[string]string{
		"US": "United States",
		"CA": "Canada", 
		"AU": "Australia",
	}
	
	fmt.Printf("Comparing holidays for %d across countries:\n\n", year)
	
	// Create country instances
	countryInstances := make(map[string]*goholidays.Country)
	for code := range countries {
		countryInstances[code] = goholidays.NewCountry(code)
	}
	
	// Display holidays by month to show seasonal differences
	for month := 1; month <= 12; month++ {
		fmt.Printf("=== %s ===\n", time.Month(month).String())
		hasHolidays := false
		
		for code, name := range countries {
			country := countryInstances[code]
			holidays := country.HolidaysForYear(year)
			
			// Filter holidays for this month
			monthHolidays := make([]time.Time, 0)
			for date := range holidays {
				if date.Month() == time.Month(month) {
					monthHolidays = append(monthHolidays, date)
				}
			}
			
			if len(monthHolidays) > 0 {
				hasHolidays = true
				fmt.Printf("%s (%s):\n", name, code)
				for _, date := range monthHolidays {
					holiday, _ := country.IsHoliday(date)
					fmt.Printf("  %-12s %s\n", date.Format("Jan 02"), holiday.Name)
				}
			}
		}
		
		if !hasHolidays {
			fmt.Println("No major holidays this month")
		}
		fmt.Println()
	}
	
	// Show seasonal comparison
	fmt.Println("=== SEASONAL COMPARISON ===")
	fmt.Println("Notice how Australia's seasons are opposite to North America:")
	fmt.Println()
	
	// Summer holidays comparison
	fmt.Println("SUMMER SEASON:")
	fmt.Println("Northern Hemisphere (US/CA): June-August")
	fmt.Println("Southern Hemisphere (AU): December-February")
	fmt.Println()
	
	// Christmas season
	fmt.Println("Christmas & New Year:")
	december := []time.Time{
		time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC),
		time.Date(year, 12, 26, 0, 0, 0, 0, time.UTC),
	}
	january := []time.Time{
		time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	
	for code, name := range countries {
		country := countryInstances[code]
		fmt.Printf("%s (%s):\n", name, code)
		
		for _, date := range december {
			if holiday, isHoliday := country.IsHoliday(date); isHoliday {
				season := "Winter"
				if code == "AU" {
					season = "Summer"
				}
				fmt.Printf("  Dec 25: %s (%s season)\n", holiday.Name, season)
			}
		}
		
		for _, date := range january {
			if holiday, isHoliday := country.IsHoliday(date); isHoliday {
				season := "Winter"
				if code == "AU" {
					season = "Summer"
				}
				fmt.Printf("  Jan 01: %s (%s season)\n", holiday.Name, season)
			}
		}
	}
	
	fmt.Println()
	
	// Unique holidays comparison
	fmt.Println("=== UNIQUE NATIONAL HOLIDAYS ===")
	uniqueHolidays := map[string][]string{
		"US": {"Independence Day (July 4)"},
		"CA": {"Canada Day (July 1)", "Thanksgiving (October)"},
		"AU": {"Australia Day (January 26)", "ANZAC Day (April 25)", "Queen's Birthday (June)", "Boxing Day (December 26)"},
	}
	
	for code, name := range countries {
		fmt.Printf("%s (%s):\n", name, code)
		for _, holiday := range uniqueHolidays[code] {
			fmt.Printf("  - %s\n", holiday)
		}
		fmt.Println()
	}
	
	// Performance demonstration
	fmt.Println("=== PERFORMANCE BENCHMARK ===")
	fmt.Println("Testing holiday lookup performance...")
	
	start := time.Now()
	totalHolidays := 0
	iterations := 1000
	
	for i := 0; i < iterations; i++ {
		for code := range countries {
			country := countryInstances[code]
			holidays := country.HolidaysForYear(year)
			totalHolidays += len(holidays)
		}
	}
	
	duration := time.Since(start)
	fmt.Printf("Processed %d holiday calculations in %v\n", totalHolidays, duration)
	fmt.Printf("Average: %.2f ns per holiday lookup\n", float64(duration.Nanoseconds())/float64(totalHolidays))
	
	fmt.Println("\n=== COUNTRY SUMMARY ===")
	for code, name := range countries {
		country := countryInstances[code]
		holidays := country.HolidaysForYear(year)
		fmt.Printf("%s (%s): %d holidays\n", name, code, len(holidays))
	}
}
