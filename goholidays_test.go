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

	if holiday.Category != "federal" {
		t.Errorf("Expected category 'federal', got '%s'", holiday.Category)
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

// TestInvalidCountry tests handling of unsupported country codes
func TestInvalidCountry(t *testing.T) {
	// Test with invalid country code
	invalid := NewCountry("ZZ")

	// Should not panic, but should have empty results
	holidays := invalid.HolidaysForYear(2024)
	if len(holidays) != 0 {
		t.Error("Unsupported country should have no holidays")
	}

	_, isHoliday := invalid.IsHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	if isHoliday {
		t.Error("Unsupported country should have no holidays")
	}
}

// TestCountryOptionsDetailed tests various country option configurations
func TestCountryOptionsDetailed(t *testing.T) {
	// Test US with California subdivision
	options := CountryOptions{
		Subdivisions: []string{"CA"},
		Language:     "es",
	}
	usCA := NewCountry("US", options)

	if len(usCA.GetSubdivisions()) != 1 || usCA.GetSubdivisions()[0] != "CA" {
		t.Error("Expected CA subdivision")
	}

	if usCA.GetLanguage() != "es" {
		t.Errorf("Expected Spanish language, got %s", usCA.GetLanguage())
	}

	// Test with empty options
	empty := NewCountry("US", CountryOptions{})
	if len(empty.GetSubdivisions()) != 0 {
		t.Error("Empty options should have no subdivisions")
	}

	if empty.GetLanguage() != "en" {
		t.Error("Default language should be English")
	}
}

// TestDateRangeEdgeCases tests edge cases for date ranges
func TestDateRangeEdgeCases(t *testing.T) {
	us := NewCountry("US")

	// Test same start and end date
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	holidays := us.HolidaysForDateRange(date, date)

	if len(holidays) != 1 {
		t.Errorf("Expected 1 holiday for single day range, got %d", len(holidays))
	}

	// Test reversed date range (end before start)
	start := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holidays = us.HolidaysForDateRange(start, end)

	if len(holidays) != 0 {
		t.Error("Reversed date range should return no holidays")
	}
}

// TestGettersAndSetters tests all getter methods
func TestGettersAndSetters(t *testing.T) {
	options := CountryOptions{
		Subdivisions: []string{"CA", "NY"},
		Categories:   []HolidayCategory{CategoryPublic, CategoryBank},
		Language:     "es",
	}
	us := NewCountry("US", options)

	// Test GetCountryCode
	if us.GetCountryCode() != "US" {
		t.Errorf("Expected country code 'US', got '%s'", us.GetCountryCode())
	}

	// Test GetSubdivisions
	subdivisions := us.GetSubdivisions()
	if len(subdivisions) != 2 {
		t.Errorf("Expected 2 subdivisions, got %d", len(subdivisions))
	}

	// Test GetCategories
	categories := us.GetCategories()
	if len(categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(categories))
	}

	// Test GetLanguage
	if us.GetLanguage() != "es" {
		t.Errorf("Expected language 'es', got '%s'", us.GetLanguage())
	}
}

// TestHolidayStructProperties tests Holiday struct functionality
func TestHolidayStructProperties(t *testing.T) {
	us := NewCountry("US")

	independenceDay := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday := us.IsHoliday(independenceDay)

	if !isHoliday {
		t.Fatal("Independence Day should be a holiday")
	}

	// Test holiday properties
	if holiday.Name == "" {
		t.Error("Holiday name should not be empty")
	}

	if holiday.Date != independenceDay {
		t.Error("Holiday date should match requested date")
	}

	if holiday.Category == "" {
		t.Error("Holiday category should not be empty")
	}

	if len(holiday.Languages) == 0 {
		t.Error("Holiday should have language translations")
	}

	// Check specific language keys
	if _, exists := holiday.Languages["en"]; !exists {
		t.Error("Holiday should have English translation")
	}
}

// TestMultipleCountries tests all supported countries
func TestMultipleCountries(t *testing.T) {
	supportedCountries := []string{"US", "GB", "CA", "AU", "NZ", "JP", "IN", "FR", "DE", "BR", "MX", "IT", "ES", "NL", "KR"}

	for _, code := range supportedCountries {
		t.Run(code, func(t *testing.T) {
			country := NewCountry(code)
			holidays := country.HolidaysForYear(2024)

			if len(holidays) == 0 {
				t.Errorf("Country %s should have holidays for 2024", code)
			}

			// Test country code retrieval
			if country.GetCountryCode() != code {
				t.Errorf("Expected country code %s, got %s", code, country.GetCountryCode())
			}
		})
	}
}

// TestHolidayLanguageSupport tests multi-language support
func TestHolidayLanguageSupport(t *testing.T) {
	// Test Spanish language for US
	options := CountryOptions{Language: "es"}
	us := NewCountry("US", options)

	newYearsDay := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, isHoliday := us.IsHoliday(newYearsDay)

	if !isHoliday {
		t.Fatal("New Year's Day should be a holiday")
	}

	if holiday.Languages["es"] == "" {
		t.Error("Holiday should have Spanish translation")
	}

	if holiday.Languages["en"] == "" {
		t.Error("Holiday should have English translation")
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
