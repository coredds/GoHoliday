package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHoliday Localization Features")
	fmt.Println("==============================")

	// 1. Multi-language Holiday Names
	fmt.Println("\n1. Holiday Names in Different Languages")

	// Create providers with different language settings
	countries := map[string]struct {
		country   *goholidays.Country
		languages []string
	}{
		"CA": {
			country:   goholidays.NewCountry("CA", goholidays.CountryOptions{Language: "fr"}),
			languages: []string{"en", "fr"},
		},
		"ES": {
			country:   goholidays.NewCountry("ES", goholidays.CountryOptions{Language: "es"}),
			languages: []string{"en", "es"},
		},
		"JP": {
			country:   goholidays.NewCountry("JP", goholidays.CountryOptions{Language: "ja"}),
			languages: []string{"en", "ja"},
		},
	}

	// Check New Year's Day in different languages
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Println("New Year's Day translations:")
	for code, info := range countries {
		if holiday, isHoliday := info.country.IsHoliday(newYear); isHoliday {
			fmt.Printf("\n%s:\n", code)
			for _, lang := range info.languages {
				fmt.Printf("- %s: %s\n", lang, holiday.Languages[lang])
			}
		}
	}

	// 2. Regional Holiday Names
	fmt.Println("\n2. Regional Holiday Variations")
	
	// Create Swiss provider with different language regions
	languages := []string{"de", "fr", "it"}
	for _, lang := range languages {
		ch := goholidays.NewCountry("CH", goholidays.CountryOptions{Language: lang})
		fmt.Printf("\nSwiss holidays in %s:\n", lang)
		
		holidays := ch.HolidaysForYear(2024)
		printSortedHolidaysWithLanguage(holidays, lang)
	}

	// 3. Country-Specific Formatting
	fmt.Println("\n3. Country-Specific Date Formats")
	sampleDate := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	
	dateFormats := map[string]string{
		"US": "January 2, 2006",
		"GB": "2 January 2006",
		"JP": "2006年1月2日",
	}

	for country, format := range dateFormats {
		fmt.Printf("%s format: %s\n", country, sampleDate.Format(format))
	}

	// 4. Language-Specific Categories
	fmt.Println("\n4. Holiday Categories in Different Languages")
	categoryTranslations := map[goholidays.HolidayCategory]map[string]string{
		goholidays.CategoryPublic: {
			"en": "Public Holiday",
			"fr": "Jour férié",
			"es": "Fiesta nacional",
			"de": "Feiertag",
		},
		goholidays.CategoryReligious: {
			"en": "Religious Holiday",
			"fr": "Fête religieuse",
			"es": "Fiesta religiosa",
			"de": "Religiöser Feiertag",
		},
	}

	for category, translations := range categoryTranslations {
		fmt.Printf("\n%s:\n", category)
		for lang, name := range translations {
			fmt.Printf("- %s: %s\n", lang, name)
		}
	}

	// 5. Regional Variations
	fmt.Println("\n5. Regional Holiday Variations")
	
	// Create US provider with state-specific holidays
	usStates := []string{"CA", "NY", "TX"}
	for _, state := range usStates {
		us := goholidays.NewCountry("US", goholidays.CountryOptions{
			Subdivisions: []string{state},
		})
		fmt.Printf("\nHolidays specific to %s:\n", state)
		stateHolidays := us.HolidaysForYear(2024)
		printSortedHolidays(stateHolidays)
	}

	fmt.Println("\nThis demonstrates GoHoliday's localization capabilities!")
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

func printSortedHolidaysWithLanguage(holidays map[time.Time]*goholidays.Holiday, lang string) {
	var dates []time.Time
	for date := range holidays {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	for _, date := range dates {
		holiday := holidays[date]
		name := holiday.Languages[lang]
		if name == "" {
			name = holiday.Name // fallback to default name
		}
		fmt.Printf("- %s: %s\n", date.Format("Jan 2"), name)
	}
}

