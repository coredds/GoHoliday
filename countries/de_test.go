package countries

import (
	"testing"
	"time"
)

func TestDEProvider_BasicHolidays(t *testing.T) {
	provider := NewDEProvider()
	holidays := provider.LoadHolidays(2024)

	// Test New Year's Day
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[newYear]; !exists {
		t.Error("Neujahr should exist")
	} else {
		if holiday.Name != "Neujahr" {
			t.Errorf("Expected 'Neujahr', got '%s'", holiday.Name)
		}
		if holiday.Category != "public" {
			t.Errorf("Expected category 'public', got '%s'", holiday.Category)
		}
		if holiday.Languages["de"] != "Neujahr" {
			t.Errorf("Expected German name 'Neujahr', got '%s'", holiday.Languages["de"])
		}
		if holiday.Languages["en"] != "New Year's Day" {
			t.Errorf("Expected English name 'New Year's Day', got '%s'", holiday.Languages["en"])
		}
	}

	// Test Labour Day
	labourDay := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[labourDay]; !exists {
		t.Error("Tag der Arbeit should exist")
	} else {
		if holiday.Name != "Tag der Arbeit" {
			t.Errorf("Expected 'Tag der Arbeit', got '%s'", holiday.Name)
		}
	}

	// Test Christmas Day
	christmas := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[christmas]; !exists {
		t.Error("1. Weihnachtsfeiertag should exist")
	} else {
		if holiday.Name != "1. Weihnachtsfeiertag" {
			t.Errorf("Expected '1. Weihnachtsfeiertag', got '%s'", holiday.Name)
		}
	}

	// Test Boxing Day
	boxingDay := time.Date(2024, 12, 26, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[boxingDay]; !exists {
		t.Error("2. Weihnachtsfeiertag should exist")
	} else {
		if holiday.Name != "2. Weihnachtsfeiertag" {
			t.Errorf("Expected '2. Weihnachtsfeiertag', got '%s'", holiday.Name)
		}
	}
}

func TestDEProvider_EasterHolidays(t *testing.T) {
	provider := NewDEProvider()
	holidays := provider.LoadHolidays(2024)

	// Easter Sunday 2024 is March 31
	easter := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

	// Test Good Friday (March 29, 2024)
	goodFriday := time.Date(2024, 3, 29, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[goodFriday]; !exists {
		t.Error("Karfreitag should exist")
	} else {
		if holiday.Name != "Karfreitag" {
			t.Errorf("Expected 'Karfreitag', got '%s'", holiday.Name)
		}
		if holiday.Category != "public" {
			t.Errorf("Expected category 'public', got '%s'", holiday.Category)
		}
	}

	// Test Easter Sunday
	if holiday, exists := holidays[easter]; !exists {
		t.Error("Ostersonntag should exist")
	} else {
		if holiday.Name != "Ostersonntag" {
			t.Errorf("Expected 'Ostersonntag', got '%s'", holiday.Name)
		}
		if holiday.Category != "religious" {
			t.Errorf("Expected category 'religious', got '%s'", holiday.Category)
		}
	}

	// Test Easter Monday (April 1, 2024)
	easterMonday := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[easterMonday]; !exists {
		t.Error("Ostermontag should exist")
	} else {
		if holiday.Name != "Ostermontag" {
			t.Errorf("Expected 'Ostermontag', got '%s'", holiday.Name)
		}
	}

	// Test Ascension Day (May 9, 2024 - 39 days after Easter)
	ascension := time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[ascension]; !exists {
		t.Error("Christi Himmelfahrt should exist")
	} else {
		if holiday.Name != "Christi Himmelfahrt" {
			t.Errorf("Expected 'Christi Himmelfahrt', got '%s'", holiday.Name)
		}
	}

	// Test Whit Monday (May 20, 2024 - 50 days after Easter)
	whitMonday := time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[whitMonday]; !exists {
		t.Error("Pfingstmontag should exist")
	} else {
		if holiday.Name != "Pfingstmontag" {
			t.Errorf("Expected 'Pfingstmontag', got '%s'", holiday.Name)
		}
	}

	// Test Corpus Christi (May 30, 2024 - 60 days after Easter)
	corpusChristi := time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[corpusChristi]; !exists {
		t.Error("Fronleichnam should exist")
	} else {
		if holiday.Name != "Fronleichnam" {
			t.Errorf("Expected 'Fronleichnam', got '%s'", holiday.Name)
		}
	}
}

func TestDEProvider_GermanUnityDay(t *testing.T) {
	provider := NewDEProvider()

	// Test that German Unity Day doesn't exist before 1990
	holidays1989 := provider.LoadHolidays(1989)
	unityDay1989 := time.Date(1989, 10, 3, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays1989[unityDay1989]; exists {
		t.Error("German Unity Day should not exist in 1989")
	}

	// Test that German Unity Day exists from 1990 onwards
	holidays1990 := provider.LoadHolidays(1990)
	unityDay1990 := time.Date(1990, 10, 3, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays1990[unityDay1990]; !exists {
		t.Error("German Unity Day should exist in 1990")
	} else {
		if holiday.Name != "Tag der Deutschen Einheit" {
			t.Errorf("Expected 'Tag der Deutschen Einheit', got '%s'", holiday.Name)
		}
		if holiday.Category != "public" {
			t.Errorf("Expected category 'public', got '%s'", holiday.Category)
		}
	}

	// Test 2024 has German Unity Day
	holidays2024 := provider.LoadHolidays(2024)
	unityDay2024 := time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC)
	if _, exists := holidays2024[unityDay2024]; !exists {
		t.Error("German Unity Day should exist in 2024")
	}
}

func TestDEProvider_RegionalHolidays(t *testing.T) {
	provider := NewDEProvider()

	// Test Bavaria (BY) - Assumption of Mary
	byHolidays := provider.GetRegionalHolidays(2024, []string{"BY"})
	assumption := time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC)
	if holiday, exists := byHolidays[assumption]; !exists {
		t.Error("Mariä Himmelfahrt should exist for Bavaria")
	} else {
		if holiday.Name != "Mariä Himmelfahrt" {
			t.Errorf("Expected 'Mariä Himmelfahrt', got '%s'", holiday.Name)
		}
		if holiday.Category != "religious" {
			t.Errorf("Expected category 'religious', got '%s'", holiday.Category)
		}
	}

	// Test Saxony (SN) - Reformation Day
	snHolidays := provider.GetRegionalHolidays(2024, []string{"SN"})
	reformation := time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC)
	if holiday, exists := snHolidays[reformation]; !exists {
		t.Error("Reformationstag should exist for Saxony")
	} else {
		if holiday.Name != "Reformationstag" {
			t.Errorf("Expected 'Reformationstag', got '%s'", holiday.Name)
		}
	}

	// Test Saxony (SN) - Repentance and Prayer Day
	// Find the correct date for Buß- und Bettag 2024 (November 20, 2024)
	repentance := time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC)
	if holiday, exists := snHolidays[repentance]; !exists {
		t.Error("Buß- und Bettag should exist for Saxony")
	} else {
		if holiday.Name != "Buß- und Bettag" {
			t.Errorf("Expected 'Buß- und Bettag', got '%s'", holiday.Name)
		}
	}
}

