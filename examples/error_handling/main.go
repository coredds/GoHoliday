package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHoliday Error Handling Example")
	fmt.Println("=================================")

	// 1. Original API (backward compatible)
	fmt.Println("\n1. Original API (still works)")
	demonstrateOriginalAPI()

	// 2. Enhanced API with error handling
	fmt.Println("\n2. Enhanced API with Error Handling")
	demonstrateEnhancedAPI()

	// 3. Context support
	fmt.Println("\n3. Context Support")
	demonstrateContextSupport()

	fmt.Println("\nExample completed!")
}

func demonstrateOriginalAPI() {
	// Original API - no changes needed
	country := goholidays.NewCountry("US")
	fmt.Printf("  Created country: %s\n", country.GetCountryCode())

	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday := country.IsHoliday(date)
	if isHoliday {
		fmt.Printf("  Holiday: %s\n", holiday.Name)
	}

	holidays := country.HolidaysForYear(2024)
	fmt.Printf("  Found %d holidays in 2024\n", len(holidays))
}

func demonstrateEnhancedAPI() {
	// Enhanced creation with validation
	fmt.Print("  Creating country with validation: ")
	country, err := goholidays.NewCountryWithError("US")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Success (%s)\n", country.GetCountryCode())

	// Enhanced holiday checking
	fmt.Print("  Checking holiday with validation: ")
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday, err := country.IsHolidayWithError(date)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if isHoliday {
		fmt.Printf("Holiday: %s\n", holiday.Name)
	}

	// Show error handling
	fmt.Print("  Testing invalid country: ")
	_, err = goholidays.NewCountryWithError("XX")
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)

		var holidayErr *goholidays.HolidayError
		if errors.As(err, &holidayErr) {
			fmt.Printf("    Error Code: %d, Country: %s\n", holidayErr.Code, holidayErr.Country)
		}
	}

	fmt.Print("  Testing invalid year: ")
	invalidDate := time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)
	_, _, err = country.IsHolidayWithError(invalidDate)
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
	}
}

func demonstrateContextSupport() {
	country := goholidays.NewCountry("US")

	// Valid context
	fmt.Print("  Using valid context: ")
	ctx := context.Background()
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	_, isHoliday, err := country.IsHolidayWithContext(ctx, date)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else if isHoliday {
		fmt.Println("Holiday found")
	}

	// Cancelled context
	fmt.Print("  Using cancelled context: ")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, _, err = country.IsHolidayWithContext(ctx, date)
	if err != nil {
		fmt.Printf("Error (expected): %v\n", err)
		if goholidays.IsContextCancelled(err) {
			fmt.Println("    Context was cancelled")
		}
	}

	// Context with timeout
	fmt.Print("  Using context with timeout: ")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	holidays, err := country.HolidaysForYearWithContext(ctx, 2024)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found %d holidays\n", len(holidays))
	}

	// Show utility functions
	fmt.Printf("  Supported countries: %d total\n", len(goholidays.GetSupportedCountries()))
	fmt.Printf("  Is 'US' valid? %v\n", goholidays.IsValidCountry("US"))
	fmt.Printf("  Is 'XX' valid? %v\n", goholidays.IsValidCountry("XX"))
}
