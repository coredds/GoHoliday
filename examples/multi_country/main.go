package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHolidays Multi-Country Comparison Demo")
	fmt.Println("========================================")

	// Create providers for different countries
	us := goholidays.NewCountry("US")
	ca := goholidays.NewCountry("CA")

	// Compare holidays for 2024
	fmt.Println("\n1. Holiday Comparison for 2024:")

	usHolidays := us.HolidaysForYear(2024)
	caHolidays := ca.HolidaysForYear(2024)

	fmt.Printf("\nUnited States (%d holidays):\n", len(usHolidays))
	printSortedHolidays(usHolidays)

	fmt.Printf("\nCanada (%d holidays):\n", len(caHolidays))
	printSortedHolidays(caHolidays)

	// Compare specific dates
	fmt.Println("\n2. Specific Date Comparisons:")

	testDates := []struct {
		date string
		desc string
	}{
		{"2024-01-01", "New Year's Day"},
		{"2024-07-01", "Canada Day (CA only)"},
		{"2024-07-04", "Independence Day (US only)"},
		{"2024-10-14", "Thanksgiving (different dates!)"},
		{"2024-11-28", "US Thanksgiving"},
		{"2024-12-25", "Christmas Day"},
	}

	for _, test := range testDates {
		date, _ := time.Parse("2006-01-02", test.date)

		fmt.Printf("\n%s (%s):\n", test.date, test.desc)

		if usHoliday, isUSHoliday := us.IsHoliday(date); isUSHoliday {
			fmt.Printf("  🇺🇸 US: %s\n", usHoliday.Name)
		} else {
			fmt.Printf("  🇺🇸 US: Not a holiday\n")
		}

		if caHoliday, isCAHoliday := ca.IsHoliday(date); isCAHoliday {
			fmt.Printf("  🇨🇦 CA: %s\n", caHoliday.Name)
		} else {
			fmt.Printf("  🇨🇦 CA: Not a holiday\n")
		}
	}

	// Business day comparison
	fmt.Println("\n3. Business Day Analysis:")

	usCalc := goholidays.NewBusinessDayCalculator(us)
	caCalc := goholidays.NewBusinessDayCalculator(ca)

	// Count business days in July 2024
	july1 := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	july31 := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)

	usBusinessDays := usCalc.BusinessDaysBetween(july1, july31)
	caBusinessDays := caCalc.BusinessDaysBetween(july1, july31)

	fmt.Printf("\nBusiness days in July 2024:\n")
	fmt.Printf("  🇺🇸 US: %d days\n", usBusinessDays)
	fmt.Printf("  🇨🇦 CA: %d days\n", caBusinessDays)

	if usBusinessDays != caBusinessDays {
		fmt.Printf("  ⚠️  Different! US has July 4th holiday, CA has July 1st\n")
	}

	// Language comparison
	fmt.Println("\n4. Multi-Language Support:")

	newYearsUS, _ := us.IsHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	newYearsCA, _ := ca.IsHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	fmt.Println("\nNew Year's Day translations:")
	fmt.Printf("  🇺🇸 US English: %s\n", newYearsUS.Languages["en"])
	fmt.Printf("  🇺🇸 US Spanish: %s\n", newYearsUS.Languages["es"])
	fmt.Printf("  🇨🇦 CA English: %s\n", newYearsCA.Languages["en"])
	fmt.Printf("  🇨🇦 CA French:  %s\n", newYearsCA.Languages["fr"])

	// Scheduling comparison
	fmt.Println("\n5. Holiday-Aware Scheduling:")

	usScheduler := goholidays.NewHolidayAwareScheduler(us)
	caScheduler := goholidays.NewHolidayAwareScheduler(ca)

	// Schedule monthly meetings starting June 1st
	startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	usMeetings := usScheduler.ScheduleRecurring(startDate, 30*24*time.Hour, 3)
	caMeetings := caScheduler.ScheduleRecurring(startDate, 30*24*time.Hour, 3)

	fmt.Printf("\nMonthly meetings starting %s:\n", startDate.Format("2006-01-02"))

	for i := 0; i < len(usMeetings) && i < len(caMeetings); i++ {
		fmt.Printf("  Meeting %d:\n", i+1)
		fmt.Printf("    🇺🇸 US: %s (%s)\n", usMeetings[i].Format("2006-01-02"), usMeetings[i].Weekday())
		fmt.Printf("    🇨🇦 CA: %s (%s)\n", caMeetings[i].Format("2006-01-02"), caMeetings[i].Weekday())

		if !usMeetings[i].Equal(caMeetings[i]) {
			fmt.Printf("    ⚠️  Different dates due to different holidays!\n")
		}
	}

	// Performance comparison
	fmt.Println("\n6. Performance Comparison:")

	start := time.Now()
	for i := 0; i < 10000; i++ {
		us.IsHoliday(time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC))
	}
	usDuration := time.Since(start)

	start = time.Now()
	for i := 0; i < 10000; i++ {
		ca.IsHoliday(time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC))
	}
	caDuration := time.Since(start)

	fmt.Printf("\n10,000 holiday checks:\n")
	fmt.Printf("  🇺🇸 US: %v (%v avg)\n", usDuration, usDuration/10000)
	fmt.Printf("  🇨🇦 CA: %v (%v avg)\n", caDuration, caDuration/10000)

	fmt.Println("\n✅ Multi-country support working perfectly!")
	fmt.Println("Ready to add more countries following the same pattern!")
}

func printSortedHolidays(holidays map[time.Time]*goholidays.Holiday) {
	// Convert to slice for sorting
	type holidayEntry struct {
		date    time.Time
		holiday *goholidays.Holiday
	}

	var sorted []holidayEntry
	for date, holiday := range holidays {
		sorted = append(sorted, holidayEntry{date, holiday})
	}

	// Simple bubble sort by date
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].date.After(sorted[j].date) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Print sorted holidays
	for _, entry := range sorted {
		fmt.Printf("  %s: %s\n", entry.date.Format("2006-01-02"), entry.holiday.Name)
	}
}
