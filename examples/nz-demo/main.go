package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("=== GoHolidays: New Zealand Holiday Demo ===")
	
	// Create New Zealand country
	nz := goholidays.NewCountry("NZ")
	
	year := 2024
	fmt.Printf("\n🇳🇿 New Zealand Public Holidays for %d:\n", year)
	fmt.Println(strings.Repeat("=", 50))
	
	// Get all holidays for 2024
	holidays := nz.HolidaysForYear(year)
	
	// Sort holidays by date for display
	type holidayPair struct {
		date    time.Time
		holiday *goholidays.Holiday
	}
	
	var sortedHolidays []holidayPair
	for date, holiday := range holidays {
		sortedHolidays = append(sortedHolidays, holidayPair{date, holiday})
	}
	
	// Simple sort by date
	for i := 0; i < len(sortedHolidays)-1; i++ {
		for j := i + 1; j < len(sortedHolidays); j++ {
			if sortedHolidays[i].date.After(sortedHolidays[j].date) {
				sortedHolidays[i], sortedHolidays[j] = sortedHolidays[j], sortedHolidays[i]
			}
		}
	}
	
	// Display holidays with bilingual support
	for _, hp := range sortedHolidays {
		date := hp.date
		holiday := hp.holiday
		
		// Format date nicely
		dateStr := date.Format("Monday, January 2")
		
		// Get Māori translation if available
		englishName := holiday.Name
		maoriName := ""
		if holiday.Languages != nil {
			if name, exists := holiday.Languages["mi"]; exists {
				maoriName = name
			}
		}
		
		if maoriName != "" && maoriName != englishName {
			fmt.Printf("📅 %-25s %s (%s)\n", dateStr, englishName, maoriName)
		} else {
			fmt.Printf("📅 %-25s %s\n", dateStr, englishName)
		}
	}
	
	// Demonstrate unique New Zealand features
	fmt.Printf("\n🌟 Special New Zealand Holiday Features:\n")
	fmt.Println(strings.Repeat("-", 40))
	
	// Check for Matariki
	matarikiFound := false
	for _, holiday := range holidays {
		if holiday.Name == "Matariki" {
			matarikiFound = true
			fmt.Printf("✨ Matariki (Māori New Year) is celebrated on %s\n", 
				holiday.Date.Format("January 2, 2006"))
			break
		}
	}
	
	if !matarikiFound {
		fmt.Println("⚠️  Matariki date not available for this year")
	}
	
	// Show bilingual support
	fmt.Println("🗣️  All holidays include English and Māori (te reo Māori) names")
	
	// Show Queen's Birthday calculation
	for _, holiday := range holidays {
		if holiday.Name == "Queen's Birthday" {
			fmt.Printf("👑 Queen's Birthday is on the first Monday in June: %s\n", 
				holiday.Date.Format("January 2"))
			break
		}
	}
	
	// Show Labour Day calculation
	for _, holiday := range holidays {
		if holiday.Name == "Labour Day" {
			fmt.Printf("⚒️  Labour Day is on the fourth Monday in October: %s\n", 
				holiday.Date.Format("January 2"))
			break
		}
	}
	
	fmt.Printf("\n📊 Total public holidays for %d: %d\n", year, len(holidays))
	
	// Demonstrate working with specific dates
	fmt.Printf("\n🔍 Holiday Lookup Examples:\n")
	fmt.Println(strings.Repeat("-", 30))
	
	testDates := []time.Time{
		time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),   // Waitangi Day
		time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),  // ANZAC Day
		time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),  // Matariki
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // Not a holiday
	}
	
	for _, testDate := range testDates {
		if holiday, exists := holidays[testDate]; exists {
			fmt.Printf("✅ %s is %s\n", testDate.Format("Jan 2"), holiday.Name)
		} else {
			fmt.Printf("❌ %s is not a public holiday\n", testDate.Format("Jan 2"))
		}
	}
	
	// Demonstrate IsHoliday method
	fmt.Printf("\n🎯 Individual Date Checking:\n")
	fmt.Println(strings.Repeat("-", 25))
	
	waitangiDay := time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC)
	if holiday, isHoliday := nz.IsHoliday(waitangiDay); isHoliday {
		fmt.Printf("✅ %s is %s\n", waitangiDay.Format("January 2"), holiday.Name)
	} else {
		fmt.Printf("❌ %s is not a holiday\n", waitangiDay.Format("January 2"))
	}
	
	// Test a regular day
	regularDay := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if holiday, isHoliday := nz.IsHoliday(regularDay); isHoliday {
		fmt.Printf("✅ %s is %s\n", regularDay.Format("January 2"), holiday.Name)
	} else {
		fmt.Printf("❌ %s is not a holiday\n", regularDay.Format("January 2"))
	}
	
	fmt.Println("\n🎉 New Zealand holidays integration successful!")
	fmt.Println("\nNext steps: Try regional holidays with subdivisions like 'AUK' (Auckland), 'WGN' (Wellington), etc.")
}
