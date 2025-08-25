package countries

import (
	"testing"
	"time"
)

func TestJPProvider_LoadHolidays(t *testing.T) {
	provider := NewJPProvider()

	// Test 2024 holidays
	holidays := provider.LoadHolidays(2024)

	if len(holidays) == 0 {
		t.Fatal("Expected holidays for Japan, got none")
	}

	// Check some key holidays for 2024
	expectedHolidays := map[string]string{
		"2024-01-01": "New Year's Day",
		"2024-02-11": "National Foundation Day",
		"2024-02-23": "Emperor's Birthday", // Emperor Naruhito's birthday
		"2024-05-03": "Constitution Memorial Day",
		"2024-05-04": "Greenery Day",
		"2024-05-05": "Children's Day",
		"2024-11-03": "Culture Day",
		"2024-11-23": "Labor Thanksgiving Day",
	}

	for dateStr, expectedName := range expectedHolidays {
		date, _ := time.Parse("2006-01-02", dateStr)
		holiday, exists := holidays[date]
		if !exists {
			t.Errorf("Expected holiday on %s, but not found", dateStr)
			continue
		}
		if holiday.Name != expectedName {
			t.Errorf("Expected holiday name '%s' on %s, got '%s'", expectedName, dateStr, holiday.Name)
		}
	}

	// Check that Coming of Age Day is second Monday of January 2024 (Jan 8)
	comingOfAge, _ := time.Parse("2006-01-02", "2024-01-08")
	if holiday, exists := holidays[comingOfAge]; !exists {
		t.Error("Expected Coming of Age Day on January 8, 2024")
	} else if holiday.Name != "Coming of Age Day" {
		t.Errorf("Expected 'Coming of Age Day', got '%s'", holiday.Name)
	}

	// Check that Marine Day is third Monday of July 2024 (July 15)
	marineDay, _ := time.Parse("2006-01-02", "2024-07-15")
	if holiday, exists := holidays[marineDay]; !exists {
		t.Error("Expected Marine Day on July 15, 2024")
	} else if holiday.Name != "Marine Day" {
		t.Errorf("Expected 'Marine Day', got '%s'", holiday.Name)
	}
}

func TestJPProvider_GetCountryCode(t *testing.T) {
	provider := NewJPProvider()
	if provider.GetCountryCode() != "JP" {
		t.Errorf("Expected country code 'JP', got '%s'", provider.GetCountryCode())
	}
}

func TestJPProvider_GetSupportedCategories(t *testing.T) {
	provider := NewJPProvider()
	categories := provider.GetSupportedCategories()
	if len(categories) != 1 || categories[0] != "public" {
		t.Errorf("Expected categories ['public'], got %v", categories)
	}
}

func TestJPProvider_EmperorBirthdayTransition(t *testing.T) {
	provider := NewJPProvider()

	// Test 2019 - Emperor Akihito's birthday (December 23)
	holidays2019 := provider.LoadHolidays(2019)
	akihitoBirthday, _ := time.Parse("2006-01-02", "2019-12-23")
	if holiday, exists := holidays2019[akihitoBirthday]; !exists {
		t.Error("Expected Emperor's Birthday on December 23, 2019")
	} else if holiday.Name != "Emperor's Birthday" {
		t.Errorf("Expected 'Emperor's Birthday', got '%s'", holiday.Name)
	}

	// Test 2020 - Emperor Naruhito's birthday (February 23)
	holidays2020 := provider.LoadHolidays(2020)
	naruhitoBirthday, _ := time.Parse("2006-01-02", "2020-02-23")
	if holiday, exists := holidays2020[naruhitoBirthday]; !exists {
		t.Error("Expected Emperor's Birthday on February 23, 2020")
	} else if holiday.Name != "Emperor's Birthday" {
		t.Errorf("Expected 'Emperor's Birthday', got '%s'", holiday.Name)
	}

	// Make sure December 23, 2020 is NOT Emperor's Birthday
	akihitoBirthday2020, _ := time.Parse("2006-01-02", "2020-12-23")
	if _, exists := holidays2020[akihitoBirthday2020]; exists {
		t.Error("December 23, 2020 should not be Emperor's Birthday")
	}
}
