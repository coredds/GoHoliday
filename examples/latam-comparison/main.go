package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("🌎 LATIN AMERICA HOLIDAY COMPARISON 🌎")
	fmt.Println("=====================================")
	fmt.Println("Comparing holiday traditions between Brazil and Mexico")

	// Create country instances
	brazil := goholidays.NewCountry("BR")
	mexico := goholidays.NewCountry("MX")

	year := 2024
	brHolidays := brazil.HolidaysForYear(year)
	mxHolidays := mexico.HolidaysForYear(year)

	fmt.Printf("\n📊 HOLIDAY STATISTICS FOR %d\n", year)
	fmt.Println("================================")
	fmt.Printf("🇧🇷 Brazil: %d holidays\n", len(brHolidays))
	fmt.Printf("🇲🇽 Mexico: %d holidays\n", len(mxHolidays))

	// Find common holidays
	fmt.Println("\n🤝 COMMON CELEBRATIONS")
	fmt.Println("======================")

	commonDates := []struct {
		date        time.Time
		description string
	}{
		{time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day"},
		{time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC), "Labour Day"},
		{time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day"},
	}

	for _, common := range commonDates {
		brHoliday, brExists := brazil.IsHoliday(common.date)
		mxHoliday, mxExists := mexico.IsHoliday(common.date)

		if brExists && mxExists {
			fmt.Printf("📅 %s:\n", common.date.Format("January 2"))
			fmt.Printf("   🇧🇷 %s\n", brHoliday.Name)
			fmt.Printf("   🇲🇽 %s\n", mxHoliday.Name)
			fmt.Printf("   🌍 %s\n\n", common.description)
		}
	}

	// Easter-based holidays
	fmt.Println("✝️ EASTER CELEBRATIONS")
	fmt.Println("======================")

	easter := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC) // Easter 2024
	goodFriday := easter.AddDate(0, 0, -2)

	brGoodFriday, brExists := brazil.IsHoliday(goodFriday)
	mxGoodFriday, mxExists := mexico.IsHoliday(goodFriday)

	if brExists && mxExists {
		fmt.Printf("🕊️ Good Friday (%s):\n", goodFriday.Format("January 2"))
		fmt.Printf("   🇧🇷 %s\n", brGoodFriday.Name)
		fmt.Printf("   🇲🇽 %s\n", mxGoodFriday.Name)
	}

	// Unique celebrations
	fmt.Println("\n🎭 UNIQUE CULTURAL CELEBRATIONS")
	fmt.Println("===============================")

	fmt.Println("🇧🇷 BRAZIL SPECIALS:")
	carnivalMonday := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brazil.IsHoliday(carnivalMonday); exists {
		fmt.Printf("   🎪 %s (%s)\n", holiday.Name, carnivalMonday.Format("Jan 2"))
		fmt.Println("      World's most famous carnival celebration")
	}

	tiradentes := time.Date(2024, 4, 21, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brazil.IsHoliday(tiradentes); exists {
		fmt.Printf("   🦾 %s (%s)\n", holiday.Name, tiradentes.Format("Jan 2"))
		fmt.Println("      Honoring Brazilian independence hero")
	}

	fmt.Println("\n🇲🇽 MEXICO SPECIALS:")
	dayOfDead := time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(dayOfDead); exists {
		fmt.Printf("   💀 %s (%s)\n", holiday.Name, dayOfDead.Format("Jan 2"))
		fmt.Println("      UNESCO World Heritage tradition")
	}

	guadalupe := time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(guadalupe); exists {
		fmt.Printf("   🌹 %s (%s)\n", holiday.Name, guadalupe.Format("Jan 2"))
		fmt.Println("      Patron saint of the Americas")
	}

	// Independence Days
	fmt.Println("\n🏛️ INDEPENDENCE CELEBRATIONS")
	fmt.Println("============================")

	brIndependence := time.Date(2024, 9, 7, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brazil.IsHoliday(brIndependence); exists {
		fmt.Printf("🇧🇷 %s (%s)\n", holiday.Name, brIndependence.Format("Jan 2"))
		fmt.Println("   Independence from Portugal (1822)")
	}

	mxIndependence := time.Date(2024, 9, 16, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mexico.IsHoliday(mxIndependence); exists {
		fmt.Printf("🇲🇽 %s (%s)\n", holiday.Name, mxIndependence.Format("Jan 2"))
		fmt.Println("   El Grito de Dolores (1810)")
	}

	// Language comparison
	fmt.Println("\n🗣️ LANGUAGE SUPPORT")
	fmt.Println("===================")

	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	if brHoliday, exists := brazil.IsHoliday(newYear); exists {
		fmt.Println("🇧🇷 Brazil (Portuguese/English):")
		fmt.Printf("   PT: %s\n", brHoliday.Name)
		if brHoliday.Languages != nil {
			fmt.Printf("   EN: %s\n", brHoliday.Languages["en"])
		}
	}

	if mxHoliday, exists := mexico.IsHoliday(newYear); exists {
		fmt.Println("🇲🇽 Mexico (Spanish/English):")
		fmt.Printf("   ES: %s\n", mxHoliday.Name)
		if mxHoliday.Languages != nil {
			fmt.Printf("   EN: %s\n", mxHoliday.Languages["en"])
		}
	}

	// Holiday categories analysis
	fmt.Println("\n📈 CATEGORY ANALYSIS")
	fmt.Println("===================")

	brCategories := countCategories(brHolidays)
	mxCategories := countCategories(mxHolidays)

	fmt.Println("🇧🇷 Brazil:")
	for category, count := range brCategories {
		fmt.Printf("   %s: %d\n", category, count)
	}

	fmt.Println("🇲🇽 Mexico:")
	for category, count := range mxCategories {
		fmt.Printf("   %s: %d\n", category, count)
	}

	fmt.Println("\n🌟 CULTURAL INSIGHTS")
	fmt.Println("===================")
	fmt.Println("Both Brazil and Mexico share strong Catholic influences with")
	fmt.Println("holidays like Good Friday and Christmas. However, each country")
	fmt.Println("has developed unique celebrations reflecting their distinct")
	fmt.Println("histories and cultural identities:")
	fmt.Println()
	fmt.Println("🇧🇷 Brazil emphasizes Carnival (African and Portuguese influences)")
	fmt.Println("🇲🇽 Mexico highlights Day of the Dead (indigenous Aztec traditions)")
	fmt.Println()
	fmt.Println("Both countries celebrate their independence from colonial rule")
	fmt.Println("with national holidays in September, just 9 days apart!")
}

func countCategories(holidays map[time.Time]*goholidays.Holiday) map[string]int {
	categories := make(map[string]int)
	for _, holiday := range holidays {
		categories[string(holiday.Category)]++
	}
	return categories
}
