package countries

import (
	"testing"
	"time"
)

func TestUAProvider_Basic(t *testing.T) {
	provider := NewUAProvider()

	// Test basic provider info
	if provider.GetCountryCode() != "UA" {
		t.Errorf("Expected country code UA, got %s", provider.GetCountryCode())
	}

	if provider.GetName() != "Ukraine" {
		t.Errorf("Expected country name Ukraine, got %s", provider.GetName())
	}

	// Test subdivisions
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) == 0 {
		t.Error("Expected subdivisions to be defined")
	}

	// Check for major oblasts
	expectedSubdivisions := []string{"KV", "LV", "OD", "DP", "KK", "30"} // Kyiv, Lviv, Odessa, Dnipro, Kharkiv, Kyiv city
	for _, expected := range expectedSubdivisions {
		if !provider.IsSubdivisionSupported(expected) {
			t.Errorf("Expected subdivision %s to be supported", expected)
		}
	}

	// Test categories
	expectedCategories := []string{"national", "orthodox", "memorial", "professional", "regional", "cultural"}
	for _, expected := range expectedCategories {
		if !provider.IsCategorySupported(expected) {
			t.Errorf("Expected category %s to be supported", expected)
		}
	}

	// Test languages
	languages := provider.GetLanguages()
	expectedLanguages := []string{"uk", "en", "ru"}
	if len(languages) != len(expectedLanguages) {
		t.Errorf("Expected %d languages, got %d", len(expectedLanguages), len(languages))
	}
}

func TestUAProvider_LoadHolidays2024(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Should have multiple holidays
	if len(holidays) == 0 {
		t.Error("Expected holidays to be loaded for 2024")
	}

	// Test specific holidays for 2024
	testCases := []struct {
		name     string
		date     time.Time
		category string
	}{
		{
			name:     "New Year's Day",
			date:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			category: "national",
		},
		{
			name:     "Orthodox Christmas",
			date:     time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC),
			category: "orthodox",
		},
		{
			name:     "International Women's Day",
			date:     time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC),
			category: "national",
		},
		{
			name:     "Labor Day",
			date:     time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			category: "national",
		},
		{
			name:     "Victory Day",
			date:     time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
			category: "memorial",
		},
		{
			name:     "Constitution Day",
			date:     time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC),
			category: "national",
		},
		{
			name:     "Day of Ukrainian Statehood",
			date:     time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC),
			category: "national",
		},
		{
			name:     "Independence Day",
			date:     time.Date(2024, 8, 24, 0, 0, 0, 0, time.UTC),
			category: "national",
		},
		{
			name:     "Defenders Day",
			date:     time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC),
			category: "memorial",
		},
		{
			name:     "Day of Ukrainian Language",
			date:     time.Date(2024, 11, 9, 0, 0, 0, 0, time.UTC),
			category: "cultural",
		},
		{
			name:     "Day of Dignity and Freedom",
			date:     time.Date(2024, 11, 21, 0, 0, 0, 0, time.UTC),
			category: "memorial",
		},
		{
			name:     "Holodomor Remembrance Day",
			date:     time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC),
			category: "memorial",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			holiday, exists := holidays[tc.date]
			if !exists {
				t.Errorf("Expected holiday %s on %s to exist", tc.name, tc.date.Format("2006-01-02"))
				return
			}

			if holiday.Category != tc.category {
				t.Errorf("Expected category %s for %s, got %s", tc.category, tc.name, holiday.Category)
			}

			// Check that Ukrainian language is present
			if holiday.Languages["uk"] == "" {
				t.Errorf("Expected Ukrainian language name for %s", tc.name)
			}

			// Check that English language is present
			if holiday.Languages["en"] == "" {
				t.Errorf("Expected English language name for %s", tc.name)
			}
		})
	}
}

func TestUAProvider_OrthodoxEaster(t *testing.T) {
	provider := NewUAProvider()

	// Test Orthodox Easter calculations for known years
	testCases := []struct {
		year         int
		expectedDate time.Time
	}{
		{2024, time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC)},  // Orthodox Easter 2024
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)}, // Orthodox Easter 2025
		{2026, time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)}, // Orthodox Easter 2026
	}

	for _, tc := range testCases {
		t.Run(string(rune(tc.year)), func(t *testing.T) {
			holidays := provider.LoadHolidays(tc.year)

			// Find Orthodox Easter
			found := false
			for date, holiday := range holidays {
				if holiday.Languages["en"] == "Orthodox Easter" {
					found = true
					if !date.Equal(tc.expectedDate) {
						t.Errorf("Expected Orthodox Easter on %s for year %d, got %s",
							tc.expectedDate.Format("2006-01-02"), tc.year, date.Format("2006-01-02"))
					}
					break
				}
			}

			if !found {
				t.Errorf("Orthodox Easter not found for year %d", tc.year)
			}
		})
	}
}

func TestUAProvider_RelatedOrthodoxHolidays(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Find Orthodox Easter to calculate related holidays
	var easter time.Time
	for date, holiday := range holidays {
		if holiday.Languages["en"] == "Orthodox Easter" {
			easter = date
			break
		}
	}

	if easter.IsZero() {
		t.Fatal("Orthodox Easter not found for 2024")
	}

	// Test Palm Sunday (1 week before Easter)
	palmSunday := easter.AddDate(0, 0, -7)
	if holiday, exists := holidays[palmSunday]; !exists {
		t.Error("Palm Sunday should exist 1 week before Orthodox Easter")
	} else if holiday.Languages["uk"] != "Вербна неділя" {
		t.Errorf("Expected Ukrainian name 'Вербна неділя' for Palm Sunday, got %s", holiday.Languages["uk"])
	}

	// Test Trinity Sunday (49 days after Easter)
	trinity := easter.AddDate(0, 0, 49)
	if holiday, exists := holidays[trinity]; !exists {
		t.Error("Trinity Sunday should exist 49 days after Orthodox Easter")
	} else if holiday.Languages["uk"] != "Трійця" {
		t.Errorf("Expected Ukrainian name 'Трійця' for Trinity Sunday, got %s", holiday.Languages["uk"])
	}
}

