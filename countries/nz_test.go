package countries

import (
	"testing"
	"time"
)

func TestNZProvider_BasicHolidays(t *testing.T) {
	provider := NewNZProvider()
	year := 2024
	holidays := provider.LoadHolidays(year)

	expectedHolidays := map[string]time.Time{
		"New Year's Day":              time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		"Day after New Year's Day":    time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		"Waitangi Day":                time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC),
		"Good Friday":                 time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
		"Easter Monday":               time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		"ANZAC Day":                   time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC),
		"Queen's Birthday":            time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC), // First Monday in June
		"Matariki":                    time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC), // Known date for 2024
		"Labour Day":                  time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC), // Fourth Monday in October
		"Christmas Day":               time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
		"Boxing Day":                  time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC),
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

func TestNZProvider_EasterCalculations(t *testing.T) {
	provider := NewNZProvider()
	
	testCases := []struct {
		year               int
		expectedGoodFriday time.Time
		expectedEasterMonday time.Time
	}{
		{
			year:               2024,
			expectedGoodFriday: time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC),
			expectedEasterMonday: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			year:               2025,
			expectedGoodFriday: time.Date(2025, 4, 18, 0, 0, 0, 0, time.UTC),
			expectedEasterMonday: time.Date(2025, 4, 21, 0, 0, 0, 0, time.UTC),
		},
	}
	
	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		if holiday, exists := holidays[tc.expectedGoodFriday]; !exists || holiday.Name != "Good Friday" {
			t.Errorf("Year %d: Good Friday not found on expected date %s", tc.year, tc.expectedGoodFriday.Format("2006-01-02"))
		}
		
		if holiday, exists := holidays[tc.expectedEasterMonday]; !exists || holiday.Name != "Easter Monday" {
			t.Errorf("Year %d: Easter Monday not found on expected date %s", tc.year, tc.expectedEasterMonday.Format("2006-01-02"))
		}
	}
}

func TestNZProvider_QueensBirthday(t *testing.T) {
	provider := NewNZProvider()
	
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC)},  // First Monday in June
		{2025, time.Date(2025, 6, 2, 0, 0, 0, 0, time.UTC)},  // First Monday in June
		{2023, time.Date(2023, 6, 5, 0, 0, 0, 0, time.UTC)},  // First Monday in June
	}
	
	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		if holiday, exists := holidays[tc.expected]; !exists || holiday.Name != "Queen's Birthday" {
			t.Errorf("Year %d: Queen's Birthday not found on expected date %s", tc.year, tc.expected.Format("2006-01-02"))
		}
		
		// Verify it's actually the first Monday
		if tc.expected.Weekday() != time.Monday {
			t.Errorf("Year %d: Queen's Birthday should be a Monday, got %s", tc.year, tc.expected.Weekday())
		}
		
		// Verify it's in June
		if tc.expected.Month() != time.June {
			t.Errorf("Year %d: Queen's Birthday should be in June, got %s", tc.year, tc.expected.Month())
		}
	}
}

func TestNZProvider_LabourDay(t *testing.T) {
	provider := NewNZProvider()
	
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 10, 28, 0, 0, 0, 0, time.UTC)}, // Fourth Monday in October
		{2025, time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)}, // Fourth Monday in October
		{2023, time.Date(2023, 10, 23, 0, 0, 0, 0, time.UTC)}, // Fourth Monday in October
	}
	
	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		if holiday, exists := holidays[tc.expected]; !exists || holiday.Name != "Labour Day" {
			t.Errorf("Year %d: Labour Day not found on expected date %s", tc.year, tc.expected.Format("2006-01-02"))
		}
		
		// Verify it's actually the fourth Monday
		if tc.expected.Weekday() != time.Monday {
			t.Errorf("Year %d: Labour Day should be a Monday, got %s", tc.year, tc.expected.Weekday())
		}
		
		// Verify it's in October
		if tc.expected.Month() != time.October {
			t.Errorf("Year %d: Labour Day should be in October, got %s", tc.year, tc.expected.Month())
		}
	}
}

