package countries

import (
	"testing"
	"time"
)

func TestNLProvider(t *testing.T) {
	provider := NewNLProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "NL" {
		t.Errorf("Expected country code 'NL', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 12 {
		t.Errorf("Expected 12 subdivisions for Netherlands, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "national", "religious", "royal"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestNLHolidays2024(t *testing.T) {
	provider := NewNLProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Dutch holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Nieuwjaarsdag", "public"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Goede Vrijdag", "religious"},      // Good Friday 2024
		{time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), "Eerste Paasdag", "religious"},     // Easter Sunday 2024
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Tweede Paasdag", "religious"},      // Easter Monday 2024
		{time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC), "Koningsdag", "royal"},             // King's Day 2024
		{time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC), "Bevrijdingsdag", "public"},         // Liberation Day
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "Hemelvaartsdag", "religious"},      // Ascension Day 2024
		{time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC), "Eerste Pinksterdag", "religious"}, // Whit Sunday 2024
		{time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), "Tweede Pinksterdag", "religious"}, // Whit Monday 2024
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Eerste Kerstdag", "religious"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Tweede Kerstdag", "religious"},
	}

	for _, tc := range testCases {
		holiday, exists := holidays[tc.date]
		if !exists {
			t.Errorf("Expected holiday on %s, but none found", tc.date.Format("2006-01-02"))
			continue
		}

		if holiday.Name != tc.name {
			t.Errorf("Expected holiday name '%s' on %s, got '%s'", tc.name, tc.date.Format("2006-01-02"), holiday.Name)
		}

		if holiday.Category != tc.category {
			t.Errorf("Expected category '%s' for %s, got '%s'", tc.category, tc.name, holiday.Category)
		}
	}

	// Check that we have a reasonable number of holidays
	if len(holidays) < 10 {
		t.Errorf("Expected at least 10 holidays for Netherlands in 2024, got %d", len(holidays))
	}
}

func TestNLKingsDay(t *testing.T) {
	provider := NewNLProvider()

	// Test King's Day when April 27 is not Sunday (2024)
	holidays2024 := provider.LoadHolidays(2024)
	kingsDay2024 := time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC) // Saturday
	if _, exists := holidays2024[kingsDay2024]; !exists {
		t.Error("King's Day 2024 should be on April 27 (Saturday)")
	}

	// Test King's Day when April 27 is Sunday (should move to April 26)
	// Find a year where April 27 is Sunday
	holidays2025 := provider.LoadHolidays(2025)
	kingsDay2025Expected := time.Date(2025, 4, 27, 0, 0, 0, 0, time.UTC) // Sunday in 2025
	kingsDayMoved := time.Date(2025, 4, 26, 0, 0, 0, 0, time.UTC)

	if kingsDay2025Expected.Weekday() == time.Sunday {
		if _, exists := holidays2025[kingsDayMoved]; !exists {
			t.Error("King's Day should move to April 26 when April 27 is Sunday")
		}
		if _, exists := holidays2025[kingsDay2025Expected]; exists {
			t.Error("King's Day should not be on April 27 when it's Sunday")
		}
	}
}

func TestNLHolidayLanguages(t *testing.T) {
	provider := NewNLProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day languages
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Dutch translation
	if dutch, ok := holiday.Languages["nl"]; !ok || dutch != "Nieuwjaarsdag" {
		t.Errorf("Expected Dutch translation 'Nieuwjaarsdag', got '%s'", dutch)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestNLEasterCalculation(t *testing.T) {
	provider := NewNLProvider()

	// Test Easter dates for known years
	testCases := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2024, time.March, 31}, // Easter 2024
		{2025, time.April, 20}, // Easter 2025
		{2026, time.April, 5},  // Easter 2026
		{2027, time.March, 28}, // Easter 2027
	}

	for _, tc := range testCases {
		easter := provider.CalculateEaster(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		if !easter.Equal(expected) {
			t.Errorf("Expected Easter %d to be %s, got %s", tc.year, expected.Format("2006-01-02"), easter.Format("2006-01-02"))
		}
	}
}

func BenchmarkNLProvider(b *testing.B) {
	provider := NewNLProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkNLEasterCalculation(b *testing.B) {
	provider := NewNLProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.CalculateEaster(2024)
	}
}
