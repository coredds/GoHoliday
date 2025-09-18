package countries

import (
	"testing"
	"time"
)

func TestCLProvider_BasicInfo(t *testing.T) {
	provider := NewCLProvider()

	if provider.GetCountryCode() != "CL" {
		t.Errorf("Expected country code CL, got %s", provider.GetCountryCode())
	}

	if provider.GetCountryName() != "Chile" {
		t.Errorf("Expected country name Chile, got %s", provider.GetCountryName())
	}

	subdivisions := provider.GetSubdivisions()
	if len(subdivisions) != 16 {
		t.Errorf("Expected 16 subdivisions, got %d", len(subdivisions))
	}

	categories := provider.GetCategories()
	expectedCategories := []string{"public", "religious", "civic", "regional"}
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestCLProvider_LoadHolidays2024(t *testing.T) {
	provider := NewCLProvider()
	holidays := provider.LoadHolidays(2024)

	if len(holidays) == 0 {
		t.Error("Expected holidays to be loaded for 2024")
	}

	// Test some key Chilean holidays for 2024
	expectedHolidays := []struct {
		date     time.Time
		name     string
		category string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "New Year's Day", "public"},
		{time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC), "Good Friday", "religious"},
		{time.Date(2024, 3, 30, 0, 0, 0, 0, time.UTC), "Holy Saturday", "religious"},
		{time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC), "Labour Day", "public"},
		{time.Date(2024, 5, 21, 0, 0, 0, 0, time.UTC), "Navy Day", "civic"},
		{time.Date(2024, 9, 18, 0, 0, 0, 0, time.UTC), "Independence Day", "public"},
		{time.Date(2024, 9, 19, 0, 0, 0, 0, time.UTC), "Army Day", "civic"},
		{time.Date(2024, 10, 12, 0, 0, 0, 0, time.UTC), "Columbus Day", "public"},
		{time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC), "All Saints' Day", "religious"},
		{time.Date(2024, 12, 8, 0, 0, 0, 0, time.UTC), "Immaculate Conception", "religious"},
		{time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Day", "religious"},
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

func TestCLProvider_VariableHolidays(t *testing.T) {
	provider := NewCLProvider()

	// Test variable holidays that move to Monday
	testCases := []struct {
		year         int
		originalDate time.Time
		expectedDate time.Time
		name         string
	}{
		// Saint Peter and Saint Paul in 2024 (Saturday June 29 -> Monday July 1)
		{2024, time.Date(2024, 6, 29, 0, 0, 0, 0, time.UTC), time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), "Saint Peter and Saint Paul"},
		// Assumption of Mary in 2024 (Thursday August 15 -> Monday August 12)
		{2024, time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 12, 0, 0, 0, 0, time.UTC), "Assumption of Mary"},
	}

	for _, tc := range testCases {
		holidays := provider.LoadHolidays(tc.year)
		
		// Check that the holiday exists on the expected (observed) date
		holiday, exists := holidays[tc.expectedDate]
		if !exists {
			t.Errorf("Expected %s to be observed on %s in %d", tc.name, tc.expectedDate.Format("2006-01-02"), tc.year)
			continue
		}

		if holiday.Name != tc.name {
			t.Errorf("Expected holiday name %s, got %s", tc.name, holiday.Name)
		}

		// Check that it doesn't exist on the original date (if different)
		if !tc.originalDate.Equal(tc.expectedDate) {
			if _, exists := holidays[tc.originalDate]; exists {
				t.Errorf("Holiday %s should not exist on original date %s, only on observed date %s", 
					tc.name, tc.originalDate.Format("2006-01-02"), tc.expectedDate.Format("2006-01-02"))
			}
		}
	}
}

