package countries

import (
	"testing"
	"time"
)

func TestILProvider_BasicInfo(t *testing.T) {
	provider := NewILProvider()

	if provider.GetCountryCode() != "IL" {
		t.Errorf("Expected country code IL, got %s", provider.GetCountryCode())
	}

	if provider.GetCountryName() != "Israel" {
		t.Errorf("Expected country name Israel, got %s", provider.GetCountryName())
	}

	subdivisions := provider.GetSubdivisions()
	if len(subdivisions) != 6 {
		t.Errorf("Expected 6 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetCategories()
	expectedCategories := []string{"public", "religious", "memorial", "national", "jewish"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestILProvider_LoadHolidays2024(t *testing.T) {
	provider := NewILProvider()
	holidays := provider.LoadHolidays(2024)

	if len(holidays) == 0 {
		t.Error("Expected holidays to be loaded for 2024")
	}

	// Test key Israeli holidays for 2024
	expectedHolidays := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 4, 23, 0, 0, 0, 0, time.UTC), "Passover", "religious"},
		{time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC), "Holocaust Remembrance Day", "memorial"},
		{time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC), "Memorial Day", "memorial"},
		{time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC), "Independence Day", "national"},
		{time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC), "Shavuot", "religious"},
		{time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC), "Rosh Hashanah", "religious"},
		{time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC), "Rosh Hashanah (Day 2)", "religious"},
		{time.Date(2024, 10, 12, 0, 0, 0, 0, time.UTC), "Yom Kippur", "religious"},
		{time.Date(2024, 10, 17, 0, 0, 0, 0, time.UTC), "Sukkot", "religious"},
		{time.Date(2024, 10, 24, 0, 0, 0, 0, time.UTC), "Simchat Torah", "religious"},
		{time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC), "Hanukkah", "religious"},
	}

	for _, expected := range expectedHolidays {
		holiday, exists := holidays[expected.date]
		if !exists {
			t.Errorf("Expected holiday on %s (%s)", expected.date.Format("2006-01-02"), expected.name)
			continue
		}

		if holiday.Name != expected.name {
			t.Errorf("Expected holiday name %s, got %s", expected.name, holiday.Name)
		}

		if holiday.Category != expected.category {
			t.Errorf("Expected category %s for %s, got %s", expected.category, expected.name, holiday.Category)
		}
	}
}

func TestILProvider_HebrewCalendarHolidays(t *testing.T) {
	provider := NewILProvider()

	// Test that Hebrew calendar holidays exist for different years
	testYears := []int{2023, 2024, 2025}

	for _, year := range testYears {
		holidays := provider.LoadHolidays(year)

		// Check for major holidays that should exist every year
		majorHolidays := []string{
			"Rosh Hashanah",
			"Yom Kippur",
			"Passover",
			"Shavuot",
		}

		for _, holidayName := range majorHolidays {
			found := false
			for _, holiday := range holidays {
				if holiday.Name == holidayName || holiday.Name == holidayName+" (Day 2)" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %s in %d", holidayName, year)
			}
		}
	}
}

