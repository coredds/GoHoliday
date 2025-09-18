package countries

import (
	"testing"
	"time"
)

func TestIEProvider_BasicInfo(t *testing.T) {
	provider := NewIEProvider()

	if provider.GetCountryCode() != "IE" {
		t.Errorf("Expected country code IE, got %s", provider.GetCountryCode())
	}

	if provider.GetCountryName() != "Ireland" {
		t.Errorf("Expected country name Ireland, got %s", provider.GetCountryName())
	}

	subdivisions := provider.GetSubdivisions()
	if len(subdivisions) != 30 { // 26 counties + 4 provinces
		t.Errorf("Expected 30 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetCategories()
	expectedCategories := []string{"public", "bank", "religious", "national", "cultural"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestIEProvider_LoadHolidays2024(t *testing.T) {
	provider := NewIEProvider()
	holidays := provider.LoadHolidays(2024)

	if len(holidays) == 0 {
		t.Error("Expected holidays to be loaded for 2024")
	}

	// Test key Irish holidays for 2024
	expectedHolidays := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day", "public"},
		{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), "Saint Brigid's Day", "cultural"},
		{time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC), "Saint Brigid's Day (Public Holiday)", "public"}, // First Monday in Feb 2024
		{time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC), "Saint Patrick's Day", "national"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Good Friday", "religious"},
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Easter Monday", "public"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "May Day", "cultural"},
		{time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC), "June Bank Holiday", "bank"}, // First Monday in June
		{time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), "Lughnasadh", "cultural"},
		{time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC), "August Bank Holiday", "bank"}, // First Monday in August
		{time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC), "October Bank Holiday", "bank"}, // Last Monday in October
		{time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC), "Samhain", "cultural"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day", "religious"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Saint Stephen's Day", "religious"},
	}

	for _, expected := range expectedHolidays {
		holiday, exists := holidays[expected.date]
		if !exists {
			t.Errorf("Expected holiday on %s (%s)", expected.date.Format("2006-01-02"), expected.name)
			continue
		}

		if holiday.Name != expected.name {
			t.Errorf("Expected holiday name %s, got %s", expected.name, holiday.Name)
		}

		if holiday.Category != expected.category {
			t.Errorf("Expected category %s for %s, got %s", expected.category, expected.name, holiday.Category)
		}
	}
}

func TestIEProvider_SaintPatricksDay(t *testing.T) {
	provider := NewIEProvider()
	holidays := provider.LoadHolidays(2024)

	stPatricksDay := time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC)
	holiday, exists := holidays[stPatricksDay]
	
	if !exists {
		t.Fatal("Saint Patrick's Day should exist")
	}

	if holiday.Name != "Saint Patrick's Day" {
		t.Errorf("Expected 'Saint Patrick's Day', got %s", holiday.Name)
	}

	if holiday.Category != "national" {
		t.Errorf("Expected category 'national', got %s", holiday.Category)
	}

	// Check Irish language name
	if holiday.Languages["ga"] != "Lá Fhéile Pádraig" {
		t.Errorf("Expected Irish name 'Lá Fhéile Pádraig', got %s", holiday.Languages["ga"])
	}
}

func TestIEProvider_BankHolidays(t *testing.T) {
	provider := NewIEProvider()

	testCases := []struct {
		year         int
		month        int
		expectedDate time.Time
		name         string
		isLast       bool // true for last Monday, false for first Monday
	}{
		{2024, 6, time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC), "June Bank Holiday", false},
		{2024, 8, time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC), "August Bank Holiday", false},
		{2024, 10, time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC), "October Bank Holiday", true},
		{2025, 6, time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC), "June Bank Holiday", false},
		{2025, 8, time.Date(2025, 8, 4, 0, 0, 0, 0, time.UTC), "August Bank Holiday", false},
		{2025, 10, time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC), "October Bank Holiday", true},
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		holiday, exists := holidays[tc.expectedDate]
		if !exists {
			t.Errorf("Expected %s on %s in %d", tc.name, tc.expectedDate.Format("2006-01-02"), tc.year)
			continue
		}

		if holiday.Name != tc.name {
			t.Errorf("Expected holiday name %s, got %s", tc.name, holiday.Name)
		}

		if holiday.Category != "bank" {
			t.Errorf("Expected category 'bank', got %s", holiday.Category)
		}

		// Verify it's actually a Monday
		if tc.expectedDate.Weekday() != time.Monday {
			t.Errorf("Bank holiday %s should be on Monday, got %s", tc.name, tc.expectedDate.Weekday())
		}
	}
}

func TestIEProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewIEProvider()

	testYears := []int{2023, 2024, 2025}

	for _, year := range testYears {
		holidays := provider.LoadHolidays(year)
		easter := EasterSunday(year)

		// Good Friday (2 days before Easter)
		goodFriday := easter.AddDate(0, 0, -2)
		if holiday, exists := holidays[goodFriday]; !exists {
			t.Errorf("Good Friday should exist in %d", year)
		} else {
			if holiday.Name != "Good Friday" {
				t.Errorf("Expected 'Good Friday', got %s", holiday.Name)
			}
			if holiday.Category != "religious" {
				t.Errorf("Expected category 'religious', got %s", holiday.Category)
			}
		}

		// Easter Monday (1 day after Easter)
		easterMonday := easter.AddDate(0, 0, 1)
		if holiday, exists := holidays[easterMonday]; !exists {
			t.Errorf("Easter Monday should exist in %d", year)
		} else {
			if holiday.Name != "Easter Monday" {
				t.Errorf("Expected 'Easter Monday', got %s", holiday.Name)
			}
			if holiday.Category != "public" {
				t.Errorf("Expected category 'public', got %s", holiday.Category)
			}
		}
	}
}

func TestIEProvider_CelticFestivals(t *testing.T) {
	provider := NewIEProvider()
	holidays := provider.LoadHolidays(2024)

	celticFestivals := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), "Saint Brigid's Day", "cultural"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "May Day", "cultural"},
		{time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), "Lughnasadh", "cultural"},
		{time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC), "Samhain", "cultural"},
	}

	for _, festival := range celticFestivals {
		holiday, exists := holidays[festival.date]
		if !exists {
			t.Errorf("Expected Celtic festival %s on %s", festival.name, festival.date.Format("2006-01-02"))
			continue
		}

		if holiday.Name != festival.name {
			t.Errorf("Expected festival name %s, got %s", festival.name, holiday.Name)
		}

		if holiday.Category != festival.category {
			t.Errorf("Expected category %s for %s, got %s", festival.category, festival.name, holiday.Category)
		}
	}
}

func TestIEProvider_BrigidsPublicHoliday(t *testing.T) {
	provider := NewIEProvider()

	// Test Saint Brigid's Day public holiday (introduced in 2023)
	t.Run("Before 2023", func(t *testing.T) {
		holidays2022 := provider.LoadHolidays(2022)
		
		// Should have cultural Brigid's Day but not public holiday
		brigidsDay := time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)
		if holiday, exists := holidays2022[brigidsDay]; exists {
			if holiday.Category == "public" {
				t.Error("Saint Brigid's Day should not be a public holiday in 2022")
			}
		}

		// Check that no public Brigid's holiday exists
		for _, holiday := range holidays2022 {
			if holiday.Name == "Saint Brigid's Day (Public Holiday)" {
				t.Error("Public Saint Brigid's Day should not exist in 2022")
			}
		}
	})

	t.Run("From 2023 onwards", func(t *testing.T) {
		testYears := []int{2023, 2024, 2025}
		
		for _, year := range testYears {
			holidays := provider.LoadHolidays(year)
			
			// Should have public Brigid's holiday
			found := false
			for _, holiday := range holidays {
				if holiday.Name == "Saint Brigid's Day (Public Holiday)" && holiday.Category == "public" {
					found = true
					// Should be on a Monday
					if holiday.Date.Weekday() != time.Monday {
						t.Errorf("Public Saint Brigid's Day should be on Monday in %d, got %s", year, holiday.Date.Weekday())
					}
					break
				}
			}
			
			if !found {
				t.Errorf("Public Saint Brigid's Day should exist in %d", year)
			}
		}
	})
}

func TestIEProvider_LanguageSupport(t *testing.T) {
	provider := NewIEProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that holidays have both English and Irish Gaelic names
	christmas := holidays[time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)]
	if christmas == nil {
		t.Fatal("Christmas Day should exist")
	}

	if christmas.Languages["en"] != "Christmas Day" {
		t.Errorf("Expected English name 'Christmas Day', got %s", christmas.Languages["en"])
	}

	if christmas.Languages["ga"] != "Lá na Nollag" {
		t.Errorf("Expected Irish name 'Lá na Nollag', got %s", christmas.Languages["ga"])
	}
}

func TestIEProvider_CategoryDistribution(t *testing.T) {
	provider := NewIEProvider()
	holidays := provider.LoadHolidays(2024)

	categories := make(map[string]int)
	for _, holiday := range holidays {
		categories[holiday.Category]++
	}

	// Should have holidays in all main categories
	expectedCategories := []string{"public", "bank", "religious", "national", "cultural"}
	for _, category := range expectedCategories {
		if count, exists := categories[category]; !exists || count == 0 {
			t.Errorf("Expected holidays in category %s", category)
		}
	}

	// Should have multiple bank holidays
	if categories["bank"] < 3 {
		t.Errorf("Expected at least 3 bank holidays, got %d", categories["bank"])
	}

	// Should have cultural holidays (Celtic festivals)
	if categories["cultural"] < 4 {
		t.Errorf("Expected at least 4 cultural holidays, got %d", categories["cultural"])
	}

	t.Logf("Holiday distribution: %+v", categories)
}
