package main

import (
	"fmt"
	"time"

	"github.com/coredds/goholiday/countries"
)

func main() {
	fmt.Println("🌍 goholiday - New Countries Demo")
	fmt.Println("==================================")
	fmt.Println()

	// Create providers for the three new countries
	providers := map[string]countries.HolidayProvider{
		"🇵🇹 Portugal": countries.NewPTProvider(),
		"🇮🇹 Italy":    countries.NewITProvider(),
		"🇮🇳 India":    countries.NewINProvider(),
	}

	year := 2024
	fmt.Printf("📅 Showing holidays for %d\n\n", year)

	for countryName, provider := range providers {
		fmt.Printf("%s (%s)\n", countryName, provider.GetCountryCode())
		fmt.Printf("├─ Subdivisions: %d\n", len(provider.GetSupportedSubdivisions()))
		fmt.Printf("├─ Categories: %v\n", provider.GetSupportedCategories())

		holidays := provider.LoadHolidays(year)
		fmt.Printf("└─ National Holidays: %d\n", len(holidays))

		// Show holidays in chronological order
		var sortedDates []time.Time
		for date := range holidays {
			sortedDates = append(sortedDates, date)
		}

		// Simple bubble sort by date
		for i := 0; i < len(sortedDates)-1; i++ {
			for j := 0; j < len(sortedDates)-i-1; j++ {
				if sortedDates[j].After(sortedDates[j+1]) {
					sortedDates[j], sortedDates[j+1] = sortedDates[j+1], sortedDates[j]
				}
			}
		}

		// Show first 5 holidays
		fmt.Println("   📋 Major Holidays:")
		count := 0
		for _, date := range sortedDates {
			if count >= 5 {
				break
			}
			holiday := holidays[date]
			fmt.Printf("   • %s - %s (%s)\n",
				date.Format("Jan 02"),
				holiday.Name,
				holiday.Category)
			count++
		}
		if len(holidays) > 5 {
			fmt.Printf("   ... and %d more holidays\n", len(holidays)-5)
		}
		fmt.Println()
	}

	// Demonstrate regional/state holidays
	fmt.Println("🏛️ Regional Holiday Examples")
	fmt.Println("============================")

	// Italy regional holidays
	itProvider := countries.NewITProvider()
	lombardyHolidays := itProvider.GetRegionalHolidays(year, "LOM")
	if len(lombardyHolidays) > 0 {
		fmt.Println("🇮🇹 Italy - Lombardy Region:")
		for _, holiday := range lombardyHolidays {
			fmt.Printf("   • %s - %s (%s)\n",
				holiday.Date.Format("Jan 02"),
				holiday.Name,
				holiday.Category)
		}
		fmt.Println()
	}

	// India state holidays
	inProvider := countries.NewINProvider()
	maharashtraHolidays := inProvider.GetStateHolidays(year, "MH")
	if len(maharashtraHolidays) > 0 {
		fmt.Println("🇮🇳 India - Maharashtra State:")
		for _, holiday := range maharashtraHolidays {
			fmt.Printf("   • %s - %s (%s)\n",
				holiday.Date.Format("Jan 02"),
				holiday.Name,
				holiday.Category)
		}
		fmt.Println()
	}

	// Show India's religious diversity
	majorFestivals := inProvider.GetMajorFestivals(year)
	if len(majorFestivals) > 0 {
		fmt.Println("🇮🇳 India - Major Religious Festivals (Approximate Dates):")

		// Group by religion
		religionFestivals := make(map[string][]string)
		for _, festival := range majorFestivals {
			religionFestivals[festival.Category] = append(religionFestivals[festival.Category],
				fmt.Sprintf("%s - %s", festival.Date.Format("Jan 02"), festival.Name))
		}

		for religion, festivals := range religionFestivals {
			fmt.Printf("   %s festivals:\n", religion)
			for _, festival := range festivals {
				fmt.Printf("     • %s\n", festival)
			}
		}
		fmt.Println()
	}

	// Language support demonstration
	fmt.Println("🌐 Multi-Language Support")
	fmt.Println("=========================")

	// Show Christmas in different countries with local names
	christmasDate := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)

	for countryName, provider := range providers {
		holidays := provider.LoadHolidays(year)
		if christmas, exists := holidays[christmasDate]; exists {
			fmt.Printf("%s Christmas:\n", countryName)
			for lang, name := range christmas.Languages {
				fmt.Printf("   %s: %s\n", lang, name)
			}
			fmt.Println()
		}
	}

	// Easter calculation demonstration
	fmt.Println("🐰 Easter Calculation Accuracy")
	fmt.Println("==============================")

	ptProvider := countries.NewPTProvider()
	easterYears := []int{2024, 2025, 2026, 2027}

	fmt.Println("Easter Sunday dates:")
	for _, testYear := range easterYears {
		holidays := ptProvider.LoadHolidays(testYear)
		for _, holiday := range holidays {
			if holiday.Name == "Easter Sunday" {
				fmt.Printf("   %d: %s\n", testYear, holiday.Date.Format("January 2, 2006"))
				break
			}
		}
	}
	fmt.Println()

	fmt.Println("✅ All three countries implemented successfully!")
	fmt.Println("📊 Total new holiday providers: 3")
	fmt.Println("🌍 Countries covered: Portugal, Italy, India")
	fmt.Println("🎯 Features: Multi-language, Regional holidays, Religious diversity")
}
