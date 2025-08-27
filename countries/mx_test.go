package countries

import (
	"fmt"
	"testing"
	"time"
)

func TestMXProvider_LoadHolidays(t *testing.T) {
	provider := NewMXProvider()

	tests := []struct {
		year int
		want int // Expected number of holidays
	}{
		{2024, 12}, // Various fixed and variable holidays
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
			} else if holiday.Name != "Año Nuevo" {
				t.Errorf("Expected Año Nuevo, got %s", holiday.Name)
			}

			independence := time.Date(tt.year, 9, 16, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[independence]; !exists {
				t.Errorf("Independence Day not found for year %d", tt.year)
			} else if holiday.Name != "Día de la Independencia" {
				t.Errorf("Expected Día de la Independencia, got %s", holiday.Name)
			}

			christmas := time.Date(tt.year, 12, 25, 0, 0, 0, 0, time.UTC)
			if holiday, exists := holidays[christmas]; !exists {
				t.Errorf("Christmas Day not found for year %d", tt.year)
			} else if holiday.Name != "Navidad" {
				t.Errorf("Expected Navidad, got %s", holiday.Name)
			}
		})
	}
}

func TestMXProvider_GetCountryCode(t *testing.T) {
	provider := NewMXProvider()
	if code := provider.GetCountryCode(); code != "MX" {
		t.Errorf("GetCountryCode() = %v, want MX", code)
	}
}

func TestMXProvider_GetSupportedSubdivisions(t *testing.T) {
	provider := NewMXProvider()
	subdivisions := provider.GetSupportedSubdivisions()

	// Check that we have all Mexican states + federal district
	expectedCount := 32 // 31 states + 1 federal district
	if len(subdivisions) != expectedCount {
		t.Errorf("GetSupportedSubdivisions() returned %d subdivisions, want %d", len(subdivisions), expectedCount)
	}

	// Check for specific states
	found := false
	for _, sub := range subdivisions {
		if sub == "CMX" { // Mexico City (Ciudad de México)
			found = true
			break
		}
	}
	if !found {
		t.Error("Mexico City (CMX) not found in supported subdivisions")
	}
}

func TestMXProvider_GetSupportedCategories(t *testing.T) {
	provider := NewMXProvider()
	categories := provider.GetSupportedCategories()

	expectedCategories := []string{"public", "national", "religious", "regional", "civic"}
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

func TestMXProvider_Languages(t *testing.T) {
	provider := NewMXProvider()
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
		if holiday.Languages["es"] == "" {
			t.Error("Spanish translation should not be empty")
		}
	}
}

func TestMXProvider_VariableHolidays(t *testing.T) {
	provider := NewMXProvider()

	// Test 2024 specifically for known dates
	holidays := provider.LoadHolidays(2024)

	// Constitution Day should be first Monday of February (February 5, 2024)
	constitutionDay := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[constitutionDay]; !exists {
		t.Error("Constitution Day not found for 2024")
	} else if holiday.Name != "Día de la Constitución" {
		t.Errorf("Expected Día de la Constitución, got %s", holiday.Name)
	}

	// Benito Juárez's Birthday should be third Monday of March (March 18, 2024)
	juarezBirthday := time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[juarezBirthday]; !exists {
		t.Error("Benito Juárez's Birthday not found for 2024")
	} else if holiday.Name != "Natalicio de Benito Juárez" {
		t.Errorf("Expected Natalicio de Benito Juárez, got %s", holiday.Name)
	}
}

func TestMXProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewMXProvider()

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
		} else if holiday.Name != "Viernes Santo" {
			t.Errorf("Expected Viernes Santo, got %s", holiday.Name)
		}

		// Check Maundy Thursday (3 days before Easter)
		maundyThursday := tc.expected.AddDate(0, 0, -3)
		if holiday, exists := holidays[maundyThursday]; !exists {
			t.Errorf("Maundy Thursday not found for year %d on expected date %v", tc.year, maundyThursday)
		} else if holiday.Name != "Jueves Santo" {
			t.Errorf("Expected Jueves Santo, got %s", holiday.Name)
		}
	}
}

func TestMXProvider_CulturalHolidays(t *testing.T) {
	provider := NewMXProvider()
	holidays := provider.LoadHolidays(2024)

	// Day of the Dead - November 2
	dayOfDead := time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[dayOfDead]; !exists {
		t.Error("Day of the Dead not found for 2024")
	} else if holiday.Name != "Día de los Muertos" {
		t.Errorf("Expected Día de los Muertos, got %s", holiday.Name)
	}

	// Our Lady of Guadalupe - December 12
	guadalupe := time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[guadalupe]; !exists {
		t.Error("Our Lady of Guadalupe not found for 2024")
	} else if holiday.Name != "Día de la Virgen de Guadalupe" {
		t.Errorf("Expected Día de la Virgen de Guadalupe, got %s", holiday.Name)
	}
}

func TestMXProvider_IsSubdivisionSupported(t *testing.T) {
	provider := NewMXProvider()

	// Test valid subdivisions
	if !provider.IsSubdivisionSupported("CMX") {
		t.Error("Mexico City (CMX) should be supported")
	}

	if !provider.IsSubdivisionSupported("JAL") {
		t.Error("Jalisco (JAL) should be supported")
	}

	// Test invalid subdivision
	if provider.IsSubdivisionSupported("XX") {
		t.Error("XX should not be supported")
	}
}

func TestMXProvider_IsCategorySupported(t *testing.T) {
	provider := NewMXProvider()

	// Test valid categories
	if !provider.IsCategorySupported("civic") {
		t.Error("civic category should be supported")
	}

	if !provider.IsCategorySupported("religious") {
		t.Error("religious category should be supported")
	}

	// Test invalid category
	if provider.IsCategorySupported("invalid") {
		t.Error("invalid category should not be supported")
	}
}

func TestMXProvider_GetNthWeekdayOfMonth(t *testing.T) {
	provider := NewMXProvider()

	// Test first Monday of February 2024 (should be February 5)
	firstMonday := provider.getFirstMondayOfFebruary(2024)
	expected := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)
	if !firstMonday.Equal(expected) {
		t.Errorf("First Monday of February 2024: got %v, want %v", firstMonday, expected)
	}

	// Test third Monday of March 2024 (should be March 18)
	thirdMonday := provider.getThirdMondayOfMarch(2024)
	expected = time.Date(2024, 3, 18, 0, 0, 0, 0, time.UTC)
	if !thirdMonday.Equal(expected) {
		t.Errorf("Third Monday of March 2024: got %v, want %v", thirdMonday, expected)
	}
}

func BenchmarkMXProvider_LoadHolidays(b *testing.B) {
	provider := NewMXProvider()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		provider.LoadHolidays(2024)
	}
}
