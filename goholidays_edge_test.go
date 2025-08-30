package goholidays

import (
	"testing"
	"time"
)

// TestHolidayEdgeCases tests various edge cases for holiday handling
func TestHolidayEdgeCases(t *testing.T) {
	// Test with all holiday categories
	t.Run("All Categories", func(t *testing.T) {
		categories := []HolidayCategory{
			CategoryPublic,
			CategoryBank,
			CategorySchool,
			CategoryGovernment,
			CategoryReligious,
			CategoryOptional,
			CategoryHalfDay,
			CategoryArmedForces,
			CategoryWorkday,
		}

		options := CountryOptions{
			Categories: categories,
		}
		us := NewCountry("US", options)

		if len(us.GetCategories()) != len(categories) {
			t.Errorf("Expected %d categories, got %d", len(categories), len(us.GetCategories()))
		}

		for _, category := range categories {
			found := false
			for _, c := range us.GetCategories() {
				if c == category {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Category %s not found in country categories", category)
			}
		}
	})

	// Test with multiple years
	t.Run("Multiple Years", func(t *testing.T) {
		us := NewCountry("US")
		years := []int{2020, 2021, 2022, 2023, 2024, 2025}

		for _, year := range years {
			holidays := us.HolidaysForYear(year)
			if len(holidays) == 0 {
				t.Errorf("Year %d should have holidays", year)
			}

			// Check specific holidays that should be consistent across years
			newYearsDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			if _, exists := holidays[newYearsDay]; !exists {
				t.Errorf("New Year's Day missing for %d", year)
			}

			independenceDay := time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC)
			if _, exists := holidays[independenceDay]; !exists {
				t.Errorf("Independence Day missing for %d", year)
			}

			christmasDay := time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC)
			if _, exists := holidays[christmasDay]; !exists {
				t.Errorf("Christmas Day missing for %d", year)
			}
		}
	})

	// Test with all supported languages
	t.Run("Supported Languages", func(t *testing.T) {
		languages := []string{"en", "es"}

		for _, lang := range languages {
			options := CountryOptions{Language: lang}
			us := NewCountry("US", options)

			if us.GetLanguage() != lang {
				t.Errorf("Expected language %s, got %s", lang, us.GetLanguage())
			}

			// Check translations exist
			holiday, exists := us.IsHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
			if exists {
				if holiday.Languages[lang] == "" {
					t.Errorf("Missing translation for language %s", lang)
				}
			}
		}
	})

	// Test with invalid dates
	t.Run("Invalid Dates", func(t *testing.T) {
		us := NewCountry("US")

		// Test with zero time - January 1, year 1 should be New Year's Day
		_, isHoliday := us.IsHoliday(time.Time{})
		if !isHoliday {
			t.Error("Zero time (January 1, year 1) should be New Year's Day holiday")
		}

		// Test with far future date - should still calculate holidays
		holidays := us.HolidaysForYear(9999)
		if len(holidays) == 0 {
			t.Error("Far future year should still have holidays")
		}

		// Test with far past date - should still calculate holidays
		holidays = us.HolidaysForYear(1)
		if len(holidays) == 0 {
			t.Error("Far past year should still have holidays")
		}
	})

	// Test with invalid date ranges
	t.Run("Invalid Date Ranges", func(t *testing.T) {
		us := NewCountry("US")

		// Test with zero times - same start and end should return one holiday if it exists
		holidays := us.HolidaysForDateRange(time.Time{}, time.Time{})
		if len(holidays) != 1 {
			t.Error("Zero time range (same start/end on New Year's Day) should return one holiday")
		}

		// Test with reversed range
		start := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		holidays = us.HolidaysForDateRange(start, end)
		if len(holidays) != 0 {
			t.Error("Reversed date range should return no holidays")
		}

		// Test with same start and end date
		date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		holidays = us.HolidaysForDateRange(date, date)
		if len(holidays) != 1 {
			t.Error("Single day range should return one holiday")
		}
	})

	// Test with all subdivisions
	t.Run("All Subdivisions", func(t *testing.T) {
		subdivisions := []string{
			"AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA",
			"HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD",
			"MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ",
			"NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC",
			"SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY",
		}

		options := CountryOptions{
			Subdivisions: subdivisions,
		}
		us := NewCountry("US", options)

		if len(us.GetSubdivisions()) != len(subdivisions) {
			t.Errorf("Expected %d subdivisions, got %d", len(subdivisions), len(us.GetSubdivisions()))
		}

		for _, subdivision := range subdivisions {
			found := false
			for _, s := range us.GetSubdivisions() {
				if s == subdivision {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Subdivision %s not found in country subdivisions", subdivision)
			}
		}
	})

	// Test concurrent access
	t.Run("Concurrent Access", func(t *testing.T) {
		us := NewCountry("US")
		done := make(chan bool)
		const workers = 10

		// Start multiple goroutines accessing holidays
		for i := 0; i < workers; i++ {
			go func() {
				defer func() { done <- true }()

				// Test different operations
				us.IsHoliday(time.Now())
				us.HolidaysForYear(2024)
				us.HolidaysForDateRange(
					time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
				)
			}()
		}

		// Wait for all goroutines
		for i := 0; i < workers; i++ {
			<-done
		}
	})

	// Test holiday properties
	t.Run("Holiday Properties", func(t *testing.T) {
		us := NewCountry("US")
		date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		holiday, exists := us.IsHoliday(date)

		if !exists {
			t.Fatal("Independence Day should exist")
		}

		// Test all holiday struct fields
		if holiday.Name == "" {
			t.Error("Holiday name should not be empty")
		}

		if holiday.Date.IsZero() {
			t.Error("Holiday date should not be zero")
		}

		if holiday.Category == "" {
			t.Error("Holiday category should not be empty")
		}

		if len(holiday.Languages) == 0 {
			t.Error("Holiday should have language translations")
		}

		if holiday.IsObserved {
			if holiday.Observed == nil {
				t.Error("Observed holiday should have Observed date set")
			}
		}
	})

	// Test with multiple options
	t.Run("Multiple Options", func(t *testing.T) {
		options := CountryOptions{
			Subdivisions: []string{"CA", "NY"},
			Categories:   []HolidayCategory{CategoryPublic, CategoryBank},
			Language:     "es",
			Years:        []int{2023, 2024, 2025},
		}
		us := NewCountry("US", options)

		// Verify all options were applied
		if len(us.GetSubdivisions()) != 2 {
			t.Error("Subdivisions not properly set")
		}

		if len(us.GetCategories()) != 2 {
			t.Error("Categories not properly set")
		}

		if us.GetLanguage() != "es" {
			t.Error("Language not properly set")
		}

		// Check that years were preloaded
		for _, year := range []int{2023, 2024, 2025} {
			holidays := us.HolidaysForYear(year)
			if len(holidays) == 0 {
				t.Errorf("Year %d should be preloaded", year)
			}
		}
	})
}

