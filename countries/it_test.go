package countries

import (
	"testing"
	"time"
)

func TestITProvider(t *testing.T) {
	provider := NewITProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "IT" {
		t.Errorf("Expected country code 'IT', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 20 {
		t.Errorf("Expected 20 subdivisions for Italy, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "national", "religious", "regional"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestITHolidays2024(t *testing.T) {
	provider := NewITProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Italian holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Capodanno", "public"},
		{time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), "Epifania", "religious"},
		{time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC), "Festa della Liberazione", "public"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Festa del Lavoro", "public"},
		{time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC), "Festa della Repubblica", "public"},
		{time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), "Assunzione di Maria", "religious"},
		{time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC), "Ognissanti", "religious"},
		{time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC), "Immacolata Concezione", "religious"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Natale", "religious"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Santo Stefano", "religious"},
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Lunedì di Pasqua", "religious"}, // Easter Monday 2024
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
		t.Errorf("Expected at least 10 holidays for Italy in 2024, got %d", len(holidays))
	}
}

func TestITHolidayLanguages(t *testing.T) {
	provider := NewITProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day languages
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Italian translation
	if italian, ok := holiday.Languages["it"]; !ok || italian != "Capodanno" {
		t.Errorf("Expected Italian translation 'Capodanno', got '%s'", italian)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestITEasterCalculation(t *testing.T) {
	provider := NewITProvider()

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

func BenchmarkITProvider(b *testing.B) {
	provider := NewITProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkITEasterCalculation(b *testing.B) {
	provider := NewITProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.CalculateEaster(2024)
	}
}