func TestDEProvider_SpecialObservances(t *testing.T) {
	provider := NewDEProvider()
	observances := provider.GetSpecialObservances(2024)

	// Easter 2024 is March 31, so Carnival Monday is February 12 (48 days before)
	carnivalMonday := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	if holiday, exists := observances[carnivalMonday]; !exists {
		t.Error("Rosenmontag should exist")
	} else {
		if holiday.Name != "Rosenmontag" {
			t.Errorf("Expected 'Rosenmontag', got '%s'", holiday.Name)
		}
		if holiday.Category != "regional" {
			t.Errorf("Expected category 'regional', got '%s'", holiday.Category)
		}
	}

	// Ash Wednesday is February 14 (46 days before Easter)
	ashWednesday := time.Date(2024, 2, 14, 0, 0, 0, 0, time.UTC)
	if holiday, exists := observances[ashWednesday]; !exists {
		t.Error("Aschermittwoch should exist")
	} else {
		if holiday.Name != "Aschermittwoch" {
			t.Errorf("Expected 'Aschermittwoch', got '%s'", holiday.Name)
		}
	}

	// Liberation Day - May 8
	liberation := time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC)
	if holiday, exists := observances[liberation]; !exists {
		t.Error("Tag der Befreiung should exist")
	} else {
		if holiday.Name != "Tag der Befreiung" {
			t.Errorf("Expected 'Tag der Befreiung', got '%s'", holiday.Name)
		}
	}
}

func TestDEProvider_RepentanceDay(t *testing.T) {
	provider := NewDEProvider()

	// Test Repentance Day calculation for different years
	testCases := []struct {
		year         int
		expectedDate time.Time
	}{
		{2024, time.Date(2024, 11, 20, 0, 0, 0, 0, time.UTC)}, // Wednesday before Nov 23
		{2025, time.Date(2025, 11, 19, 0, 0, 0, 0, time.UTC)}, // Wednesday before Nov 23
		{2026, time.Date(2026, 11, 18, 0, 0, 0, 0, time.UTC)}, // Wednesday before Nov 23
	}

	for _, tc := range testCases {
		calculated := provider.getRepentanceDay(tc.year)
		if !calculated.Equal(tc.expectedDate) {
			t.Errorf("Year %d: Expected Repentance Day %v, got %v", 
				tc.year, tc.expectedDate, calculated)
		}
		
		// Verify it's a Wednesday
		if calculated.Weekday() != time.Wednesday {
			t.Errorf("Year %d: Repentance Day should be a Wednesday, got %v", 
				tc.year, calculated.Weekday())
		}
	}
}

func TestDEProvider_ProviderInfo(t *testing.T) {
	provider := NewDEProvider()

	// Test country code
	if provider.GetCountryCode() != "DE" {
		t.Errorf("Expected country code 'DE', got '%s'", provider.GetCountryCode())
	}

	// Test subdivisions (16 German states)
	subdivisions := provider.GetSupportedSubdivisions()
	expectedCount := 16
	if len(subdivisions) != expectedCount {
		t.Errorf("Expected %d subdivisions, got %d", expectedCount, len(subdivisions))
	}

	// Test categories
	expectedCategories := []string{"public", "religious", "regional"}
	categories := provider.GetSupportedCategories()
	if len(categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(categories))
	}
}

func TestDEProvider_BilingualSupport(t *testing.T) {
	provider := NewDEProvider()
	holidays := provider.LoadHolidays(2024)

	// Test that holidays have both German and English names
	newYear := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[newYear]; exists {
		if holiday.Languages["de"] != "Neujahr" {
			t.Errorf("Expected German name 'Neujahr', got '%s'", holiday.Languages["de"])
		}
		if holiday.Languages["en"] != "New Year's Day" {
			t.Errorf("Expected English name 'New Year's Day', got '%s'", holiday.Languages["en"])
		}
	}

	labourDay := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	if holiday, exists := holidays[labourDay]; exists {
		if holiday.Languages["de"] != "Tag der Arbeit" {
			t.Errorf("Expected German name 'Tag der Arbeit', got '%s'", holiday.Languages["de"])
		}
		if holiday.Languages["en"] != "Labour Day" {
			t.Errorf("Expected English name 'Labour Day', got '%s'", holiday.Languages["en"])
		}
	}
}
