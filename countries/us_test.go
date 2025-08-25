package countries

import (
	"testing"
	"time"
)

func TestUSProvider_FederalHolidays(t *testing.T) {
	provider := NewUSProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[newYear]; !exists {
		t.Error("New Year's Day should exist")
	} else {
		if holiday.Name != "New Year's Day" {
			t.Errorf("Expected 'New Year's Day', got '%s'", holiday.Name)
		}
		if holiday.Category != "federal" {
			t.Errorf("Expected category 'federal', got '%s'", holiday.Category)
		}
	}

	// Test Independence Day
	independenceDay := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[independenceDay]; !exists {
		t.Error("Independence Day should exist")
	} else {
		if holiday.Name != "Independence Day" {
			t.Errorf("Expected 'Independence Day', got '%s'", holiday.Name)
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

	// Test Veterans Day
	veteransDay := time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[veteransDay]; !exists {
		t.Error("Veterans Day should exist")
	} else {
		if holiday.Name != "Veterans Day" {
			t.Errorf("Expected 'Veterans Day', got '%s'", holiday.Name)
		}
	}
}

func TestUSProvider_VariableHolidays(t *testing.T) {
	provider := NewUSProvider()
	holidays := provider.LoadHolidays(2024)

	// MLK Day - 3rd Monday in January 2024 (January 15)
	mlkDay := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[mlkDay]; !exists {
		t.Error("Martin Luther King Jr. Day should exist")
	} else {
		if holiday.Name != "Martin Luther King Jr. Day" {
			t.Errorf("Expected 'Martin Luther King Jr. Day', got '%s'", holiday.Name)
		}
	}

	// Presidents' Day - 3rd Monday in February 2024 (February 19)
	presidentsDay := time.Date(2024, 2, 19, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[presidentsDay]; !exists {
		t.Error("Presidents' Day should exist")
	} else {
		if holiday.Name != "Presidents' Day" {
			t.Errorf("Expected 'Presidents' Day', got '%s'", holiday.Name)
		}
	}

	// Memorial Day - Last Monday in May 2024 (May 27)
	memorialDay := time.Date(2024, 5, 27, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[memorialDay]; !exists {
		t.Error("Memorial Day should exist")
	} else {
		if holiday.Name != "Memorial Day" {
			t.Errorf("Expected 'Memorial Day', got '%s'", holiday.Name)
		}
	}

	// Labor Day - 1st Monday in September 2024 (September 2)
	laborDay := time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[laborDay]; !exists {
		t.Error("Labor Day should exist")
	} else {
		if holiday.Name != "Labor Day" {
			t.Errorf("Expected 'Labor Day', got '%s'", holiday.Name)
		}
	}

	// Thanksgiving - 4th Thursday in November 2024 (November 28)
	thanksgiving := time.Date(2024, 11, 28, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[thanksgiving]; !exists {
		t.Error("Thanksgiving Day should exist")
	} else {
		if holiday.Name != "Thanksgiving Day" {
			t.Errorf("Expected 'Thanksgiving Day', got '%s'", holiday.Name)
		}
	}
}

func TestUSProvider_JuneteenthHistory(t *testing.T) {
	provider := NewUSProvider()

	// Test that Juneteenth doesn't exist before 2021
	holidays2020 := provider.LoadHolidays(2020)
	juneteenth2020 := time.Date(2020, 6, 19, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2020[juneteenth2020]; exists {
		t.Error("Juneteenth should not exist in 2020")
	}

	// Test that Juneteenth exists from 2021 onwards
	holidays2021 := provider.LoadHolidays(2021)
	juneteenth2021 := time.Date(2021, 6, 19, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays2021[juneteenth2021]; !exists {
		t.Error("Juneteenth should exist in 2021")
	} else {
		if holiday.Name != "Juneteenth" {
			t.Errorf("Expected 'Juneteenth', got '%s'", holiday.Name)
		}
		if holiday.Category != "federal" {
			t.Errorf("Expected category 'federal', got '%s'", holiday.Category)
		}
	}

	// Test 2024 has Juneteenth
	holidays2024 := provider.LoadHolidays(2024)
	juneteenth2024 := time.Date(2024, 6, 19, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2024[juneteenth2024]; !exists {
		t.Error("Juneteenth should exist in 2024")
	}
}

func TestUSProvider_MLKHistory(t *testing.T) {
	provider := NewUSProvider()

	// Test that MLK Day doesn't exist before 1983
	holidays1982 := provider.LoadHolidays(1982)
	mlk1982 := NthWeekdayOfMonth(1982, 1, time.Monday, 3)
	if _, exists := holidays1982[mlk1982]; exists {
		t.Error("MLK Day should not exist in 1982")
	}

	// Test that MLK Day exists from 1983 onwards
	holidays1983 := provider.LoadHolidays(1983)
	mlk1983 := NthWeekdayOfMonth(1983, 1, time.Monday, 3)
	if holiday, exists := holidays1983[mlk1983]; !exists {
		t.Error("MLK Day should exist in 1983")
	} else {
		if holiday.Name != "Martin Luther King Jr. Day" {
			t.Errorf("Expected 'Martin Luther King Jr. Day', got '%s'", holiday.Name)
		}
	}
}

func TestUSProvider_StateHolidays(t *testing.T) {
	provider := NewUSProvider()

	// Test California state holidays
	caHolidays := provider.GetStateHolidays(2024, []string{"CA"})
	chavezDay := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	if holiday, exists := caHolidays[chavezDay]; !exists {
		t.Error("Cesar Chavez Day should exist for California")
	} else {
		if holiday.Name != "Cesar Chavez Day" {
			t.Errorf("Expected 'Cesar Chavez Day', got '%s'", holiday.Name)
		}
	}

	// Test Texas state holidays
	txHolidays := provider.GetStateHolidays(2024, []string{"TX"})
	texasIndependence := time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC)
	if holiday, exists := txHolidays[texasIndependence]; !exists {
		t.Error("Texas Independence Day should exist for Texas")
	} else {
		if holiday.Name != "Texas Independence Day" {
			t.Errorf("Expected 'Texas Independence Day', got '%s'", holiday.Name)
		}
	}

	// Test Massachusetts state holidays
	maHolidays := provider.GetStateHolidays(2024, []string{"MA"})
	patriotsDay := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC) // 3rd Monday in April 2024
	if holiday, exists := maHolidays[patriotsDay]; !exists {
		t.Error("Patriots' Day should exist for Massachusetts")
	} else {
		if holiday.Name != "Patriots' Day" {
			t.Errorf("Expected 'Patriots' Day', got '%s'", holiday.Name)
		}
	}
}

func TestUSProvider_MultipleStates(t *testing.T) {
	provider := NewUSProvider()

	// Test multiple states
	stateHolidays := provider.GetStateHolidays(2024, []string{"CA", "TX", "MA"})
	
	// Should have holidays from all three states
	expectedCount := 3 // Chavez Day, Texas Independence, Patriots' Day
	if len(stateHolidays) != expectedCount {
		t.Errorf("Expected %d state holidays, got %d", expectedCount, len(stateHolidays))
	}
}

func TestUSProvider_ProviderInfo(t *testing.T) {
	provider := NewUSProvider()

	// Test country code
	if provider.GetCountryCode() != "US" {
		t.Errorf("Expected country code 'US', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (should include all 50 states + DC + territories)
	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) < 51 { // At least 50 states + DC
		t.Errorf("Expected at least 51 subdivisions, got %d", len(subdivisions))
	}

	// Test categories
	expectedCategories := []string{"federal", "state", "religious", "observance"}
	categories := provider.GetSupportedCategories()
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestUSProvider_FederalHolidayCount(t *testing.T) {
	provider := NewUSProvider()

	// Test 2024 - should have 11 federal holidays (including Juneteenth)
	holidays2024 := provider.LoadHolidays(2024)
	federalCount := 0
	for _, holiday := range holidays2024 {
		if holiday.Category == "federal" {
			federalCount++
		}
	}
	
	expectedFederal := 11 // All current federal holidays including Juneteenth
	if federalCount != expectedFederal {
		t.Errorf("Expected %d federal holidays in 2024, got %d", expectedFederal, federalCount)
	}

	// Test 2020 - should have 10 federal holidays (no Juneteenth)
	holidays2020 := provider.LoadHolidays(2020)
	federalCount2020 := 0
	for _, holiday := range holidays2020 {
		if holiday.Category == "federal" {
			federalCount2020++
		}
	}
	
	expectedFederal2020 := 10 // Federal holidays before Juneteenth
	if federalCount2020 != expectedFederal2020 {
		t.Errorf("Expected %d federal holidays in 2020, got %d", expectedFederal2020, federalCount2020)
	}

	// Test 1982 - should have 9 federal holidays (no MLK Day, no Juneteenth)
	holidays1982 := provider.LoadHolidays(1982)
	federalCount1982 := 0
	for _, holiday := range holidays1982 {
		if holiday.Category == "federal" {
			federalCount1982++
		}
	}
	
	expectedFederal1982 := 9 // Federal holidays before MLK Day
	if federalCount1982 != expectedFederal1982 {
		t.Errorf("Expected %d federal holidays in 1982, got %d", expectedFederal1982, federalCount1982)
	}
}

func TestUSProvider_BilingualSupport(t *testing.T) {
	provider := NewUSProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that major holidays have Spanish translations
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[newYear]; exists {
		if holiday.Languages["es"] != "Año Nuevo" {
			t.Errorf("Expected Spanish translation 'Año Nuevo', got '%s'", holiday.Languages["es"])
		}
	}

	independence := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[independence]; exists {
		if holiday.Languages["es"] != "Día de la Independencia" {
			t.Errorf("Expected Spanish translation 'Día de la Independencia', got '%s'", holiday.Languages["es"])
		}
	}
}
