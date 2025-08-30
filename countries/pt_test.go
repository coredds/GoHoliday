package countries

import (
	"testing"
	"time"
)

func TestPTProvider_Creation(t *testing.T) {
	provider := NewPTProvider()

	if provider.GetCountryCode() != "PT" {
		t.Errorf("Expected country code PT, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 20 { // 18 districts + 2 autonomous regions
		t.Errorf("Expected 20 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "religious", "regional", "municipal"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestPTProvider_LoadHolidays2024(t *testing.T) {
	provider := NewPTProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key holidays for 2024
	expectedHolidays := map[string]time.Time{
		"New Year's Day":              time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"Freedom Day":                 time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),
		"Labour Day":                  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		"Portugal Day":                time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC),
		"Assumption of Mary":          time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
		"Republic Day":                time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC),
		"All Saints' Day":             time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		"Restoration of Independence": time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		"Immaculate Conception":       time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC),
		"Christmas Day":               time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
	}

	for name, expectedDate := range expectedHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == name && holiday.Date.Equal(expectedDate) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected holiday %s on %s not found", name, expectedDate.Format("2006-01-02"))
		}
	}

	// Test Easter-based holidays for 2024 (Easter was March 31, 2024)
	easterBasedHolidays := map[string]time.Time{
		"Good Friday":      time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), // Easter - 2 days
		"Easter Sunday":    time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), // Easter
		"Carnival Tuesday": time.Date(2024, 2, 13, 0, 0, 0, 0, time.UTC), // Easter - 47 days
		"Corpus Christi":   time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC), // Easter + 60 days
	}

	for name, expectedDate := range easterBasedHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == name && holiday.Date.Equal(expectedDate) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected Easter-based holiday %s on %s not found", name, expectedDate.Format("2006-01-02"))
		}
	}

	// Check total number of holidays
	if len(holidays) < 14 {
		t.Errorf("Expected at least 14 holidays, got %d", len(holidays))
	}
}

func TestPTProvider_EasterCalculation(t *testing.T) {
	provider := NewPTProvider()

	// Test Easter dates for known years
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)},
		{2027, time.Date(2027, 3, 28, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		easter := provider.calculateEaster(tc.year)
		if !easter.Equal(tc.expected) {
			t.Errorf("Easter %d: expected %s, got %s", tc.year, tc.expected.Format("2006-01-02"), easter.Format("2006-01-02"))
		}
	}
}

func TestPTProvider_HolidayLanguages(t *testing.T) {
	provider := NewPTProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that holidays have both Portuguese and English names
	for _, holiday := range holidays {
		if holiday.Languages == nil {
			t.Errorf("Holiday %s missing languages", holiday.Name)
			continue
		}

		if _, hasEn := holiday.Languages["en"]; !hasEn {
			t.Errorf("Holiday %s missing English translation", holiday.Name)
		}

		if _, hasPt := holiday.Languages["pt"]; !hasPt {
			t.Errorf("Holiday %s missing Portuguese translation", holiday.Name)
		}
	}
}

func TestPTProvider_HolidayCategories(t *testing.T) {
	provider := NewPTProvider()
	holidays := provider.LoadHolidays(2024)

	validCategories := map[string]bool{
		"public":    true,
		"religious": true,
		"regional":  true,
		"municipal": true,
	}

	for _, holiday := range holidays {
		if !validCategories[holiday.Category] {
			t.Errorf("Holiday %s has invalid category: %s", holiday.Name, holiday.Category)
		}
	}
}
