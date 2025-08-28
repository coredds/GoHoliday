package countries

import (
	"testing"
	"time"
)

func TestSEProvider(t *testing.T) {
	provider := NewSEProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "SE" {
		t.Errorf("Expected country code 'SE', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (21 counties)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 21 {
		t.Errorf("Expected 21 subdivisions for Sweden, got %d", len(subdivisions))
	}

	// Test categories
	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"public", "religious", "cultural", "traditional"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestSEHolidays2024(t *testing.T) {
	provider := NewSEProvider()
	holidays := provider.LoadHolidays(2024)

	// Test some key Swedish holidays for 2024
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Nyårsdagen", "public"},
		{time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), "Trettondedag jul", "religious"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Långfredagen", "religious"},          // Good Friday 2024
		{time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), "Påskdagen", "religious"},             // Easter Sunday 2024
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Annandag påsk", "religious"},          // Easter Monday 2024
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Första maj", "public"},                // Labour Day
		{time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), "Kristi himmelsfärdsdag", "religious"}, // Ascension 2024
		{time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC), "Pingstdagen", "religious"},           // Whit Sunday 2024
		{time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC), "Sveriges nationaldag", "public"},
		{time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC), "Midsommarafton", "traditional"}, // Midsummer Eve 2024
		{time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC), "Midsommardagen", "traditional"}, // Midsummer Day 2024
		{time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC), "Alla helgons dag", "religious"}, // All Saints 2024
		{time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC), "Julafton", "traditional"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Juldagen", "religious"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Annandag jul", "religious"},
		{time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), "Nyårsafton", "traditional"},
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
	if len(holidays) < 15 {
		t.Errorf("Expected at least 15 holidays for Sweden in 2024, got %d", len(holidays))
	}
}

func TestSESwedishLanguageSupport(t *testing.T) {
	provider := NewSEProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day in Swedish and English
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("New Year's Day not found")
	}

	// Check Swedish translation
	if swedish, ok := holiday.Languages["sv"]; !ok || swedish != "Nyårsdagen" {
		t.Errorf("Expected Swedish translation 'Nyårsdagen', got '%s'", swedish)
	}

	// Check English translation
	if english, ok := holiday.Languages["en"]; !ok || english != "New Year's Day" {
		t.Errorf("Expected English translation 'New Year's Day', got '%s'", english)
	}
}

func TestSEMidsummer(t *testing.T) {
	provider := NewSEProvider()

	// Test Midsummer calculations for multiple years
	testCases := []struct {
		year            int
		expectedEveDate time.Time
		expectedDayDate time.Time
	}{
		{2024, time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC), time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC)}, // Friday/Saturday
		{2025, time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC), time.Date(2025, 6, 21, 0, 0, 0, 0, time.UTC)}, // Friday/Saturday
		{2026, time.Date(2026, 6, 19, 0, 0, 0, 0, time.UTC), time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC)}, // Friday/Saturday
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)

		// Check Midsummer Eve
		if holiday, exists := holidays[tc.expectedEveDate]; exists {
			if holiday.Name != "Midsommarafton" {
				t.Errorf("Expected Midsummer Eve name 'Midsommarafton' for %d, got '%s'", tc.year, holiday.Name)
			}
			if tc.expectedEveDate.Weekday() != time.Friday {
				t.Errorf("Expected Midsummer Eve %d to be on Friday, got %s", tc.year, tc.expectedEveDate.Weekday())
			}
		} else {
			t.Errorf("Midsummer Eve not found for %d on %s", tc.year, tc.expectedEveDate.Format("2006-01-02"))
		}

		// Check Midsummer Day
		if holiday, exists := holidays[tc.expectedDayDate]; exists {
			if holiday.Name != "Midsommardagen" {
				t.Errorf("Expected Midsummer Day name 'Midsommardagen' for %d, got '%s'", tc.year, holiday.Name)
			}
			if tc.expectedDayDate.Weekday() != time.Saturday {
				t.Errorf("Expected Midsummer Day %d to be on Saturday, got %s", tc.year, tc.expectedDayDate.Weekday())
			}
		} else {
			t.Errorf("Midsummer Day not found for %d on %s", tc.year, tc.expectedDayDate.Format("2006-01-02"))
		}
	}
}

func TestSEAllSaintsDay(t *testing.T) {
	provider := NewSEProvider()

	// Test All Saints' Day calculations for multiple years
	testCases := []struct {
		year         int
		expectedDate time.Time
	}{
		{2024, time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC)},  // Saturday
		{2025, time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC)},  // Saturday
		{2026, time.Date(2026, 10, 31, 0, 0, 0, 0, time.UTC)}, // Saturday
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)

		if holiday, exists := holidays[tc.expectedDate]; exists {
			if holiday.Name != "Alla helgons dag" {
				t.Errorf("Expected All Saints' Day name 'Alla helgons dag' for %d, got '%s'", tc.year, holiday.Name)
			}
			if tc.expectedDate.Weekday() != time.Saturday {
				t.Errorf("Expected All Saints' Day %d to be on Saturday, got %s", tc.year, tc.expectedDate.Weekday())
			}
		} else {
			t.Errorf("All Saints' Day not found for %d on %s", tc.year, tc.expectedDate.Format("2006-01-02"))
		}
	}
}

func TestSETraditionalHolidays(t *testing.T) {
	provider := NewSEProvider()
	holidays := provider.LoadHolidays(2024)

	// Count traditional holidays
	traditionalCount := 0
	traditionalHolidays := []string{}
	for _, holiday := range holidays {
		if holiday.Category == "traditional" {
			traditionalCount++
			traditionalHolidays = append(traditionalHolidays, holiday.Name)
		}
	}

	// Should have Midsummer Eve, Midsummer Day, Christmas Eve, New Year's Eve = 4 traditional holidays
	if traditionalCount != 4 {
		t.Errorf("Expected 4 traditional holidays for Sweden in 2024, got %d: %v", traditionalCount, traditionalHolidays)
	}
}

func TestSEEasterCalculation(t *testing.T) {
	provider := NewSEProvider()

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

func BenchmarkSEProvider(b *testing.B) {
	provider := NewSEProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(2024)
	}
}

func BenchmarkSEMidsummerCalculation(b *testing.B) {
	provider := NewSEProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.calculateMidsummerEve(2024)
	}
}
