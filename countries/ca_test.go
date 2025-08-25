package countries

import (
	"testing"
	"time"
)

func TestCAProvider(t *testing.T) {
	ca := NewCAProvider()
	
	if ca.GetCountryCode() != "CA" {
		t.Errorf("Expected country code 'CA', got '%s'", ca.GetCountryCode())
	}
	
	subdivisions := ca.GetSupportedSubdivisions()
	if len(subdivisions) != 13 {
		t.Errorf("Expected 13 provinces/territories, got %d", len(subdivisions))
	}
}

func TestCAHolidays2024(t *testing.T) {
	ca := NewCAProvider()
	holidays := ca.LoadHolidays(2024)
	
	// Test fixed holidays
	testCases := []struct {
		date time.Time
		name string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day"},
		{time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), "Canada Day"},
		{time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC), "Remembrance Day"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Boxing Day"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Good Friday"}, // Good Friday 2024
		{time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), "Easter Monday"}, // Easter Monday 2024
	}
	
	for _, tc := range testCases {
		holiday, exists := holidays[tc.date]
		if !exists {
			t.Errorf("Holiday %s on %s should exist", tc.name, tc.date.Format("2006-01-02"))
			continue
		}
		if holiday.Name != tc.name {
			t.Errorf("Expected holiday name '%s', got '%s'", tc.name, holiday.Name)
		}
	}
}

func TestCAVariableHolidays2024(t *testing.T) {
	ca := NewCAProvider()
	holidays := ca.LoadHolidays(2024)
	
	// Test variable holidays for 2024
	testCases := []struct {
		date time.Time
		name string
	}{
		{time.Date(2024, 2, 19, 0, 0, 0, 0, time.UTC), "Family Day"}, // 3rd Monday in February
		{time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), "Victoria Day"}, // Monday before May 25
		{time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC), "Labour Day"}, // 1st Monday in September
		{time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC), "Thanksgiving Day"}, // 2nd Monday in October
	}
	
	for _, tc := range testCases {
		holiday, exists := holidays[tc.date]
		if !exists {
			t.Errorf("Variable holiday %s on %s should exist", tc.name, tc.date.Format("2006-01-02"))
			continue
		}
		if holiday.Name != tc.name {
			t.Errorf("Expected holiday name '%s', got '%s'", tc.name, holiday.Name)
		}
	}
}

func TestVictoriaDay(t *testing.T) {
	ca := NewCAProvider()
	
	// Test Victoria Day calculation for multiple years
	testCases := []struct {
		year int
		expected time.Time
	}{
		{2024, time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)}, // Monday before May 25
		{2025, time.Date(2025, 5, 19, 0, 0, 0, 0, time.UTC)}, // Monday before May 25
		{2026, time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)}, // Monday before May 25
	}
	
	for _, tc := range testCases {
		actual := ca.getVictoriaDay(tc.year)
		if !actual.Equal(tc.expected) {
			t.Errorf("Victoria Day %d: expected %s, got %s", 
				tc.year, tc.expected.Format("2006-01-02"), actual.Format("2006-01-02"))
		}
		
		// Verify it's always a Monday
		if actual.Weekday() != time.Monday {
			t.Errorf("Victoria Day %d should be a Monday, got %s", tc.year, actual.Weekday())
		}
	}
}

func TestCAMultiLanguageSupport(t *testing.T) {
	ca := NewCAProvider()
	holidays := ca.LoadHolidays(2024)
	
	newYears := holidays[time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)]
	
	if newYears.Languages["en"] != "New Year's Day" {
		t.Errorf("Expected English name 'New Year's Day', got '%s'", newYears.Languages["en"])
	}
	
	if newYears.Languages["fr"] != "Jour de l'An" {
		t.Errorf("Expected French name 'Jour de l'An', got '%s'", newYears.Languages["fr"])
	}
}

func TestCAProvincialHolidays(t *testing.T) {
	ca := NewCAProvider()
	
	// Test Quebec-specific holidays
	qcHolidays := ca.GetProvincialHolidays(2024, []string{"QC"})
	
	stJean := time.Date(2024, 6, 24, 0, 0, 0, 0, time.UTC)
	holiday, exists := qcHolidays[stJean]
	if !exists {
		t.Error("St. Jean Baptiste Day should exist in Quebec")
	}
	if holiday != nil && holiday.Name != "St. Jean Baptiste Day" {
		t.Errorf("Expected 'St. Jean Baptiste Day', got '%s'", holiday.Name)
	}
	
	// Test Newfoundland-specific holidays
	nlHolidays := ca.GetProvincialHolidays(2024, []string{"NL"})
	
	// St. Patrick's Day 2024 falls on Sunday, so should be observed on Monday March 18
	stPatricks := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)
	holiday, exists = nlHolidays[stPatricks]
	if !exists {
		t.Error("St. Patrick's Day should exist in Newfoundland and Labrador")
	}
	if holiday != nil && holiday.Name != "St. Patrick's Day" {
		t.Errorf("Expected 'St. Patrick's Day', got '%s'", holiday.Name)
	}
}

func TestCAObservedDates(t *testing.T) {
	ca := NewCAProvider()
	
	// Test Saturday holiday moved to Monday
	saturday := time.Date(2024, 3, 16, 0, 0, 0, 0, time.UTC) // Saturday
	observed := ca.getObservedDate(saturday)
	expected := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC) // Monday
	
	if !observed.Equal(expected) {
		t.Errorf("Saturday holiday should be observed on Monday: expected %s, got %s",
			expected.Format("2006-01-02"), observed.Format("2006-01-02"))
	}
	
	// Test Sunday holiday moved to Monday
	sunday := time.Date(2024, 3, 17, 0, 0, 0, 0, time.UTC) // Sunday
	observed = ca.getObservedDate(sunday)
	expected = time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC) // Monday
	
	if !observed.Equal(expected) {
		t.Errorf("Sunday holiday should be observed on Monday: expected %s, got %s",
			expected.Format("2006-01-02"), observed.Format("2006-01-02"))
	}
	
	// Test weekday holiday unchanged
	tuesday := time.Date(2024, 3, 19, 0, 0, 0, 0, time.UTC)
	observed = ca.getObservedDate(tuesday)
	
	if !observed.Equal(tuesday) {
		t.Errorf("Weekday holiday should not change: expected %s, got %s",
			tuesday.Format("2006-01-02"), observed.Format("2006-01-02"))
	}
}

func BenchmarkCALoadHolidays(b *testing.B) {
	ca := NewCAProvider()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.LoadHolidays(2024)
	}
}

func BenchmarkCAProvincialHolidays(b *testing.B) {
	ca := NewCAProvider()
	provinces := []string{"QC", "ON", "BC", "AB"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ca.GetProvincialHolidays(2024, provinces)
	}
}
