package countries

import (
	"testing"
	"time"
)

func TestITProvider_Creation(t *testing.T) {
	provider := NewITProvider()

	if provider.GetCountryCode() != "IT" {
		t.Errorf("Expected country code IT, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 20 { // 20 regions
		t.Errorf("Expected 20 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "religious", "regional", "patron"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestITProvider_LoadHolidays2024(t *testing.T) {
	provider := NewITProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key holidays for 2024
	expectedHolidays := map[string]time.Time{
		"New Year's Day":        time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"Epiphany":              time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
		"Liberation Day":        time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),
		"Labour Day":            time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
		"Republic Day":          time.Date(2024, 6, 2, 0, 0, 0, 0, time.UTC),
		"Assumption of Mary":    time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
		"All Saints' Day":       time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		"Immaculate Conception": time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC),
		"Christmas Day":         time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		"St. Stephen's Day":     time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
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
		"Easter Monday": time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), // Easter + 1 day
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
	if len(holidays) < 11 {
		t.Errorf("Expected at least 11 holidays, got %d", len(holidays))
	}
}

func TestITProvider_EasterCalculation(t *testing.T) {
	provider := NewITProvider()

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

func TestITProvider_RegionalHolidays(t *testing.T) {
	provider := NewITProvider()

	// Test Lombardy regional holidays
	lombardyHolidays := provider.GetRegionalHolidays(2024, "LOM")
	if len(lombardyHolidays) == 0 {
		t.Error("Expected Lombardy to have regional holidays")
	}

	// Check for St. Ambrose Day in Lombardy
	stAmbroseDate := time.Date(2024, 12, 7, 0, 0, 0, 0, time.UTC)
	found := false
	for _, holiday := range lombardyHolidays {
		if holiday.Date.Equal(stAmbroseDate) && holiday.Name == "St. Ambrose Day" {
			found = true
			if holiday.Category != "patron" {
				t.Errorf("Expected St. Ambrose Day to be patron category, got %s", holiday.Category)
			}
			if len(holiday.Subdivisions) == 0 || holiday.Subdivisions[0] != "LOM" {
				t.Errorf("Expected St. Ambrose Day to be specific to LOM subdivision")
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find St. Ambrose Day in Lombardy holidays")
	}

	// Test Veneto regional holidays
	venetoHolidays := provider.GetRegionalHolidays(2024, "VEN")
	if len(venetoHolidays) == 0 {
		t.Error("Expected Veneto to have regional holidays")
	}

	// Test unknown region
	unknownHolidays := provider.GetRegionalHolidays(2024, "XXX")
	if len(unknownHolidays) != 0 {
		t.Error("Expected unknown region to have no holidays")
	}
}

func TestITProvider_HolidayLanguages(t *testing.T) {
	provider := NewITProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that holidays have both Italian and English names
	for _, holiday := range holidays {
		if holiday.Languages == nil {
			t.Errorf("Holiday %s missing languages", holiday.Name)
			continue
		}

		if _, hasEn := holiday.Languages["en"]; !hasEn {
			t.Errorf("Holiday %s missing English translation", holiday.Name)
		}

		if _, hasIt := holiday.Languages["it"]; !hasIt {
			t.Errorf("Holiday %s missing Italian translation", holiday.Name)
		}
	}
}

func TestITProvider_HolidayCategories(t *testing.T) {
	provider := NewITProvider()
	holidays := provider.LoadHolidays(2024)

	validCategories := map[string]bool{
		"public":    true,
		"religious": true,
		"regional":  true,
		"patron":    true,
	}

	for _, holiday := range holidays {
		if !validCategories[holiday.Category] {
			t.Errorf("Holiday %s has invalid category: %s", holiday.Name, holiday.Category)
		}
	}
}

func TestITProvider_UniqueHolidays(t *testing.T) {
	provider := NewITProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that Italy has some unique holidays not found in other countries
	uniqueHolidays := []string{
		"Epiphany",
		"Liberation Day",
		"Republic Day",
		"St. Stephen's Day",
	}

	for _, uniqueName := range uniqueHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == uniqueName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find unique Italian holiday: %s", uniqueName)
		}
	}
}
