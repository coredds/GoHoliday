package countries

import (
	"fmt"
	"testing"
	"time"
)

func TestATProvider_NewATProvider(t *testing.T) {
	provider := NewATProvider()

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.GetCountryCode() != "AT" {
		t.Errorf("Expected country code AT, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	expectedSubdivisions := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	if len(subdivisions) != len(expectedSubdivisions) {
		t.Errorf("Expected %d subdivisions, got %d", len(expectedSubdivisions), len(subdivisions))
	}
}

func TestATProvider_LoadHolidays(t *testing.T) {
	provider := NewATProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test some key Austrian holidays
	testCases := []struct {
		date     time.Time
		name     string
		category string
	}{
		{
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			name:     "Neujahr",
			category: "public",
		},
		{
			date:     time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC),
			name:     "Heilige Drei Könige",
			category: "religious",
		},
		{
			date:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			name:     "Staatsfeiertag",
			category: "public",
		},
		{
			date:     time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
			name:     "Mariä Himmelfahrt",
			category: "religious",
		},
		{
			date:     time.Date(2024, 10, 26, 0, 0, 0, 0, time.UTC),
			name:     "Nationalfeiertag",
			category: "public",
		},
		{
			date:     time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			name:     "Allerheiligen",
			category: "religious",
		},
		{
			date:     time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC),
			name:     "Mariä Empfängnis",
			category: "religious",
		},
		{
			date:     time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			name:     "Christtag",
			category: "public",
		},
		{
			date:     time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
			name:     "Stefanitag",
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

func TestATProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewATProvider()
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
			name: "Ostersonntag",
		},
		{
			date: easter.AddDate(0, 0, 1), // Easter Monday
			name: "Ostermontag",
		},
		{
			date: easter.AddDate(0, 0, 39), // Ascension Day
			name: "Christi Himmelfahrt",
		},
		{
			date: easter.AddDate(0, 0, 49), // Whit Sunday
			name: "Pfingstsonntag",
		},
		{
			date: easter.AddDate(0, 0, 50), // Whit Monday
			name: "Pfingstmontag",
		},
		{
			date: easter.AddDate(0, 0, 60), // Corpus Christi
			name: "Fronleichnam",
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

func TestATProvider_NationalDayHistoricalChange(t *testing.T) {
	provider := NewATProvider()

	// Test that National Day is not present before 1965
	holidays1964 := provider.LoadHolidays(1964)
	nationalDay1964 := time.Date(1964, 10, 26, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays1964[nationalDay1964]; exists {
		t.Error("Expected National Day to not exist in 1964 (before it was established)")
	}

	// Test that National Day is present from 1965 onwards
	holidays1965 := provider.LoadHolidays(1965)
	nationalDay1965 := time.Date(1965, 10, 26, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays1965[nationalDay1965]; !exists {
		t.Error("Expected National Day to exist in 1965 (when it was established)")
	} else if holiday.Name != "Nationalfeiertag" {
		t.Errorf("Expected holiday name Nationalfeiertag, got %s", holiday.Name)
	}
}

func TestATProvider_GermanLanguageSupport(t *testing.T) {
	provider := NewATProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Test National Day German/English support
	nationalDay := time.Date(2024, 10, 26, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[nationalDay]
	if !exists {
		t.Fatal("Expected National Day to exist")
	}

	expectedLanguages := map[string]string{
		"de": "Nationalfeiertag",
		"en": "Austrian National Day",
	}

	for lang, expectedName := range expectedLanguages {
		if name, exists := holiday.Languages[lang]; !exists {
			t.Errorf("Expected language %s to exist for National Day", lang)
		} else if name != expectedName {
			t.Errorf("Expected %s name %s, got %s", lang, expectedName, name)
		}
	}
}

func TestATProvider_GetRegionalHolidays(t *testing.T) {
	provider := NewATProvider()
	year := 2024

	testCases := []struct {
		state        string
		expectedDate time.Time
		expectedName string
	}{
		{
			state:        "1", // Burgenland
			expectedDate: time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
			expectedName: "Martinstag",
		},
		{
			state:        "2", // Carinthia
			expectedDate: time.Date(2024, 10, 10, 0, 0, 0, 0, time.UTC),
			expectedName: "Tag der Volksabstimmung",
		},
		{
			state:        "8", // Vorarlberg
			expectedDate: time.Date(2024, 3, 19, 0, 0, 0, 0, time.UTC),
			expectedName: "Josefitag",
		},
	}

	for _, tc := range testCases {
		t.Run("state_"+tc.state, func(t *testing.T) {
			regionalHolidays := provider.GetRegionalHolidays(year, []string{tc.state})

			holiday, exists := regionalHolidays[tc.expectedDate]
			if !exists {
				t.Errorf("Expected regional holiday for state %s on %s to exist", tc.state, tc.expectedDate.Format("2006-01-02"))
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

func TestATProvider_TyrolSacredHeartCalculation(t *testing.T) {
	provider := NewATProvider()
	year := 2024

	// Test Tyrol's Sacred Heart calculation (Friday after Corpus Christi)
	regionalHolidays := provider.GetRegionalHolidays(year, []string{"7"}) // Tyrol

	// Easter 2024 is March 31, Corpus Christi is 60 days later (May 30)
	// Sacred Heart should be the Friday after that
	easter := EasterSunday(2024)
	corpusChristi := easter.AddDate(0, 0, 60) // May 30, 2024 (Thursday)

	// Find the Friday after Corpus Christi
	daysToFriday := (5 - int(corpusChristi.Weekday()) + 7) % 7
	if daysToFriday == 0 {
		daysToFriday = 7
	}
	expectedSacredHeart := corpusChristi.AddDate(0, 0, daysToFriday)

	holiday, exists := regionalHolidays[expectedSacredHeart]
	if !exists {
		t.Errorf("Expected Sacred Heart holiday on %s to exist for Tyrol", expectedSacredHeart.Format("2006-01-02"))
		return
	}

	if holiday.Name != "Herz-Jesu-Fest" {
		t.Errorf("Expected holiday name Herz-Jesu-Fest, got %s", holiday.Name)
	}
}

func TestATProvider_GetSpecialObservances(t *testing.T) {
	provider := NewATProvider()
	year := 2024

	observances := provider.GetSpecialObservances(year)

	// Test Saint Nicholas Day
	stNicholas := time.Date(2024, 12, 6, 0, 0, 0, 0, time.UTC)
	holiday, exists := observances[stNicholas]
	if !exists {
		t.Error("Expected Saint Nicholas Day to exist in special observances")
		return
	}

	if holiday.Name != "Nikolaus" {
		t.Errorf("Expected holiday name Nikolaus, got %s", holiday.Name)
	}

	if holiday.Category != "cultural" {
		t.Errorf("Expected category cultural, got %s", holiday.Category)
	}

	// Test Austrian State Treaty Day (should exist for 2024 since it's after 1955)
	stateTreaty := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)
	holiday, exists = observances[stateTreaty]
	if !exists {
		t.Error("Expected Austrian State Treaty Day to exist in special observances for 2024")
		return
	}

	if holiday.Name != "Staatsvertragsunterzeichnung" {
		t.Errorf("Expected holiday name Staatsvertragsunterzeichnung, got %s", holiday.Name)
	}

	// Test Good Friday (Easter-based)
	easter := EasterSunday(2024) // March 31, 2024
	expectedGoodFriday := easter.AddDate(0, 0, -2)
	holiday, exists = observances[expectedGoodFriday]
	if !exists {
		t.Errorf("Expected Good Friday to exist on %s", expectedGoodFriday.Format("2006-01-02"))
		return
	}

	if holiday.Name != "Karfreitag" {
		t.Errorf("Expected holiday name Karfreitag, got %s", holiday.Name)
	}
}

func TestATProvider_HolidayCount(t *testing.T) {
	provider := NewATProvider()
	year := 2024

	holidays := provider.LoadHolidays(year)

	// Austria should have 15 main holidays (9 fixed + 6 Easter-based)
	expectedCount := 15
	if len(holidays) != expectedCount {
		t.Errorf("Expected %d holidays for Austria in %d, got %d", expectedCount, year, len(holidays))
	}
}

func TestATProvider_CorpusChristiCalculation(t *testing.T) {
	provider := NewATProvider()

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

			if holiday.Name != "Fronleichnam" {
				t.Errorf("Expected holiday name Fronleichnam, got %s", holiday.Name)
			}

			if holiday.Category != "public" {
				t.Errorf("Expected category public, got %s", holiday.Category)
			}
		})
	}
}

// Benchmark tests
func BenchmarkATProvider_LoadHolidays(b *testing.B) {
	provider := NewATProvider()
	year := 2024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(year)
	}
}