func TestCLProvider_RegionalHolidays(t *testing.T) {
	provider := NewCLProvider()

	// Test Battle of Arica (introduced in 2020)
	t.Run("Battle of Arica", func(t *testing.T) {
		// Should not exist before 2020
		holidays2019 := provider.LoadHolidays(2019)
		aricaDate2019 := time.Date(2019, 6, 7, 0, 0, 0, 0, time.UTC)
		if _, exists := holidays2019[aricaDate2019]; exists {
			t.Error("Battle of Arica should not exist in 2019")
		}

		// Should exist from 2020 onwards
		holidays2020 := provider.LoadHolidays(2020)
		aricaDate2020 := time.Date(2020, 6, 7, 0, 0, 0, 0, time.UTC)
		holiday, exists := holidays2020[aricaDate2020]
		if !exists {
			t.Error("Battle of Arica should exist in 2020")
		} else {
			if holiday.Name != "Battle of Arica" {
				t.Errorf("Expected 'Battle of Arica', got %s", holiday.Name)
			}
			if holiday.Category != "regional" {
				t.Errorf("Expected category 'regional', got %s", holiday.Category)
			}
		}
	})

	// Test Chillán Foundation Day (introduced in 2019)
	t.Run("Chillán Foundation Day", func(t *testing.T) {
		// Should not exist before 2019
		holidays2018 := provider.LoadHolidays(2018)
		chillanDate2018 := time.Date(2018, 8, 20, 0, 0, 0, 0, time.UTC)
		if _, exists := holidays2018[chillanDate2018]; exists {
			t.Error("Chillán Foundation Day should not exist in 2018")
		}

		// Should exist from 2019 onwards
		holidays2019 := provider.LoadHolidays(2019)
		chillanDate2019 := time.Date(2019, 8, 20, 0, 0, 0, 0, time.UTC)
		holiday, exists := holidays2019[chillanDate2019]
		if !exists {
			t.Error("Chillán Foundation Day should exist in 2019")
		} else {
			if holiday.Name != "Chillán Foundation Day" {
				t.Errorf("Expected 'Chillán Foundation Day', got %s", holiday.Name)
			}
			if holiday.Category != "regional" {
				t.Errorf("Expected category 'regional', got %s", holiday.Category)
			}
		}
	})
}

func TestCLProvider_EasterBasedHolidays(t *testing.T) {
	provider := NewCLProvider()

	// Test Easter-based holidays for different years
	testYears := []int{2023, 2024, 2025}

	for _, year := range testYears {
		holidays := provider.LoadHolidays(year)
		easter := EasterSunday(year)

		// Good Friday (2 days before Easter)
		goodFriday := easter.AddDate(0, 0, -2)
		if holiday, exists := holidays[goodFriday]; !exists {
			t.Errorf("Good Friday should exist in %d", year)
		} else {
			if holiday.Name != "Good Friday" {
				t.Errorf("Expected 'Good Friday', got %s", holiday.Name)
			}
			if holiday.Category != "religious" {
				t.Errorf("Expected category 'religious', got %s", holiday.Category)
			}
		}

		// Holy Saturday (1 day before Easter)
		holySaturday := easter.AddDate(0, 0, -1)
		if holiday, exists := holidays[holySaturday]; !exists {
			t.Errorf("Holy Saturday should exist in %d", year)
		} else {
			if holiday.Name != "Holy Saturday" {
				t.Errorf("Expected 'Holy Saturday', got %s", holiday.Name)
			}
		}
	}
}

func TestCLProvider_LanguageSupport(t *testing.T) {
	provider := NewCLProvider()
	holidays := provider.LoadHolidays(2024)

	// Check that holidays have both English and Spanish names
	newYears := holidays[time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)]
	if newYears == nil {
		t.Fatal("New Year's Day should exist")
	}

	if newYears.Languages["en"] != "New Year's Day" {
		t.Errorf("Expected English name 'New Year's Day', got %s", newYears.Languages["en"])
	}

	if newYears.Languages["es"] != "Año Nuevo" {
		t.Errorf("Expected Spanish name 'Año Nuevo', got %s", newYears.Languages["es"])
	}
}

func TestCLProvider_CategoryDistribution(t *testing.T) {
	provider := NewCLProvider()
	holidays := provider.LoadHolidays(2024)

	categories := make(map[string]int)
	for _, holiday := range holidays {
		categories[holiday.Category]++
	}

	// Should have holidays in all main categories
	expectedCategories := []string{"public", "religious", "civic"}
	for _, category := range expectedCategories {
		if count, exists := categories[category]; !exists || count == 0 {
			t.Errorf("Expected holidays in category %s", category)
		}
	}

	t.Logf("Holiday distribution: %+v", categories)
}
