package countries

import (
	"fmt"
	"testing"
	"time"
)

func TestPLProvider_NewPLProvider(t *testing.T) {
	provider := NewPLProvider()

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.GetCountryCode() != "PL" {
		t.Errorf("Expected country code PL, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	expectedSubdivisions := []string{
		"DS", "KP", "LB", "LD", "LU", "MA", "MZ", "OP", "PK", "PD", "PM", "SL", "SK", "WN", "WP", "ZP",
	}
	if len(subdivisions) != len(expectedSubdivisions) {
		t.Errorf("Expected %d subdivisions, got %d", len(expectedSubdivisions), len(subdivisions))
	}
}

func TestPLProvider_LoadHolidays(t *testing.T) {
	provider := NewPLProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test some key Polish holidays
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			name:     "Nowy Rok",
			category: "public",
		},
		{
			date:     time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			name:     "Święto Trzech Króli",
			category: "religious",
		},
		{
			date:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			name:     "Święto Pracy",
			category: "public",
		},
		{
			date:     time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC),
			name:     "Święto Konstytucji 3 Maja",
			category: "national",
		},
		{
			date:     time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			name:     "Wniebowzięcie Najświętszej Maryi Panny",
			category: "religious",
		},
		{
			date:     time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			name:     "Wszystkich Świętych",
			category: "religious",
		},
		{
			date:     time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
			name:     "Narodowe Święto Niepodległości",
			category: "national",
		},
		{
			date:     time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			name:     "Boże Narodzenie",
			category: "public",
		},
		{
			date:     time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
			name:     "Drugi dzień Bożego Narodzenia",
			category: "public",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			holiday, exists := holidays[tc.date]
			if !exists {
				t.Errorf("Expected holiday %s on %s to exist", tc.name, tc.date.Format("2006-01-02"))
				return
			}

			if holiday.Name != tc.name {
				t.Errorf("Expected holiday name %s, got %s", tc.name, holiday.Name)
			}

			if holiday.Category != tc.category {
				t.Errorf("Expected category %s, got %s", tc.category, holiday.Category)
			}
		})
	}
}

func TestPLProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewPLProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Easter 2024 is March 31
	easter := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		date time.Time
		name string
	}{
		{
			date: easter,
			name: "Niedziela Wielkanocna",
		},
		{
			date: easter.AddDate(0, 0, 1), // Easter Monday
			name: "Poniedziałek Wielkanocny",
		},
		{
			date: easter.AddDate(0, 0, 49), // Whit Sunday
			name: "Zielone Świątki",
		},
		{
			date: easter.AddDate(0, 0, 60), // Corpus Christi
			name: "Boże Ciało",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			holiday, exists := holidays[tc.date]
			if !exists {
				t.Errorf("Expected Easter-based holiday %s on %s to exist", tc.name, tc.date.Format("2006-01-02"))
				return
			}

			if holiday.Name != tc.name {
				t.Errorf("Expected holiday name %s, got %s", tc.name, holiday.Name)
			}
		})
	}
}

