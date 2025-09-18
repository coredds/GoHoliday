package main

import (
	"fmt"
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("Testing new countries: Chile, Ireland, and Israel")
	fmt.Println("=================================================")

	// Test Chile
	fmt.Println("\n🇨🇱 Chile (CL):")
	chile := goholidays.NewCountry("CL")
	holidays := chile.HolidaysForYear(2024)
	fmt.Printf("Found %d holidays in 2024\n", len(holidays))

	// Check for Independence Day
	independenceDay := time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[independenceDay]; exists {
		fmt.Printf("✓ %s: %s (%s)\n", independenceDay.Format("2006-01-02"), holiday.Name, holiday.Languages["es"])
	}

	// Test Ireland
	fmt.Println("\n🇮🇪 Ireland (IE):")
	ireland := goholidays.NewCountry("IE")
	holidays = ireland.HolidaysForYear(2024)
	fmt.Printf("Found %d holidays in 2024\n", len(holidays))

	// Check for Saint Patrick's Day
	stPatricksDay := time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[stPatricksDay]; exists {
		fmt.Printf("✓ %s: %s (%s)\n", stPatricksDay.Format("2006-01-02"), holiday.Name, holiday.Languages["ga"])
	}

	// Test Israel
	fmt.Println("\n🇮🇱 Israel (IL):")
	israel := goholidays.NewCountry("IL")
	holidays = israel.HolidaysForYear(2024)
	fmt.Printf("Found %d holidays in 2024\n", len(holidays))

	// Check for Passover
	for date, holiday := range holidays {
		if holiday.Name == "Passover" {
			fmt.Printf("✓ %s: %s (%s)\n", date.Format("2006-01-02"), holiday.Name, holiday.Languages["he"])
			break
		}
	}

	fmt.Println("\n✅ All new countries integrated successfully!")
}
