package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	// Test Japan holidays
	jp := goholidays.NewCountry("JP")

	// Test New Year's Day 2024
	holiday, exists := jp.IsHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	if exists {
		fmt.Printf("Found Japan holiday: %s (%s)\n", holiday.Name, holiday.Languages["ja"])
	} else {
		fmt.Println("New Year's Day not found")
	}

	// Test Coming of Age Day 2024 (should be January 8, 2024)
	holiday, exists = jp.IsHoliday(time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC))
	if exists {
		fmt.Printf("Found Japan holiday: %s (%s)\n", holiday.Name, holiday.Languages["ja"])
	} else {
		fmt.Println("Coming of Age Day not found")
	}

	// Test Emperor's Birthday 2024 (should be February 23, 2024)
	holiday, exists = jp.IsHoliday(time.Date(2024, 2, 23, 0, 0, 0, 0, time.UTC))
	if exists {
		fmt.Printf("Found Japan holiday: %s (%s)\n", holiday.Name, holiday.Languages["ja"])
	} else {
		fmt.Println("Emperor's Birthday not found")
	}

	// Test a non-holiday
	holiday, exists = jp.IsHoliday(time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC))
	if exists {
		fmt.Printf("Unexpected holiday found: %s\n", holiday.Name)
	} else {
		fmt.Println("June 15, 2024 correctly identified as non-holiday")
	}
}