func TestUAProvider_YearSpecificHolidays(t *testing.T) {
	provider := NewUAProvider()

	testCases := []struct {
		year        int
		holidayName string
		shouldExist bool
		reason      string
	}{
		{2013, "Day of Dignity and Freedom", false, "introduced in 2014"},
		{2015, "Day of Dignity and Freedom", true, "should exist from 2014"},
		{2014, "Defenders Day", false, "introduced in 2015"},
		{2015, "Defenders Day", true, "should exist from 2015"},
		{2018, "Day of Ukrainian Language", false, "introduced in 2019"},
		{2019, "Day of Ukrainian Language", true, "should exist from 2019"},
		{2020, "Day of Ukrainian Statehood", false, "introduced in 2021"},
		{2021, "Day of Ukrainian Statehood", true, "should exist from 2021"},
	}

	for _, tc := range testCases {
		t.Run(string(rune(tc.year))+"_"+tc.holidayName, func(t *testing.T) {
			holidays := provider.LoadHolidays(tc.year)

			found := false
			for _, holiday := range holidays {
				if holiday.Languages["en"] == tc.holidayName {
					found = true
					break
				}
			}

			if found != tc.shouldExist {
				if tc.shouldExist {
					t.Errorf("Expected %s to exist in %d (%s), but it was not found",
						tc.holidayName, tc.year, tc.reason)
				} else {
					t.Errorf("Expected %s not to exist in %d (%s), but it was found",
						tc.holidayName, tc.year, tc.reason)
				}
			}
		})
	}
}

func TestUAProvider_CulturalHolidays(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Old New Year
	oldNewYear := time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[oldNewYear]; !exists {
		t.Error("Old New Year should exist on January 14")
	} else {
		if holiday.Languages["uk"] != "Старий Новий рік" {
			t.Errorf("Expected Ukrainian name 'Старий Новий рік' for Old New Year, got %s", holiday.Languages["uk"])
		}
		if holiday.Category != "cultural" {
			t.Errorf("Expected category 'cultural' for Old New Year, got %s", holiday.Category)
		}
	}
}

func TestUAProvider_MemorialHolidays(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Test Holodomor Remembrance Day
	holodomorDay := time.Date(2024, 11, 25, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[holodomorDay]; !exists {
		t.Error("Holodomor Remembrance Day should exist on November 25")
	} else {
		if holiday.Languages["uk"] != "День пам'яті жертв голодомору" {
			t.Errorf("Expected Ukrainian name for Holodomor Remembrance Day, got %s", holiday.Languages["uk"])
		}
		if holiday.Category != "memorial" {
			t.Errorf("Expected category 'memorial' for Holodomor Remembrance Day, got %s", holiday.Category)
		}
	}

	// Test Day of Remembrance of Victims of Political Repressions
	repressionDay := time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[repressionDay]; !exists {
		t.Error("Day of Remembrance of Victims of Political Repressions should exist on May 19")
	} else {
		if holiday.Category != "memorial" {
			t.Errorf("Expected category 'memorial' for Day of Remembrance, got %s", holiday.Category)
		}
	}
}

func TestUAProvider_LanguageSupport(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that all holidays have proper language support
	for date, holiday := range holidays {
		// Every holiday should have Ukrainian name
		if holiday.Languages["uk"] == "" {
			t.Errorf("Holiday on %s is missing Ukrainian name", date.Format("2006-01-02"))
		}

		// Every holiday should have English name
		if holiday.Languages["en"] == "" {
			t.Errorf("Holiday on %s is missing English name", date.Format("2006-01-02"))
		}

		// Most holidays should have Russian name (historical context)
		if holiday.Languages["ru"] == "" {
			t.Logf("Holiday on %s is missing Russian name: %s", date.Format("2006-01-02"), holiday.Languages["en"])
		}
	}
}

func TestUAProvider_CategoryDistribution(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Count holidays by category
	categoryCount := make(map[string]int)
	for _, holiday := range holidays {
		categoryCount[holiday.Category]++
	}

	// Should have holidays in multiple categories
	expectedCategories := []string{"national", "orthodox", "memorial", "cultural"}
	for _, category := range expectedCategories {
		if count := categoryCount[category]; count == 0 {
			t.Errorf("Expected at least one holiday in category %s, got %d", category, count)
		}
	}

	// Log category distribution for verification
	t.Logf("Holiday category distribution for 2024:")
	for category, count := range categoryCount {
		t.Logf("  %s: %d holidays", category, count)
	}
}

func TestUAProvider_NoConflictingDates(t *testing.T) {
	provider := NewUAProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that Defenders Day and Day of Ukrainian Cossacks don't conflict in 2015+
	defendersDay := time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC)

	holidaysOnDate := 0
	var holidayNames []string
	for date, holiday := range holidays {
		if date.Equal(defendersDay) {
			holidaysOnDate++
			holidayNames = append(holidayNames, holiday.Languages["en"])
		}
	}

	// Should only be Defenders Day on October 14 for years 2015+
	if holidaysOnDate != 1 {
		t.Errorf("Expected exactly 1 holiday on October 14, 2024, got %d: %v", holidaysOnDate, holidayNames)
	}

	if len(holidayNames) > 0 && holidayNames[0] != "Defenders Day" {
		t.Errorf("Expected Defenders Day on October 14, 2024, got %s", holidayNames[0])
	}
}
