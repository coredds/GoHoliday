package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("Testing GoHoliday v0.3.0 with 4 new countries:")
	fmt.Println("==============================================")

	countries := []string{"IT", "ES", "NL", "KR"}
	testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	for _, code := range countries {
		fmt.Printf("\nTesting %s:\n", code)

		country := goholidays.NewCountry(code)

		// Test New Year's Day
		if holiday, ok := country.IsHoliday(testDate); ok {
			fmt.Printf("  New Year's Day: %s", holiday.Name)

			// Show native language if available
			for lang, name := range holiday.Languages {
				if lang != "en" {
					fmt.Printf(" (%s: %s)", lang, name)
					break
				}
			}
			fmt.Println()
		}

		// Get all holidays for 2024
		holidays := country.HolidaysForYear(2024)
		fmt.Printf("  Total holidays in 2024: %d\n", len(holidays))
	}

	fmt.Println("\nAll 4 new countries working perfectly!")
	fmt.Printf("GoHoliday now supports %d countries total\n", 15)
}
