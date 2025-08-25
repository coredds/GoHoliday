package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
	"github.com/coredds/GoHoliday/updater"
)

func main() {
	fmt.Println("GoHolidays Library Demo")
	fmt.Println("======================")

	// Basic usage example
	fmt.Println("\n1. Basic Holiday Checking:")

	// Create a US holiday provider
	us := goholidays.NewCountry("US")

	// Check some specific dates
	testDates := []string{
		"2024-01-01", // New Year's Day
		"2024-07-04", // Independence Day
		"2024-12-25", // Christmas
		"2024-03-15", // Random day
	}

	for _, dateStr := range testDates {
		date, _ := parseDate(dateStr)
		if holiday, isHoliday := us.IsHoliday(date); isHoliday {
			fmt.Printf("  %s: %s (%s)\n", dateStr, holiday.Name, holiday.Category)
		} else {
			fmt.Printf("  %s: Not a holiday\n", dateStr)
		}
	}

	// Get all holidays for a year
	fmt.Println("\n2. All US Holidays for 2024:")
	holidays := us.HolidaysForYear(2024)

	// Convert to slice for sorting
	type holidayEntry struct {
		date string
		name string
	}

	var sortedHolidays []holidayEntry
	for date, holiday := range holidays {
		sortedHolidays = append(sortedHolidays, holidayEntry{
			date: date.Format("2006-01-02"),
			name: holiday.Name,
		})
	}

	// Simple bubble sort by date
	for i := 0; i < len(sortedHolidays); i++ {
		for j := i + 1; j < len(sortedHolidays); j++ {
			if sortedHolidays[i].date > sortedHolidays[j].date {
				sortedHolidays[i], sortedHolidays[j] = sortedHolidays[j], sortedHolidays[i]
			}
		}
	}

	for _, holiday := range sortedHolidays {
		fmt.Printf("  %s: %s\n", holiday.date, holiday.name)
	}

	// Multi-language support example
	fmt.Println("\n3. Multi-language Support:")

	independenceDay, _ := parseDate("2024-07-04")
	if holiday, isHoliday := us.IsHoliday(independenceDay); isHoliday {
		fmt.Printf("  English: %s\n", holiday.Languages["en"])
		fmt.Printf("  Spanish: %s\n", holiday.Languages["es"])
	}

	// Country with subdivisions example
	fmt.Println("\n4. Country with Subdivisions:")

	usWithStates := goholidays.NewCountry("US", goholidays.CountryOptions{
		Subdivisions: []string{"CA", "TX"},
		Language:     "en",
	})

	fmt.Printf("  Country: %s\n", usWithStates.GetCountryCode())
	fmt.Printf("  Subdivisions: %v\n", usWithStates.GetSubdivisions())
	fmt.Printf("  Language: %s\n", usWithStates.GetLanguage())

	// Date range example
	fmt.Println("\n5. Holidays in Date Range:")

	start, _ := parseDate("2024-06-01")
	end, _ := parseDate("2024-08-31")
	summerHolidays := us.HolidaysForDateRange(start, end)

	fmt.Printf("  Holidays between %s and %s:\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
	for date, holiday := range summerHolidays {
		fmt.Printf("    %s: %s\n", date.Format("2006-01-02"), holiday.Name)
	}

	// Update system demo
	fmt.Println("\n6. Update System Demo:")

	// Create updater (this would normally sync from Python holidays repo)
	dataDir := "holiday_data"
	sync := updater.NewPythonHolidaysSync(dataDir)

	fmt.Println("  Checking for updates...")
	if hasUpdates, err := sync.CheckForUpdates(); err != nil {
		fmt.Printf("  Error checking for updates: %v\n", err)
	} else if hasUpdates {
		fmt.Println("  Updates available! Run sync to get latest data.")
	} else {
		fmt.Println("  Holiday data is up to date.")
	}

	// Performance demo
	fmt.Println("\n7. Performance Test:")

	performanceTest(us)

	fmt.Println("\nDemo completed! Check out the CLI tool with: go run cmd/goholidays/main.go -help")
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func performanceTest(country *goholidays.Country) {
	start := time.Now()

	// Test 1000 holiday checks
	testDate, _ := parseDate("2024-07-04")
	for i := 0; i < 1000; i++ {
		country.IsHoliday(testDate)
	}

	duration := time.Since(start)
	fmt.Printf("  1000 holiday checks took: %v\n", duration)
	fmt.Printf("  Average per check: %v\n", duration/1000)
}
