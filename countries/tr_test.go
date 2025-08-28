package countries

import (
	"testing"
	"time"
)

func TestTRProvider(t *testing.T) {
	provider := NewTRProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "TR" {
		t.Errorf("Expected country code 'TR', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (81 provinces)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 81 {
		t.Errorf("Expected 81 subdivisions for Turkey, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "religious", "commemorative", "seasonal"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestTRHolidays2024(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Turkish holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Yılbaşı", "national"},
		{time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC), "Ulusal Egemenlik ve Çocuk Bayramı", "national"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Emek ve Dayanışma Günü", "national"},
		{time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC), "Atatürk'ü Anma, Gençlik ve Spor Bayramı", "commemorative"},
		{time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), "Demokrasi ve Milli Birlik Günü", "commemorative"},
		{time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC), "Zafer Bayramı", "national"},
		{time.Date(2024, 10, 29, 0, 0, 0, 0, time.UTC), "Cumhuriyet Bayramı", "national"},
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

	// Check that we have a reasonable number of holidays (should include Islamic holidays)
	if len(holidays) < 10 {
		t.Errorf("Expected at least 10 holidays for Turkey in 2024, got %d", len(holidays))
	}
}

func TestTRLanguageSupport(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Republic Day in Turkish and English
	republicDay := time.Date(2024, 10, 29, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[republicDay]
	if !exists {
		t.Fatal("Republic Day not found")
	}

	// Check Turkish translation
	if turkish, ok := holiday.Languages["tr"]; !ok || turkish != "Cumhuriyet Bayramı" {
		t.Errorf("Expected Turkish translation 'Cumhuriyet Bayramı', got '%s'", turkish)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "Republic Day" {
		t.Errorf("Expected English translation 'Republic Day', got '%s'", english)
	}
}

func TestTRRepublicDay(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Republic Day (October 29) - Turkey's most important national holiday
	republicDay := time.Date(2024, 10, 29, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[republicDay]
	if !exists {
		t.Fatal("Republic Day not found")
	}

	if holiday.Name != "Cumhuriyet Bayramı" {
		t.Errorf("Expected Republic Day name 'Cumhuriyet Bayramı', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Republic Day to be national category, got '%s'", holiday.Category)
	}
}

func TestTRAtaturkDay(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Atatürk Commemoration Day (May 19)
	ataturkDay := time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[ataturkDay]
	if !exists {
		t.Fatal("Atatürk Commemoration Day not found")
	}

	expectedName := "Atatürk'ü Anma, Gençlik ve Spor Bayramı"
	if holiday.Name != expectedName {
		t.Errorf("Expected Atatürk Day name '%s', got '%s'", expectedName, holiday.Name)
	}

	if holiday.Category != "commemorative" {
		t.Errorf("Expected Atatürk Day to be commemorative category, got '%s'", holiday.Category)
	}
}

func TestTRDemocracyDay(t *testing.T) {
	provider := NewTRProvider()

	// Test that Democracy Day exists for years >= 2017
	holidays2017 := provider.LoadHolidays(2017)
	democracyDay2017 := time.Date(2017, 7, 15, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2017[democracyDay2017]; !exists {
		t.Error("Democracy Day should exist in 2017")
	}

	// Test that Democracy Day doesn't exist for years < 2017
	holidays2016 := provider.LoadHolidays(2016)
	democracyDay2016 := time.Date(2016, 7, 15, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2016[democracyDay2016]; exists {
		t.Error("Democracy Day should not exist in 2016")
	}
}

func TestTRIslamicHolidays(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Count Islamic holidays (should have multiple days for each festival)
	islamicCount := 0
	for _, holiday := range holidays {
		if holiday.Category == "religious" {
			islamicCount++
		}
	}

	// Should have 7 days total (3 Ramadan + 4 Sacrifice)
	if islamicCount != 7 {
		t.Errorf("Expected 7 Islamic holidays for Turkey in 2024, got %d", islamicCount)
	}

	// Test that Ramadan Festival exists
	ramadanFound := false
	sacrificeFound := false
	for _, holiday := range holidays {
		if holiday.Category == "religious" {
			if holiday.Languages["en"] == "Ramadan Festival Day 1" {
				ramadanFound = true
			}
			if holiday.Languages["en"] == "Sacrifice Festival Day 1" {
				sacrificeFound = true
			}
		}
	}

	if !ramadanFound {
		t.Error("Ramadan Festival should be found in holidays")
	}
	if !sacrificeFound {
		t.Error("Sacrifice Festival should be found in holidays")
	}
}

func TestTRNationalSovereigntyDay(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test National Sovereignty and Children's Day (April 23)
	sovereigntyDay := time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[sovereigntyDay]
	if !exists {
		t.Fatal("National Sovereignty and Children's Day not found")
	}

	expectedName := "Ulusal Egemenlik ve Çocuk Bayramı"
	if holiday.Name != expectedName {
		t.Errorf("Expected sovereignty day name '%s', got '%s'", expectedName, holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected sovereignty day to be national category, got '%s'", holiday.Category)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "National Sovereignty and Children's Day" {
		t.Errorf("Expected proper English translation for National Sovereignty Day")
	}
}

func TestTRVictoryDay(t *testing.T) {
	provider := NewTRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Victory Day (August 30)
	victoryDay := time.Date(2024, 8, 30, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[victoryDay]
	if !exists {
		t.Fatal("Victory Day not found")
	}

	if holiday.Name != "Zafer Bayramı" {
		t.Errorf("Expected Victory Day name 'Zafer Bayramı', got '%s'", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected Victory Day to be national category, got '%s'", holiday.Category)
	}
}

func BenchmarkTRProvider(b *testing.B) {
	provider := NewTRProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkTRIslamicHolidays(b *testing.B) {
	provider := NewTRProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		holidays := make(map[time.Time]*Holiday)
		provider.addIslamicHolidays(holidays, 2024)
	}
}
