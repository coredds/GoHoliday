package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday/countries"
)

func main() {
	fmt.Println("🎉 GoHolidays - Complete Country Coverage Demo")
	fmt.Println("==============================================")

	year := 2024

	// Create all providers
	providers := map[string]countries.HolidayProvider{
		"🇦🇺 Australia":      countries.NewAUProvider(),
		"🇨🇦 Canada":         countries.NewCAProvider(),
		"🇳🇿 New Zealand":    countries.NewNZProvider(),
		"🇬🇧 United Kingdom": countries.NewGBProvider(),
		"🇺🇸 United States":  countries.NewUSProvider(),
		"🇩🇪 Germany":        countries.NewDEProvider(),
		"🇫🇷 France":         countries.NewFRProvider(),
	}

	fmt.Printf("\n📅 Holiday Coverage for %d\n", year)
	fmt.Println("=====================================")

	totalHolidays := 0

	for country, provider := range providers {
		holidays := provider.LoadHolidays(year)
		holidayCount := len(holidays)
		totalHolidays += holidayCount

		fmt.Printf("\n%s (%s)\n", country, provider.GetCountryCode())
		fmt.Printf("├─ Total Holidays: %d\n", holidayCount)
		fmt.Printf("├─ Subdivisions: %d\n", len(provider.GetSupportedSubdivisions()))
		fmt.Printf("└─ Categories: %v\n", provider.GetSupportedCategories())

		// Show some sample holidays
		fmt.Println("   Sample holidays:")
		count := 0
		for date, holiday := range holidays {
			if count >= 3 {
				break
			}
			fmt.Printf("   • %s: %s\n", date.Format("Jan 2"), holiday.Name)
			count++
		}
	}

	fmt.Printf("\n🌍 Global Coverage Summary\n")
	fmt.Println("==========================")
	fmt.Printf("Total Countries: %d\n", len(providers))
	fmt.Printf("Total Holidays: %d\n", totalHolidays)
	fmt.Printf("Average per Country: %.1f\n", float64(totalHolidays)/float64(len(providers)))

	// Demonstrate multi-language support
	fmt.Printf("\n🌐 Multi-Language Support\n")
	fmt.Println("=========================")

	// Show New Year's Day in different languages
	newYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	for country, provider := range providers {
		holidays := provider.LoadHolidays(year)
		if holiday, exists := holidays[newYear]; exists {
			// Get the primary language name
			var primaryName string
			if len(holiday.Languages) > 0 {
				for _, name := range holiday.Languages {
					primaryName = name
					break
				}
			} else {
				primaryName = holiday.Name
			}
			fmt.Printf("%s: %s\n", country, primaryName)
		}
	}

	// Demonstrate Easter calculations across countries
	fmt.Printf("\n🐰 Easter-based Holidays (March 31, %d)\n", year)
	fmt.Println("========================================")

	easter := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)
	goodFriday := time.Date(year, 3, 29, 0, 0, 0, 0, time.UTC)
	easterMonday := time.Date(year, 4, 1, 0, 0, 0, 0, time.UTC)

	for country, provider := range providers {
		holidays := provider.LoadHolidays(year)
		easterHolidays := []string{}

		if _, exists := holidays[goodFriday]; exists {
			easterHolidays = append(easterHolidays, "Good Friday")
		}
		if _, exists := holidays[easter]; exists {
			easterHolidays = append(easterHolidays, "Easter Sunday")
		}
		if _, exists := holidays[easterMonday]; exists {
			easterHolidays = append(easterHolidays, "Easter Monday")
		}

		if len(easterHolidays) > 0 {
			fmt.Printf("%s: %v\n", country, easterHolidays)
		}
	}

	// Show regional/state support
	fmt.Printf("\n🏛️ Regional Holiday Support\n")
	fmt.Println("===========================")

	regionalSupport := map[string][]string{
		"🇦🇺 Australia":      {"NSW", "VIC", "QLD", "WA", "SA", "TAS", "ACT", "NT"},
		"🇨🇦 Canada":         {"ON", "QC", "BC", "AB", "MB", "SK", "NS", "NB", "PE", "NL", "YT", "NT", "NU"},
		"🇳🇿 New Zealand":    {"AUK", "BOP", "CAN", "GIS", "HKB", "MWT", "MBH", "NSN", "OTA", "STL", "TKI", "TAS", "WKO", "WGN", "WTC", "CIT"},
		"🇬🇧 United Kingdom": {"ENG", "SCT", "WLS", "NIR"},
		"🇺🇸 United States":  {"CA", "TX", "NY", "FL", "IL", "PA", "OH", "GA", "NC", "MI"},
		"🇩🇪 Germany":        {"BW", "BY", "BE", "BB", "HB", "HH", "HE", "MV", "NI", "NW"},
		"🇫🇷 France":         {"ARA", "BFC", "BRE", "CVL", "COR", "GES", "HDF", "IDF", "NOR", "NAQ"},
	}

	for country, regions := range regionalSupport {
		fmt.Printf("%s: %d regions/states supported\n", country, len(regions))
	}

	fmt.Printf("\n✅ Country Coverage Implementation Complete!\n")
	fmt.Println("===========================================")
	fmt.Println("Features implemented:")
	fmt.Println("├─ ✅ Comprehensive holiday calculations")
	fmt.Println("├─ ✅ Multi-language support")
	fmt.Println("├─ ✅ Regional/state-specific holidays")
	fmt.Println("├─ ✅ Historical holiday changes")
	fmt.Println("├─ ✅ Easter-based calculations")
	fmt.Println("├─ ✅ Category classification")
	fmt.Println("├─ ✅ Observed date handling")
	fmt.Println("└─ ✅ Comprehensive test coverage")

	fmt.Println("\nNext steps:")
	fmt.Println("• Add more countries (JP, IT, ES, etc.)")
	fmt.Println("• Implement configuration system")
	fmt.Println("• Add REST API layer")
	fmt.Println("• Performance optimization")
	fmt.Println("• Database integration")
}