// TestHolidayCalculations tests various holiday calculation methods
func TestHolidayCalculations(t *testing.T) {
	us := NewCountry("US")

	// Test fixed date holidays
	t.Run("Fixed Date Holidays", func(t *testing.T) {
		fixedDates := []struct {
			name string
			date time.Time
		}{
			{"New Year's Day", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			{"Independence Day", time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)},
			{"Christmas Day", time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)},
		}

		for _, fd := range fixedDates {
			holiday, exists := us.IsHoliday(fd.date)
			if !exists {
				t.Errorf("%s should exist", fd.name)
			}
			if exists && holiday.Name != fd.name {
				t.Errorf("Expected %s, got %s", fd.name, holiday.Name)
			}
		}
	})

	// Test observed holidays
	t.Run("Observed Holidays", func(t *testing.T) {
		// July 4th, 2026 falls on Saturday
		date := time.Date(2026, 7, 4, 0, 0, 0, 0, time.UTC)
		holiday, exists := us.IsHoliday(date)

		if !exists {
			t.Fatal("Independence Day should exist")
		}

		if !holiday.IsObserved {
			t.Error("Saturday holiday should be observed")
		}

		if holiday.Observed == nil {
			t.Error("Observed date should be set")
		}

		// Check observed date is Friday
		expectedObserved := time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC)
		if holiday.Observed != nil && !holiday.Observed.Equal(expectedObserved) {
			t.Errorf("Expected observed date %v, got %v", expectedObserved, holiday.Observed)
		}
	})

	// Test holiday cache
	t.Run("Holiday Cache", func(t *testing.T) {
		// Access same year multiple times
		year := 2024
		firstAccess := us.HolidaysForYear(year)
		secondAccess := us.HolidaysForYear(year)

		if len(firstAccess) != len(secondAccess) {
			t.Error("Cached results should be consistent")
		}

		// Compare specific holidays
		for date, holiday := range firstAccess {
			if cached, exists := secondAccess[date]; !exists {
				t.Errorf("Holiday %s missing in cached result", holiday.Name)
			} else if cached.Name != holiday.Name {
				t.Errorf("Cached holiday name mismatch: expected %s, got %s", holiday.Name, cached.Name)
			}
		}
	})
}

// TestHolidayPerformance tests performance characteristics
func TestHolidayPerformance(t *testing.T) {
	// Test memory usage
	t.Run("Memory Usage", func(t *testing.T) {
		country := NewCountry("US")
		years := make([]int, 100)
		for i := range years {
			years[i] = 2000 + i
		}

		// Load many years of holidays
		for _, year := range years {
			holidays := country.HolidaysForYear(year)
			if len(holidays) == 0 {
				t.Errorf("Year %d should have holidays", year)
			}
		}
	})

	// Test date range performance
	t.Run("Date Range Performance", func(t *testing.T) {
		country := NewCountry("US")
		ranges := []struct {
			name  string
			start time.Time
			end   time.Time
		}{
			{
				"Short Range",
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			},
			{
				"Medium Range",
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
			},
			{
				"Long Range",
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
			},
			{
				"Multi-Year Range",
				time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
			},
		}

		for _, r := range ranges {
			holidays := country.HolidaysForDateRange(r.start, r.end)
			if len(holidays) == 0 {
				t.Errorf("%s should have holidays", r.name)
			}
		}
	})
}
