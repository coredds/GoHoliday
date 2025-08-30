package countries

import (
	"testing"
	"time"
)

func TestINProvider_Creation(t *testing.T) {
	provider := NewINProvider()

	if provider.GetCountryCode() != "IN" {
		t.Errorf("Expected country code IN, got %s", provider.GetCountryCode())
	}

	subdivisions := provider.GetSupportedSubdivisions()
	if len(subdivisions) != 36 { // 28 states + 8 union territories
		t.Errorf("Expected 36 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetSupportedCategories()
	expectedCategories := []string{"national", "gazetted", "restricted", "hindu", "muslim", "christian", "sikh", "buddhist", "jain", "regional"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestINProvider_LoadHolidays2024(t *testing.T) {
	provider := NewINProvider()
	holidays := provider.LoadHolidays(2024)

	// Test national holidays for 2024
	expectedHolidays := map[string]time.Time{
		"Republic Day":     time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC),
		"Independence Day": time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC),
		"Gandhi Jayanti":   time.Date(2024, 10, 2, 0, 0, 0, 0, time.UTC),
		"Christmas Day":    time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
	}

	for name, expectedDate := range expectedHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == name && holiday.Date.Equal(expectedDate) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected national holiday %s on %s not found", name, expectedDate.Format("2006-01-02"))
		}
	}

	// Test Christian holidays for 2024 (Easter was March 31, 2024)
	christianHolidays := map[string]time.Time{
		"Good Friday":   time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), // Easter - 2 days
		"Easter Sunday": time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC), // Easter
	}

	for name, expectedDate := range christianHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == name && holiday.Date.Equal(expectedDate) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected Christian holiday %s on %s not found", name, expectedDate.Format("2006-01-02"))
		}
	}

	// Check minimum number of holidays
	if len(holidays) < 6 {
		t.Errorf("Expected at least 6 holidays, got %d", len(holidays))
	}
}

func TestINProvider_NationalHolidays(t *testing.T) {
	provider := NewINProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that national holidays have proper categories
	nationalHolidays := []string{"Republic Day", "Independence Day", "Gandhi Jayanti"}

	for _, name := range nationalHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == name {
				found = true
				if holiday.Category != "national" {
					t.Errorf("Expected %s to be national category, got %s", name, holiday.Category)
				}
				// Check for Hindi translation
				if holiday.Languages == nil || holiday.Languages["hi"] == "" {
					t.Errorf("Expected %s to have Hindi translation", name)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected to find national holiday: %s", name)
		}
	}
}

func TestINProvider_StateHolidays(t *testing.T) {
	provider := NewINProvider()

	// Test Maharashtra state holidays
	maharashtraHolidays := provider.GetStateHolidays(2024, "MH")
	if len(maharashtraHolidays) == 0 {
		t.Error("Expected Maharashtra to have state holidays")
	}

	// Check for Maharashtra Day
	maharashtraDay := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	found := false
	for _, holiday := range maharashtraHolidays {
		if holiday.Date.Equal(maharashtraDay) && holiday.Name == "Maharashtra Day" {
			found = true
			if holiday.Category != "regional" {
				t.Errorf("Expected Maharashtra Day to be regional category, got %s", holiday.Category)
			}
			if len(holiday.Subdivisions) == 0 || holiday.Subdivisions[0] != "MH" {
				t.Errorf("Expected Maharashtra Day to be specific to MH subdivision")
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find Maharashtra Day in Maharashtra holidays")
	}

	// Test Gujarat state holidays
	gujaratHolidays := provider.GetStateHolidays(2024, "GJ")
	if len(gujaratHolidays) == 0 {
		t.Error("Expected Gujarat to have state holidays")
	}

	// Test unknown state
	unknownHolidays := provider.GetStateHolidays(2024, "XX")
	if len(unknownHolidays) != 0 {
		t.Error("Expected unknown state to have no holidays")
	}
}

func TestINProvider_MajorFestivals(t *testing.T) {
	provider := NewINProvider()
	festivals := provider.GetMajorFestivals(2024)

	// Check that major festivals are included
	expectedFestivals := []string{"Holi", "Dussehra", "Diwali", "Janmashtami", "Ram Navami"}

	for _, festivalName := range expectedFestivals {
		found := false
		for _, festival := range festivals {
			if festival.Name == festivalName {
				found = true
				if festival.Category != "hindu" {
					t.Errorf("Expected %s to be hindu category, got %s", festivalName, festival.Category)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected to find major festival: %s", festivalName)
		}
	}

	// Check for non-Hindu festivals
	nonHinduFestivals := map[string]string{
		"Buddha Purnima":     "buddhist",
		"Mahavir Jayanti":    "jain",
		"Guru Nanak Jayanti": "sikh",
	}

	for festivalName, expectedCategory := range nonHinduFestivals {
		found := false
		for _, festival := range festivals {
			if festival.Name == festivalName {
				found = true
				if festival.Category != expectedCategory {
					t.Errorf("Expected %s to be %s category, got %s", festivalName, expectedCategory, festival.Category)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected to find festival: %s", festivalName)
		}
	}
}

func TestINProvider_EasterCalculation(t *testing.T) {
	provider := NewINProvider()

	// Test Easter dates for known years
	testCases := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		easter := provider.calculateEaster(tc.year)
		if !easter.Equal(tc.expected) {
			t.Errorf("Easter %d: expected %s, got %s", tc.year, tc.expected.Format("2006-01-02"), easter.Format("2006-01-02"))
		}
	}
}

func TestINProvider_HolidayCategories(t *testing.T) {
	provider := NewINProvider()
	holidays := provider.LoadHolidays(2024)

	validCategories := map[string]bool{
		"national":   true,
		"gazetted":   true,
		"restricted": true,
		"hindu":      true,
		"muslim":     true,
		"christian":  true,
		"sikh":       true,
		"buddhist":   true,
		"jain":       true,
		"regional":   true,
	}

	for _, holiday := range holidays {
		if !validCategories[holiday.Category] {
			t.Errorf("Holiday %s has invalid category: %s", holiday.Name, holiday.Category)
		}
	}
}

func TestINProvider_DiverseReligions(t *testing.T) {
	provider := NewINProvider()
	holidays := provider.LoadHolidays(2024)
	festivals := provider.GetMajorFestivals(2024)

	// Combine all holidays
	allHolidays := make(map[time.Time]*Holiday)
	for date, holiday := range holidays {
		allHolidays[date] = holiday
	}
	for date, festival := range festivals {
		allHolidays[date] = festival
	}

	// Check that multiple religions are represented
	religionCategories := map[string]bool{
		"hindu":     false,
		"christian": false,
		"buddhist":  false,
		"sikh":      false,
		"jain":      false,
	}

	for _, holiday := range allHolidays {
		if _, exists := religionCategories[holiday.Category]; exists {
			religionCategories[holiday.Category] = true
		}
	}

	for religion, found := range religionCategories {
		if !found {
			t.Errorf("Expected to find holidays from %s tradition", religion)
		}
	}
}

func TestINProvider_UniqueIndianHolidays(t *testing.T) {
	provider := NewINProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that India has unique holidays
	uniqueHolidays := []string{
		"Republic Day",
		"Independence Day",
		"Gandhi Jayanti",
	}

	for _, uniqueName := range uniqueHolidays {
		found := false
		for _, holiday := range holidays {
			if holiday.Name == uniqueName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find unique Indian holiday: %s", uniqueName)
		}
	}
}
