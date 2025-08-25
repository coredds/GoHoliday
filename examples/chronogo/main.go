package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday/chronogo"
)

func main() {
	fmt.Println("GoHolidays + ChronoGo Integration Demo")
	fmt.Println("====================================")

	// Example 1: Basic Holiday Checking
	fmt.Println("\n1. Basic Holiday Checking:")
	checker := chronogo.Checker("US")

	testDates := []struct {
		name string
		date time.Time
	}{
		{"New Year's Day 2024", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
		{"Regular Day", time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)},
		{"Independence Day 2024", time.Date(2024, time.July, 4, 0, 0, 0, 0, time.UTC)},
		{"Christmas Day 2024", time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range testDates {
		isHoliday := checker.IsHoliday(test.date)
		holidayName := checker.GetHolidayName(test.date)
		fmt.Printf("  %-25s: %v", test.name, isHoliday)
		if holidayName != "" {
			fmt.Printf(" (%s)", holidayName)
		}
		fmt.Println()
	}

	// Example 2: Multi-Country Comparison
	fmt.Println("\n2. Multi-Country Holiday Comparison:")
	countries := []string{"US", "CA", "GB"}
	checkers := make(map[string]*chronogo.FastCountryChecker)

	for _, country := range countries {
		checkers[country] = chronogo.Checker(country)
	}

	comparisonDates := []struct {
		name string
		date time.Time
	}{
		{"Canada Day", time.Date(2024, time.July, 1, 0, 0, 0, 0, time.UTC)},
		{"Independence Day", time.Date(2024, time.July, 4, 0, 0, 0, 0, time.UTC)},
		{"Christmas Day", time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)},
	}

	for _, test := range comparisonDates {
		fmt.Printf("  %-20s:", test.name)
		for _, country := range countries {
			isHoliday := checkers[country].IsHoliday(test.date)
			fmt.Printf(" %s=%v", country, isHoliday)
		}
		fmt.Println()
	}

	// Example 3: Batch Holiday Processing
	fmt.Println("\n3. Batch Holiday Processing:")
	batchDates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year's
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),  // MLK Day
		time.Date(2024, 2, 19, 0, 0, 0, 0, time.UTC),  // Presidents Day
		time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC),  // Memorial Day
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // Independence Day
		time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC),   // Labor Day
		time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC), // Thanksgiving
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
	}

	results := checker.AreHolidays(batchDates)
	for i, date := range batchDates {
		status := "Regular Day"
		if results[i] {
			status = "Holiday"
		}
		fmt.Printf("  %s: %s\n", date.Format("Jan 2, 2006"), status)
	}

	// Example 4: Holiday Range Analysis
	fmt.Println("\n4. Holiday Range Analysis:")
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	holidayCount := checker.CountHolidaysInRange(start, end)
	holidays := checker.GetHolidaysInRange(start, end)

	fmt.Printf("  US Federal Holidays in 2024: %d\n", holidayCount)
	fmt.Printf("  Holiday List:\n")

	// Sort holidays by date for display
	sortedDates := make([]time.Time, 0, len(holidays))
	for date := range holidays {
		sortedDates = append(sortedDates, date)
	}

	// Simple sort by comparing Unix timestamps
	for i := 0; i < len(sortedDates); i++ {
		for j := i + 1; j < len(sortedDates); j++ {
			if sortedDates[i].After(sortedDates[j]) {
				sortedDates[i], sortedDates[j] = sortedDates[j], sortedDates[i]
			}
		}
	}

	for _, date := range sortedDates {
		fmt.Printf("    %s: %s\n", date.Format("Jan 2, 2006"), holidays[date])
	}

	// Example 5: Performance Demonstration
	fmt.Println("\n5. Performance Demonstration:")
	iterations := 10000
	start = time.Now()

	for i := 0; i < iterations; i++ {
		checker.IsHoliday(time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC))
	}

	duration := time.Since(start)
	fmt.Printf("  %d holiday checks completed in: %v\n", iterations, duration)
	fmt.Printf("  Average per check: %v\n", duration/time.Duration(iterations))
	fmt.Printf("  Checks per second: %.0f\n", float64(iterations)/duration.Seconds())

	// Example 6: Business Day Calculation Helper
	fmt.Println("\n6. Business Day Calculation Helper:")
	businessStart := time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC)
	businessEnd := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	businessHolidays := checker.CountHolidaysInRange(businessStart, businessEnd)
	totalDays := int(businessEnd.Sub(businessStart).Hours()/24) + 1

	// Count weekends
	weekendDays := 0
	for d := businessStart; !d.After(businessEnd); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			weekendDays++
		}
	}

	businessDays := totalDays - weekendDays - businessHolidays

	fmt.Printf("  Date Range: %s to %s\n",
		businessStart.Format("Jan 2, 2006"),
		businessEnd.Format("Jan 2, 2006"))
	fmt.Printf("  Total Days: %d\n", totalDays)
	fmt.Printf("  Weekend Days: %d\n", weekendDays)
	fmt.Printf("  Holiday Days: %d\n", businessHolidays)
	fmt.Printf("  Business Days: %d\n", businessDays)

	fmt.Println("\n✅ ChronoGo Integration Demo Complete!")
	fmt.Println("\n📚 Usage in ChronoGo:")
	fmt.Println("   // Create a fast holiday checker")
	fmt.Println("   holidayChecker := chronogo.Checker(\"US\")")
	fmt.Println("   ")
	fmt.Println("   // Business day calculations with GoHolidays")
	fmt.Println("   isHoliday := holidayChecker.IsHoliday(someDate)")
	fmt.Println("   holidayName := holidayChecker.GetHolidayName(someDate)")
	fmt.Println("   holidayCount := holidayChecker.CountHolidaysInRange(start, end)")
}
