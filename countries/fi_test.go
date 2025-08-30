package countries

import (
	"testing"
	"time"
)

func TestFIHolidays(t *testing.T) {
	provider := NewFIProvider()

	// Test basic provider properties
	if provider.GetCountryCode() != "FI" {
		t.Errorf("Expected country code FI, got %s", provider.GetCountryCode())
	}

	// Test 2024 holidays
	holidays := provider.LoadHolidays(2024)

	testCases := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"New Year's Day", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{"Epiphany", time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), true},
		{"Good Friday", time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), true},
		{"Easter Sunday", time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), true},
		{"Easter Monday", time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC), true},
		{"May Day", time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), true},
		{"Ascension Day", time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC), true},
		{"Pentecost", time.Date(2024, 5, 19, 0, 0, 0, 0, time.UTC), true},
		{"Midsummer Eve", time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC), true},
		{"Midsummer Day", time.Date(2024, 6, 22, 0, 0, 0, 0, time.UTC), true},
		{"All Saints' Day", time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC), true},
		{"Independence Day", time.Date(2024, 12, 6, 0, 0, 0, 0, time.UTC), true},
		{"Christmas Eve", time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC), true},
		{"Christmas Day", time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), true},
		{"Boxing Day", time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), true},

		// Negative test cases
		{"Random Day", time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC), false},
		{"Random Day 2", time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, exists := holidays[tc.date]
			if exists != tc.expected {
				if tc.expected {
					t.Errorf("%s should be a holiday", tc.name)
				} else {
					t.Errorf("%s should not be a holiday", tc.name)
				}
			}
		})
	}

	// Test specific holiday properties
	christmas, exists := holidays[time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)]
	if !exists {
		t.Fatal("Christmas should be a holiday")
	}

	if christmas.Name != "Joulupäivä" {
		t.Errorf("Expected name 'Joulupäivä', got '%s'", christmas.Name)
	}

	if christmas.Category != "public" {
		t.Errorf("Expected category 'public', got '%s'", christmas.Category)
	}

	// Test language translations
	if christmas.Languages["fi"] != "Joulupäivä" {
		t.Errorf("Expected Finnish name 'Joulupäivä', got '%s'", christmas.Languages["fi"])
	}

	if christmas.Languages["en"] != "Christmas Day" {
		t.Errorf("Expected English name 'Christmas Day', got '%s'", christmas.Languages["en"])
	}

	if christmas.Languages["sv"] != "Juldagen" {
		t.Errorf("Expected Swedish name 'Juldagen', got '%s'", christmas.Languages["sv"])
	}
}