func TestNZProvider_Matariki(t *testing.T) {
	provider := NewNZProvider()
	
	// Test known Matariki dates
	testCases := []struct {
		year     int
		expected time.Time
		hasHoliday bool
	}{
		{2024, time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC), true},
		{2025, time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC), true},
		{2023, time.Date(2023, 7, 14, 0, 0, 0, 0, time.UTC), true},
		{2020, time.Time{}, false}, // Year not in our table
	}
	
	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		if tc.hasHoliday {
			if holiday, exists := holidays[tc.expected]; !exists || holiday.Name != "Matariki" {
				t.Errorf("Year %d: Matariki not found on expected date %s", tc.year, tc.expected.Format("2006-01-02"))
			}
		} else {
			// Check that no Matariki holiday exists for this year
			foundMatariki := false
			for _, holiday := range holidays {
				if holiday.Name == "Matariki" {
					foundMatariki = true
					break
				}
			}
			if foundMatariki {
				t.Errorf("Year %d: Unexpected Matariki holiday found", tc.year)
			}
		}
	}
}

func TestNZProvider_RegionalHolidays(t *testing.T) {
	provider := NewNZProvider()
	year := 2024

	// Test Auckland Anniversary
	aucklandHolidays := provider.GetRegionalHolidays(year, []string{"AUK"})
	
	expectedAuckland := time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC) // Monday closest to Jan 29
	if holiday, exists := aucklandHolidays[expectedAuckland]; !exists || holiday.Name != "Auckland Anniversary Day" {
		t.Errorf("Auckland Anniversary Day not found on expected date %s", expectedAuckland.Format("2006-01-02"))
	}

	// Test Wellington Anniversary
	wellingtonHolidays := provider.GetRegionalHolidays(year, []string{"WGN"})
	
	expectedWellington := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC) // Monday closest to Jan 22
	if holiday, exists := wellingtonHolidays[expectedWellington]; !exists || holiday.Name != "Wellington Anniversary Day" {
		t.Errorf("Wellington Anniversary Day not found on expected date %s", expectedWellington.Format("2006-01-02"))
	}

	// Test Canterbury Anniversary (Show Day)
	canterburyHolidays := provider.GetRegionalHolidays(year, []string{"CAN"})
	
	// Canterbury Show Day is the Friday after the first Tuesday in November
	firstTuesday := NthWeekdayOfMonth(2024, 11, time.Tuesday, 1) // Nov 5, 2024
	expectedCanterbury := firstTuesday.AddDate(0, 0, 3) // Friday after (Nov 8, 2024)
	
	if holiday, exists := canterburyHolidays[expectedCanterbury]; !exists || holiday.Name != "Canterbury Anniversary Day" {
		t.Errorf("Canterbury Anniversary Day not found on expected date %s", expectedCanterbury.Format("2006-01-02"))
	}

	// Test Southland Anniversary (Easter Tuesday)
	southlandHolidays := provider.GetRegionalHolidays(year, []string{"STL"})
	
	easter := EasterSunday(2024) // March 31, 2024
	expectedSouthland := easter.AddDate(0, 0, 2) // Tuesday after Easter (April 2, 2024)
	
	if holiday, exists := southlandHolidays[expectedSouthland]; !exists || holiday.Name != "Southland Anniversary Day" {
		t.Errorf("Southland Anniversary Day not found on expected date %s", expectedSouthland.Format("2006-01-02"))
	}
}

