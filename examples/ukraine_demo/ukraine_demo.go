package main

import (
	"fmt"
	"sort"
	"time"

	goholidays "github.com/coredds/GoHoliday"
	"github.com/coredds/GoHoliday/countries"
)

func main() {
	fmt.Println("🇺🇦 GoHoliday - Ukraine Demo")
	fmt.Println("=============================")
	fmt.Println()

	// Create Ukraine provider directly and through main library
	ua := goholidays.NewCountry("UA")
	provider := countries.NewUAProvider()

	fmt.Printf("Country: %s (%s)\n", provider.GetName(), provider.GetCountryCode())
	fmt.Printf("Subdivisions: %d oblasts and cities\n", len(provider.GetSupportedSubdivisions()))
	fmt.Printf("Categories: %v\n", provider.GetSupportedCategories())
	fmt.Printf("Languages: %v\n", provider.GetLanguages())
	fmt.Println()

	// Demonstrate key Ukrainian holidays for 2024
	year := 2024
	holidays := ua.HolidaysForYear(year)
	fmt.Printf("📅 Ukrainian holidays in %d: %d total\n\n", year, len(holidays))

	// Sort holidays by date for better display
	var sortedDates []time.Time
	for date := range holidays {
		sortedDates = append(sortedDates, date)
	}
	sort.Slice(sortedDates, func(i, j int) bool {
		return sortedDates[i].Before(sortedDates[j])
	})

	// Group holidays by category
	categoryHolidays := make(map[string][]time.Time)
	for _, date := range sortedDates {
		holiday := holidays[date]
		categoryHolidays[string(holiday.Category)] = append(categoryHolidays[string(holiday.Category)], date)
	}

	// Display holidays by category
	categories := []string{"national", "orthodox", "memorial", "cultural"}

	for _, category := range categories {
		if dates, exists := categoryHolidays[category]; exists {
			fmt.Printf("🏛️  %s holidays (%d):\n", category, len(dates))
			for _, date := range dates {
				holiday := holidays[date]
				fmt.Printf("   %s - %s\n", date.Format("Jan 02"), holiday.Name)
				if holiday.Languages["uk"] != "" {
					fmt.Printf("      🇺🇦 %s\n", holiday.Languages["uk"])
				}
				if holiday.Languages["ru"] != "" && holiday.Languages["ru"] != holiday.Languages["uk"] {
					fmt.Printf("      🇷🇺 %s\n", holiday.Languages["ru"])
				}
			}
			fmt.Println()
		}
	}

	// Demonstrate Orthodox Easter calculations
	fmt.Println("🔬 Orthodox Easter Calculations:")
	fmt.Println("================================")
	for year := 2024; year <= 2027; year++ {
		yearHolidays := ua.HolidaysForYear(year)
		for date, holiday := range yearHolidays {
			if holiday.Languages["en"] == "Orthodox Easter" {
				fmt.Printf("%d: %s (Orthodox Easter)\n", year, date.Format("January 2"))

				// Show related Orthodox holidays
				palmSunday := date.AddDate(0, 0, -7)
				trinity := date.AddDate(0, 0, 49)

				if palmHoliday, exists := yearHolidays[palmSunday]; exists {
					fmt.Printf("      Palm Sunday: %s - %s\n", palmSunday.Format("Jan 2"), palmHoliday.Languages["uk"])
				}
				if trinityHoliday, exists := yearHolidays[trinity]; exists {
					fmt.Printf("      Trinity: %s - %s\n", trinity.Format("Jan 2"), trinityHoliday.Languages["uk"])
				}
				break
			}
		}
	}
	fmt.Println()

	// Demonstrate historical context - holidays that change over time
	fmt.Println("📊 Historical Context:")
	fmt.Println("======================")

	testYears := []int{2013, 2014, 2015, 2019, 2021, 2024}
	holidayIntroductions := map[string]int{
		"Day of Dignity and Freedom": 2014,
		"Defenders Day":              2015,
		"Day of Ukrainian Language":  2019,
		"Day of Ukrainian Statehood": 2021,
	}

	for holidayName, introYear := range holidayIntroductions {
		fmt.Printf("%s (introduced %d):\n", holidayName, introYear)
		for _, year := range testYears {
			yearHolidays := ua.HolidaysForYear(year)
			found := false
			for _, holiday := range yearHolidays {
				if holiday.Languages["en"] == holidayName {
					fmt.Printf("   %d: ✅ Present\n", year)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("   %d: ❌ Not yet introduced\n", year)
			}
		}
		fmt.Println()
	}

	// Demonstrate multi-language support
	fmt.Println("🌍 Multi-Language Support:")
	fmt.Println("===========================")

	keyHolidays := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year
		time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),   // Orthodox Christmas
		time.Date(2024, 8, 24, 0, 0, 0, 0, time.UTC),  // Independence Day
		time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC), // Holodomor Remembrance
	}

	for _, date := range keyHolidays {
		if holiday, exists := holidays[date]; exists {
			fmt.Printf("%s:\n", date.Format("January 2"))
			fmt.Printf("   🇬🇧 English: %s\n", holiday.Languages["en"])
			fmt.Printf("   🇺🇦 Ukrainian: %s\n", holiday.Languages["uk"])
			if holiday.Languages["ru"] != "" {
				fmt.Printf("   🇷🇺 Russian: %s\n", holiday.Languages["ru"])
			}
			fmt.Printf("   Category: %s\n", holiday.Category)
			fmt.Println()
		}
	}

	fmt.Println("🎯 Integration Examples:")
	fmt.Println("=========================")

	// Business day calculation example
	businessDays := 0
	startDate := time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC) // Week before Independence Day
	endDate := time.Date(2024, 8, 26, 0, 0, 0, 0, time.UTC)   // Week after Independence Day

	fmt.Printf("Business days from %s to %s:\n", startDate.Format("Jan 2"), endDate.Format("Jan 2"))
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		isWeekend := d.Weekday() == time.Saturday || d.Weekday() == time.Sunday
		_, isHoliday := ua.IsHoliday(d)

		status := "✅ Business day"
		if isWeekend {
			status = "🚫 Weekend"
		} else if isHoliday {
			status = "🎉 Holiday"
		} else {
			businessDays++
		}

		fmt.Printf("   %s (%s): %s\n", d.Format("Jan 2"), d.Weekday().String()[:3], status)
	}
	fmt.Printf("Total business days: %d\n", businessDays)
}
