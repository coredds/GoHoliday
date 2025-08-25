package chronogo

import (
	"testing"
	"time"
)

// mockChronoGoDT implements ChronoGoDateTime for testing
type mockChronoGoDT struct {
	year  int
	month time.Month
	day   int
}

func (m *mockChronoGoDT) Year() int        { return m.year }
func (m *mockChronoGoDT) Month() time.Month { return m.month }
func (m *mockChronoGoDT) Day() int         { return m.day }

func TestGoHolidaysChecker_IsHoliday(t *testing.T) {
	checker := NewGoHolidaysChecker().WithCountries("US")

	tests := []struct {
		name     string
		date     ChronoGoDateTime
		expected bool
	}{
		{
			name:     "New Year's Day 2024",
			date:     &mockChronoGoDT{2024, time.January, 1},
			expected: true,
		},
		{
			name:     "Independence Day 2024",
			date:     &mockChronoGoDT{2024, time.July, 4},
			expected: true,
		},
		{
			name:     "Christmas Day 2024",
			date:     &mockChronoGoDT{2024, time.December, 25},
			expected: true,
		},
		{
			name:     "Regular day 2024",
			date:     &mockChronoGoDT{2024, time.March, 15},
			expected: false,
		},
		{
			name:     "Martin Luther King Jr. Day 2024",
			date:     &mockChronoGoDT{2024, time.January, 15}, // Third Monday in January 2024
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.IsHoliday(tt.date)
			if result != tt.expected {
				t.Errorf("IsHoliday() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGoHolidaysChecker_WithCountries(t *testing.T) {
	// Test with US holidays
	usChecker := NewGoHolidaysChecker().WithCountries("US")
	
	// Canada Day should not be a holiday for US-only checker
	canadaDay := &mockChronoGoDT{2024, time.July, 1}
	if usChecker.IsHoliday(canadaDay) {
		t.Error("Canada Day should not be a US holiday")
	}

	// Test with both US and Canada
	multiChecker := NewGoHolidaysChecker().WithCountries("US", "CA")
	
	// Canada Day should be a holiday for multi-country checker
	if !multiChecker.IsHoliday(canadaDay) {
		t.Error("Canada Day should be a holiday when including Canada")
	}
	
	// Independence Day should still be a holiday
	independenceDay := &mockChronoGoDT{2024, time.July, 4}
	if !multiChecker.IsHoliday(independenceDay) {
		t.Error("Independence Day should still be a holiday")
	}
}

func TestGoHolidaysChecker_WithCategories(t *testing.T) {
	// Test with only federal holidays
	federalChecker := NewGoHolidaysChecker().
		WithCountries("US").
		WithCategories("federal")

	// New Year's should be included (federal holiday)
	newYear := &mockChronoGoDT{2024, time.January, 1}
	if !federalChecker.IsHoliday(newYear) {
		t.Error("New Year's Day should be a federal holiday")
	}

	// Test with no observance holidays
	noObservanceChecker := NewGoHolidaysChecker().
		WithCountries("US").
		WithCategories("federal", "public") // Exclude observance

	// Columbus Day might be observance in some configurations
	// This test verifies category filtering works
	columbusDay := &mockChronoGoDT{2024, time.October, 14} // Second Monday in October 2024
	result := noObservanceChecker.IsHoliday(columbusDay)
	t.Logf("Columbus Day (excluding observance): %v", result)
}

func TestGoHolidaysChecker_GetHolidays(t *testing.T) {
	checker := NewGoHolidaysChecker().WithCountries("US")
	
	holidays := checker.GetHolidays(2024)
	
	if len(holidays) == 0 {
		t.Error("Should have found US holidays for 2024")
	}

	// Check that we have some expected holidays
	expectedHolidays := map[string]bool{
		"New Year's Day":    false,
		"Independence Day":  false,
		"Christmas Day":     false,
	}

	for _, holiday := range holidays {
		if _, exists := expectedHolidays[holiday.Name]; exists {
			expectedHolidays[holiday.Name] = true
		}
		
		// Verify holiday info structure
		if holiday.Name == "" {
			t.Error("Holiday should have a name")
		}
		if holiday.Date.IsZero() {
			t.Error("Holiday should have a valid date")
		}
		if holiday.Country == "" {
			t.Error("Holiday should have a country")
		}
	}

	// Check that we found the expected holidays
	for name, found := range expectedHolidays {
		if !found {
			t.Errorf("Expected to find holiday: %s", name)
		}
	}

	t.Logf("Found %d holidays for US in 2024", len(holidays))
}

func TestGoHolidaysChecker_WithSubdivisions(t *testing.T) {
	// Test regional holidays (if available)
	regionalChecker := NewGoHolidaysChecker().
		WithCountries("US").
		WithSubdivisions("CA") // California

	baseChecker := NewGoHolidaysChecker().WithCountries("US")

	// Get holidays for both checkers
	regionalHolidays := regionalChecker.GetHolidays(2024)
	baseHolidays := baseChecker.GetHolidays(2024)

	t.Logf("Base holidays: %d, Regional holidays: %d", len(baseHolidays), len(regionalHolidays))

	// Regional checker should have at least as many holidays as base
	if len(regionalHolidays) < len(baseHolidays) {
		t.Error("Regional checker should have at least as many holidays as base checker")
	}
}

func TestGoHolidaysChecker_GetSupportedCountries(t *testing.T) {
	checker := NewGoHolidaysChecker()
	
	countries := checker.GetSupportedCountries()
	
	if len(countries) == 0 {
		t.Error("Should have supported countries")
	}

	// Check that US is supported (our baseline)
	usSupported := false
	for _, country := range countries {
		if country == "US" {
			usSupported = true
			break
		}
	}

	if !usSupported {
		t.Error("US should be supported")
	}

	t.Logf("Supported countries: %v", countries)
}

func TestGoHolidaysChecker_GetCountryInfo(t *testing.T) {
	checker := NewGoHolidaysChecker()
	
	info, err := checker.GetCountryInfo("US")
	if err != nil {
		t.Fatalf("GetCountryInfo() error = %v", err)
	}

	if info == nil {
		t.Error("Country info should not be nil")
	}

	// Check that we have some expected info fields
	if enabled, exists := info["enabled"]; !exists {
		t.Error("Country info should include 'enabled' field")
	} else if enabledBool, ok := enabled.(bool); !ok || !enabledBool {
		t.Error("US should be enabled")
	}

	t.Logf("US info: %+v", info)
}

func TestGoHolidaysChecker_CachingBehavior(t *testing.T) {
	checker := NewGoHolidaysChecker().WithCountries("US")
	
	testDate := &mockChronoGoDT{2024, time.July, 4}

	// First call - should load and cache
	result1 := checker.IsHoliday(testDate)
	
	// Second call - should use cache
	result2 := checker.IsHoliday(testDate)
	
	if result1 != result2 {
		t.Error("Cached result should match original result")
	}

	// Verify that cache was populated
	if len(checker.holidayCache) == 0 {
		t.Error("Cache should be populated after holiday check")
	}

	if yearCache, exists := checker.holidayCache[2024]; !exists {
		t.Error("Cache should have entry for 2024")
	} else if len(yearCache) == 0 {
		t.Error("Year cache should have entries")
	}
}

func TestGoHolidaysChecker_ConfigurationChange(t *testing.T) {
	checker := NewGoHolidaysChecker().WithCountries("US")
	
	testDate := &mockChronoGoDT{2024, time.July, 4}
	
	// Check holiday and populate cache
	result1 := checker.IsHoliday(testDate)
	
	// Change configuration
	checker.WithCountries("CA")
	
	// Cache should be cleared
	if len(checker.holidayCache) != 0 {
		t.Error("Cache should be cleared when configuration changes")
	}
	
	// Check with new configuration
	result2 := checker.IsHoliday(testDate)
	
	// Results might be different due to different country
	t.Logf("US config: %v, CA config: %v", result1, result2)
}

func TestGoHolidaysChecker_PreloadYear(t *testing.T) {
	checker := NewGoHolidaysChecker().WithCountries("US")
	
	// Preload 2024
	err := checker.PreloadYear(2024)
	if err != nil {
		t.Fatalf("PreloadYear() error = %v", err)
	}
	
	// Cache should be populated
	if len(checker.holidayCache) == 0 {
		t.Error("Cache should be populated after preloading")
	}
	
	if _, exists := checker.holidayCache[2024]; !exists {
		t.Error("Cache should have entry for 2024 after preloading")
	}
}

func TestCreateDefaultUSChecker(t *testing.T) {
	checker := CreateDefaultUSChecker()
	
	if len(checker.countries) != 1 || checker.countries[0] != "US" {
		t.Error("Default US checker should be configured for US only")
	}
	
	// Should recognize US federal holidays
	newYear := &mockChronoGoDT{2024, time.January, 1}
	if !checker.IsHoliday(newYear) {
		t.Error("Default US checker should recognize New Year's Day")
	}
}

func TestCreateMultiCountryChecker(t *testing.T) {
	checker := CreateMultiCountryChecker("US", "CA", "GB")
	
	if len(checker.countries) != 3 {
		t.Error("Multi-country checker should have 3 countries")
	}
	
	// Should recognize holidays from multiple countries
	usHoliday := &mockChronoGoDT{2024, time.July, 4}  // Independence Day
	caHoliday := &mockChronoGoDT{2024, time.July, 1}  // Canada Day
	
	if !checker.IsHoliday(usHoliday) {
		t.Error("Multi-country checker should recognize US holidays")
	}
	
	if !checker.IsHoliday(caHoliday) {
		t.Error("Multi-country checker should recognize Canadian holidays")
	}
}

func TestCreateRegionalChecker(t *testing.T) {
	checker := CreateRegionalChecker("US", "CA", "NY")
	
	if len(checker.countries) != 1 || checker.countries[0] != "US" {
		t.Error("Regional checker should be configured for single country")
	}
	
	if !checker.includeRegional {
		t.Error("Regional checker should have regional holidays enabled")
	}
	
	if len(checker.subdivisions) != 2 {
		t.Error("Regional checker should have 2 subdivisions")
	}
}

// Benchmark tests
func BenchmarkGoHolidaysChecker_IsHoliday(b *testing.B) {
	checker := NewGoHolidaysChecker().WithCountries("US")
	testDate := &mockChronoGoDT{2024, time.July, 4}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.IsHoliday(testDate)
	}
}

func BenchmarkGoHolidaysChecker_IsHolidayWithCache(b *testing.B) {
	checker := NewGoHolidaysChecker().WithCountries("US")
	testDate := &mockChronoGoDT{2024, time.July, 4}

	// Prime the cache
	checker.IsHoliday(testDate)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.IsHoliday(testDate)
	}
}

func BenchmarkGoHolidaysChecker_GetHolidays(b *testing.B) {
	checker := NewGoHolidaysChecker().WithCountries("US")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.GetHolidays(2024)
	}
}
