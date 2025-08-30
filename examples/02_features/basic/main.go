package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHoliday Basic Features")
	fmt.Println("=======================")

	// Create a holiday provider
	us := goholidays.NewCountry("US", goholidays.CountryOptions{
		Categories: []goholidays.HolidayCategory{
			goholidays.CategoryPublic,
			goholidays.CategoryBank,
		},
		Language: "en",
	})

	// 1. Holiday Categories
	fmt.Println("\n1. Holiday Categories")
	fmt.Printf("Supported categories: %v\n", us.GetCategories())

	// 2. Check specific dates
	fmt.Println("\n2. Holiday Checks")
	checkDates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year's Day
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // Independence Day
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
		time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),  // Regular day
	}

	for _, date := range checkDates {
		if holiday, isHoliday := us.IsHoliday(date); isHoliday {
			fmt.Printf("%s: %s (%s)\n", date.Format("Jan 2"), holiday.Name, holiday.Category)
		} else {
			fmt.Printf("%s: Not a holiday\n", date.Format("Jan 2"))
		}
	}

	// 3. Date Range Features
	fmt.Println("\n3. Holiday Ranges")
	start := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Holidays between %s and %s:\n", start.Format("Jan 2"), end.Format("Jan 2"))

	holidays := us.HolidaysForDateRange(start, end)
	printSortedHolidays(holidays)

	// 4. Year Overview
	fmt.Println("\n4. Full Year Overview")
	yearHolidays := us.HolidaysForYear(2024)
	fmt.Printf("Total holidays in 2024: %d\n", len(yearHolidays))

	// Group by month
	monthlyCount := make(map[time.Month]int)
	for date := range yearHolidays {
		monthlyCount[date.Month()]++
	}

	// Print monthly distribution
	for month := time.January; month <= time.December; month++ {
		if count := monthlyCount[month]; count > 0 {
			fmt.Printf("%s: %d holiday(s)\n", month, count)
		}
	}

	fmt.Println("\nCheck out other examples for more advanced features!")
}

func printSortedHolidays(holidays map[time.Time]*goholidays.Holiday) {
	// Convert to slice for sorting
	type holidayEntry struct {
		date    time.Time
		holiday *goholidays.Holiday
	}

	var entries []holidayEntry
	for date, holiday := range holidays {
		entries = append(entries, holidayEntry{date, holiday})
	}

	// Sort by date
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].date.Before(entries[j].date)
	})

	// Print sorted holidays
	for _, entry := range entries {
		fmt.Printf("- %s: %s\n", entry.date.Format("Jan 2"), entry.holiday.Name)
	}
}
