package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	// Test Brazil
	fmt.Println("Testing Brazil (BR)...")
	brCountry := goholidays.NewCountry("BR")
	brHolidays := brCountry.HolidaysForYear(2024)
	fmt.Printf("Brazil has %d holidays in 2024:\n", len(brHolidays))

	// Check Independence Day
	independenceDay := time.Date(2024, 9, 7, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brCountry.IsHoliday(independenceDay); exists {
		fmt.Printf("✓ %s: %s\n", independenceDay.Format("2006-01-02"), holiday.Name)
	}

	// Check Carnival Monday
	carnivalMonday := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := brCountry.IsHoliday(carnivalMonday); exists {
		fmt.Printf("✓ %s: %s\n", carnivalMonday.Format("2006-01-02"), holiday.Name)
	}

	fmt.Println()

	// Test Mexico
	fmt.Println("Testing Mexico (MX)...")
	mxCountry := goholidays.NewCountry("MX")
	mxHolidays := mxCountry.HolidaysForYear(2024)
	fmt.Printf("Mexico has %d holidays in 2024:\n", len(mxHolidays))

	// Check Independence Day
	independenceDay = time.Date(2024, 9, 16, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mxCountry.IsHoliday(independenceDay); exists {
		fmt.Printf("✓ %s: %s\n", independenceDay.Format("2006-01-02"), holiday.Name)
	}

	// Check Constitution Day (first Monday of February)
	constitutionDay := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)
	if holiday, exists := mxCountry.IsHoliday(constitutionDay); exists {
		fmt.Printf("✓ %s: %s\n", constitutionDay.Format("2006-01-02"), holiday.Name)
	}
}
