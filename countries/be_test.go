package countries

import (
	"testing"
	"time"
)

func TestBEProvider_NewBEProvider(t *testing.T) {
	provider := NewBEProvider()

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.GetCountryCode() != "BE" {
		t.Errorf("Expected country code BE, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	expectedSubdivisions := []string{"BRU", "VLG", "WAL"}
	if len(subdivisions) != len(expectedSubdivisions) {
		t.Errorf("Expected %d subdivisions, got %d", len(expectedSubdivisions), len(subdivisions))
	}
}

func TestBEProvider_LoadHolidays(t *testing.T) {
	provider := NewBEProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test some key Belgian holidays
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			name:     "Nieuwjaar",
			category: "public",
		},
		{
			date:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			name:     "Dag van de Arbeid",
			category: "public",
		},
		{
			date:     time.Date(2024, 7, 21, 0, 0, 0, 0, time.UTC),
			name:     "Nationale Feestdag",
			category: "public",
		},
		{
			date:     time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			name:     "Onze-Lieve-Vrouw-Hemelvaart",
			category: "religious",
		},
		{
			date:     time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			name:     "Allerheiligen",
			category: "religious",
		},
		{
			date:     time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
			name:     "Wapenstilstand",
			category: "public",
		},
		{
			date:     time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			name:     "Kerstmis",
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

func TestBEProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewBEProvider()
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
			name: "Pasen",
		},
		{
			date: easter.AddDate(0, 0, 1), // Easter Monday
			name: "Paasmaandag",
		},
		{
			date: easter.AddDate(0, 0, 39), // Ascension Day
			name: "Onze-Lieve-Heer-Hemelvaart",
		},
		{
			date: easter.AddDate(0, 0, 49), // Whit Sunday
			name: "Pinksteren",
		},
		{
			date: easter.AddDate(0, 0, 50), // Whit Monday
			name: "Pinkstermaandag",
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

func TestBEProvider_MultilingualSupport(t *testing.T) {
	provider := NewBEProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test New Year's Day multilingual support
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[newYear]
	if !exists {
		t.Fatal("Expected New Year's Day to exist")
	}

	expectedLanguages := map[string]string{
		"nl": "Nieuwjaar",
		"fr": "Nouvel An",
		"de": "Neujahr",
		"en": "New Year's Day",
	}

	for lang, expectedName := range expectedLanguages {
		if name, exists := holiday.Languages[lang]; !exists {
			t.Errorf("Expected language %s to exist for New Year's Day", lang)
		} else if name != expectedName {
			t.Errorf("Expected %s name %s, got %s", lang, expectedName, name)
		}
	}
}

func TestBEProvider_GetRegionalHolidays(t *testing.T) {
	provider := NewBEProvider()
	year := 2024

	testCases := []struct {
		region       string
		expectedDate time.Time
		expectedName string
	}{
		{
			region:       "VLG",
			expectedDate: time.Date(2024, 7, 11, 0, 0, 0, 0, time.UTC),
			expectedName: "Feest van de Vlaamse Gemeenschap",
		},
		{
			region:       "WAL",
			expectedDate: time.Date(2024, 9, 27, 0, 0, 0, 0, time.UTC),
			expectedName: "Fête de la Communauté française",
		},
		{
			region:       "BRU",
			expectedDate: time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
			expectedName: "Feest van het Brussels Hoofdstedelijk Gewest",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.region, func(t *testing.T) {
			regionalHolidays := provider.GetRegionalHolidays(year, []string{tc.region})

			holiday, exists := regionalHolidays[tc.expectedDate]
			if !exists {
				t.Errorf("Expected regional holiday for %s on %s to exist", tc.region, tc.expectedDate.Format("2006-01-02"))
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

func TestBEProvider_GetSpecialObservances(t *testing.T) {
	provider := NewBEProvider()
	year := 2024

	observances := provider.GetSpecialObservances(year)

	// Test Saint Nicholas Day
	stNicholas := time.Date(2024, 12, 6, 0, 0, 0, 0, time.UTC)
	holiday, exists := observances[stNicholas]
	if !exists {
		t.Error("Expected Saint Nicholas Day to exist in special observances")
		return
	}

	if holiday.Name != "Sinterklaas" {
		t.Errorf("Expected holiday name Sinterklaas, got %s", holiday.Name)
	}

	if holiday.Category != "religious" {
		t.Errorf("Expected category religious, got %s", holiday.Category)
	}

	// Test King's Day (should exist for 2024 since Philippe became king in 2013)
	kingsDay := time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC)
	holiday, exists = observances[kingsDay]
	if !exists {
		t.Error("Expected King's Day to exist in special observances for 2024")
		return
	}

	if holiday.Name != "Koningsdag" {
		t.Errorf("Expected holiday name Koningsdag, got %s", holiday.Name)
	}
}

func TestBEProvider_HolidayCount(t *testing.T) {
	provider := NewBEProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Belgium should have 12 main holidays (7 fixed + 5 Easter-based)
	expectedCount := 12
	if len(holidays) != expectedCount {
		t.Errorf("Expected %d holidays for Belgium in %d, got %d", expectedCount, year, len(holidays))
	}
}

// Benchmark tests
func BenchmarkBEProvider_LoadHolidays(b *testing.B) {
	provider := NewBEProvider()
	year := 2024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(year)
	}
}
