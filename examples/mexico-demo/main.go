package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("🇲🇽 MEXICO HOLIDAYS SHOWCASE 🇲🇽")
	fmt.Println("===============================")

	// Create Mexico country instance
	mexico := goholidays.NewCountry("MX")

	// Get holidays for 2024
	holidays := mexico.HolidaysForYear(2024)

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

	fmt.Printf("\n📅 MEXICAN HOLIDAYS FOR 2024 (%d total)\n", len(holidays))
	fmt.Println("=======================================")

	for _, hi := range sortedHolidays {
		// Get both Spanish and English names
		esName := hi.holiday.Name
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
			esName,
			enName,
			string(hi.holiday.Category))
	}

	fmt.Println("\n🏛️ CONSTITUTIONAL REFORMS")
	fmt.Println("=========================")
	fmt.Println("Since 2006, Mexico moved several holidays to Mondays to create")
	fmt.Println("long weekends, boosting tourism and family time:")

	// Check variable Monday holidays
	constitution := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(constitution); exists {
		fmt.Printf("📜 %s: %s (1st Monday of February)\n", holiday.Name, holiday.Languages["en"])
	}

	juarez := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(juarez); exists {
		fmt.Printf("👨‍⚖️ %s: %s (3rd Monday of March)\n", holiday.Name, holiday.Languages["en"])
	}

	revolution := time.Date(2024, 11, 18, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(revolution); exists {
		fmt.Printf("⚔️ %s: %s (3rd Monday of November)\n", holiday.Name, holiday.Languages["en"])
	}

	fmt.Println("\n🎨 CULTURAL CELEBRATIONS")
	fmt.Println("========================")

	// Highlight Independence Day
	independence := time.Date(2024, 9, 16, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(independence); exists {
		fmt.Printf("🇲🇽 %s: %s\n", holiday.Name, holiday.Languages["en"])
		fmt.Println("   El Grito de Dolores - The cry for independence (1810)")
	}

	// Highlight Day of the Dead
	dayOfDead := time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(dayOfDead); exists {
		fmt.Printf("💀 %s: %s\n", holiday.Name, holiday.Languages["en"])
		fmt.Println("   UNESCO World Heritage celebration honoring ancestors")
	}

	// Highlight Our Lady of Guadalupe
	guadalupe := time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(guadalupe); exists {
		fmt.Printf("🌹 %s: %s\n", holiday.Name, holiday.Languages["en"])
		fmt.Println("   Patron saint of Mexico and the Americas")
	}

	fmt.Println("\n🙏 HOLY WEEK TRADITIONS")
	fmt.Println("=======================")

	// Check Easter-based holidays
	maundyThursday := time.Date(2024, 3, 28, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(maundyThursday); exists {
		fmt.Printf("🕊️ %s: %s\n", holiday.Name, holiday.Languages["en"])
	}

	goodFriday := time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(goodFriday); exists {
		fmt.Printf("✝️ %s: %s\n", holiday.Name, holiday.Languages["en"])
	}

	fmt.Println("\n📊 HOLIDAY STATISTICS")
	fmt.Println("====================")
	fmt.Printf("Total Mexican states/territories supported: %d\n", len(mexico.GetSubdivisions()))
	fmt.Printf("Holiday categories: %v\n", mexico.GetCategories())
	fmt.Printf("Primary language: %s\n", mexico.GetLanguage())

	// Count holidays by category
	categoryCount := make(map[string]int)
	for _, hi := range sortedHolidays {
		categoryCount[string(hi.holiday.Category)]++
	}

	fmt.Println("\nHolidays by category:")
	for category, count := range categoryCount {
		fmt.Printf("  %s: %d\n", category, count)
	}

	fmt.Println("\n🌎 ABOUT MEXICAN HOLIDAYS")
	fmt.Println("=========================")
	fmt.Println("Mexico's holiday calendar reflects its rich cultural heritage,")
	fmt.Println("blending indigenous traditions, Catholic influence, and modern")
	fmt.Println("constitutional reforms. The 2006 labor law changes moved several")
	fmt.Println("holidays to Mondays, creating puentes (bridges) - long weekends")
	fmt.Println("that boost domestic tourism and family gatherings across all")
	fmt.Println("32 Mexican federal entities.")
}

func getCategoryIcon(category string) string {
	switch category {
	case "public":
		return "🏛️"
	case "national":
		return "🇲🇽"
	case "religious":
		return "✝️"
	case "civic":
		return "⚖️"
	case "regional":
		return "🏞️"
	default:
		return "📅"
	}
}
