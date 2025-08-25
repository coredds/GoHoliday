package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday/chronogo"
)

// This demo shows how ChronoGo would use GoHoliday's FastCountryChecker
// for efficient business day calculations

func main() {
	fmt.Println("ChronoGo + GoHoliday Integration Demo")
	fmt.Println("====================================")

	// Create fast holiday checkers for different countries
	usChecker := chronogo.Checker("US")
	caChecker := chronogo.Checker("CA")
	gbChecker := chronogo.Checker("GB")

	// Example 1: Single holiday checks
	fmt.Println("\n1. Single Holiday Checks:")
	testDate := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)

	fmt.Printf("July 4, 2024:\n")
	fmt.Printf("  US: %v (%s)\n", usChecker.IsHoliday(testDate), usChecker.GetHolidayName(testDate))
	fmt.Printf("  CA: %v (%s)\n", caChecker.IsHoliday(testDate), caChecker.GetHolidayName(testDate))
	fmt.Printf("  GB: %v (%s)\n", gbChecker.IsHoliday(testDate), gbChecker.GetHolidayName(testDate))

	// Example 2: Business day calculation helper
	fmt.Println("\n2. Business Day Range Analysis:")
	start := time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	holidayCount := usChecker.CountHolidaysInRange(start, end)
	totalDays := int(end.Sub(start).Hours()/24) + 1
	weekendDays := countWeekendDays(start, end)
	businessDays := totalDays - weekendDays - holidayCount

	fmt.Printf("Date Range: %s to %s\n", start.Format("Jan 2, 2006"), end.Format("Jan 2, 2006"))
	fmt.Printf("  Total Days: %d\n", totalDays)
	fmt.Printf("  Weekend Days: %d\n", weekendDays)
	fmt.Printf("  US Holidays: %d\n", holidayCount)
	fmt.Printf("  Business Days: %d\n", businessDays)

	// Example 3: Batch holiday checking (efficient for ChronoGo)
	fmt.Println("\n3. Batch Holiday Checking:")
	dates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year's
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),   // Day after
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // Independence Day
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
		time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),  // Random day
	}

	results := usChecker.AreHolidays(dates)
	for i, date := range dates {
		status := "business day"
		if results[i] {
			status = "holiday"
		}
		fmt.Printf("  %s: %s\n", date.Format("Jan 2, 2006"), status)
	}

	// Example 4: Holiday listings for calendar integration
	fmt.Println("\n4. Holiday Listings (January 2024):")
	janStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	janEnd := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	holidays := usChecker.GetHolidaysInRange(janStart, janEnd)
	for date, name := range holidays {
		fmt.Printf("  %s: %s\n", date.Format("Jan 2, 2006"), name)
	}

	// Example 5: Performance comparison simulation
	fmt.Println("\n5. Performance Demo (1000 holiday checks):")
	start = time.Now()
	for i := 0; i < 1000; i++ {
		usChecker.IsHoliday(testDate)
	}
	duration := time.Since(start)
	fmt.Printf("  1000 checks completed in: %v\n", duration)
	fmt.Printf("  Average per check: %v\n", duration/1000)

	// Example 6: Multi-country business operations
	fmt.Println("\n6. Multi-Country Business Day Check:")
	testDates := []time.Time{
		time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC),  // Canada Day
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),  // US Independence Day
		time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC), // UK Spring Bank Holiday (last Monday in May)
	}

	for _, date := range testDates {
		fmt.Printf("  %s:\n", date.Format("Jan 2, 2006"))
		fmt.Printf("    US: %v, CA: %v, GB: %v\n",
			usChecker.IsHoliday(date),
			caChecker.IsHoliday(date),
			gbChecker.IsHoliday(date))
	}
}

// Helper function to count weekend days in a range
func countWeekendDays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			count++
		}
	}
	return count
}
