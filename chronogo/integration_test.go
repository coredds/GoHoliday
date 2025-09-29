package chronogo

import (
	"testing"
	"time"
)

func TestFastCountryChecker_IsHoliday(t *testing.T) {
	checker := Checker("US")

	// Test known US holidays for 2024
	testCases := []struct {
		date        time.Time
		isHoliday   bool
		description string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), true, "New Year's Day"},
		{time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC), true, "Independence Day"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), true, "Christmas Day"},
		{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), false, "Regular day"},
		{time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false, "Regular day"},
	}

	for _, tc := range testCases {
		result := checker.IsHoliday(tc.date)
		if result != tc.isHoliday {
			t.Errorf("%s: expected %v, got %v", tc.description, tc.isHoliday, result)
		}
	}
}

func TestFastCountryChecker_GetHolidayName(t *testing.T) {
	checker := Checker("US")

	// Test getting holiday names
	testCases := []struct {
		date         time.Time
		expectedName string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day"},
		{time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC), "Independence Day"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day"},
		{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), ""}, // Not a holiday
	}

	for _, tc := range testCases {
		result := checker.GetHolidayName(tc.date)
		if result != tc.expectedName {
			t.Errorf("Date %s: expected '%s', got '%s'",
				tc.date.Format("2006-01-02"), tc.expectedName, result)
		}
	}
}

func TestFastCountryChecker_AreHolidays(t *testing.T) {
	checker := Checker("US")

	dates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year's Day
		time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),   // Regular day
		time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC),   // Independence Day
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas Day
	}

	expected := []bool{true, false, true, true}
	results := checker.AreHolidays(dates)

	if len(results) != len(expected) {
		t.Fatalf("Expected %d results, got %d", len(expected), len(results))
	}

	for i, expectedResult := range expected {
		if results[i] != expectedResult {
			t.Errorf("Date %s: expected %v, got %v",
				dates[i].Format("2006-01-02"), expectedResult, results[i])
		}
	}
}

func TestFastCountryChecker_GetHolidaysInRange(t *testing.T) {
	checker := Checker("US")

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	holidays := checker.GetHolidaysInRange(start, end)

	// Should include New Year's Day and MLK Day in January 2024
	if len(holidays) < 1 {
		t.Errorf("Expected at least 1 holiday in January 2024, got %d", len(holidays))
	}

	// Check if New Year's Day is included
	newYears := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holidayName, exists := holidays[newYears]; !exists {
		t.Error("New Year's Day should be included in January 2024 holidays")
	} else if holidayName != "New Year's Day" {
		t.Errorf("Expected 'New Year's Day', got '%s'", holidayName)
	}
}

func TestFastCountryChecker_CountHolidaysInRange(t *testing.T) {
	checker := Checker("US")

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	count := checker.CountHolidaysInRange(start, end)

	// US should have approximately 11 federal holidays in 2024
	if count < 10 || count > 15 {
		t.Errorf("Expected 10-15 holidays in 2024, got %d", count)
	}
}

func TestFastCountryChecker_MultipleCountries(t *testing.T) {
	usChecker := Checker("US")
	caChecker := Checker("CA")
	gbChecker := Checker("GB")
	clChecker := Checker("CL")
	ieChecker := Checker("IE")
	ilChecker := Checker("IL")

	// Test a date that might be a holiday in some countries but not others
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)

	usResult := usChecker.IsHoliday(date) // Should be true (Independence Day)
	caResult := caChecker.IsHoliday(date) // Should be false
	gbResult := gbChecker.IsHoliday(date) // Should be false
	clResult := clChecker.IsHoliday(date) // Should be false
	ieResult := ieChecker.IsHoliday(date) // Should be false
	ilResult := ilChecker.IsHoliday(date) // Should be false

	if !usResult {
		t.Error("July 4th should be a holiday in the US")
	}
	if caResult {
		t.Error("July 4th should not be a holiday in Canada")
	}
	if gbResult {
		t.Error("July 4th should not be a holiday in Great Britain")
	}
	if clResult {
		t.Error("July 4th should not be a holiday in Chile")
	}
	if ieResult {
		t.Error("July 4th should not be a holiday in Ireland")
	}
	if ilResult {
		t.Error("July 4th should not be a holiday in Israel")
	}
}

