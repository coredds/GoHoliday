package countries

import (
	"fmt"
	"testing"
	"time"
)

func TestBRProvider_LoadHolidays(t *testing.T) {
	provider := NewBRProvider()

	tests := []struct {
		year int
		want int // Expected number of holidays
	}{
		{2024, 12}, // 8 fixed + 4 Easter-based holidays
		{2025, 12},
		{2026, 12},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("year_%d", tt.year), func(t *testing.T) {
			holidays := provider.LoadHolidays(tt.year)

			if len(holidays) != tt.want {
				t.Errorf("LoadHolidays(%d) returned %d holidays, want %d", tt.year, len(holidays), tt.want)
			}

			// Test specific national holidays
			newYear := time.Date(tt.year, 1, 1, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[newYear]; !exists {
				t.Errorf("New Year's Day not found for year %d", tt.year)
			} else if holiday.Name != "Confraternização Universal" {
				t.Errorf("Expected Confraternização Universal, got %s", holiday.Name)
			}

			independence := time.Date(tt.year, 9, 7, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[independence]; !exists {
				t.Errorf("Independence Day not found for year %d", tt.year)
			} else if holiday.Name != "Independência do Brasil" {
				t.Errorf("Expected Independência do Brasil, got %s", holiday.Name)
			}

			christmas := time.Date(tt.year, 12, 25, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[christmas]; !exists {
				t.Errorf("Christmas Day not found for year %d", tt.year)
			} else if holiday.Name != "Natal" {
				t.Errorf("Expected Natal, got %s", holiday.Name)
			}
		})
	}
}

func TestBRProvider_GetCountryCode(t *testing.T) {
	provider := NewBRProvider()
	if code := provider.GetCountryCode(); code != "BR" {
		t.Errorf("GetCountryCode() = %v, want BR", code)
	}
}

func TestBRProvider_GetSupportedSubdivisions(t *testing.T) {
	provider := NewBRProvider()
	subdivisions := provider.GetSupportedSubdivisions()

	// Check that we have all Brazilian states + federal district
	expectedCount := 27 // 26 states + 1 federal district
	if len(subdivisions) != expectedCount {
		t.Errorf("GetSupportedSubdivisions() returned %d subdivisions, want %d", len(subdivisions), expectedCount)
	}

	// Check for specific states
	found := false
	for _, sub := range subdivisions {
		if sub == "SP" { // São Paulo
			found = true
			break
		}
	}
	if !found {
		t.Error("São Paulo (SP) not found in supported subdivisions")
	}
}

func TestBRProvider_GetSupportedCategories(t *testing.T) {
	provider := NewBRProvider()
	categories := provider.GetSupportedCategories()

	expectedCategories := []string{"public", "national", "religious", "regional", "carnival"}
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

func TestBRProvider_Languages(t *testing.T) {
	provider := NewBRProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that holidays have multilingual names
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[newYear]; exists {
		if holiday.Languages == nil {
			t.Error("New Year's Day should have language translations")
		}
		if holiday.Languages["en"] != "New Year's Day" {
			t.Errorf("Expected English name 'New Year's Day', got %s", holiday.Languages["en"])
		}
		if holiday.Languages["pt"] == "" {
			t.Error("Portuguese translation should not be empty")
		}
	}
}

func TestBRProvider_CarnivalHolidays(t *testing.T) {
	provider := NewBRProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that Carnival holidays exist
	carnivalFound := false
	for _, holiday := range holidays {
		if holiday.Category == "carnival" {
			carnivalFound = true
			break
		}
	}

	if !carnivalFound {
		t.Error("No Carnival holidays found")
	}
}

func TestBRProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewBRProvider()

	// Test known Easter dates
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)}, // March 31, 2024
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)}, // April 20, 2025
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)

		// Check Good Friday (2 days before Easter)
		goodFriday := tc.expected.AddDate(0, 0, -2)
		if holiday, exists := holidays[goodFriday]; !exists {
			t.Errorf("Good Friday not found for year %d on expected date %v", tc.year, goodFriday)
		} else if holiday.Name != "Sexta-feira Santa" {
			t.Errorf("Expected Sexta-feira Santa, got %s", holiday.Name)
		}

		// Check Corpus Christi (60 days after Easter)
		corpusChristi := tc.expected.AddDate(0, 0, 60)
		if holiday, exists := holidays[corpusChristi]; !exists {
			t.Errorf("Corpus Christi not found for year %d on expected date %v", tc.year, corpusChristi)
		} else if holiday.Name != "Corpus Christi" {
			t.Errorf("Expected Corpus Christi, got %s", holiday.Name)
		}
	}
}

func TestBRProvider_IsSubdivisionSupported(t *testing.T) {
	provider := NewBRProvider()

	// Test valid subdivisions
	if !provider.IsSubdivisionSupported("SP") {
		t.Error("São Paulo (SP) should be supported")
	}

	if !provider.IsSubdivisionSupported("RJ") {
		t.Error("Rio de Janeiro (RJ) should be supported")
	}

	// Test invalid subdivision
	if provider.IsSubdivisionSupported("XX") {
		t.Error("XX should not be supported")
	}
}

func TestBRProvider_IsCategorySupported(t *testing.T) {
	provider := NewBRProvider()

	// Test valid categories
	if !provider.IsCategorySupported("carnival") {
		t.Error("carnival category should be supported")
	}

	if !provider.IsCategorySupported("religious") {
		t.Error("religious category should be supported")
	}

	// Test invalid category
	if provider.IsCategorySupported("invalid") {
		t.Error("invalid category should not be supported")
	}
}

func BenchmarkBRProvider_LoadHolidays(b *testing.B) {
	provider := NewBRProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(2024)
	}
}
