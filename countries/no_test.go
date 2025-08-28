package countries

import (
	"testing"
	"time"
)

func TestNOProvider(t *testing.T) {
	provider := NewNOProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "NO" {
		t.Errorf("Expected country code 'NO', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (11 counties after 2020 reform)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 11 {
		t.Errorf("Expected 11 subdivisions for Norway, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "religious", "traditional", "royal"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestNOHolidays2024(t *testing.T) {
	provider := NewNOProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Norwegian holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Nyttårsdag", "national"},
		{time.Date(2024, 3, 28, 0, 0, 0, 0, time.UTC), "Skjærtorsdag", "religious"},         // Maundy Thursday 2024
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Langfredag", "religious"},           // Good Friday 2024
		{time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), "Første påskedag", "religious"},      // Easter Sunday 2024
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Andre påskedag", "religious"},        // Easter Monday 2024
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Arbeidernes dag", "national"},        // Labour Day
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "Kristi himmelfartsdag", "religious"}, // Ascension Day 2024
		{time.Date(2024, 5, 17, 0, 0, 0, 0, time.UTC), "Grunnlovsdag", "national"},          // Constitution Day
		{time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC), "Første pinsedag", "religious"},      // Whit Sunday 2024
		{time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), "Andre pinsedag", "religious"},       // Whit Monday 2024
		{time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC), "Julaften", "traditional"},          // Christmas Eve
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Første juledag", "religious"},      // Christmas Day
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Andre juledag", "traditional"},     // Boxing Day
		{time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), "Nyttårsaften", "traditional"},      // New Year's Eve
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
		t.Errorf("Expected at least 12 holidays for Norway in 2024, got %d", len(holidays))
	}
}

func TestNOLanguageSupport(t *testing.T) {
	provider := NewNOProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Constitution Day in Norwegian and English
	constitutionDay := time.Date(2024, 5, 17, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[constitutionDay]
	if !exists {
		t.Fatal("Constitution Day not found")
	}

	// Check Norwegian translation
	if norwegian, ok := holiday.Languages["no"]; !ok || norwegian != "Grunnlovsdag" {
		t.Errorf("Expected Norwegian translation 'Grunnlovsdag', got '%s'", norwegian)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "Constitution Day" {
		t.Errorf("Expected English translation 'Constitution Day', got '%s'", english)
	}
}

func TestNOEasterBasedHolidays(t *testing.T) {
	provider := NewNOProvider()
	holidays := provider.LoadHolidays(2024)

	// Count Easter-based holidays
	easterBasedCount := 0
	easterBasedHolidays := []string{}
	easterBasedNames := []string{
		"Skjærtorsdag", "Langfredag", "Første påskedag", "Andre påskedag",
		"Kristi himmelfartsdag", "Første pinsedag", "Andre pinsedag",
	}

	for _, holiday := range holidays {
		for _, name := range easterBasedNames {
			if holiday.Name == name {
				easterBasedCount++
				easterBasedHolidays = append(easterBasedHolidays, holiday.Name)
				break
			}
		}
	}

	// Should have 7 Easter-based holidays
	if easterBasedCount != 7 {
		t.Errorf("Expected 7 Easter-based holidays for Norway in 2024, got %d: %v", easterBasedCount, easterBasedHolidays)
	}
}

func TestNOConstitutionDay(t *testing.T) {
	provider := NewNOProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Constitution Day (May 17) - Norway's National Day
	constitutionDay := time.Date(2024, 5, 17, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[constitutionDay]
	if !exists {
		t.Fatal("Constitution Day not found")
	}

	if holiday.Name != "Grunnlovsdag" {
		t.Errorf("Expected Constitution Day name 'Grunnlovsdag', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Constitution Day to be national category, got '%s'", holiday.Category)
	}

	// Check that it's Norway's most important national holiday
	if norwegian, ok := holiday.Languages["no"]; !ok || norwegian != "Grunnlovsdag" {
		t.Errorf("Expected proper Norwegian name for Constitution Day")
	}
}

func TestNOChristmasHolidays(t *testing.T) {
	provider := NewNOProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Christmas period holidays
	christmasHolidays := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC), "Julaften", "traditional"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Første juledag", "religious"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Andre juledag", "traditional"},
	}

	for _, tc := range christmasHolidays {
		holiday, exists := holidays[tc.date]
		if !exists {
			t.Errorf("Christmas holiday not found on %s", tc.date.Format("2006-01-02"))
			continue
		}

		if holiday.Name != tc.name {
			t.Errorf("Expected Christmas holiday name '%s', got '%s'", tc.name, holiday.Name)
		}

		if holiday.Category != tc.category {
			t.Errorf("Expected Christmas holiday category '%s', got '%s'", tc.category, holiday.Category)
		}
	}
}

func TestNOEasterCalculation(t *testing.T) {
	provider := NewNOProvider()

	// Test Easter dates for known years (Western/Gregorian Easter)
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
		holidays := provider.LoadHolidays(tc.year)
		expected := time.Date(tc.year, tc.month, tc.day, 0, 0, 0, 0, time.UTC)

		// Find Easter Sunday in holidays
		var easterFound bool
		for date, holiday := range holidays {
			if holiday.Name == "Første påskedag" && date.Equal(expected) {
				easterFound = true
				break
			}
		}

		if !easterFound {
			t.Errorf("Expected Easter %d to be %s, but not found in holidays", tc.year, expected.Format("2006-01-02"))
		}
	}
}

func BenchmarkNOProvider(b *testing.B) {
	provider := NewNOProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkNOEasterCalculation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EasterSunday(2024)
	}
}
