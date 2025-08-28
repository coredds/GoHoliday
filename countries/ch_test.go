package countries

import (
	"testing"
	"time"
)

func TestCHProvider(t *testing.T) {
	provider := NewCHProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "CH" {
		t.Errorf("Expected country code 'CH', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (26 cantons)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 26 {
		t.Errorf("Expected 26 subdivisions for Switzerland, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"federal", "cantonal", "religious", "cultural"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestCHHolidays2024(t *testing.T) {
	provider := NewCHProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Swiss holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Neujahr", "federal"},
		{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), "Berchtoldstag", "cantonal"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Karfreitag", "federal"},     // Good Friday 2024
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Ostermontag", "cantonal"},    // Easter Monday 2024
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Tag der Arbeit", "cantonal"}, // Labour Day
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "Auffahrt", "federal"},        // Ascension 2024
		{time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), "Pfingstmontag", "federal"},  // Whit Monday 2024
		{time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC), "Fronleichnam", "cantonal"},  // Corpus Christi 2024
		{time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), "Schweizer Nationalfeiertag", "federal"},
		{time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), "Mariä Himmelfahrt", "cantonal"},
		{time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC), "Allerheiligen", "cantonal"},
		{time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC), "Mariä Empfängnis", "cantonal"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Weihnachten", "federal"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Stephanstag", "cantonal"},
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
	if len(holidays) < 12 {
		t.Errorf("Expected at least 12 holidays for Switzerland in 2024, got %d", len(holidays))
	}
}

func TestCHMultilingualSupport(t *testing.T) {
	provider := NewCHProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day in all 4 official languages + English
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check German (Deutsch)
	if german, ok := holiday.Languages["de"]; !ok || german != "Neujahr" {
		t.Errorf("Expected German translation 'Neujahr', got '%s'", german)
	}

	// Check French (Français)
	if french, ok := holiday.Languages["fr"]; !ok || french != "Nouvel An" {
		t.Errorf("Expected French translation 'Nouvel An', got '%s'", french)
	}

	// Check Italian (Italiano)
	if italian, ok := holiday.Languages["it"]; !ok || italian != "Capodanno" {
		t.Errorf("Expected Italian translation 'Capodanno', got '%s'", italian)
	}

	// Check Romansh (Rumantsch)
	if romansh, ok := holiday.Languages["rm"]; !ok || romansh != "Niev on" {
		t.Errorf("Expected Romansh translation 'Niev on', got '%s'", romansh)
	}

	// Check English
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestCHFederalVsCantonal(t *testing.T) {
	provider := NewCHProvider()
	holidays := provider.LoadHolidays(2024)

	// Count federal vs cantonal holidays
	federalCount := 0
	cantonalCount := 0
	for _, holiday := range holidays {
		switch holiday.Category {
		case "federal":
			federalCount++
		case "cantonal":
			cantonalCount++
		}
	}

	// Should have several federal holidays
	if federalCount < 5 {
		t.Errorf("Expected at least 5 federal holidays for Switzerland in 2024, got %d", federalCount)
	}

	// Should have several cantonal holidays
	if cantonalCount < 5 {
		t.Errorf("Expected at least 5 cantonal holidays for Switzerland in 2024, got %d", cantonalCount)
	}

	// Test specific federal holidays
	federalHolidays := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),   // New Year
		time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),  // Good Friday
		time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),   // Ascension
		time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),  // Whit Monday
		time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC),   // National Day
		time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), // Christmas
	}

	for _, date := range federalHolidays {
		if holiday, exists := holidays[date]; exists {
			if holiday.Category != "federal" {
				t.Errorf("Expected holiday on %s to be federal, got '%s'", date.Format("2006-01-02"), holiday.Category)
			}
		}
	}
}

func TestCHEasterBasedHolidays(t *testing.T) {
	provider := NewCHProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Easter-based holidays (Easter 2024 is March 31)
	easterHolidays := map[string]time.Time{
		"Karfreitag":    time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), // Good Friday (-2 days)
		"Ostermontag":   time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),  // Easter Monday (+1 day)
		"Auffahrt":      time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),  // Ascension (+39 days)
		"Pfingstmontag": time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), // Whit Monday (+50 days)
		"Fronleichnam":  time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC), // Corpus Christi (+60 days)
	}

	for name, expectedDate := range easterHolidays {
		if holiday, exists := holidays[expectedDate]; exists {
			if holiday.Name != name {
				t.Errorf("Expected holiday name '%s' on %s, got '%s'", name, expectedDate.Format("2006-01-02"), holiday.Name)
			}
		} else {
			t.Errorf("Expected Easter-based holiday '%s' on %s", name, expectedDate.Format("2006-01-02"))
		}
	}
}

func TestCHEasterCalculation(t *testing.T) {
	provider := NewCHProvider()

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

func BenchmarkCHProvider(b *testing.B) {
	provider := NewCHProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkCHEasterCalculation(b *testing.B) {
	provider := NewCHProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.CalculateEaster(2024)
	}
}