func TestFastCountryChecker_ClearCache(t *testing.T) {
	checker := Checker("US")

	// Load some data
	checker.IsHoliday(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	// Verify cache has data
	checker.mutex.RLock()
	cacheSize := len(checker.yearCache)
	checker.mutex.RUnlock()

	if cacheSize == 0 {
		t.Error("Cache should have data after holiday check")
	}

	// Clear cache
	checker.ClearCache()

	// Verify cache is empty
	checker.mutex.RLock()
	cacheSize = len(checker.yearCache)
	checker.mutex.RUnlock()

	if cacheSize != 0 {
		t.Error("Cache should be empty after ClearCache()")
	}
}

func TestFastCountryChecker_GetCountryCode(t *testing.T) {
	testCases := []string{
		"AR", "AT", "AU", "BE", "BR", "CA", "CH", "CL", "CN", "DE",
		"ES", "FI", "FR", "GB", "ID", "IE", "IL", "IN", "IT", "JP",
		"KR", "MX", "NL", "NO", "NZ", "PL", "PT", "RU", "SE", "SG",
		"TH", "TR", "UA", "US",
	}

	for _, countryCode := range testCases {
		checker := Checker(countryCode)
		if checker.GetCountryCode() != countryCode {
			t.Errorf("Expected country code %s, got %s", countryCode, checker.GetCountryCode())
		}
	}
}

func TestFastCountryChecker_UnsupportedCountry(t *testing.T) {
	// Test with unsupported country code - should fallback to US
	checker := Checker("XX")

	if checker.GetCountryCode() != "XX" {
		t.Errorf("Expected country code XX, got %s", checker.GetCountryCode())
	}

	// Should still work (using US provider as fallback)
	result := checker.IsHoliday(time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC))
	if !result {
		t.Error("Fallback to US provider should recognize July 4th as a holiday")
	}
}

// Benchmark tests for performance optimization
func BenchmarkFastCountryChecker_IsHoliday(b *testing.B) {
	checker := Checker("US")
	date := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.IsHoliday(date)
	}
}

func BenchmarkFastCountryChecker_AreHolidays_100Dates(b *testing.B) {
	checker := Checker("US")

	// Create 100 dates spanning a year
	dates := make([]time.Time, 100)
	for i := 0; i < 100; i++ {
		dates[i] = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i*3)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.AreHolidays(dates)
	}
}

func BenchmarkFastCountryChecker_CountHolidaysInRange(b *testing.B) {
	checker := Checker("US")
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.CountHolidaysInRange(start, end)
	}
}

// TestFastCountryChecker_NewCountries tests the newly added countries
func TestFastCountryChecker_NewCountries(t *testing.T) {
	testCases := []struct {
		countryCode string
		date        time.Time
		isHoliday   bool
		description string
	}{
		// Chile - Independence Day
		{"CL", time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), true, "Chile Independence Day"},
		// Ireland - Saint Patrick's Day
		{"IE", time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC), true, "Ireland Saint Patrick's Day"},
		// Israel - Independence Day (varies by Hebrew calendar, but test a known date)
		{"IL", time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC), true, "Israel Independence Day"},
		// Brazil - Independence Day
		{"BR", time.Date(2024, 9, 7, 0, 0, 0, 0, time.UTC), true, "Brazil Independence Day"},
		// Argentina - Independence Day
		{"AR", time.Date(2024, 7, 9, 0, 0, 0, 0, time.UTC), true, "Argentina Independence Day"},
		// Regular days that should not be holidays
		{"CL", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false, "Chile regular day"},
		{"IE", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false, "Ireland regular day"},
		{"IL", time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), false, "Israel regular day"},
	}

	for _, tc := range testCases {
		checker := Checker(tc.countryCode)
		result := checker.IsHoliday(tc.date)
		if result != tc.isHoliday {
			t.Errorf("%s: expected %v, got %v", tc.description, tc.isHoliday, result)
		}
	}
}

// TestFastCountryChecker_AllCountriesSupported tests that all 34 countries are supported
func TestFastCountryChecker_AllCountriesSupported(t *testing.T) {
	supportedCountries := []string{
		"AR", "AT", "AU", "BE", "BR", "CA", "CH", "CL", "CN", "DE",
		"ES", "FI", "FR", "GB", "ID", "IE", "IL", "IN", "IT", "JP",
		"KR", "MX", "NL", "NO", "NZ", "PL", "PT", "RU", "SE", "SG",
		"TH", "TR", "UA", "US",
	}

	for _, countryCode := range supportedCountries {
		checker := Checker(countryCode)

		// Test that the checker was created successfully
		if checker == nil {
			t.Errorf("Failed to create checker for country %s", countryCode)
			continue
		}

		// Test that it returns the correct country code
		if checker.GetCountryCode() != countryCode {
			t.Errorf("Country %s: expected country code %s, got %s",
				countryCode, countryCode, checker.GetCountryCode())
		}

		// Test that it can perform basic operations without panicking
		testDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		_ = checker.IsHoliday(testDate)
		_ = checker.GetHolidayName(testDate)
	}

	t.Logf("Successfully tested all %d supported countries", len(supportedCountries))
}
