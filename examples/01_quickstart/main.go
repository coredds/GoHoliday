package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHoliday Quickstart Example")
	fmt.Println("==========================")

	// Create a holiday provider for the United States
	us := goholidays.NewCountry("US")

	// 1. Basic holiday check
	fmt.Println("\n1. Is today a holiday?")
	today := time.Now()
	if holiday, isHoliday := us.IsHoliday(today); isHoliday {
		fmt.Printf("Yes! Today is %s\n", holiday.Name)
	} else {
		fmt.Printf("No, today is not a holiday\n")
	}

	// 2. Get upcoming holidays
	fmt.Println("\n2. Next 3 holidays:")
	holidays := us.HolidaysForYear(today.Year())
	count := 0
	for date, holiday := range holidays {
		if date.After(today) {
			fmt.Printf("- %s: %s\n", date.Format("Jan 2"), holiday.Name)
			count++
			if count >= 3 {
				break
			}
		}
	}

	// 3. Multi-language support
	fmt.Println("\n3. Holiday names in different languages:")
	christmas := time.Date(today.Year(), 12, 25, 0, 0, 0, 0, time.UTC)
	if holiday, isHoliday := us.IsHoliday(christmas); isHoliday {
		fmt.Printf("Christmas in:\n")
		fmt.Printf("- English: %s\n", holiday.Languages["en"])
		fmt.Printf("- Spanish: %s\n", holiday.Languages["es"])
	}

	// 4. Business day calculation
	fmt.Println("\n4. Business days:")
	calc := goholidays.NewBusinessDayCalculator(us)
	nextWeek := today.AddDate(0, 0, 7)
	businessDays := calc.BusinessDaysBetween(today, nextWeek)
	fmt.Printf("Business days in the next week: %d\n", businessDays)

	fmt.Println("\nThat's it! Check out the other examples for more features.")
}
