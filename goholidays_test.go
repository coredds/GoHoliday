package goholidays

import (
	"testing"
	"time"
)

func TestNewCountry(t *testing.T) {
	// Test basic country creation
	us := NewCountry("US")
	if us.GetCountryCode() != "US" {
		t.Errorf("Expected country code 'US', got '%s'", us.GetCountryCode())
	}

	// Test with options
	options := CountryOptions{
		Subdivisions: []string{"CA", "NY"},
		Categories:   []HolidayCategory{CategoryPublic, CategoryBank},
		Language:     "es",
	}
	usWithOptions := NewCountry("US", options)

	if len(usWithOptions.GetSubdivisions()) != 2 {
		t.Errorf("Expected 2 subdivisions, got %d", len(usWithOptions.GetSubdivisions()))
	}

	if usWithOptions.GetLanguage() != "es" {
		t.Errorf("Expected language 'es', got '%s'", usWithOptions.GetLanguage())
	}
}

func TestIsHoliday(t *testing.T) {
	us := NewCountry("US")

	// Test Independence Day
	independenceDay := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday := us.IsHoliday(independenceDay)

	if !isHoliday {
		t.Error("July 4th should be a holiday in the US")
	}

	if holiday.Name != "Independence Day" {
		t.Errorf("Expected 'Independence Day', got '%s'", holiday.Name)
	}

	// Test non-holiday
	randomDay := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	_, isHoliday = us.IsHoliday(randomDay)

	if isHoliday {
		t.Error("March 15th should not be a holiday in the US")
	}
}

func TestHolidaysForYear(t *testing.T) {
	us := NewCountry("US")
	holidays := us.HolidaysForYear(2024)

	if len(holidays) == 0 {
		t.Error("Should have holidays for 2024")
	}

	// Check for specific holidays
	newYearsDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[newYearsDay]; !exists {
		t.Error("New Year's Day should be in 2024 holidays")
	}

	christmasDay := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[christmasDay]; !exists {
		t.Error("Christmas Day should be in 2024 holidays")
	}
}

func TestHolidaysForDateRange(t *testing.T) {
	us := NewCountry("US")

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 7, 31, 0, 0, 0, 0, time.UTC)

	holidays := us.HolidaysForDateRange(start, end)

	// Should include New Year's Day and Independence Day
	newYearsDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	independenceDay := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	christmasDay := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)

	if _, exists := holidays[newYearsDay]; !exists {
		t.Error("New Year's Day should be in date range")
	}

	if _, exists := holidays[independenceDay]; !exists {
		t.Error("Independence Day should be in date range")
	}

	if _, exists := holidays[christmasDay]; exists {
		t.Error("Christmas Day should not be in date range")
	}
}

func TestHolidayCategories(t *testing.T) {
	us := NewCountry("US")

	independenceDay := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday := us.IsHoliday(independenceDay)

	if !isHoliday {
		t.Fatal("Independence Day should be a holiday")
	}

	if holiday.Category != CategoryPublic {
		t.Errorf("Expected category '%s', got '%s'", CategoryPublic, holiday.Category)
	}
}

func TestMultiLanguageSupport(t *testing.T) {
	us := NewCountry("US")

	newYearsDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday := us.IsHoliday(newYearsDay)

	if !isHoliday {
		t.Fatal("New Year's Day should be a holiday")
	}

	if holiday.Languages["en"] != "New Year's Day" {
		t.Errorf("Expected English name 'New Year's Day', got '%s'", holiday.Languages["en"])
	}

	if holiday.Languages["es"] != "Año Nuevo" {
		t.Errorf("Expected Spanish name 'Año Nuevo', got '%s'", holiday.Languages["es"])
	}
}

func BenchmarkIsHoliday(b *testing.B) {
	us := NewCountry("US")
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		us.IsHoliday(date)
	}
}

func BenchmarkHolidaysForYear(b *testing.B) {
	us := NewCountry("US")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		us.HolidaysForYear(2024)
	}
}
