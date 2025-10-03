package main

import (
	"fmt"
	"sort"
	"time"

	goholidays "github.com/coredds/goholiday"
)

func main() {
	fmt.Println("goholiday Multi-Country Features")
	fmt.Println("================================")

	// Create providers for different countries
	countries := map[string]*goholidays.Country{
		"US": goholidays.NewCountry("US"),
		"GB": goholidays.NewCountry("GB"),
		"JP": goholidays.NewCountry("JP"),
		"AU": goholidays.NewCountry("AU"),
	}

	// 1. Compare specific dates across countries
	fmt.Println("\n1. International Holiday Comparison")
	checkDates := []struct {
		date time.Time
		desc string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "May Day/Labor Day"},
	}

	for _, check := range checkDates {
		fmt.Printf("\n%s (%s):\n", check.date.Format("Jan 2"), check.desc)
		for code, country := range countries {
			if holiday, isHoliday := country.IsHoliday(check.date); isHoliday {
				fmt.Printf("- %s: %s\n", code, holiday.Name)
			} else {
				fmt.Printf("- %s: Not a holiday\n", code)
			}
		}
	}

	// 2. Country-specific holidays
	fmt.Println("\n2. Unique National Holidays (2024)")
	uniqueHolidays := map[string]time.Time{
		"US": time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),  // Independence Day
		"GB": time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),  // Queen's Birthday
		"JP": time.Date(2024, 2, 11, 0, 0, 0, 0, time.UTC), // National Foundation Day
		"AU": time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), // Australia Day
	}

	for code, date := range uniqueHolidays {
		if holiday, isHoliday := countries[code].IsHoliday(date); isHoliday {
			fmt.Printf("%s: %s (%s)\n", code, holiday.Name, date.Format("Jan 2"))
		}
	}

	// 3. Holiday Count Comparison
	fmt.Println("\n3. Total Holidays per Country (2024)")
	for code, country := range countries {
		holidays := country.HolidaysForYear(2024)
		fmt.Printf("%s: %d holidays\n", code, len(holidays))
	}

	// 4. Multi-language Support
	fmt.Println("\n4. Multi-language Holiday Names")
	christmas := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	for code, country := range countries {
		if holiday, isHoliday := country.IsHoliday(christmas); isHoliday {
			fmt.Printf("\n%s Christmas names:\n", code)
			for lang, name := range holiday.Languages {
				fmt.Printf("- %s: %s\n", lang, name)
			}
		}
	}

	// 5. Business Day Analysis
	fmt.Println("\n5. Business Days Comparison")
	// Compare business days in December 2024
	start := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	for code, country := range countries {
		calc := goholidays.NewBusinessDayCalculator(country)
		days := calc.BusinessDaysBetween(start, end)
		fmt.Printf("%s: %d business days in December 2024\n", code, days)
	}

	// 6. Holiday Categories
	fmt.Println("\n6. Holiday Categories by Country")
	for code, country := range countries {
		fmt.Printf("\n%s Categories:\n", code)
		holidays := country.HolidaysForYear(2024)
		categories := make(map[goholidays.HolidayCategory]int)
		for _, holiday := range holidays {
			categories[holiday.Category]++
		}
		for category, count := range categories {
			fmt.Printf("- %s: %d holiday(s)\n", category, count)
		}
	}

	fmt.Println("\nThis demonstrates the power of goholiday's multi-country support!")
}

func printSortedHolidays(holidays map[time.Time]*goholidays.Holiday) {
	var dates []time.Time
	for date := range holidays {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	for _, date := range dates {
		fmt.Printf("- %s: %s\n", date.Format("Jan 2"), holidays[date].Name)
	}
}
