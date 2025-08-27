package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("🇧🇷 BRAZIL HOLIDAYS SHOWCASE 🇧🇷")
	fmt.Println("================================")

	// Create Brazil country instance
	brazil := goholidays.NewCountry("BR")

	// Get holidays for 2024
	holidays := brazil.HolidaysForYear(2024)

	// Convert to slice and sort by date
	type holidayInfo struct {
		date    time.Time
		holiday *goholidays.Holiday
	}

	var sortedHolidays []holidayInfo
	for date, holiday := range holidays {
		sortedHolidays = append(sortedHolidays, holidayInfo{date, holiday})
	}

	sort.Slice(sortedHolidays, func(i, j int) bool {
		return sortedHolidays[i].date.Before(sortedHolidays[j].date)
	})

	fmt.Printf("\n📅 BRAZILIAN HOLIDAYS FOR 2024 (%d total)\n", len(holidays))
	fmt.Println("==========================================")

	for _, hi := range sortedHolidays {
		// Get both Portuguese and English names
		ptName := hi.holiday.Name
		enName := ""
		if hi.holiday.Languages != nil {
			if en, exists := hi.holiday.Languages["en"]; exists {
				enName = fmt.Sprintf(" (%s)", en)
			}
		}

		// Format date nicely
		dateStr := hi.date.Format("Monday, January 2")
		categoryIcon := getCategoryIcon(string(hi.holiday.Category))

		fmt.Printf("%-20s %s %s%s %s\n",
			dateStr,
			categoryIcon,
			ptName,
			enName,
			string(hi.holiday.Category))
	}

	fmt.Println("\n🎭 CULTURAL HIGHLIGHTS")
	fmt.Println("=====================")

	// Highlight Carnival
	carnivalMonday := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brazil.IsHoliday(carnivalMonday); exists {
		fmt.Printf("🎪 Carnival Monday: %s\n", holiday.Name)
		fmt.Println("   Brazil's most famous cultural celebration!")
	}

	// Highlight Independence Day
	independence := time.Date(2024, 9, 7, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brazil.IsHoliday(independence); exists {
		fmt.Printf("🇧🇷 Independence Day: %s\n", holiday.Name)
		fmt.Println("   Celebrating Brazil's independence from Portugal (1822)")
	}

	// Highlight Day of the Dead
	finados := time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brazil.IsHoliday(finados); exists {
		fmt.Printf("🕯️  Day of the Dead: %s\n", holiday.Name)
		fmt.Println("   A day to honor deceased loved ones")
	}

	fmt.Println("\n📊 HOLIDAY STATISTICS")
	fmt.Println("====================")
	fmt.Printf("Total Brazilian states/territories supported: %d\n", len(brazil.GetSubdivisions()))
	fmt.Printf("Holiday categories: %v\n", brazil.GetCategories())
	fmt.Printf("Primary language: %s\n", brazil.GetLanguage())

	// Count holidays by category
	categoryCount := make(map[string]int)
	for _, hi := range sortedHolidays {
		categoryCount[string(hi.holiday.Category)]++
	}

	fmt.Println("\nHolidays by category:")
	for category, count := range categoryCount {
		fmt.Printf("  %s: %d\n", category, count)
	}

	fmt.Println("\n🌎 ABOUT BRAZILIAN HOLIDAYS")
	fmt.Println("===========================")
	fmt.Println("Brazil has a rich mixture of public, religious, and cultural holidays.")
	fmt.Println("The country celebrates both Catholic traditions (like Good Friday and")
	fmt.Println("Corpus Christi) and uniquely Brazilian events (like Carnival and")
	fmt.Println("Tiradentes Day). This implementation includes all major federal")
	fmt.Println("holidays celebrated across Brazil's 27 states and federal district.")
}

func getCategoryIcon(category string) string {
	switch category {
	case "public":
		return "🏛️"
	case "religious":
		return "✝️"
	case "carnival":
		return "🎭"
	case "national":
		return "🇧🇷"
	default:
		return "📅"
	}
}
