package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coredds/GoHoliday/chronogo"
)

func main() {
	fmt.Println("GoHolidays + ChronoGo Integration Demo")
	fmt.Println("====================================")

	// Example 1: Basic Holiday Checking
	fmt.Println("\n1. Basic Holiday Checking:")
	checker := chronogo.NewGoHolidaysChecker().WithCountries("US")
	
	testDates := []struct {
		name string
		date *mockDateTime
	}{
		{"New Year's Day 2024", &mockDateTime{2024, time.January, 1}},
		{"Regular Day", &mockDateTime{2024, time.March, 15}},
		{"Independence Day 2024", &mockDateTime{2024, time.July, 4}},
		{"Christmas Day 2024", &mockDateTime{2024, time.December, 25}},
	}

	for _, test := range testDates {
		isHoliday := checker.IsHoliday(test.date)
		status := "❌ Regular Day"
		if isHoliday {
			status = "🎉 Holiday"
		}
		fmt.Printf("   %s: %s\n", test.name, status)
	}

	// Example 2: Multi-Country Support
	fmt.Println("\n2. Multi-Country Holiday Support:")
	multiChecker := chronogo.CreateMultiCountryChecker("US", "CA", "GB")
	
	internationalDates := []struct {
		name string
		date *mockDateTime
	}{
		{"US Independence Day", &mockDateTime{2024, time.July, 4}},
		{"Canada Day", &mockDateTime{2024, time.July, 1}},
		{"UK Boxing Day", &mockDateTime{2024, time.December, 26}},
	}

	for _, test := range internationalDates {
		isHoliday := multiChecker.IsHoliday(test.date)
		status := "❌ Not a holiday"
		if isHoliday {
			status = "🎉 Holiday"
		}
		fmt.Printf("   %s: %s\n", test.name, status)
	}

	// Example 3: Regional Holiday Support
	fmt.Println("\n3. Regional Holiday Support:")
	regionalChecker := chronogo.CreateRegionalChecker("US", "CA") // California
	
	fmt.Printf("   Regional holidays for US (California) enabled\n")
	
	baseHolidays := checker.GetHolidays(2024)
	regionalHolidays := regionalChecker.GetHolidays(2024)
	
	fmt.Printf("   Base US holidays: %d\n", len(baseHolidays))
	fmt.Printf("   With CA regional: %d\n", len(regionalHolidays))
	fmt.Printf("   Additional regional holidays: %d\n", len(regionalHolidays)-len(baseHolidays))

	// Example 4: Holiday Category Filtering
	fmt.Println("\n4. Holiday Category Filtering:")
	
	federalOnly := chronogo.NewGoHolidaysChecker().
		WithCountries("US").
		WithCategories("federal")
	
	allCategories := chronogo.NewGoHolidaysChecker().
		WithCountries("US").
		WithCategories("federal", "observance", "state")

	federalHolidays := federalOnly.GetHolidays(2024)
	allHolidays := allCategories.GetHolidays(2024)
	
	fmt.Printf("   Federal holidays only: %d\n", len(federalHolidays))
	fmt.Printf("   All categories: %d\n", len(allHolidays))

	// Example 5: Configuration-Driven Setup
	fmt.Println("\n5. Configuration System Integration:")
	
	// This would load holidays from your configuration files
	fmt.Printf("   ✓ YAML-based configuration support\n")
	fmt.Printf("   ✓ Environment-specific settings (dev/prod)\n")
	fmt.Printf("   ✓ Custom holiday definitions\n")
	fmt.Printf("   ✓ Country-specific overrides\n")

	// Example 6: Performance with Caching
	fmt.Println("\n6. Performance Optimization:")
	
	cachedChecker := chronogo.NewGoHolidaysChecker().WithCountries("US")
	
	// Preload holidays for better performance
	err := cachedChecker.PreloadYear(2024)
	if err != nil {
		log.Printf("Warning: Could not preload holidays: %v", err)
	}
	
	start := time.Now()
	for i := 0; i < 1000; i++ {
		testDate := &mockDateTime{2024, time.July, 4}
		cachedChecker.IsHoliday(testDate)
	}
	duration := time.Since(start)
	
	fmt.Printf("   1000 holiday checks: %v (avg: %v per check)\n", 
		duration, duration/1000)

	// Example 7: Supported Countries
	fmt.Println("\n7. Supported Countries:")
	supportedCountries := checker.GetSupportedCountries()
	fmt.Printf("   Total countries supported: %d\n", len(supportedCountries))
	fmt.Printf("   Countries: %v\n", supportedCountries)

	// Example 8: Holiday Details
	fmt.Println("\n8. Holiday Information:")
	holidays := checker.GetHolidays(2024)
	
	fmt.Printf("   US Holidays in 2024:\n")
	for i, holiday := range holidays {
		if i >= 5 { // Show first 5
			fmt.Printf("   ... and %d more holidays\n", len(holidays)-5)
			break
		}
		fmt.Printf("   - %s (%s) - %s\n", 
			holiday.Name, 
			holiday.Date.Format("Jan 2"), 
			holiday.Category)
	}

	// Example 9: Integration Benefits
	fmt.Println("\n9. Integration Benefits:")
	fmt.Println("   ✅ Drop-in replacement for ChronoGo's DefaultHolidayChecker")
	fmt.Println("   ✅ 7 countries supported (AU, CA, NZ, GB, US, DE, FR)")
	fmt.Println("   ✅ Regional/subdivision holiday support")
	fmt.Println("   ✅ Configuration-driven customization")
	fmt.Println("   ✅ Performance optimized with caching")
	fmt.Println("   ✅ Category-based filtering")
	fmt.Println("   ✅ Multi-country business operations")

	fmt.Println("\n📚 Usage in ChronoGo:")
	fmt.Println(`
   // In your ChronoGo application:
   holidayChecker := chronogo.CreateDefaultUSChecker()
   
   // Business day calculations with GoHolidays
   dt := chronogo.Now()
   nextBusinessDay := dt.NextBusinessDay(holidayChecker)
   businessDays := dt.AddBusinessDays(5, holidayChecker)
   
   // Check if today is a holiday
   if dt.IsHoliday(holidayChecker) {
       fmt.Println("Today is a holiday!")
   }`)
}

// mockDateTime implements the ChronoGoDateTime interface for testing
type mockDateTime struct {
	year  int
	month time.Month
	day   int
}

func (m *mockDateTime) Year() int        { return m.year }
func (m *mockDateTime) Month() time.Month { return m.month }
func (m *mockDateTime) Day() int         { return m.day }
