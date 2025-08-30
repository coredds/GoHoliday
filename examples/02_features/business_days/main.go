package main

import (
	"fmt"
	"time"

	goholidays "github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHoliday Business Day Features")
	fmt.Println("==============================")

	// Create a US holiday provider with business day settings
	us := goholidays.NewCountry("US", goholidays.CountryOptions{
		Categories: []goholidays.HolidayCategory{
			goholidays.CategoryPublic,
			goholidays.CategoryBank,
		},
	})

	// Create business day calculator
	calc := goholidays.NewBusinessDayCalculator(us)

	// 1. Basic Business Day Calculation
	fmt.Println("\n1. Basic Business Day Count")
	start := time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC) // Monday before Christmas
	end := time.Date(2024, 12, 27, 0, 0, 0, 0, time.UTC)   // Friday after Christmas
	days := calc.BusinessDaysBetween(start, end)
	fmt.Printf("Business days between %s and %s: %d\n",
		start.Format("Jan 2"), end.Format("Jan 2"), days)

	// 2. Add Business Days
	fmt.Println("\n2. Adding Business Days")
	date := time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC) // Day before Independence Day
	for i := 1; i <= 3; i++ {
		next := calc.AddBusinessDays(date, i)
		fmt.Printf("Adding %d business day(s) to %s: %s\n",
			i, date.Format("Jan 2"), next.Format("Jan 2"))
	}

	// 3. Previous Business Days
	fmt.Println("\n3. Finding Previous Business Days")
	date = time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC) // Day after Independence Day
	for i := 1; i <= 3; i++ {
		current := date
		for j := 0; j < i; j++ {
			current = calc.PreviousBusinessDay(current)
		}
		fmt.Printf("Going back %d business day(s) from %s: %s\n",
			i, date.Format("Jan 2"), current.Format("Jan 2"))
	}

	// 4. Next Business Day
	fmt.Println("\n4. Finding Next Business Day")
	holidays := []time.Time{
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // Independence Day
	}

	for _, holiday := range holidays {
		next := calc.NextBusinessDay(holiday)
		fmt.Printf("Next business day after %s (%s): %s (%s)\n",
			holiday.Format("Jan 2"), holiday.Weekday(),
			next.Format("Jan 2"), next.Weekday())
	}

	// 5. Previous Business Day
	fmt.Println("\n5. Finding Previous Business Day")
	for _, holiday := range holidays {
		prev := calc.PreviousBusinessDay(holiday)
		fmt.Printf("Previous business day before %s (%s): %s (%s)\n",
			holiday.Format("Jan 2"), holiday.Weekday(),
			prev.Format("Jan 2"), prev.Weekday())
	}

	// 6. Business Day Check
	fmt.Println("\n6. Checking Business Days")
	checkDates := []time.Time{
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas (Holiday)
		time.Date(2024, 12, 28, 0, 0, 0, 0, time.UTC), // Saturday
		time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC), // Tuesday
	}

	for _, date := range checkDates {
		isBusinessDay := calc.IsBusinessDay(date)
		fmt.Printf("%s (%s): %v\n", date.Format("Jan 2"), date.Weekday(),
			map[bool]string{true: "Business day", false: "Non-business day"}[isBusinessDay])
	}

	// 7. Business Quarter Analysis
	fmt.Println("\n7. Business Days in Q4 2024")
	q4Start := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	q4End := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	q4Days := calc.BusinessDaysBetween(q4Start, q4End)
	fmt.Printf("Total business days in Q4 2024: %d\n", q4Days)

	// List all holidays in Q4
	fmt.Println("Q4 Holidays:")
	q4Holidays := us.HolidaysForDateRange(q4Start, q4End)
	for date, holiday := range q4Holidays {
		fmt.Printf("- %s: %s\n", date.Format("Jan 2"), holiday.Name)
	}

	fmt.Println("\nThis demonstrates GoHoliday's business day calculation features!")
}
