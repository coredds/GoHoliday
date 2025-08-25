package main

import (
	"fmt"
	"time"

	"github.com/coredds/GoHoliday"
)

func main() {
	fmt.Println("GoHolidays Business Day Features Demo")
	fmt.Println("=====================================")

	// Create a US holiday provider
	us := goholidays.NewCountry("US")
	
	// Business Day Calculator Demo
	fmt.Println("\n1. Business Day Calculator:")
	calc := goholidays.NewBusinessDayCalculator(us)
	
	testDates := []string{
		"2024-07-03", // Wednesday before July 4th
		"2024-07-04", // Independence Day (holiday)
		"2024-07-05", // Friday after July 4th
		"2024-07-06", // Saturday
		"2024-07-07", // Sunday
		"2024-07-08", // Monday
	}
	
	for _, dateStr := range testDates {
		date, _ := time.Parse("2006-01-02", dateStr)
		isBusinessDay := calc.IsBusinessDay(date)
		dayType := "Business Day"
		if !isBusinessDay {
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				dayType = "Weekend"
			} else {
				dayType = "Holiday"
			}
		}
		fmt.Printf("  %s (%s): %s\n", dateStr, date.Weekday().String(), dayType)
	}
	
	// Next Business Day Demo
	fmt.Println("\n2. Next Business Day Calculation:")
	
	// Starting from July 3rd (Wednesday)
	startDate := time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC)
	fmt.Printf("  Starting from: %s (%s)\n", startDate.Format("2006-01-02"), startDate.Weekday())
	
	for i := 1; i <= 5; i++ {
		nextBusiness := calc.AddBusinessDays(startDate, i)
		fmt.Printf("  +%d business days: %s (%s)\n", i, nextBusiness.Format("2006-01-02"), nextBusiness.Weekday())
	}
	
	// Business Days Between Demo
	fmt.Println("\n3. Counting Business Days:")
	start := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)
	businessDays := calc.BusinessDaysBetween(start, end)
	fmt.Printf("  Business days in July 2024: %d\n", businessDays)
	
	// Holiday-Aware Scheduler Demo
	fmt.Println("\n4. Holiday-Aware Scheduling:")
	scheduler := goholidays.NewHolidayAwareScheduler(us)
	
	// Schedule weekly meetings starting July 1st
	meetingStart := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	meetings := scheduler.ScheduleRecurring(meetingStart, 7*24*time.Hour, 4)
	
	fmt.Printf("  Weekly meetings starting %s:\n", meetingStart.Format("2006-01-02"))
	for i, meeting := range meetings {
		fmt.Printf("    Meeting %d: %s (%s)\n", i+1, meeting.Format("2006-01-02"), meeting.Weekday())
	}
	
	// End-of-month scheduling
	fmt.Println("\n5. End-of-Month Business Day Scheduling:")
	eomStart := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	eomSchedule := scheduler.ScheduleMonthlyEndOfMonth(eomStart, 6)
	
	fmt.Println("  Last business day of each month:")
	for i, date := range eomSchedule {
		fmt.Printf("    Month %d: %s (%s)\n", i+1, date.Format("2006-01-02"), date.Weekday())
	}
	
	// Holiday Calendar Demo
	fmt.Println("\n6. Holiday Calendar for July 2024:")
	calendar := goholidays.NewHolidayCalendar(us)
	calendar.PrintMonth(2024, 7)
	
	// Custom Weekend Configuration
	fmt.Println("\n7. Custom Weekend Configuration (Middle East style):")
	customCalc := goholidays.NewBusinessDayCalculator(us)
	customCalc.SetWeekends([]time.Weekday{time.Friday, time.Saturday})
	
	fmt.Println("  Business days with Friday-Saturday weekends:")
	customTestDates := []string{
		"2024-07-04", // Thursday (Independence Day - still holiday)
		"2024-07-05", // Friday (weekend)
		"2024-07-06", // Saturday (weekend)
		"2024-07-07", // Sunday (business day!)
		"2024-07-08", // Monday (business day)
	}
	
	for _, dateStr := range customTestDates {
		date, _ := time.Parse("2006-01-02", dateStr)
		isBusinessDay := customCalc.IsBusinessDay(date)
		status := "Business Day"
		if !isBusinessDay {
			if date.Weekday() == time.Friday || date.Weekday() == time.Saturday {
				status = "Weekend"
			} else {
				status = "Holiday"
			}
		}
		fmt.Printf("    %s (%s): %s\n", dateStr, date.Weekday().String(), status)
	}
	
	// Performance demonstration
	fmt.Println("\n8. Performance Test:")
	businessPerformanceTest(calc)
	
	fmt.Println("\nDemo completed!")
}

func businessPerformanceTest(calc *goholidays.BusinessDayCalculator) {
	start := time.Now()
	
	// Test 10,000 business day calculations
	testDate := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 10000; i++ {
		calc.IsBusinessDay(testDate.AddDate(0, 0, i%365))
	}
	
	duration := time.Since(start)
	fmt.Printf("  10,000 business day checks took: %v\n", duration)
	fmt.Printf("  Average per check: %v\n", duration/10000)
	
	// Test business day additions
	start = time.Now()
	baseDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 1000; i++ {
		calc.AddBusinessDays(baseDate, 10)
	}
	
	duration = time.Since(start)
	fmt.Printf("  1,000 'add 10 business days' operations took: %v\n", duration)
	fmt.Printf("  Average per operation: %v\n", duration/1000)
}
