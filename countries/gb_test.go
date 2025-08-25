package countries

import (
	"testing"
	"time"
)

func TestGBProvider_BasicHolidays(t *testing.T) {
	provider := NewGBProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[newYear]; !exists {
		t.Error("New Year's Day should exist")
	} else {
		if holiday.Name != "New Year's Day" {
			t.Errorf("Expected 'New Year's Day', got '%s'", holiday.Name)
		}
		if holiday.Category != "public" {
			t.Errorf("Expected category 'public', got '%s'", holiday.Category)
		}
	}

	// Test Christmas Day
	christmas := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[christmas]; !exists {
		t.Error("Christmas Day should exist")
	} else {
		if holiday.Name != "Christmas Day" {
			t.Errorf("Expected 'Christmas Day', got '%s'", holiday.Name)
		}
	}

	// Test Boxing Day
	boxingDay := time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[boxingDay]; !exists {
		t.Error("Boxing Day should exist")
	} else {
		if holiday.Name != "Boxing Day" {
			t.Errorf("Expected 'Boxing Day', got '%s'", holiday.Name)
		}
	}
}

func TestGBProvider_EasterHolidays(t *testing.T) {
	provider := NewGBProvider()
	holidays := provider.LoadHolidays(2024)

	// Easter Sunday 2024 is March 31
	easter := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

	// Test Good Friday (March 29, 2024)
	goodFriday := time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[goodFriday]; !exists {
		t.Error("Good Friday should exist")
	} else {
		if holiday.Name != "Good Friday" {
			t.Errorf("Expected 'Good Friday', got '%s'", holiday.Name)
		}
	}

	// Test Easter Monday (April 1, 2024)
	easterMonday := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[easterMonday]; !exists {
		t.Error("Easter Monday should exist")
	} else {
		if holiday.Name != "Easter Monday" {
			t.Errorf("Expected 'Easter Monday', got '%s'", holiday.Name)
		}
	}

	// Verify Easter calculation is correct
	calculatedEaster := EasterSunday(2024)
	if !calculatedEaster.Equal(easter) {
		t.Errorf("Easter calculation incorrect. Expected %v, got %v", easter, calculatedEaster)
	}
}