func TestNZProvider_GetClosestMonday(t *testing.T) {
	provider := NewNZProvider()
	
	testCases := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "already monday",
			input:    time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC), // Monday
			expected: time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "tuesday - previous monday",
			input:    time.Date(2024, 1, 23, 0, 0, 0, 0, time.UTC), // Tuesday
			expected: time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC), // Monday
		},
		{
			name:     "friday - next monday",
			input:    time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC), // Friday
			expected: time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC), // Monday
		},
		{
			name:     "sunday - next monday",
			input:    time.Date(2024, 1, 28, 0, 0, 0, 0, time.UTC), // Sunday
			expected: time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC), // Monday
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := provider.getClosestMonday(tc.input)
			if !result.Equal(tc.expected) {
				t.Errorf("Expected %s, got %s", tc.expected.Format("2006-01-02"), result.Format("2006-01-02"))
			}
			
			// Verify result is actually a Monday
			if result.Weekday() != time.Monday {
				t.Errorf("Result should be a Monday, got %s", result.Weekday())
			}
		})
	}
}

func TestNZProvider_Seasons(t *testing.T) {
	provider := NewNZProvider()
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

func TestNZProvider_MultipleRegions(t *testing.T) {
	provider := NewNZProvider()
	year := 2024
	
	// Test multiple regions at once
	holidays := provider.GetRegionalHolidays(year, []string{"AUK", "WGN", "CAN"})
	
	// Should have holidays from all three regions
	expectedCount := 3
	if len(holidays) != expectedCount {
		t.Errorf("Expected %d regional holidays, got %d", expectedCount, len(holidays))
	}
	
	// Verify specific holidays exist
	aucklandDay := time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[aucklandDay]; !exists {
		t.Error("Auckland Anniversary Day not found")
	}
	
	wellingtonDay := time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays[wellingtonDay]; !exists {
		t.Error("Wellington Anniversary Day not found")
	}
}

func TestNZProvider_ProviderInfo(t *testing.T) {
	provider := NewNZProvider()
	
	if provider.GetCountryCode() != "NZ" {
		t.Errorf("Expected country code 'NZ', got '%s'", provider.GetCountryCode())
	}
	
	subdivisions := provider.GetSupportedSubdivisions()
	expectedSubdivisions := []string{"AUK", "BOP", "CAN", "GIS", "HKB", "MWT", "MBH", "NSN", "NTL", "OTA", "STL", "TKI", "TAS", "WKO", "WGN", "WTC", "CIT"}
	
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

func TestNZProvider_MaoriLanguageSupport(t *testing.T) {
	provider := NewNZProvider()
	year := 2024
	holidays := provider.LoadHolidays(year)
	
	// Check that key holidays have Māori translations
	waitangiDay := time.Date(2024, 2, 6, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[waitangiDay]; exists {
		if maoriName, hasMaori := holiday.Languages["mi"]; !hasMaori {
			t.Error("Waitangi Day should have Māori translation")
		} else if maoriName != "Te Rā o Waitangi" {
			t.Errorf("Expected Māori name 'Te Rā o Waitangi', got '%s'", maoriName)
		}
	}
	
	// Check Matariki has Māori name
	matariki := time.Date(2024, 6, 28, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[matariki]; exists {
		if maoriName, hasMaori := holiday.Languages["mi"]; !hasMaori {
			t.Error("Matariki should have Māori translation")
		} else if maoriName != "Matariki" {
			t.Errorf("Expected Māori name 'Matariki', got '%s'", maoriName)
		}
	}
}

// Benchmark New Zealand holiday calculations
func BenchmarkNZProvider_LoadHolidays(b *testing.B) {
	provider := NewNZProvider()
	year := 2024
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.LoadHolidays(year)
	}
}

func BenchmarkNZProvider_RegionalHolidays(b *testing.B) {
	provider := NewNZProvider()
	year := 2024
	regions := []string{"AUK", "WGN", "CAN", "OTA"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.GetRegionalHolidays(year, regions)
	}
}

func BenchmarkNZProvider_GetClosestMonday(b *testing.B) {
	provider := NewNZProvider()
	testDate := time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC) // Thursday
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = provider.getClosestMonday(testDate)
	}
}
