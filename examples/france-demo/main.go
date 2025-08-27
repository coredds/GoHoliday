package main

import (
	"fmt"
	"log"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("🇫🇷 France Holiday Demo - GoHoliday Library")
	fmt.Println("============================================")

	// Create France holiday provider
	france := goholidays.NewCountry("FR")

	// Test current year
	currentYear := time.Now().Year()
	fmt.Printf("\n📅 France Holidays for %d:\n", currentYear)
	fmt.Println("-----------------------------")

	holidays := france.HolidaysForYear(currentYear)
	if len(holidays) == 0 {
		log.Printf("No holidays found for France in %d", currentYear)
		return
	}

	// Sort and display holidays
	for date, holiday := range holidays {
		fmt.Printf("%-20s | %s", date.Format("January 2, 2006"), holiday.Name)
		
		// Show French translation
		if holiday.Languages != nil && holiday.Languages["fr"] != "" {
			fmt.Printf(" (%s)", holiday.Languages["fr"])
		}
		
		// Show category
		fmt.Printf(" [%s]", holiday.Category)
		fmt.Println()
	}

	// Test specific dates
	fmt.Println("\n🔍 Testing Specific French Holidays:")
	fmt.Println("-------------------------------------")

	testDates := []struct {
		date        time.Time
		description string
	}{
		{time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day"},
		{time.Date(currentYear, 5, 1, 0, 0, 0, 0, time.UTC), "Labour Day"},
		{time.Date(currentYear, 5, 8, 0, 0, 0, 0, time.UTC), "Victory in Europe Day"},
		{time.Date(currentYear, 7, 14, 0, 0, 0, 0, time.UTC), "Bastille Day"},
		{time.Date(currentYear, 8, 15, 0, 0, 0, 0, time.UTC), "Assumption of Mary"},
		{time.Date(currentYear, 11, 1, 0, 0, 0, 0, time.UTC), "All Saints' Day"},
		{time.Date(currentYear, 11, 11, 0, 0, 0, 0, time.UTC), "Armistice Day"},
		{time.Date(currentYear, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day"},
		{time.Date(currentYear, 6, 21, 0, 0, 0, 0, time.UTC), "Random Date (should not be holiday)"},
	}

	for _, test := range testDates {
		holiday, isHoliday := france.IsHoliday(test.date)
		if isHoliday {
			fmt.Printf("✅ %s is %s", test.date.Format("January 2"), holiday.Name)
			if holiday.Languages != nil && holiday.Languages["fr"] != "" {
				fmt.Printf(" (%s)", holiday.Languages["fr"])
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
		yearHolidays := france.HolidaysForYear(year)
		fmt.Printf("Year %d: %d holidays\n", year, len(yearHolidays))
	}

	// Show cultural context
	fmt.Println("\n🎭 French Holiday Categories:")
	fmt.Println("-----------------------------")
	
	holidays2024 := france.HolidaysForYear(2024)
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

	// Show bilingual features
	fmt.Println("\n🌐 Bilingual Holiday Names (French/English):")
	fmt.Println("---------------------------------------------")
	
	for _, holiday := range holidays2024 {
		if len(holiday.Languages) > 1 {
			fmt.Printf("🇫🇷 %-20s | 🇬🇧 %s\n", 
				holiday.Languages["fr"], 
				holiday.Languages["en"])
		}
	}

	fmt.Println("\n✨ France holiday integration completed successfully!")
	fmt.Printf("📝 France now has %d holidays with full bilingual support!\n", len(holidays2024))
	fmt.Println("🇫🇷 Liberté, Égalité, Fraternité - French holidays at your service!")
}