func TestGBProvider_BankHolidays(t *testing.T) {
	provider := NewGBProvider()
	holidays := provider.LoadHolidays(2024)

	// Early May Bank Holiday - 1st Monday in May 2024 (May 6)
	earlyMay := time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[earlyMay]; !exists {
		t.Error("Early May Bank Holiday should exist")
	} else {
		if holiday.Name != "Early May Bank Holiday" {
			t.Errorf("Expected 'Early May Bank Holiday', got '%s'", holiday.Name)
		}
		if holiday.Category != "bank" {
			t.Errorf("Expected category 'bank', got '%s'", holiday.Category)
		}
	}

	// Spring Bank Holiday - Last Monday in May 2024 (May 27)
	springBank := time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[springBank]; !exists {
		t.Error("Spring Bank Holiday should exist")
	} else {
		if holiday.Name != "Spring Bank Holiday" {
			t.Errorf("Expected 'Spring Bank Holiday', got '%s'", holiday.Name)
		}
	}

	// Summer Bank Holiday - Last Monday in August 2024 (August 26)
	summerBank := time.Date(2024, 8, 26, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[summerBank]; !exists {
		t.Error("Summer Bank Holiday should exist")
	} else {
		if holiday.Name != "Summer Bank Holiday" {
			t.Errorf("Expected 'Summer Bank Holiday', got '%s'", holiday.Name)
		}
	}
}

func TestGBProvider_SpecialHolidays(t *testing.T) {
	provider := NewGBProvider()

	// Test 2022 Platinum Jubilee
	holidays2022 := provider.LoadHolidays(2022)
	platinumJubilee := time.Date(2022, 6, 3, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays2022[platinumJubilee]; !exists {
		t.Error("Platinum Jubilee should exist in 2022")
	} else {
		if holiday.Name != "Platinum Jubilee" {
			t.Errorf("Expected 'Platinum Jubilee', got '%s'", holiday.Name)
		}
	}

	// Test 2023 Coronation
	holidays2023 := provider.LoadHolidays(2023)
	coronation := time.Date(2023, 5, 8, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays2023[coronation]; !exists {
		t.Error("Coronation should exist in 2023")
	} else {
		if holiday.Name != "Coronation of King Charles III" {
			t.Errorf("Expected 'Coronation of King Charles III', got '%s'", holiday.Name)
		}
	}

	// Test normal year has no special holidays
	holidays2024 := provider.LoadHolidays(2024)
	if len(holidays2024) > 8 { // Should have exactly 8 standard holidays
		t.Errorf("2024 should have exactly 8 holidays, got %d", len(holidays2024))
	}
}

func TestGBProvider_RegionalHolidays(t *testing.T) {
	provider := NewGBProvider()

	// Test Scotland
	scotlandHolidays := provider.GetRegionalHolidays(2024, []string{"SCT"})
	stAndrews := time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC)
	if holiday, exists := scotlandHolidays[stAndrews]; !exists {
		t.Error("St. Andrew's Day should exist for Scotland")
	} else {
		if holiday.Name != "St. Andrew's Day" {
			t.Errorf("Expected 'St. Andrew's Day', got '%s'", holiday.Name)
		}
	}

	// Test Wales
	walesHolidays := provider.GetRegionalHolidays(2024, []string{"WLS"})
	stDavids := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := walesHolidays[stDavids]; !exists {
		t.Error("St. David's Day should exist for Wales")
	} else {
		if holiday.Name != "St. David's Day" {
			t.Errorf("Expected 'St. David's Day', got '%s'", holiday.Name)
		}
	}

	// Test Northern Ireland
	nirHolidays := provider.GetRegionalHolidays(2024, []string{"NIR"})
	stPatricks := time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC)
	if holiday, exists := nirHolidays[stPatricks]; !exists {
		t.Error("St. Patrick's Day should exist for Northern Ireland")
	} else {
		if holiday.Name != "St. Patrick's Day" {
			t.Errorf("Expected 'St. Patrick's Day', got '%s'", holiday.Name)
		}
	}

	battleOfBoyne := time.Date(2024, 7, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := nirHolidays[battleOfBoyne]; !exists {
		t.Error("Battle of the Boyne should exist for Northern Ireland")
	} else {
		if holiday.Name != "Battle of the Boyne" {
			t.Errorf("Expected 'Battle of the Boyne', got '%s'", holiday.Name)
		}
	}
}

func TestGBProvider_MultipleRegions(t *testing.T) {
	provider := NewGBProvider()

	// Test multiple regions
	regionalHolidays := provider.GetRegionalHolidays(2024, []string{"SCT", "WLS", "NIR"})

	// Should have all regional holidays
	expectedCount := 4 // St. Andrew's, St. David's, St. Patrick's, Battle of Boyne
	if len(regionalHolidays) != expectedCount {
		t.Errorf("Expected %d regional holidays, got %d", expectedCount, len(regionalHolidays))
	}
}

func TestGBProvider_ProviderInfo(t *testing.T) {
	provider := NewGBProvider()

	// Test country code
	if provider.GetCountryCode() != "GB" {
		t.Errorf("Expected country code 'GB', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions
	expectedSubdivisions := []string{"ENG", "SCT", "WLS", "NIR"}
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != len(expectedSubdivisions) {
		t.Errorf("Expected %d subdivisions, got %d", len(expectedSubdivisions), len(subdivisions))
	}

	// Test categories
	expectedCategories := []string{"public", "bank", "government"}
	categories := provider.GetSupportedCategories()
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestGBProvider_HolidayDateCalculations(t *testing.T) {
	provider := NewGBProvider()

	// Test multiple years to ensure calculations are correct
	testYears := []int{2023, 2024, 2025, 2026}

	for _, year := range testYears {
		holidays := provider.LoadHolidays(year)

		// Verify we have the expected number of standard holidays
		expectedStandardHolidays := 8
		standardCount := 0
		for _, holiday := range holidays {
			if holiday.Category == "public" || holiday.Category == "bank" {
				standardCount++
			}
		}

		if standardCount < expectedStandardHolidays {
			t.Errorf("Year %d: Expected at least %d standard holidays, got %d",
				year, expectedStandardHolidays, standardCount)
		}

		// Verify Easter-based holidays are calculated correctly
		easter := EasterSunday(year)
		goodFriday := easter.AddDate(0, 0, -2)
		easterMonday := easter.AddDate(0, 0, 1)

		if _, exists := holidays[goodFriday]; !exists {
			t.Errorf("Year %d: Good Friday missing for Easter %v", year, easter)
		}

		if _, exists := holidays[easterMonday]; !exists {
			t.Errorf("Year %d: Easter Monday missing for Easter %v", year, easter)
		}
	}
}
