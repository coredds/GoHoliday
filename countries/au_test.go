package countries

import (
	"testing"
	"time"
)

func TestAUProvider_BasicHolidays(t *testing.T) {
	provider := NewAUProvider()
	year := 2024
	holidays := provider.LoadHolidays(year)

	expectedHolidays := map[string]time.Time{
		"New Year's Day":   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"Australia Day":    time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
		"Good Friday":      time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
		"Easter Saturday":  time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC),
		"Easter Monday":    time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		"ANZAC Day":        time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),
		"Queen's Birthday": time.Date(2024, 6, 10, 0, 0, 0, 0, time.UTC), // 2nd Monday in June
		"Labour Day":       time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC),  // 1st Monday in October
		"Christmas Day":    time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		"Boxing Day":       time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
	}

	for name, expectedDate := range expectedHolidays {
		holiday, exists := holidays[expectedDate]
		if !exists {
			t.Errorf("Expected holiday %s on %s not found", name, expectedDate.Format("2006-01-02"))
			continue
		}
		if holiday.Name != name {
			t.Errorf("Holiday name mismatch. Expected: %s, Got: %s", name, holiday.Name)
		}
	}

	// Verify we have the expected number of base holidays
	if len(holidays) != len(expectedHolidays) {
		t.Errorf("Expected %d holidays, got %d", len(expectedHolidays), len(holidays))
	}
}

func TestAUProvider_EasterCalculations(t *testing.T) {
	provider := NewAUProvider()
	
	testCases := []struct {
		year         int
		expectedGoodFriday   time.Time
		expectedEasterSaturday time.Time
		expectedEasterMonday time.Time
	}{
		{
			year:                2024,
			expectedGoodFriday:  time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			expectedEasterSaturday: time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC),
			expectedEasterMonday: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			year:                2025,
			expectedGoodFriday:  time.Date(2025, 4, 18, 0, 0, 0, 0, time.UTC),
			expectedEasterSaturday: time.Date(2025, 4, 19, 0, 0, 0, 0, time.UTC),
			expectedEasterMonday: time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC),
		},
	}
	
	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		if holiday, exists := holidays[tc.expectedGoodFriday]; !exists || holiday.Name != "Good Friday" {
			t.Errorf("Year %d: Good Friday not found on expected date %s", tc.year, tc.expectedGoodFriday.Format("2006-01-02"))
		}
		
		if holiday, exists := holidays[tc.expectedEasterSaturday]; !exists || holiday.Name != "Easter Saturday" {
			t.Errorf("Year %d: Easter Saturday not found on expected date %s", tc.year, tc.expectedEasterSaturday.Format("2006-01-02"))
		}
		
		if holiday, exists := holidays[tc.expectedEasterMonday]; !exists || holiday.Name != "Easter Monday" {
			t.Errorf("Year %d: Easter Monday not found on expected date %s", tc.year, tc.expectedEasterMonday.Format("2006-01-02"))
		}
	}
}

func TestAUProvider_StateSpecificHolidays(t *testing.T) {
	provider := NewAUProvider()
	year := 2024

	// Test Victoria specific holidays
	vicHolidays := provider.GetStateHolidays(year, []string{"VIC"})
	
	expectedMelbourneCup := time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC) // First Tuesday in November
	if holiday, exists := vicHolidays[expectedMelbourneCup]; !exists || holiday.Name != "Melbourne Cup Day" {
		t.Errorf("Melbourne Cup Day not found on expected date %s", expectedMelbourneCup.Format("2006-01-02"))
	}
	
	expectedVicLabourDay := time.Date(2024, 3, 11, 0, 0, 0, 0, time.UTC) // 2nd Monday in March
	if holiday, exists := vicHolidays[expectedVicLabourDay]; !exists || holiday.Name != "Labour Day" {
		t.Errorf("Victoria Labour Day not found on expected date %s", expectedVicLabourDay.Format("2006-01-02"))
	}

	// Test Western Australia specific holidays
	waHolidays := provider.GetStateHolidays(year, []string{"WA"})
	
	expectedWADay := time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC) // 1st Monday in June
	if holiday, exists := waHolidays[expectedWADay]; !exists || holiday.Name != "Western Australia Day" {
		t.Errorf("Western Australia Day not found on expected date %s", expectedWADay.Format("2006-01-02"))
	}

	// Test Queensland specific holidays
	qldHolidays := provider.GetStateHolidays(year, []string{"QLD"})
	
	expectedQLDLabourDay := time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC) // 1st Monday in May
	if holiday, exists := qldHolidays[expectedQLDLabourDay]; !exists || holiday.Name != "Labour Day" {
		t.Errorf("Queensland Labour Day not found on expected date %s", expectedQLDLabourDay.Format("2006-01-02"))
	}
}