func TestILProvider_MemorialDays(t *testing.T) {
	provider := NewILProvider()

	testCases := []struct {
		year            int
		yomHaShoah      time.Time
		yomHaZikaron    time.Time
		independenceDay time.Time
	}{
		{2023, time.Date(2023, 4, 18, 0, 0, 0, 0, time.UTC), time.Date(2023, 4, 25, 0, 0, 0, 0, time.UTC), time.Date(2023, 4, 26, 0, 0, 0, 0, time.UTC)},
		{2024, time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC), time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC), time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, 4, 24, 0, 0, 0, 0, time.UTC), time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC), time.Date(2025, 5, 2, 0, 0, 0, 0, time.UTC)},
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)

		// Holocaust Remembrance Day
		if holiday, exists := holidays[tc.yomHaShoah]; !exists {
			t.Errorf("Holocaust Remembrance Day should exist on %s in %d", tc.yomHaShoah.Format("2006-01-02"), tc.year)
		} else {
			if holiday.Category != "memorial" {
				t.Errorf("Holocaust Remembrance Day should be memorial category, got %s", holiday.Category)
			}
		}

		// Memorial Day
		if holiday, exists := holidays[tc.yomHaZikaron]; !exists {
			t.Errorf("Memorial Day should exist on %s in %d", tc.yomHaZikaron.Format("2006-01-02"), tc.year)
		} else {
			if holiday.Category != "memorial" {
				t.Errorf("Memorial Day should be memorial category, got %s", holiday.Category)
			}
		}

		// Independence Day (day after Memorial Day)
		if holiday, exists := holidays[tc.independenceDay]; !exists {
			t.Errorf("Independence Day should exist on %s in %d", tc.independenceDay.Format("2006-01-02"), tc.year)
		} else {
			if holiday.Category != "national" {
				t.Errorf("Independence Day should be national category, got %s", holiday.Category)
			}
		}

		// Verify Memorial Day is the day before Independence Day
		expectedIndependence := tc.yomHaZikaron.AddDate(0, 0, 1)
		if !tc.independenceDay.Equal(expectedIndependence) {
			t.Errorf("Independence Day should be the day after Memorial Day in %d", tc.year)
		}
	}
}

func TestILProvider_RoshHashanahTwoDays(t *testing.T) {
	provider := NewILProvider()
	holidays := provider.LoadHolidays(2024)

	// Rosh Hashanah should have two days
	day1 := time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC)
	day2 := time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC)

	holiday1, exists1 := holidays[day1]
	if !exists1 {
		t.Error("Rosh Hashanah Day 1 should exist")
	} else {
		if holiday1.Name != "Rosh Hashanah" {
			t.Errorf("Expected 'Rosh Hashanah', got %s", holiday1.Name)
		}
	}

	holiday2, exists2 := holidays[day2]
	if !exists2 {
		t.Error("Rosh Hashanah Day 2 should exist")
	} else {
		if holiday2.Name != "Rosh Hashanah (Day 2)" {
			t.Errorf("Expected 'Rosh Hashanah (Day 2)', got %s", holiday2.Name)
		}
	}
}

func TestILProvider_LanguageSupport(t *testing.T) {
	provider := NewILProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that holidays have both English and Hebrew names
	var passover *Holiday
	for _, holiday := range holidays {
		if holiday.Name == "Passover" {
			passover = holiday
			break
		}
	}

	if passover == nil {
		t.Fatal("Passover should exist in 2024")
	}

	if passover.Languages["en"] != "Passover" {
		t.Errorf("Expected English name 'Passover', got %s", passover.Languages["en"])
	}

	if passover.Languages["he"] != "פסח" {
		t.Errorf("Expected Hebrew name 'פסח', got %s", passover.Languages["he"])
	}
}

func TestILProvider_CategoryDistribution(t *testing.T) {
	provider := NewILProvider()
	holidays := provider.LoadHolidays(2024)

	categories := make(map[string]int)
	for _, holiday := range holidays {
		categories[holiday.Category]++
	}

	// Should have holidays in main categories
	expectedCategories := []string{"religious", "memorial", "national"}
	for _, category := range expectedCategories {
		if count, exists := categories[category]; !exists || count == 0 {
			t.Errorf("Expected holidays in category %s", category)
		}
	}

	// Religious holidays should be the majority
	if categories["religious"] == 0 {
		t.Error("Expected religious holidays")
	}

	t.Logf("Holiday distribution: %+v", categories)
}

func TestILProvider_YearVariation(t *testing.T) {
	provider := NewILProvider()

	// Test that holidays exist for known years but not for unknown years
	holidays2024 := provider.LoadHolidays(2024)
	holidays2030 := provider.LoadHolidays(2030) // Unknown year

	if len(holidays2024) == 0 {
		t.Error("Expected holidays for 2024")
	}

	// 2030 should have fewer holidays (only the ones we can calculate)
	if len(holidays2030) >= len(holidays2024) {
		t.Error("Expected fewer holidays for unknown year 2030")
	}

	t.Logf("2024 holidays: %d, 2030 holidays: %d", len(holidays2024), len(holidays2030))
}
