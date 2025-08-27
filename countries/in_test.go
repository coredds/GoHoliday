package countries

import (
	"fmt"
	"testing"
	"time"
)

func TestINProvider_LoadHolidays(t *testing.T) {
	provider := NewINProvider()
	
	tests := []struct {
		year int
		want int // Expected number of holidays
	}{
		{2024, 8}, // Approximate count for major holidays
		{2025, 8},
		{2026, 8},
	}
	
	for _, tt := range tests {
		t.Run(fmt.Sprintf("year_%d", tt.year), func(t *testing.T) {
			holidays := provider.LoadHolidays(tt.year)
			
			if len(holidays) < 6 { // At least the fixed national holidays + some religious ones
				t.Errorf("LoadHolidays(%d) returned %d holidays, want at least 6", tt.year, len(holidays))
			}
			
			// Test specific national holidays
			republicDay := time.Date(tt.year, 1, 26, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[republicDay]; !exists {
				t.Errorf("Republic Day not found for year %d", tt.year)
			} else if holiday.Name != "Republic Day" {
				t.Errorf("Expected Republic Day, got %s", holiday.Name)
			}
			
			independenceDay := time.Date(tt.year, 8, 15, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[independenceDay]; !exists {
				t.Errorf("Independence Day not found for year %d", tt.year)
			} else if holiday.Name != "Independence Day" {
				t.Errorf("Expected Independence Day, got %s", holiday.Name)
			}
			
			gandhiJayanti := time.Date(tt.year, 10, 2, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[gandhiJayanti]; !exists {
				t.Errorf("Gandhi Jayanti not found for year %d", tt.year)
			} else if holiday.Name != "Gandhi Jayanti" {
				t.Errorf("Expected Gandhi Jayanti, got %s", holiday.Name)
			}
			
			christmas := time.Date(tt.year, 12, 25, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[christmas]; !exists {
				t.Errorf("Christmas Day not found for year %d", tt.year)
			} else if holiday.Name != "Christmas Day" {
				t.Errorf("Expected Christmas Day, got %s", holiday.Name)
			}
		})
	}
}

func TestINProvider_GetCountryCode(t *testing.T) {
	provider := NewINProvider()
	if code := provider.GetCountryCode(); code != "IN" {
		t.Errorf("GetCountryCode() = %v, want IN", code)
	}
}

func TestINProvider_GetSupportedSubdivisions(t *testing.T) {
	provider := NewINProvider()
	subdivisions := provider.GetSupportedSubdivisions()
	
	// Check that we have all major Indian states/territories
	expectedCount := 36 // As per the list in NewINProvider
	if len(subdivisions) != expectedCount {
		t.Errorf("GetSupportedSubdivisions() returned %d subdivisions, want %d", len(subdivisions), expectedCount)
	}
	
	// Check for specific states
	found := false
	for _, sub := range subdivisions {
		if sub == "DL" { // Delhi
			found = true
			break
		}
	}
	if !found {
		t.Error("Delhi (DL) not found in supported subdivisions")
	}
}

func TestINProvider_GetSupportedCategories(t *testing.T) {
	provider := NewINProvider()
	categories := provider.GetSupportedCategories()
	
	expectedCategories := []string{"public", "national", "religious", "hindu", "islamic", "christian"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("GetSupportedCategories() returned %d categories, want %d", len(categories), len(expectedCategories))
	}
	
	// Check for specific categories
	categoryMap := make(map[string]bool)
	for _, cat := range categories {
		categoryMap[cat] = true
	}
	
	for _, expected := range expectedCategories {
		if !categoryMap[expected] {
			t.Errorf("Expected category %s not found", expected)
		}
	}
}

func TestINProvider_Languages(t *testing.T) {
	provider := NewINProvider()
	holidays := provider.LoadHolidays(2024)
	
	// Test that holidays have multilingual names
	republicDay := time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[republicDay]; exists {
		if holiday.Languages == nil {
			t.Error("Republic Day should have language translations")
		}
		if holiday.Languages["en"] != "Republic Day" {
			t.Errorf("Expected English name 'Republic Day', got %s", holiday.Languages["en"])
		}
		if holiday.Languages["hi"] == "" {
			t.Error("Hindi translation should not be empty")
		}
	}
}

func BenchmarkINProvider_LoadHolidays(b *testing.B) {
	provider := NewINProvider()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(2024)
	}
}

func TestINProvider_EasterCalculation(t *testing.T) {
	provider := NewINProvider()
	
	// Test known Easter dates
	testCases := []struct {
		year int
		expected time.Time
	}{
		{2024, time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)}, // March 31, 2024
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)}, // April 20, 2025
	}
	
	for _, tc := range testCases {
		easter := provider.calculateEaster(tc.year)
		if !easter.Equal(tc.expected) {
			t.Errorf("calculateEaster(%d) = %v, want %v", tc.year, easter, tc.expected)
		}
	}
}