func TestAUProvider_MelbourneCupDay(t *testing.T) {
	provider := NewAUProvider()
	
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC)},  // First Tuesday
		{2025, time.Date(2025, 11, 4, 0, 0, 0, 0, time.UTC)},  // First Tuesday
		{2023, time.Date(2023, 11, 7, 0, 0, 0, 0, time.UTC)},  // First Tuesday
	}
	
	for _, tc := range testCases {
		result := provider.getMelbourneCupDay(tc.year)
		if !result.Equal(tc.expected) {
			t.Errorf("Year %d: Expected Melbourne Cup Day on %s, got %s", 
				tc.year, tc.expected.Format("2006-01-02"), result.Format("2006-01-02"))
		}
		
		// Verify it's actually a Tuesday
		if result.Weekday() != time.Tuesday {
			t.Errorf("Year %d: Melbourne Cup Day should be a Tuesday, got %s", tc.year, result.Weekday())
		}
	}
}

func TestAUProvider_Seasons(t *testing.T) {
	provider := NewAUProvider()
	seasons := provider.GetSeasons(2024)
	
	expectedSeasons := map[string][2]time.Time{
		"Summer": {
			time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		"Autumn": {
			time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
		},
		"Winter": {
			time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC),
		},
		"Spring": {
			time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC),
		},
	}
	
	for seasonName, expectedDates := range expectedSeasons {
		actualDates, exists := seasons[seasonName]
		if !exists {
			t.Errorf("Season %s not found", seasonName)
			continue
		}
		
		if len(actualDates) != 2 {
			t.Errorf("Season %s should have 2 dates, got %d", seasonName, len(actualDates))
			continue
		}
		
		if !actualDates[0].Equal(expectedDates[0]) {
			t.Errorf("Season %s start date mismatch. Expected: %s, Got: %s",
				seasonName, expectedDates[0].Format("2006-01-02"), actualDates[0].Format("2006-01-02"))
		}
		
		if !actualDates[1].Equal(expectedDates[1]) {
			t.Errorf("Season %s end date mismatch. Expected: %s, Got: %s",
				seasonName, expectedDates[1].Format("2006-01-02"), actualDates[1].Format("2006-01-02"))
		}
	}
}

func TestAUProvider_MultipleStates(t *testing.T) {
	provider := NewAUProvider()
	year := 2024
	
	// Test multiple states at once
	holidays := provider.GetStateHolidays(year, []string{"VIC", "QLD", "WA"})
	
	// Should have holidays from all three states
	expectedCount := 6 // 2 from VIC, 2 from QLD, 2 from WA
	if len(holidays) != expectedCount {
		t.Errorf("Expected %d state holidays, got %d", expectedCount, len(holidays))
	}
	
	// Verify specific holidays exist
	melbourneCup := time.Date(2024, 11, 5, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[melbourneCup]; !exists {
		t.Error("Melbourne Cup Day from VIC not found")
	}
	
	qldLabourDay := time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[qldLabourDay]; !exists {
		t.Error("Queensland Labour Day not found")
	}
	
	waDay := time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[waDay]; !exists {
		t.Error("Western Australia Day not found")
	}
}

func TestAUProvider_ProviderInfo(t *testing.T) {
	provider := NewAUProvider()
	
	if provider.GetCountryCode() != "AU" {
		t.Errorf("Expected country code 'AU', got '%s'", provider.GetCountryCode())
	}
	
	subdivisions := provider.GetSupportedSubdivisions()
	expectedSubdivisions := []string{"NSW", "VIC", "QLD", "SA", "WA", "TAS", "NT", "ACT"}
	
	if len(subdivisions) != len(expectedSubdivisions) {
		t.Errorf("Expected %d subdivisions, got %d", len(expectedSubdivisions), len(subdivisions))
	}
	
	for _, expected := range expectedSubdivisions {
		found := false
		for _, actual := range subdivisions {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subdivision '%s' not found", expected)
		}
	}
}

func TestAUProvider_HolidayCategories(t *testing.T) {
	provider := NewAUProvider()
	year := 2024
	holidays := provider.LoadHolidays(year)
	
	// Check that all holidays have the correct category
	for _, holiday := range holidays {
		if holiday.Category != "public" {
			t.Errorf("Holiday %s has category '%s', expected 'public'", holiday.Name, holiday.Category)
		}
	}
	
	// Check localization exists
	for _, holiday := range holidays {
		if len(holiday.Languages) == 0 {
			t.Errorf("Holiday %s has no languages", holiday.Name)
		}
		
		if englishName, exists := holiday.Languages["en"]; !exists {
			t.Errorf("Holiday %s has no English localization", holiday.Name)
		} else if englishName != holiday.Name {
			t.Errorf("Holiday %s English localization mismatch: %s", holiday.Name, englishName)
		}
	}
}

// Benchmark Australia holiday calculations
func BenchmarkAUProvider_LoadHolidays(b *testing.B) {
	provider := NewAUProvider()
	year := 2024
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(year)
	}
}

func BenchmarkAUProvider_StateHolidays(b *testing.B) {
	provider := NewAUProvider()
	year := 2024
	states := []string{"VIC", "QLD", "WA", "SA"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.GetStateHolidays(year, states)
	}
}
