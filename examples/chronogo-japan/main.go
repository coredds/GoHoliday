package main

import (
	"fmt"
	"time"
	
	"github.com/coredds/GoHoliday/chronogo"
)

func main() {
	// Test ChronoGo with Japan
	jpChecker := chronogo.Checker("JP")
	if jpChecker == nil {
		fmt.Println("Error: Could not create Japan checker")
		return
	}
	
	fmt.Println("Testing ChronoGo Japan Integration:")
	fmt.Printf("Country: %s\n\n", jpChecker.GetCountryCode())
	
	// Test some key dates
	testDates := []struct {
		date time.Time
		desc string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day 2024"},
		{time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), "Coming of Age Day 2024"},
		{time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC), "Emperor's Birthday 2024"},
		{time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC), "Constitution Memorial Day 2024"},
		{time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC), "Children's Day 2024"},
		{time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), "Regular day (should not be holiday)"},
		{time.Date(2019, 12, 23, 0, 0, 0, 0, time.UTC), "Emperor Akihito's Birthday 2019"},
		{time.Date(2020, 2, 23, 0, 0, 0, 0, time.UTC), "Emperor Naruhito's Birthday 2020"},
	}
	
	for _, test := range testDates {
		isHoliday := jpChecker.IsHoliday(test.date)
		if isHoliday {
			name := jpChecker.GetHolidayName(test.date)
			fmt.Printf("✓ %s: %s (HOLIDAY)\n", test.desc, name)
		} else {
			fmt.Printf("  %s: Not a holiday\n", test.desc)
		}
	}
	
	// Test batch operation
	fmt.Println("\nBatch testing multiple dates:")
	dates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC),
	}
	
	results := jpChecker.AreHolidays(dates)
	for i, date := range dates {
		if results[i] {
			fmt.Printf("✓ %s: Holiday\n", date.Format("2006-01-02"))
		} else {
			fmt.Printf("  %s: Not a holiday\n", date.Format("2006-01-02"))
		}
	}
	
	// Test range counting
	fmt.Println("\nHoliday count for Q1 2024:")
	q1Start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	q1End := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	count := jpChecker.CountHolidaysInRange(q1Start, q1End)
	fmt.Printf("Q1 2024 has %d holidays\n", count)
}
