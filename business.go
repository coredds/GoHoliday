package goholidays

import (
	"fmt"
	"time"
)

// BusinessDayCalculator provides business day calculations with holiday awareness
type BusinessDayCalculator struct {
	country *Country
	weekends []time.Weekday
}

// NewBusinessDayCalculator creates a new business day calculator
func NewBusinessDayCalculator(country *Country) *BusinessDayCalculator {
	return &BusinessDayCalculator{
		country:  country,
		weekends: []time.Weekday{time.Saturday, time.Sunday}, // Default weekends
	}
}

// SetWeekends sets custom weekend days
func (bdc *BusinessDayCalculator) SetWeekends(weekends []time.Weekday) {
	bdc.weekends = weekends
}

// IsBusinessDay checks if a date is a business day (not weekend or holiday)
func (bdc *BusinessDayCalculator) IsBusinessDay(date time.Time) bool {
	// Check if it's a weekend
	for _, weekend := range bdc.weekends {
		if date.Weekday() == weekend {
			return false
		}
	}
	
	// Check if it's a holiday
	_, isHoliday := bdc.country.IsHoliday(date)
	return !isHoliday
}

// NextBusinessDay returns the next business day after the given date
func (bdc *BusinessDayCalculator) NextBusinessDay(date time.Time) time.Time {
	next := date.AddDate(0, 0, 1)
	for !bdc.IsBusinessDay(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

// PreviousBusinessDay returns the previous business day before the given date
func (bdc *BusinessDayCalculator) PreviousBusinessDay(date time.Time) time.Time {
	prev := date.AddDate(0, 0, -1)
	for !bdc.IsBusinessDay(prev) {
		prev = prev.AddDate(0, 0, -1)
	}
	return prev
}

// AddBusinessDays adds a specified number of business days to a date
func (bdc *BusinessDayCalculator) AddBusinessDays(date time.Time, days int) time.Time {
	if days == 0 {
		return date
	}
	
	current := date
	if days > 0 {
		for i := 0; i < days; i++ {
			current = bdc.NextBusinessDay(current)
		}
	} else {
		for i := 0; i < -days; i++ {
			current = bdc.PreviousBusinessDay(current)
		}
	}
	
	return current
}

// BusinessDaysBetween calculates the number of business days between two dates
func (bdc *BusinessDayCalculator) BusinessDaysBetween(start, end time.Time) int {
	if start.After(end) {
		return -bdc.BusinessDaysBetween(end, start)
	}
	
	count := 0
	current := start
	
	for current.Before(end) {
		if bdc.IsBusinessDay(current) {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}
	
	return count
}

// IsEndOfMonth checks if a date is the last business day of the month
func (bdc *BusinessDayCalculator) IsEndOfMonth(date time.Time) bool {
	if !bdc.IsBusinessDay(date) {
		return false
	}
	
	// Get the last day of the month
	nextMonth := time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, date.Location())
	lastDay := nextMonth.AddDate(0, 0, -1)
	
	// Find the last business day of the month
	for !bdc.IsBusinessDay(lastDay) {
		lastDay = lastDay.AddDate(0, 0, -1)
	}
	
	return date.Year() == lastDay.Year() && 
		   date.Month() == lastDay.Month() && 
		   date.Day() == lastDay.Day()
}

// HolidayAwareScheduler provides scheduling functionality with holiday awareness
type HolidayAwareScheduler struct {
	calculator *BusinessDayCalculator
}

// NewHolidayAwareScheduler creates a new scheduler
func NewHolidayAwareScheduler(country *Country) *HolidayAwareScheduler {
	return &HolidayAwareScheduler{
		calculator: NewBusinessDayCalculator(country),
	}
}

// ScheduleRecurring schedules recurring events avoiding holidays and weekends
func (has *HolidayAwareScheduler) ScheduleRecurring(start time.Time, frequency time.Duration, count int) []time.Time {
	var schedule []time.Time
	current := start
	
	for i := 0; i < count; i++ {
		// Ensure the scheduled date is a business day
		if !has.calculator.IsBusinessDay(current) {
			current = has.calculator.NextBusinessDay(current)
		}
		
		schedule = append(schedule, current)
		current = current.Add(frequency)
	}
	
	return schedule
}

// ScheduleMonthlyEndOfMonth schedules events for the last business day of each month
func (has *HolidayAwareScheduler) ScheduleMonthlyEndOfMonth(start time.Time, months int) []time.Time {
	var schedule []time.Time
	current := start
	
	for i := 0; i < months; i++ {
		// Get the last day of the current month
		nextMonth := time.Date(current.Year(), current.Month()+1, 1, 0, 0, 0, 0, current.Location())
		lastDay := nextMonth.AddDate(0, 0, -1)
		
		// Find the last business day
		for !has.calculator.IsBusinessDay(lastDay) {
			lastDay = lastDay.AddDate(0, 0, -1)
		}
		
		schedule = append(schedule, lastDay)
		
		// Move to next month
		current = time.Date(current.Year(), current.Month()+1, 1, 0, 0, 0, 0, current.Location())
	}
	
	return schedule
}

// HolidayCalendar provides a calendar view with holiday information
type HolidayCalendar struct {
	country *Country
}

// NewHolidayCalendar creates a new holiday calendar
func NewHolidayCalendar(country *Country) *HolidayCalendar {
	return &HolidayCalendar{country: country}
}

// CalendarEntry represents a single day in the calendar
type CalendarEntry struct {
	Date        time.Time `json:"date"`
	IsHoliday   bool      `json:"is_holiday"`
	IsWeekend   bool      `json:"is_weekend"`
	Holiday     *Holiday  `json:"holiday,omitempty"`
	IsBusinessDay bool    `json:"is_business_day"`
}

// GenerateMonth generates a calendar for a specific month
func (hc *HolidayCalendar) GenerateMonth(year int, month time.Month) []CalendarEntry {
	var entries []CalendarEntry
	
	// Get first day of month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	
	// Get first day of next month to determine when to stop
	nextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	
	current := firstDay
	for current.Before(nextMonth) {
		holiday, isHoliday := hc.country.IsHoliday(current)
		isWeekend := current.Weekday() == time.Saturday || current.Weekday() == time.Sunday
		
		entry := CalendarEntry{
			Date:          current,
			IsHoliday:     isHoliday,
			IsWeekend:     isWeekend,
			IsBusinessDay: !isHoliday && !isWeekend,
		}
		
		if isHoliday {
			entry.Holiday = holiday
		}
		
		entries = append(entries, entry)
		current = current.AddDate(0, 0, 1)
	}
	
	return entries
}

// PrintMonth prints a formatted calendar for a month
func (hc *HolidayCalendar) PrintMonth(year int, month time.Month) {
	entries := hc.GenerateMonth(year, month)
	
	fmt.Printf("\n%s %d\n", month.String(), year)
	fmt.Println("Su Mo Tu We Th Fr Sa")
	
	// Get first day of month to calculate starting position
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	startPos := int(firstDay.Weekday())
	
	// Print leading spaces
	for i := 0; i < startPos; i++ {
		fmt.Print("   ")
	}
	
	// Print days
	for _, entry := range entries {
		dayStr := fmt.Sprintf("%2d", entry.Date.Day())
		
		if entry.IsHoliday {
			fmt.Printf("*%s", dayStr[1:]) // Mark holidays with *
		} else {
			fmt.Print(dayStr)
		}
		
		// New line after Saturday
		if entry.Date.Weekday() == time.Saturday {
			fmt.Println()
		} else {
			fmt.Print(" ")
		}
	}
	
	fmt.Println()
	fmt.Println("* = Holiday")
}
