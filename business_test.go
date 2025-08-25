package goholidays

import (
	"testing"
	"time"
)

func TestBusinessDayCalculator(t *testing.T) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)

	// Test a regular Tuesday (should be business day)
	tuesday := time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)
	if !calc.IsBusinessDay(tuesday) {
		t.Error("Tuesday should be a business day")
	}

	// Test Saturday (should not be business day)
	saturday := time.Date(2024, 3, 9, 0, 0, 0, 0, time.UTC)
	if calc.IsBusinessDay(saturday) {
		t.Error("Saturday should not be a business day")
	}

	// Test Independence Day (holiday - should not be business day)
	independenceDay := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	if calc.IsBusinessDay(independenceDay) {
		t.Error("Independence Day should not be a business day")
	}
}

func TestNextBusinessDay(t *testing.T) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)

	// Test from Friday to next Monday
	friday := time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC)
	nextBusiness := calc.NextBusinessDay(friday)
	expected := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC) // Monday

	if !nextBusiness.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, nextBusiness)
	}
}

func TestAddBusinessDays(t *testing.T) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)

	// Add 5 business days from Monday
	monday := time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	result := calc.AddBusinessDays(monday, 5)
	expected := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC) // Next Monday

	if !result.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestBusinessDaysBetween(t *testing.T) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)

	// Count business days in a week (Monday to Friday)
	monday := time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC)
	friday := time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC)

	count := calc.BusinessDaysBetween(monday, friday)
	expected := 4 // Mon, Tue, Wed, Thu (not including Friday)

	if count != expected {
		t.Errorf("Expected %d business days, got %d", expected, count)
	}
}

func TestHolidayAwareScheduler(t *testing.T) {
	us := NewCountry("US")
	scheduler := NewHolidayAwareScheduler(us)

	// Schedule 3 events starting from July 3rd (day before Independence Day)
	start := time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC)
	schedule := scheduler.ScheduleRecurring(start, 24*time.Hour, 3)

	if len(schedule) != 3 {
		t.Errorf("Expected 3 scheduled events, got %d", len(schedule))
	}

	// First event should be July 3rd (business day)
	if !schedule[0].Equal(start) {
		t.Errorf("First event should be July 3rd, got %v", schedule[0])
	}

	// Second event should skip July 4th (holiday) and go to July 5th
	expected := time.Date(2024, 7, 5, 0, 0, 0, 0, time.UTC)
	if !schedule[1].Equal(expected) {
		t.Errorf("Second event should be July 5th, got %v", schedule[1])
	}
}

func TestHolidayCalendar(t *testing.T) {
	us := NewCountry("US")
	calendar := NewHolidayCalendar(us)

	// Generate calendar for July 2024
	entries := calendar.GenerateMonth(2024, 7)

	if len(entries) != 31 {
		t.Errorf("July should have 31 days, got %d", len(entries))
	}

	// Find July 4th in the entries
	var july4th *CalendarEntry
	for i, entry := range entries {
		if entry.Date.Day() == 4 {
			july4th = &entries[i]
			break
		}
	}

	if july4th == nil {
		t.Fatal("Could not find July 4th in calendar entries")
	}

	if !july4th.IsHoliday {
		t.Error("July 4th should be marked as a holiday")
	}

	if july4th.Holiday == nil || july4th.Holiday.Name != "Independence Day" {
		t.Error("July 4th should be Independence Day")
	}

	// Check that July 4th is not a business day
	if july4th.IsBusinessDay {
		t.Error("July 4th should not be a business day")
	}
}

func TestCustomWeekends(t *testing.T) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)

	// Set custom weekends (Friday and Saturday for some Middle Eastern countries)
	calc.SetWeekends([]time.Weekday{time.Friday, time.Saturday})

	// Test Friday (should not be business day with custom weekends)
	friday := time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC)
	if calc.IsBusinessDay(friday) {
		t.Error("Friday should not be a business day with custom weekends")
	}

	// Test Sunday (should be business day with custom weekends)
	sunday := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	if !calc.IsBusinessDay(sunday) {
		t.Error("Sunday should be a business day with custom weekends")
	}
}

func BenchmarkBusinessDayCalculation(b *testing.B) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.IsBusinessDay(date)
	}
}

func BenchmarkAddBusinessDays(b *testing.B) {
	us := NewCountry("US")
	calc := NewBusinessDayCalculator(us)
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.AddBusinessDays(start, 10)
	}
}
