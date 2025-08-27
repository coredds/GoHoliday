package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("🇮🇳 India Holiday Demo - GoHoliday Library")
	fmt.Println("==========================================")

	// Create India holiday provider
	india := goholidays.NewCountry("IN")

	// Test current year
	currentYear := time.Now().Year()
	fmt.Printf("\n📅 India Holidays for %d:\n", currentYear)
	fmt.Println("----------------------------")

	holidays := india.HolidaysForYear(currentYear)
	if len(holidays) == 0 {
		log.Printf("No holidays found for India in %d", currentYear)
		return
	}

	// Sort and display holidays
	for date, holiday := range holidays {
		fmt.Printf("%-20s | %s", date.Format("January 2, 2006"), holiday.Name)
		
		// Show Hindi translation if available
		if holiday.Languages != nil && holiday.Languages["hi"] != "" {
			fmt.Printf(" (%s)", holiday.Languages["hi"])
		}
		
		// Show category
		fmt.Printf(" [%s]", holiday.Category)
		fmt.Println()
	}

	// Test specific dates
	fmt.Println("\n🔍 Testing Specific Dates:")
	fmt.Println("---------------------------")

	testDates := []struct {
		date        time.Time
		description string
	}{
		{time.Date(currentYear, 1, 26, 0, 0, 0, 0, time.UTC), "Republic Day"},
		{time.Date(currentYear, 8, 15, 0, 0, 0, 0, time.UTC), "Independence Day"},
		{time.Date(currentYear, 10, 2, 0, 0, 0, 0, time.UTC), "Gandhi Jayanti"},
		{time.Date(currentYear, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas"},
		{time.Date(currentYear, 7, 4, 0, 0, 0, 0, time.UTC), "Random Date (should not be holiday)"},
	}

	for _, test := range testDates {
		holiday, isHoliday := india.IsHoliday(test.date)
		if isHoliday {
			fmt.Printf("✅ %s is %s", test.date.Format("January 2"), holiday.Name)
			if holiday.Languages != nil && holiday.Languages["hi"] != "" {
				fmt.Printf(" (%s)", holiday.Languages["hi"])
			}
			fmt.Println()
		} else {
			fmt.Printf("❌ %s is not a holiday\n", test.date.Format("January 2"))
		}
	}

	// Test multiple years
	fmt.Println("\n📊 Holiday Count by Year:")
	fmt.Println("-------------------------")
	for year := currentYear; year <= currentYear+2; year++ {
		yearHolidays := india.HolidaysForYear(year)
		fmt.Printf("Year %d: %d holidays\n", year, len(yearHolidays))
	}

	// Show cultural diversity
	fmt.Println("\n🎭 Cultural Holiday Categories:")
	fmt.Println("-------------------------------")
	
	holidays2024 := india.HolidaysForYear(2024)
	categories := make(map[string][]string)
	
	for _, holiday := range holidays2024 {
		categoryName := string(holiday.Category)
		if categories[categoryName] == nil {
			categories[categoryName] = []string{}
		}
		categories[categoryName] = append(categories[categoryName], holiday.Name)
	}
	
	for category, holidayList := range categories {
		fmt.Printf("%s: %v\n", category, holidayList)
	}

	fmt.Println("\n✨ India holiday integration completed successfully!")
	fmt.Println("📝 Note: Religious festival dates are approximated. Production systems")
	fmt.Println("     should use proper lunar calendar calculations for accuracy.")
}
