package countries

import (
	"testing"
	"time"
)

func TestKRProvider(t *testing.T) {
	provider := NewKRProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "KR" {
		t.Errorf("Expected country code 'KR', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 17 {
		t.Errorf("Expected 17 subdivisions for South Korea, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "national", "traditional", "commemorative"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestKRHolidays2024(t *testing.T) {
	provider := NewKRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key South Korean holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "신정", "public"},
		{time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), "설날", "traditional"}, // Lunar New Year 2024
		{time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC), "삼일절", "national"},
		{time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC), "어린이날", "public"},
		{time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC), "부처님 오신 날", "traditional"}, // Buddha's Birthday 2024
		{time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC), "현충일", "commemorative"},
		{time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), "광복절", "national"},
		{time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC), "추석", "traditional"}, // Chuseok 2024
		{time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC), "개천절", "national"},
		{time.Date(2024, 10, 9, 0, 0, 0, 0, time.UTC), "한글날", "national"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "성탄절", "public"},
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
		t.Errorf("Expected at least 12 holidays for South Korea in 2024, got %d", len(holidays))
	}
}

func TestKRSeollalHolidays(t *testing.T) {
	provider := NewKRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that Seollal (Lunar New Year) spans 3 days in 2024
	seollalDates := []time.Time{
		time.Date(2024, 2, 9, 0, 0, 0, 0, time.UTC),  // Day before
		time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC), // Actual day
		time.Date(2024, 2, 11, 0, 0, 0, 0, time.UTC), // Day after
	}

	for _, date := range seollalDates {
		if _, exists := holidays[date]; !exists {
			t.Errorf("Expected Seollal holiday on %s", date.Format("2006-01-02"))
		}
	}

	// Check the main Seollal day
	mainSeollal := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[mainSeollal]; exists {
		if holiday.Name != "설날" {
			t.Errorf("Expected main Seollal day to be named '설날', got '%s'", holiday.Name)
		}
	}
}

func TestKRChuseokHolidays(t *testing.T) {
	provider := NewKRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that Chuseok spans 3 days in 2024
	chuseokDates := []time.Time{
		time.Date(2024, 9, 16, 0, 0, 0, 0, time.UTC), // Day before
		time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC), // Actual day
		time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), // Day after
	}

	for _, date := range chuseokDates {
		if _, exists := holidays[date]; !exists {
			t.Errorf("Expected Chuseok holiday on %s", date.Format("2006-01-02"))
		}
	}

	// Check the main Chuseok day
	mainChuseok := time.Date(2024, 9, 17, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[mainChuseok]; exists {
		if holiday.Name != "추석" {
			t.Errorf("Expected main Chuseok day to be named '추석', got '%s'", holiday.Name)
		}
	}
}

func TestKRHolidayLanguages(t *testing.T) {
	provider := NewKRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day languages
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Korean translation
	if korean, ok := holiday.Languages["ko"]; !ok || korean != "신정" {
		t.Errorf("Expected Korean translation '신정', got '%s'", korean)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestKRTraditionalHolidays(t *testing.T) {
	provider := NewKRProvider()
	holidays := provider.LoadHolidays(2024)

	// Count traditional holidays
	traditionalCount := 0
	for _, holiday := range holidays {
		if holiday.Category == "traditional" {
			traditionalCount++
		}
	}

	// Should have Seollal (3 days), Buddha's Birthday, and Chuseok (3 days) = 7 traditional holidays
	if traditionalCount != 7 {
		t.Errorf("Expected 7 traditional holidays for South Korea in 2024, got %d", traditionalCount)
	}
}

func BenchmarkKRProvider(b *testing.B) {
	provider := NewKRProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}
