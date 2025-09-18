package goholidays

import (
	"context"
	"errors"
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

// ============================================================================
// Error Handling and Enhanced API Tests
// ============================================================================

func TestHolidayError(t *testing.T) {
	t.Run("Basic Error", func(t *testing.T) {
		err := NewHolidayError(ErrInvalidCountry, "invalid country")

		if err.Code != ErrInvalidCountry {
			t.Errorf("Expected error code %d, got %d", ErrInvalidCountry, err.Code)
		}

		if err.Error() != "invalid country" {
			t.Errorf("Expected error message 'invalid country', got '%s'", err.Error())
		}
	})

	t.Run("Error With Cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := NewHolidayErrorWithCause(ErrDataLoadFailed, "failed to load", cause)

		if err.Unwrap() != cause {
			t.Error("Expected unwrapped error to match cause")
		}

		expected := "failed to load: underlying error"
		if err.Error() != expected {
			t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("Error Is Method", func(t *testing.T) {
		err1 := NewHolidayError(ErrInvalidCountry, "test")
		err2 := NewHolidayError(ErrInvalidCountry, "different message")
		err3 := NewHolidayError(ErrInvalidYear, "test")

		if !errors.Is(err1, err2) {
			t.Error("Expected errors with same code to match")
		}

		if errors.Is(err1, err3) {
			t.Error("Expected errors with different codes not to match")
		}
	})
}

func TestValidation(t *testing.T) {
	t.Run("ValidateCountryCode", func(t *testing.T) {
		// Valid countries
		validCodes := []string{"US", "GB", "DE", "UA"}
		for _, code := range validCodes {
			if err := ValidateCountryCode(code); err != nil {
				t.Errorf("Expected %s to be valid, got error: %v", code, err)
			}
		}

		// Invalid countries
		invalidCodes := []string{"XX", "", "us"}
		for _, code := range invalidCodes {
			if err := ValidateCountryCode(code); err == nil {
				t.Errorf("Expected %s to be invalid", code)
			}
		}
	})

	t.Run("ValidateYear", func(t *testing.T) {
		// Valid years
		validYears := []int{1900, 2024, 2200}
		for _, year := range validYears {
			if err := ValidateYear(year); err != nil {
				t.Errorf("Expected year %d to be valid, got error: %v", year, err)
			}
		}

		// Invalid years
		invalidYears := []int{1800, 2300, -1}
		for _, year := range invalidYears {
			if err := ValidateYear(year); err == nil {
				t.Errorf("Expected year %d to be invalid", year)
			}
		}
	})
}

func TestEnhancedAPI(t *testing.T) {
	t.Run("NewCountryWithError", func(t *testing.T) {
		// Valid country
		country, err := NewCountryWithError("US")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if country.GetCountryCode() != "US" {
			t.Errorf("Expected US, got %s", country.GetCountryCode())
		}

		// Invalid country
		country, err = NewCountryWithError("XX")
		if err == nil {
			t.Fatal("Expected error for invalid country")
		}
		if country != nil {
			t.Error("Expected nil country for invalid code")
		}
	})

	t.Run("IsHolidayWithError", func(t *testing.T) {
		country := NewCountry("US")

		// Valid date
		date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		holiday, isHoliday, err := country.IsHolidayWithError(date)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !isHoliday {
			t.Error("Expected July 4th to be a holiday")
		}
		if holiday == nil {
			t.Error("Expected holiday object")
		}

		// Invalid year
		date = time.Date(1800, 7, 4, 0, 0, 0, 0, time.UTC)
		holiday, isHoliday, err = country.IsHolidayWithError(date)
		if err == nil {
			t.Fatal("Expected error for invalid year")
		}
		if isHoliday || holiday != nil {
			t.Error("Expected no holiday for invalid year")
		}
	})

	t.Run("Context Support", func(t *testing.T) {
		country := NewCountry("US")

		// Valid context
		ctx := context.Background()
		date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		_, isHoliday, err := country.IsHolidayWithContext(ctx, date)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !isHoliday {
			t.Error("Expected July 4th to be a holiday")
		}

		// Cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, _, err = country.IsHolidayWithContext(ctx, date)
		if err == nil {
			t.Fatal("Expected error for cancelled context")
		}
		if !IsContextCancelled(err) {
			t.Error("Expected context cancellation error")
		}
	})

	t.Run("Backward Compatibility", func(t *testing.T) {
		// Original API still works exactly the same
		country := NewCountry("US")

		date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		holiday, isHoliday := country.IsHoliday(date)
		if !isHoliday {
			t.Error("Original API should still work")
		}
		if holiday == nil {
			t.Error("Expected holiday object")
		}

		holidays := country.HolidaysForYear(2024)
		if len(holidays) == 0 {
			t.Error("Original API should still work")
		}

		// Enhanced API provides additional safety
		holidayWithErr, isHolidayWithErr, err := country.IsHolidayWithError(date)
		if err != nil {
			t.Fatalf("Enhanced API failed: %v", err)
		}

		// Results should be identical
		if isHoliday != isHolidayWithErr {
			t.Error("Results should match between APIs")
		}
		if holiday.Name != holidayWithErr.Name {
			t.Error("Holiday names should match")
		}
	})
}

func TestUtilityFunctions(t *testing.T) {
	t.Run("IsValidCountry", func(t *testing.T) {
		if !IsValidCountry("US") {
			t.Error("Expected US to be valid")
		}

		if IsValidCountry("XX") {
			t.Error("Expected XX to be invalid")
		}
	})

	t.Run("GetSupportedCountries", func(t *testing.T) {
		countries := GetSupportedCountries()

		if len(countries) != len(SupportedCountries) {
			t.Errorf("Expected %d countries, got %d", len(SupportedCountries), len(countries))
		}

		// Check that all returned countries are in the map
		for _, country := range countries {
			if !SupportedCountries[country] {
				t.Errorf("Country %s returned but not in SupportedCountries map", country)
			}
		}
	})

	t.Run("IsContextCancelled", func(t *testing.T) {
		if !IsContextCancelled(context.Canceled) {
			t.Error("Expected context.Canceled to be detected")
		}

		if !IsContextCancelled(context.DeadlineExceeded) {
			t.Error("Expected context.DeadlineExceeded to be detected")
		}

		err := NewHolidayError(ErrCancelled, "cancelled")
		if !IsContextCancelled(err) {
			t.Error("Expected HolidayError with ErrCancelled to be detected")
		}

		err = NewHolidayError(ErrInvalidCountry, "invalid")
		if IsContextCancelled(err) {
			t.Error("Expected non-cancelled error not to be detected")
		}
	})
}