func TestPLProvider_EpiphanyHistoricalChange(t *testing.T) {
	provider := NewPLProvider()

	// Test that Epiphany is not present before 2011
	holidays2010 := provider.LoadHolidays(2010)
	epiphany2010 := time.Date(2010, 1, 6, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2010[epiphany2010]; exists {
		t.Error("Expected Epiphany to not exist in 2010 (before it became a public holiday)")
	}

	// Test that Epiphany is present from 2011 onwards
	holidays2011 := provider.LoadHolidays(2011)
	epiphany2011 := time.Date(2011, 1, 6, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays2011[epiphany2011]; !exists {
		t.Error("Expected Epiphany to exist in 2011 (when it became a public holiday)")
	} else if holiday.Name != "Święto Trzech Króli" {
		t.Errorf("Expected holiday name Święto Trzech Króli, got %s", holiday.Name)
	}
}

func TestPLProvider_BilingualSupport(t *testing.T) {
	provider := NewPLProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test Constitution Day bilingual support
	constitutionDay := time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[constitutionDay]
	if !exists {
		t.Fatal("Expected Constitution Day to exist")
	}

	expectedLanguages := map[string]string{
		"pl": "Święto Konstytucji 3 Maja",
		"en": "Constitution Day",
	}

	for lang, expectedName := range expectedLanguages {
		if name, exists := holiday.Languages[lang]; !exists {
			t.Errorf("Expected language %s to exist for Constitution Day", lang)
		} else if name != expectedName {
			t.Errorf("Expected %s name %s, got %s", lang, expectedName, name)
		}
	}
}

func TestPLProvider_GetRegionalHolidays(t *testing.T) {
	provider := NewPLProvider()
	year := 2024

	testCases := []struct {
		voivodeship  string
		expectedDate time.Time
		expectedName string
	}{
		{
			voivodeship:  "SL",
			expectedDate: time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC),
			expectedName: "Dzień Powstań Śląskich",
		},
		{
			voivodeship:  "MA",
			expectedDate: time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
			expectedName: "Święty Stanisław",
		},
		{
			voivodeship:  "PD",
			expectedDate: time.Date(2024, 9, 8, 0, 0, 0, 0, time.UTC),
			expectedName: "Narodzenie Najświętszej Maryi Panny",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.voivodeship, func(t *testing.T) {
			regionalHolidays := provider.GetRegionalHolidays(year, []string{tc.voivodeship})

			holiday, exists := regionalHolidays[tc.expectedDate]
			if !exists {
				t.Errorf("Expected regional holiday for %s on %s to exist", tc.voivodeship, tc.expectedDate.Format("2006-01-02"))
				return
			}

			if holiday.Name != tc.expectedName {
				t.Errorf("Expected holiday name %s, got %s", tc.expectedName, holiday.Name)
			}

			if holiday.Category != "regional" {
				t.Errorf("Expected category regional, got %s", holiday.Category)
			}
		})
	}
}

func TestPLProvider_GetSpecialObservances(t *testing.T) {
	provider := NewPLProvider()
	year := 2024

	observances := provider.GetSpecialObservances(year)

	// Test Saint Nicholas Day
	stNicholas := time.Date(2024, 12, 6, 0, 0, 0, 0, time.UTC)
	holiday, exists := observances[stNicholas]
	if !exists {
		t.Error("Expected Saint Nicholas Day to exist in special observances")
		return
	}

	if holiday.Name != "Mikołajki" {
		t.Errorf("Expected holiday name Mikołajki, got %s", holiday.Name)
	}

	if holiday.Category != "cultural" {
		t.Errorf("Expected category cultural, got %s", holiday.Category)
	}

	// Test Warsaw Uprising Day (should exist for 2024 since it's after 1944)
	warsawUprising := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists = observances[warsawUprising]
	if !exists {
		t.Error("Expected Warsaw Uprising Day to exist in special observances for 2024")
		return
	}

	if holiday.Name != "Dzień Powstania Warszawskiego" {
		t.Errorf("Expected holiday name Dzień Powstania Warszawskiego, got %s", holiday.Name)
	}

	// Test Fat Thursday (Easter-based)
	easter := EasterSunday(2024) // March 31, 2024
	expectedFatThursday := easter.AddDate(0, 0, -52)
	holiday, exists = observances[expectedFatThursday]
	if !exists {
		t.Errorf("Expected Fat Thursday to exist on %s", expectedFatThursday.Format("2006-01-02"))
		return
	}

	if holiday.Name != "Tłusty Czwartek" {
		t.Errorf("Expected holiday name Tłusty Czwartek, got %s", holiday.Name)
	}
}

func TestPLProvider_HolidayCount(t *testing.T) {
	provider := NewPLProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Poland should have 14 main holidays (9 fixed + 4 Easter-based + 1 Christmas Eve)
	expectedCount := 14
	if len(holidays) != expectedCount {
		t.Errorf("Expected %d holidays for Poland in %d, got %d", expectedCount, year, len(holidays))
	}
}

func TestPLProvider_CorpusChristiCalculation(t *testing.T) {
	provider := NewPLProvider()

	testCases := []struct {
		year                  int
		expectedEaster        time.Time
		expectedCorpusChristi time.Time
	}{
		{
			year:                  2024,
			expectedEaster:        time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC),
			expectedCorpusChristi: time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC), // 60 days after Easter
		},
		{
			year:                  2025,
			expectedEaster:        time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC),
			expectedCorpusChristi: time.Date(2025, 6, 19, 0, 0, 0, 0, time.UTC), // 60 days after Easter
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("year_%d", tc.year), func(t *testing.T) {
			holidays := provider.LoadHolidays(tc.year)

			// Verify Easter calculation
			easter := EasterSunday(tc.year)
			if !easter.Equal(tc.expectedEaster) {
				t.Errorf("Expected Easter %s, got %s", tc.expectedEaster.Format("2006-01-02"), easter.Format("2006-01-02"))
			}

			// Verify Corpus Christi
			holiday, exists := holidays[tc.expectedCorpusChristi]
			if !exists {
				t.Errorf("Expected Corpus Christi on %s to exist", tc.expectedCorpusChristi.Format("2006-01-02"))
				return
			}

			if holiday.Name != "Boże Ciało" {
				t.Errorf("Expected holiday name Boże Ciało, got %s", holiday.Name)
			}

			if holiday.Category != "public" {
				t.Errorf("Expected category public, got %s", holiday.Category)
			}
		})
	}
}

// Benchmark tests
func BenchmarkPLProvider_LoadHolidays(b *testing.B) {
	provider := NewPLProvider()
	year := 2024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(year)
	}
}
