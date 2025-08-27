package countries

import (
	"testing"
	"time"
)

func TestESProvider(t *testing.T) {
	provider := NewESProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "ES" {
		t.Errorf("Expected country code 'ES', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 19 {
		t.Errorf("Expected 19 subdivisions for Spain, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "national", "religious", "regional"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestESHolidays2024(t *testing.T) {
	provider := NewESProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Spanish holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Año Nuevo", "public"},
		{time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), "Día de los Reyes Magos", "religious"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Día del Trabajador", "public"},
		{time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), "Asunción de la Virgen", "religious"},
		{time.Date(2024, 10, 12, 0, 0, 0, 0, time.UTC), "Fiesta Nacional de España", "public"},
		{time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC), "Día de Todos los Santos", "religious"},
		{time.Date(2024, 12, 6, 0, 0, 0, 0, time.UTC), "Día de la Constitución", "public"},
		{time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC), "Inmaculada Concepción", "religious"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Navidad", "religious"},
		{time.Date(2024, 3, 28, 0, 0, 0, 0, time.UTC), "Jueves Santo", "religious"},  // Maundy Thursday 2024
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Viernes Santo", "religious"}, // Good Friday 2024
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
		t.Errorf("Expected at least 10 holidays for Spain in 2024, got %d", len(holidays))
	}
}

func TestESHolidayLanguages(t *testing.T) {
	provider := NewESProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day languages
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Spanish translation
	if spanish, ok := holiday.Languages["es"]; !ok || spanish != "Año Nuevo" {
		t.Errorf("Expected Spanish translation 'Año Nuevo', got '%s'", spanish)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestESEasterCalculation(t *testing.T) {
	provider := NewESProvider()

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

func BenchmarkESProvider(b *testing.B) {
	provider := NewESProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkESEasterCalculation(b *testing.B) {
	provider := NewESProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.CalculateEaster(2024)
	}
}
