package countries

import (
	"testing"
	"time"
)

func TestSGProvider(t *testing.T) {
	provider := NewSGProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "SG" {
		t.Errorf("Expected country code 'SG', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 5 {
		t.Errorf("Expected 5 subdivisions for Singapore, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "religious", "cultural", "national"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestSGHolidays2024(t *testing.T) {
	provider := NewSGProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Singaporean holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day", "public"},
		{time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), "Chinese New Year", "cultural"},
		{time.Date(2024, 2, 11, 0, 0, 0, 0, time.UTC), "Chinese New Year (Day 2)", "cultural"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Good Friday", "religious"}, // Good Friday 2024
		{time.Date(2024, 4, 10, 0, 0, 0, 0, time.UTC), "Hari Raya Puasa", "religious"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Labour Day", "public"},
		{time.Date(2024, 5, 22, 0, 0, 0, 0, time.UTC), "Vesak Day", "religious"},
		{time.Date(2024, 6, 17, 0, 0, 0, 0, time.UTC), "Hari Raya Haji", "religious"},
		{time.Date(2024, 8, 9, 0, 0, 0, 0, time.UTC), "National Day", "national"},
		{time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC), "Deepavali", "religious"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day", "religious"},
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
		t.Errorf("Expected at least 10 holidays for Singapore in 2024, got %d", len(holidays))
	}
}

func TestSGMultilingualSupport(t *testing.T) {
	provider := NewSGProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day in all 4 official languages
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check English (official)
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}

	// Check Chinese (官方语言)
	if chinese, ok := holiday.Languages["zh"]; !ok || chinese != "元旦" {
		t.Errorf("Expected Chinese translation '元旦', got '%s'", chinese)
	}

	// Check Malay (bahasa rasmi)
	if malay, ok := holiday.Languages["ms"]; !ok || malay != "Hari Tahun Baru" {
		t.Errorf("Expected Malay translation 'Hari Tahun Baru', got '%s'", malay)
	}

	// Check Tamil (அதிகாரப்பூர்வ மொழி)
	if tamil, ok := holiday.Languages["ta"]; !ok || tamil != "புத்தாண்டு" {
		t.Errorf("Expected Tamil translation 'புத்தாண்டு', got '%s'", tamil)
	}
}

func TestSGCulturalHolidays(t *testing.T) {
	provider := NewSGProvider()
	holidays := provider.LoadHolidays(2024)

	// Count cultural holidays (Chinese New Year)
	culturalCount := 0
	for _, holiday := range holidays {
		if holiday.Category == "cultural" {
			culturalCount++
		}
	}

	// Should have Chinese New Year (2 days) = 2 cultural holidays
	if culturalCount != 2 {
		t.Errorf("Expected 2 cultural holidays for Singapore in 2024, got %d", culturalCount)
	}

	// Test Chinese New Year specifically
	cny1 := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[cny1]; exists {
		if holiday.Category != "cultural" {
			t.Errorf("Expected Chinese New Year to be cultural category, got '%s'", holiday.Category)
		}
		if holiday.Languages["zh"] != "农历新年" {
			t.Errorf("Expected Chinese translation '农历新年', got '%s'", holiday.Languages["zh"])
		}
	} else {
		t.Error("Chinese New Year Day 1 not found")
	}
}

func TestSGReligiousHolidays(t *testing.T) {
	provider := NewSGProvider()
	holidays := provider.LoadHolidays(2024)

	// Count religious holidays
	religiousCount := 0
	religiousHolidays := []string{}
	for _, holiday := range holidays {
		if holiday.Category == "religious" {
			religiousCount++
			religiousHolidays = append(religiousHolidays, holiday.Name)
		}
	}

	// Should have Good Friday, Hari Raya Puasa, Vesak Day, Hari Raya Haji, Deepavali, Christmas = 6 religious holidays
	if religiousCount != 6 {
		t.Errorf("Expected 6 religious holidays for Singapore in 2024, got %d: %v", religiousCount, religiousHolidays)
	}
}

func TestSGEasterCalculation(t *testing.T) {
	provider := NewSGProvider()

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

func BenchmarkSGProvider(b *testing.B) {
	provider := NewSGProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkSGEasterCalculation(b *testing.B) {
	provider := NewSGProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.CalculateEaster(2024)
	}
}
